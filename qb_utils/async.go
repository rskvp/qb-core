package qb_utils

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rskvp/qb-core/qb_stopwatch"
)

type AsyncHelper struct {
}

var Async *AsyncHelper

func init() {
	Async = new(AsyncHelper)
}

func (instance *AsyncHelper) NewAsyncTask() *AsyncTask {
	task := &AsyncTask{
		watch:   qb_stopwatch.New(),
		timeout: 0,
	}
	return task
}

func (instance *AsyncHelper) NewAsyncTimedTask(timeout time.Duration, handler func()) *AsyncTask {
	task := &AsyncTask{
		watch:          qb_stopwatch.New(),
		timeout:        timeout,
		handlerTimeout: handler,
	}
	return task
}

func (instance *AsyncHelper) NewConcurrentPool(limit int) *ConcurrentPool {
	if limit < 1 {
		limit = DefaultMaxConcurrent
	}

	// allocate a limiter instance
	pool := new(ConcurrentPool)
	pool.limit = limit
	pool.tickets = make(chan int, limit) // buffered channel with a limit
	pool.numOfRunningRoutines = 0

	// allocate the tickets:
	for i := 0; i < pool.limit; i++ {
		pool.tickets <- i
	}

	return pool
}

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const (
	// DefaultMaxConcurrent is max number of concurrent routines
	DefaultMaxConcurrent = 100
)

type GoRoutineWrapper func() error
type GoRoutineWrapperArgs func(args ...interface{}) error

//----------------------------------------------------------------------------------------------------------------------
//	ConcurrentPool
//----------------------------------------------------------------------------------------------------------------------

type ConcurrentPool struct {
	limit                int
	tickets              chan int
	numOfRunningRoutines int32
	wg                   sync.WaitGroup
	errMux               sync.Mutex
	err                  *Error
}

func (instance *ConcurrentPool) Limit() int {
	if nil != instance {
		return instance.limit
	}
	return 0
}

func (instance *ConcurrentPool) Count() int32 {
	if nil != instance {
		return atomic.LoadInt32(&instance.numOfRunningRoutines)
	}
	return 0
}

func (instance *ConcurrentPool) Run(f GoRoutineWrapper) int {
	return instance.RunArgs(func(args ...interface{}) error {
		return f()
	})
}

// RunArgs Execute adds a function to the execution queue.
// If num of go routines allocated by this instance is < limit
// launch a new go routine to execute job
// else wait until a go routine becomes available
func (instance *ConcurrentPool) RunArgs(f GoRoutineWrapperArgs, args ...interface{}) int {
	if nil != instance {
		// pop a ticket
		ticket := <-instance.tickets
		atomic.AddInt32(&instance.numOfRunningRoutines, 1)

		instance.wg.Add(1)

		go func() {

			defer func() {
				// push a ticket
				instance.tickets <- ticket
				atomic.AddInt32(&instance.numOfRunningRoutines, -1)
			}()

			// run the job
			defer instance.wg.Done()
			if err := f(args...); err != nil {
				instance.errMux.Lock()
				instance.err = Errors.Append(instance.err, err)
				instance.errMux.Unlock()
			}
		}()

		return ticket
	}
	return -1
}

// Wait all jobs are executed, if any in queue
func (instance *ConcurrentPool) Wait() *Error {
	if nil != instance {
		instance.wg.Wait()
		instance.errMux.Lock()
		defer instance.errMux.Unlock()
		return instance.err
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	AsyncContext
//----------------------------------------------------------------------------------------------------------------------

type AsyncContext struct {
	ctx    context.Context
	stop   context.CancelFunc
	closed bool
}

func NewAsyncContext() *AsyncContext {
	instance := new(AsyncContext)
	instance.ctx, instance.stop = context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-instance.ctx.Done(): // closes when the caller cancels the ctx
				instance.closed = true
				break
			}
		}
	}()
	return instance
}

func (instance *AsyncContext) IsCancelled() bool {
	if nil != instance {
		return instance.closed
	}
	return true
}

func (instance *AsyncContext) Done() <-chan struct{} {
	if nil != instance {
		return instance.ctx.Done()
	}
	return nil
}

func (instance *AsyncContext) Cancel() {
	if nil != instance && nil != instance.stop {
		instance.stop()
	}
}

func (instance *AsyncContext) Err() error {
	if nil != instance && nil != instance.ctx {
		return instance.ctx.Err()
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	AsyncTask
//----------------------------------------------------------------------------------------------------------------------

type AsyncTaskGoRoutine func(ctx *AsyncContext, args ...interface{}) (interface{}, error)

type AsyncTask struct {
	wg             *sync.WaitGroup
	err            error
	response       interface{}
	watch          *qb_stopwatch.StopWatch
	handlerTimeout func()
	handlerSuccess func(response interface{})
	handlerError   func(err error)
	handlerFinish  func(response interface{}, err error)
	timeout        time.Duration
	timer          *time.Timer
	finished       bool
}

func (instance *AsyncTask) SetTimeout(timeout time.Duration) *AsyncTask {
	if nil != instance {
		instance.timeout = timeout
	}
	return instance
}

func (instance *AsyncTask) OnTimeout(callback func()) *AsyncTask {
	if nil != instance {
		instance.handlerTimeout = callback
	}
	return instance
}

func (instance *AsyncTask) OnSuccess(callback func(response interface{})) *AsyncTask {
	if nil != instance {
		instance.handlerSuccess = callback
	}
	return instance
}

func (instance *AsyncTask) OnError(callback func(err error)) *AsyncTask {
	if nil != instance {
		instance.handlerError = callback
	}
	return instance
}

func (instance *AsyncTask) OnFinish(callback func(response interface{}, err error)) *AsyncTask {
	if nil != instance {
		instance.handlerFinish = callback
	}
	return instance
}

func (instance *AsyncTask) ElapsedMs() int {
	if nil != instance && nil != instance.watch {
		return instance.watch.Milliseconds()
	}
	return 0
}

func (instance *AsyncTask) Response() interface{} {
	if nil != instance {
		return instance.response
	}
	return nil
}

func (instance *AsyncTask) ResponseString() string {
	if nil != instance {
		return Convert.ToString(instance.response)
	}
	return ""
}

func (instance *AsyncTask) ResponseBool() bool {
	if nil != instance {
		return Convert.ToBool(instance.response)
	}
	return false
}

func (instance *AsyncTask) ResponseInt() int {
	if nil != instance {
		return Convert.ToInt(instance.response)
	}
	return -1
}

func (instance *AsyncTask) ResponseInt8() int8 {
	if nil != instance {
		return Convert.ToInt8(instance.response)
	}
	return -1
}

func (instance *AsyncTask) ResponseInt32() int32 {
	if nil != instance {
		return Convert.ToInt32(instance.response)
	}
	return -1
}

func (instance *AsyncTask) ResponseInt64() int64 {
	if nil != instance {
		return Convert.ToInt64(instance.response)
	}
	return -1
}

func (instance *AsyncTask) ResponseFloat32() float32 {
	if nil != instance {
		return Convert.ToFloat32(instance.response)
	}
	return -1
}

func (instance *AsyncTask) ResponseFloat64() float64 {
	if nil != instance {
		return Convert.ToFloat64(instance.response)
	}
	return -1
}

func (instance *AsyncTask) Error() error {
	if nil != instance {
		return instance.err
	}
	return nil
}

func (instance *AsyncTask) RunSync(f AsyncTaskGoRoutine, args ...interface{}) (elapsed int, response interface{}, err error) {
	if nil != instance && nil != f {
		instance.Run(f, args...)
		elapsed, response, err = instance.Wait()
	}
	return
}

func (instance *AsyncTask) Run(f AsyncTaskGoRoutine, args ...interface{}) *AsyncTask {
	if nil != instance && nil != f {
		ctx := NewAsyncContext()
		instance.wg = new(sync.WaitGroup)
		instance.wg.Add(1)
		if instance.timeout > 0 {
			instance.timer = time.NewTimer(instance.timeout)
		}
		// run
		go instance.run(ctx, f, args...)
		// wait timer
		go instance.checkTimeout(ctx)
	}
	return instance
}

func (instance *AsyncTask) Close() {
	if nil != instance {
		instance.wg.Done()
	}
}

func (instance *AsyncTask) Wait() (elapsed int, response interface{}, err error) {
	elapsed = -1
	if nil != instance {
		instance.wg.Wait()
		// exiting process
		elapsed = instance.doFinish()
		response = instance.response
		err = instance.err
	}
	return
}

func (instance *AsyncTask) run(ctx *AsyncContext, f AsyncTaskGoRoutine, args ...interface{}) {
	if r := recover(); r != nil {
		// recovered from panic
		instance.err = errors.New(fmt.Sprintf("%v", r))
	}
	response, err := f(ctx, args...)
	instance.response = response
	if nil == instance.err {
		instance.err = err
	}
	if nil != instance.timer {
		instance.timer.Stop()
		instance.timer = nil
	}
	_ = instance.doFinish()
}

func (instance *AsyncTask) checkTimeout(ctx *AsyncContext) {
	if nil != instance.timer {
		<-instance.timer.C
		instance.timer = nil
		ctx.Cancel()
		if nil == instance.handlerTimeout {
			instance.err = errors.New("timeout_error")
		} else {
			instance.handlerTimeout()
		}
	}
}

func (instance *AsyncTask) doFinish() (elapsed int) {
	if nil != instance {
		// handle internal watch
		if nil != instance.watch {
			if instance.watch.IsRunning() {
				instance.watch.Stop()
			}
			elapsed = instance.watch.Milliseconds()
		}
		// finish
		if !instance.finished {
			instance.finished = true
			instance.wg.Done()

			if nil == instance.err {
				if nil != instance.handlerSuccess {
					instance.handlerSuccess(instance.response)
				}
			} else {
				if nil != instance.handlerError {
					instance.handlerError(instance.err)
				}
			}
			if nil != instance.handlerFinish {
				instance.handlerFinish(instance.response, instance.err)
			}
		}
	}
	return
}

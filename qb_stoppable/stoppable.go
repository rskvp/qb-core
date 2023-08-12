package qb_stoppable

import (
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
	"time"

	qbCore "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_"
)

type ShutdownCallback func() error

var stoppableInstances = make([]*Stoppable, 0)
var mux sync.Mutex

// Stoppable object
type Stoppable struct {
	name               string
	index              int
	stopChan           chan bool
	mux                sync.Mutex
	waiting            bool
	onStart            func()
	onStop             func()
	logger             qb_.ILogger
	shutdownOperations map[string]ShutdownCallback
}

func NewStoppable() *Stoppable {
	instance := new(Stoppable)
	instance.name = qbCore.Rnd.Uuid()
	instance.shutdownOperations = make(map[string]ShutdownCallback)
	    qbCore.Sys.OnSignal(instance.onSignal,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGKILL,
		syscall.SIGQUIT,
	)

	// add to internal list
	mux.Lock()
	defer mux.Unlock()
	stoppableInstances = append(stoppableInstances, instance)

	instance.index = len(stoppableInstances)

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *Stoppable) String() string {
	if nil != instance {
		return qbCore.JSON.Stringify(map[string]interface{}{
			"id":      instance.ItemId(),
			"name":    instance.name,
			"actions": instance.OperationsName(),
		})
	}
	return ""
}

func (instance *Stoppable) ItemId() string {
	if nil != instance {
		return fmt.Sprintf("item#%v", instance.index)
	}
	return ""
}

func (instance *Stoppable) SetLogger(logger qb_.ILogger) *Stoppable {
	if nil != instance {
		instance.logger = logger
	}
	return instance
}

func (instance *Stoppable) SetName(name string) *Stoppable {
	if nil != instance {
		instance.name = name
	}
	return instance
}

func (instance *Stoppable) AddStopOperation(name string, callback ShutdownCallback) *Stoppable {
	if nil != instance && nil != instance.shutdownOperations {
		instance.shutdownOperations[name] = callback
	}
	return instance
}

func (instance *Stoppable) OperationsName() []string {
	response := make([]string, 0)
	if nil != instance && nil != instance.shutdownOperations {
		for name, _ := range instance.shutdownOperations {
			response = append(response, name)
		}
	}
	return response
}

func (instance *Stoppable) Start() bool {
	if nil != instance && nil == instance.stopChan {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		instance.stopChan = make(chan bool, 1)
		instance.doStart()

		return true
	}
	return false // not executed
}

func (instance *Stoppable) Stop() bool {
	if nil != instance && nil != instance.stopChan {
		instance.onSignal(syscall.SIGQUIT)
		if nil != instance.logger {
			// wait a little to allow file writing
			time.Sleep(1 * time.Second)
		}
		return true
	}
	return false // not executed
}

func (instance *Stoppable) IsStopped() bool {
	if nil != instance && nil != instance.stopChan {
		return false
	}
	return true
}

func (instance *Stoppable) IsJoined() bool {
	return nil != instance && instance.waiting
}

func (instance *Stoppable) Join() {
	if nil != instance && nil != instance.stopChan && !instance.waiting {

		instance.waiting = true

		// wait exit
		<-instance.stopChan
		// reset channel
		instance.stopChan = nil
		instance.waiting = false
	}
}

func (instance *Stoppable) JoinTimeout(d time.Duration) {
	go func() {
		time.Sleep(d)
		instance.Stop()
	}()
	instance.Join()
}

// ---------------------------------------------------------------------------------------------------------------------
//	e v e n t s
// ---------------------------------------------------------------------------------------------------------------------

func (instance *Stoppable) OnStart(f func()) {
	instance.onStart = f
}

func (instance *Stoppable) OnStop(f func()) {
	instance.onStop = f
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *Stoppable) doStop() {
	if nil != instance && nil != instance.onStop {
		instance.onStop()
	}
}

func (instance *Stoppable) doStart() {
	if nil != instance && nil != instance.onStart {
		instance.onStart()
	}
}

func (instance *Stoppable) onSignal(s os.Signal) {
	if nil != instance {
		// log.Println("onSignal", s)
		logger := instance.logger
		wait := make(chan struct{})

		// intercepted close signal
		if len(instance.shutdownOperations) > 0 {
			msg := fmt.Sprintf("STARTING SHUTDOWN for '%s' from signal '%v'...", instance.ItemId(), s)
			if nil != logger {
				logger.Info(msg)
			} else {
				log.Println(msg)
			}

			// set timeout for the ops to be done to prevent system hang
			timeout := time.Second * 3 // max 3 seconds to shutdown
			timeoutFunc := time.AfterFunc(timeout, func() {
				msg := fmt.Sprintf("[%s] timeout %d ms has been elapsed, force exit",
					instance.ItemId(), timeout.Milliseconds())
				if nil != logger {
					logger.Info(msg)
				} else {
					log.Println(msg)
				}
				os.Exit(0)
			})
			defer timeoutFunc.Stop()

			var wg sync.WaitGroup
			// Do the operations asynchronously to save time
			for key, op := range instance.shutdownOperations {
				wg.Add(1)
				innerOp := op
				innerKey := key
				go func() {
					defer wg.Done()

					msg := fmt.Sprintf("\t[%s] cleaning up: %s", instance.ItemId(), innerKey)
					if nil != logger {
						logger.Info(msg)
					} else {
						log.Println(msg)
					}
					if err := innerOp(); err != nil {
						msg = fmt.Sprintf("\t[%s] %s: clean up failed: %s", instance.ItemId(), innerKey, err.Error())
						if nil != logger {
							logger.Info(msg)
						} else {
							log.Println(msg)
						}
						return
					}

					msg = fmt.Sprintf("\t[%s] %s was shutdown gracefully", instance.ItemId(), innerKey)
					if nil != logger {
						logger.Info(msg)
					} else {
						log.Println(msg)
					}
				}()
			}
			wg.Wait()

			// close wait channel
			close(wait)
			msg = fmt.Sprintf("TERMINATED SHUTDOWN for '%s'.", instance.ItemId())
			if nil != logger {
				logger.Info(msg)
			} else {
				log.Println(msg)
			}

			_ = instance.free()
		}
	}
}

func (instance *Stoppable) free() bool {
	if nil != instance && nil != instance.stopChan {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		instance.stopChan <- true
		instance.stopChan = nil
		instance.doStop()

		return true
	}
	return false // not executed
}

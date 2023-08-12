package shell

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/rskvp/qb-core/qb_events"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_utils"
)

const (
	eventOnOut = "on_out"
	eventOnErr = "on_err"
)

// ShellCmd wrap a command
type ShellCmd struct {
	// properties
	isBackground bool
	command      string
	args         []string
	logger       qb_log.ILogger

	out            bytes.Buffer
	err            bytes.Buffer
	chanCmd        chan bool // command executed
	chanQuit       chan bool // command terminated
	ended          bool
	cancelFunc     context.CancelFunc
	mux            sync.Mutex
	events         *qb_events.Emitter
	lastOutputTime time.Time

	_c   *exec.Cmd
	_dir string
	_pid int
	_err error
}

func NewShellCmd(command string, args ...string) (instance *ShellCmd) {
	instance = new(ShellCmd)
	instance.events = qb_events.Events.NewEmitter()
	instance.command = command
	instance.args = make([]string, 0)
	if len(args) > 0 {
		instance.args = append(instance.args, args...)
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p r o p e r t i e s
//----------------------------------------------------------------------------------------------------------------------

func (instance *ShellCmd) SetLogger(value qb_log.ILogger) *ShellCmd {
	if nil != instance {
		instance.logger = value
	}
	return instance
}

func (instance *ShellCmd) SetBackground(value bool) *ShellCmd {
	if nil != instance {
		instance.isBackground = value
	}
	return instance
}

func (instance *ShellCmd) IsBackground() bool {
	if nil != instance {
		return instance.isBackground
	}
	return false
}

func (instance *ShellCmd) SetCommand(value string) *ShellCmd {
	if nil != instance {
		instance.command = value
	}
	return instance
}

func (instance *ShellCmd) GetCommand() string {
	if nil != instance {
		return instance.command
	}
	return ""
}

func (instance *ShellCmd) SetArgs(values ...string) *ShellCmd {
	if nil != instance {
		instance.args = values
	}
	return instance
}

func (instance *ShellCmd) GetArgs() []string {
	if nil != instance {
		return instance.args
	}
	return make([]string, 0)
}

//----------------------------------------------------------------------------------------------------------------------
//	e v e n t s
//----------------------------------------------------------------------------------------------------------------------

func (instance *ShellCmd) OnOut(callback func(event *qb_events.Event)) {
	if nil != instance && nil != instance.events {
		instance.events.On(eventOnOut, callback)
	}
}

func (instance *ShellCmd) OnErr(callback func(event *qb_events.Event)) {
	if nil != instance && nil != instance.events {
		instance.events.On(eventOnErr, callback)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Pid return current or last pid
func (instance *ShellCmd) Pid() int {
	if nil != instance {
		return instance._pid
	}
	return 0
}

func (instance *ShellCmd) IsRunning() bool {
	if nil != instance && nil != instance._c {
		return !instance.ended
	}
	return false
}

func (instance *ShellCmd) StdOut() string {
	if nil != instance {
		return instance.out.String()
	}
	return ""
}

func (instance *ShellCmd) StdErr() string {
	if nil != instance {
		return instance.err.String()
	}
	return ""
}

func (instance *ShellCmd) OutWasInactiveFor(d time.Duration) bool {
	if nil != instance && !instance.lastOutputTime.IsZero() {
		dd := qb_utils.Dates.Sub(time.Now(), instance.lastOutputTime)
		return dd > d
	}
	return false
}

func (instance *ShellCmd) Run() (err error) {
	if nil != instance {
		// reset
		instance._c = nil
		instance._dir = ""
		instance._pid = 0
		instance._err = nil
		instance.ended = false
		instance.chanCmd = make(chan bool, 1)
		instance.chanQuit = make(chan bool, 1)

		// run async
		go instance.run()
		go instance.checkOutput()

		// wait command run
		<-instance.chanCmd

		// get error from execution
		err = instance._err
	}
	return
}

func (instance *ShellCmd) Wait() (err error) {
	if nil != instance && nil != instance._c {
		if nil != instance.chanQuit {
			// wait command terminated
			<-instance.chanQuit
		}
		err = instance._err
	}
	return
}

func (instance *ShellCmd) Kill() error {
	if nil != instance._c && !instance.ended {
		if nil != instance._c.Process {
			instance.quit()
			err := instance._c.Process.Kill()
			return err
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ShellCmd) cmd() (c *exec.Cmd) {
	if nil != instance && len(instance.command) > 0 {

		if instance.isBackground {
			ctx, cancelFunc := context.WithCancel(context.Background())
			instance.cancelFunc = cancelFunc
			// ctx := context.Background()
			c = exec.CommandContext(ctx, instance.command, instance.args...)
		} else {
			c = exec.Command(instance.command, instance.args...)
		}
	}

	c.Stdout = &instance.out
	c.Stderr = &instance.err

	return
}

func (instance *ShellCmd) checkOutput() {
	var lastOut string
	var lastErr string

	for {
		if nil != instance && nil != instance._c && !instance.ended && nil != instance.events {
			time.Sleep(300 * time.Millisecond)
			out := instance.out.String()
			if lastOut != out {
				diff := strings.Replace(out, lastOut, "", 1)
				if len(diff) > 0 {
					instance.lastOutputTime = time.Now()
					lastOut = out
					instance.events.Emit(eventOnOut, diff)
					instance.log(qb_log.InfoLevel, diff)
				}
			}
			err := instance.err.String()
			if lastErr != err {
				diff := strings.Replace(err, lastErr, "", 1)
				if len(diff) > 0 {
					instance.lastOutputTime = time.Now()
					lastErr = err
					instance.events.Emit(eventOnErr, diff)
					instance.log(qb_log.ErrorLevel, diff)
				}
			}
		}
	}
}

func (instance *ShellCmd) log(level qb_log.Level, text string) {
	if nil != instance && nil != instance.logger {
		rows := strings.Split(text, "\n")
		for _, row := range rows {
			if len(row) > 0 {
				switch level {
				case qb_log.ErrorLevel:
					instance.logger.Error(row)
				case qb_log.InfoLevel:
					instance.logger.Info(row)
				default:
					instance.logger.Debug(row)
				}
			}
		}
	}
}

func (instance *ShellCmd) run() {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	// get the cmd and run
	instance._c = instance.cmd()
	instance._err = instance._c.Start() // async in background

	if nil != instance._err {
		instance._c = nil // clear instance
		return
	}
	instance._dir = instance._c.Dir
	if nil != instance._c.Process {
		instance._pid = instance._c.Process.Pid
	}

	// just the time to let os create child process
	time.Sleep(500 * time.Millisecond)

	// command run
	instance.chanCmd <- true

	// now wait execution is terminated
	instance._err = instance._c.Wait()

	// confirm finish
	instance.quit()

}

func (instance *ShellCmd) quit() {
	if nil != instance {
		// instance.stopWatch.Stop()
		instance.ended = true
		instance.chanQuit <- true
		instance.chanQuit = nil
		// instance.chanOsInterrupt = nil
		if nil != instance.cancelFunc {
			instance.cancelFunc()
		}
	}
}

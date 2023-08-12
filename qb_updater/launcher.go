package qb_updater

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/rskvp/qb-core/qb_events"
	"github.com/rskvp/qb-core/qb_exec"
	"github.com/rskvp/qb-core/qb_utils"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t a n t s
//----------------------------------------------------------------------------------------------------------------------

const (
	onQuit    = "on_quit"
	onStart   = "on_start"
	onStarted = "on_started"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type Launcher struct {
	keepHandle       bool // keep session handler
	cmd              *exec.Cmd
	command          string
	err              error
	out              *bytes.Buffer
	pid              int
	chanCmd          chan bool // command executed
	chanQuit         chan bool // command terminated
	ended            bool
	events           *qb_events.Emitter
	onQuitHandler    func(command string, pid int)
	onStartHandler   func(command string)
	onStartedHandler func(command string, pid int)
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewLauncher(keepHandle bool) *Launcher {
	instance := new(Launcher)
	instance.keepHandle = keepHandle
	instance.ended = true
	instance.pid = -1
	instance.events = qb_events.Events.NewEmitter()
	if nil != instance.events {
		instance.events.On(onQuit, func(event *qb_events.Event) {
			if nil != instance && nil != instance.onQuitHandler {
				instance.onQuitHandler(event.ArgumentAsString(0), event.ArgumentAsInt(1))
			}
		})
		instance.events.On(onStart, func(event *qb_events.Event) {
			if nil != instance && nil != instance.onStartHandler {
				instance.onStartHandler(event.ArgumentAsString(0))
			}
		})
		instance.events.On(onStarted, func(event *qb_events.Event) {
			if nil != instance && nil != instance.onStartedHandler {
				instance.onStartedHandler(event.ArgumentAsString(0), event.ArgumentAsInt(1))
			}
		})
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Launcher) String() string {
	if nil != instance {
		return instance.GoString()
	}
	return ""
}

func (instance *Launcher) GoString() string {
	if nil != instance {
		info := map[string]interface{}{
			"pid": instance.pid,
			"out": instance.Output(),
		}
		return qb_utils.JSON.Stringify(info)
	}
	return ""
}

func (instance *Launcher) Run(command string) error {
	if nil != instance {
		instance.runCommand(command)

		// wait command run
		<-instance.chanCmd

		return instance.err
	}
	return nil
}

// Wait
// wait command terminated
func (instance *Launcher) Wait() error {
	if nil != instance {
		// wait command terminated
		<-instance.chanQuit

		return instance.err
	}
	return nil
}

func (instance *Launcher) Kill() error {
	if nil != instance.cmd && !instance.ended {
		if nil != instance.cmd.Process {
			instance.quit(true)
			err := instance.cmd.Process.Kill()
			return err
		}
	}
	return nil
}

func (instance *Launcher) IsKillable() bool {
	if nil != instance.cmd && !instance.ended {
		if nil != instance.cmd.Process {
			return true
		}
	}
	return false
}

func (instance *Launcher) Pid() int {
	if nil != instance {
		return instance.pid
	}
	return -1
}

func (instance *Launcher) Output() string {
	if nil != instance.cmd && nil != instance.out {
		return instance.out.String()
	}
	return ""
}

func (instance *Launcher) Error() error {
	if nil != instance {
		return instance.err
	}
	return nil
}

func (instance *Launcher) OnQuit(callback func(command string, pid int)) {
	if nil != instance {
		instance.onQuitHandler = callback
	}
}

func (instance *Launcher) OnStart(callback func(command string)) {
	if nil != instance {
		instance.onStartHandler = callback
	}
}

func (instance *Launcher) OnStarted(callback func(command string, pid int)) {
	if nil != instance {
		instance.onStartedHandler = callback
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Launcher) init() {
	if nil != instance {
		instance.pid = -1
		instance.ended = false
		instance.err = nil
		instance.chanCmd = make(chan bool, 1)
		instance.chanQuit = make(chan bool, 1)
		if nil != instance.out {
			instance.out.Reset()
		}
	}
}

func (instance *Launcher) runCommand(rawCommand string) {
	if nil != instance {

		instance.init()

		instance.command = rawCommand

		if len(rawCommand) > 0 {

			tokens := strings.Split(rawCommand, " ")
			command := tokens[0]
			params := make([]string, 0)
			if len(tokens) > 1 {
				params = tokens[1:]
			}

			if instance.keepHandle {
				go instance.start(command, params...)
			} else {
				// run and forget: does not wait response to avoid screen lock
				go instance.runBackground(command, params...)
			}

			return
		}
		instance.chanCmd <- false
		instance.quit(false)
	}
}

func (instance *Launcher) emit(eventName string, args ...interface{}) {
	if nil != instance && nil != instance.events {
		instance.events.Emit(eventName, args...)
	}
}

func (instance *Launcher) runBackground(command string, params ...string) {
	if nil != instance {
		instance.emit(onStart, instance.command)
		_, _ = qb_exec.Exec.RunBackground(command, params...)

		// command run
		instance.chanCmd <- true
		instance.emit(onStarted, instance.command)

		instance.quit(true)
	}
}

func (instance *Launcher) start(command string, params ...string) {
	if nil != instance {
		instance.emit(onStart, instance.command)
		cmd, buff, _ := qb_exec.Exec.Command(command, params...)
		instance.cmd = cmd
		instance.out = buff

		// start
		err := cmd.Start()
		if nil != err {
			instance.err = err
		} else {
			instance.pid = cmd.Process.Pid
		}
		// command run
		instance.chanCmd <- true
		instance.emit(onStarted, instance.command, instance.pid)

		// wait
		if nil == err {
			err = cmd.Wait() // try to wat, but some nohup programs are not "waitable"
			if nil != err {
				instance.err = err
			}
		}

		// notify exit
		instance.quit(true)
	}
}

func (instance *Launcher) quit(value bool) {
	if nil != instance {
		pid := instance.pid
		cmd := instance.command

		instance.pid = -1
		instance.ended = true
		instance.chanQuit <- value
		instance.emit(onQuit, cmd, pid)
	}
}

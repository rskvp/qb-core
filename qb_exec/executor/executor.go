package executor

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/rskvp/qb-core/qb_stopwatch"
	"github.com/rskvp/qb-core/qb_utils"
)

//----------------------------------------------------------------------------------------------------------------------
//	Executor
//----------------------------------------------------------------------------------------------------------------------

type Executor struct {
	execDir         string
	command         string
	cmd             *exec.Cmd
	cancelFunc      context.CancelFunc
	currPid         int
	lastPid         int
	chanCmd         chan bool // command executed
	chanQuit        chan bool // command terminated
	chanOsInterrupt chan os.Signal
	ended           bool
	inputs          []string
	outWriters      []io.Writer
	errWriters      []io.Writer
	stdout          bytes.Buffer
	stderr          bytes.Buffer
	err             error // internal errors
	stopWatch       *qb_stopwatch.StopWatch
}

func NewExecutor(cmd string) *Executor {
	instance := new(Executor)
	instance.command = cmd
	instance.inputs = make([]string, 0)
	instance.ended = true
	instance.currPid = -1
	instance.lastPid = -1
	instance.outWriters = make([]io.Writer, 0)
	instance.errWriters = make([]io.Writer, 0)

	instance.OutWriterAppend(&instance.stdout)
	instance.ErrorWriterAppend(&instance.stderr)

	return instance
}

func NewExecutorWithDir(cmd, execDir string) *Executor {
	instance := NewExecutor(cmd)
	instance.execDir = execDir

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Executor) String() string {
	if nil != instance {
		return instance.GoString()
	}
	return ""
}

func (instance *Executor) GoString() string {
	if nil != instance {
		info := map[string]interface{}{
			"current-pid": instance.currPid,
			"last-pid":    instance.lastPid,
			"out":         instance.StdOut(),
		}
		return qb_utils.JSON.Stringify(info)
	}
	return ""
}

func (instance *Executor) Inputs() []string {
	return instance.inputs
}

func (instance *Executor) InputsAppend(text string) {
	instance.inputs = append(instance.inputs, text)
}

func (instance *Executor) OutWriterAppend(w io.Writer) {
	instance.outWriters = append(instance.outWriters, w)
}

func (instance *Executor) ErrorWriterAppend(w io.Writer) {
	instance.errWriters = append(instance.errWriters, w)
}

func (instance *Executor) Elapsed() int {
	if nil != instance && nil != instance.stopWatch {
		return instance.stopWatch.Milliseconds()
	}
	return -1
}

// Wait wait command terminated
func (instance *Executor) Wait() error {
	if nil != instance {
		if nil != instance.chanQuit {
			// wait command terminated
			<-instance.chanQuit
		}
		return instance.err
	}
	return nil
}

// Close kill executable using os.Interrupt signal
func (instance *Executor) Close() {
	if nil != instance && nil != instance.chanOsInterrupt {
		signal.Notify(instance.chanOsInterrupt, os.Interrupt)
	}
}

func (instance *Executor) Kill() error {
	if nil != instance.cmd && !instance.ended {
		if nil != instance.cmd.Process {
			instance.quit(true)
			err := instance.cmd.Process.Kill()
			return err
		}
	}
	return nil
}

func (instance *Executor) IsKillable() bool {
	if nil != instance.cmd && !instance.ended {
		if nil != instance.cmd.Process {
			return true
		}
	}
	return false
}

func (instance *Executor) PidCurrent() int {
	return instance.currPid
}

func (instance *Executor) PidLatest() int {
	return instance.lastPid
}

func (instance *Executor) StdOut() string {
	if nil != instance && nil != instance.cmd {
		return string(instance.stdout.Bytes())
	}
	return ""
}

func (instance *Executor) StdErr() string {
	if nil != instance && nil != instance.cmd {
		return string(instance.stderr.Bytes())
	}
	return ""
}

func (instance *Executor) StdOutLines() []string {
	output := instance.StdOut()
	return strings.Split(output, "\n")
}

func (instance *Executor) StdOutMap() map[string]interface{} {
	var response map[string]interface{}
	output := strings.TrimSpace(instance.StdOut())
	_ = qb_utils.JSON.Read(output, &response)
	return response
}

func (instance *Executor) Run(args ...string) error {

	instance.ended = false
	instance.chanCmd = make(chan bool, 1)
	instance.chanQuit = make(chan bool, 1)
	instance.chanOsInterrupt = make(chan os.Signal, 1)
	instance.stopWatch = qb_stopwatch.New()
	instance.stopWatch.Start()

	// creates background executor
	ctx, cancelFunc := context.WithCancel(context.Background())
	instance.cancelFunc = cancelFunc
	instance.cmd = exec.CommandContext(ctx, instance.command, args...)
	// instance.cmd.SysProcAttr = &syscall.SysProcAttr{}

	// signal.Notify(instance.chanOsInterrupt, os.Interrupt)
	go func() {
		for sig := range instance.chanOsInterrupt {
			// sig is a ^C, handle it
			if sig == os.Interrupt {
				_ = instance.Kill()
			}
		}
	}()

	go instance.run()

	// wait command run
	<-instance.chanCmd

	return instance.err
}

func (instance *Executor) IsRunning() bool {
	if nil != instance {
		if nil != instance.cmd && nil != instance.cmd.Process && instance.cmd.Process.Pid > 0 {
			return true
		}
	}
	return false
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Executor) run() {

	if len(instance.execDir) > 0 {
		instance.cmd.Dir = instance.execDir
	}

	instance.cmd.Stdout = io.MultiWriter(instance.outWriters...)
	instance.cmd.Stderr = io.MultiWriter(instance.errWriters...)

	// concatenate inputs
	inBuffer := bytes.NewBuffer([]byte(strings.Join(instance.inputs, "\n")))
	instance.cmd.Stdin = io.MultiReader(inBuffer)
	//stdin, pipeerr := instance.cmd.StdinPipe()

	// start
	err := instance.cmd.Start()
	if nil != err {
		instance.err = err
	} else {
		instance.currPid = instance.cmd.Process.Pid
		instance.lastPid = instance.cmd.Process.Pid

		//stdin.Write([]byte(strings.Join(instance.inputs, "\n")))
		// _ = stdin.Close()
	}

	// command run
	instance.chanCmd <- true

	// wait
	if nil == err {
		err := instance.cmd.Wait()
		if nil != err {
			instance.err = err
		}
	}

	// notify exit
	instance.quit(true)
}

func (instance *Executor) quit(value bool) {
	if nil != instance {
		instance.stopWatch.Stop()
		instance.currPid = -1
		instance.ended = true
		instance.chanQuit <- value
		instance.chanQuit = nil
		instance.chanOsInterrupt = nil
		if nil != instance.cancelFunc {
			instance.cancelFunc()
		}
	}
}

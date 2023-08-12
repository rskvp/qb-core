package executor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rskvp/qb-core/qb_rnd"
	"github.com/rskvp/qb-core/qb_sys"
	"github.com/rskvp/qb-core/qb_utils"
)

var (
	NotExistsFileError = errors.New("not_exists_file")
	NotInstalledError  = errors.New("not_installed")
)

type ConsoleProgram struct {
	command  string
	filename string
	dir      string
	args     []string

	inputs     []string
	outWriters []io.Writer
	errWriters []io.Writer
}

func NewConsoleProgram(command string, args ...string) *ConsoleProgram {
	instance := new(ConsoleProgram)
	instance.command = command
	instance.args = args
	instance.inputs = make([]string, 0)
	instance.outWriters = make([]io.Writer, 0)
	instance.errWriters = make([]io.Writer, 0)

	return instance
}

func NewConsoleProgramWithFile(command string, filename string, args ...string) *ConsoleProgram {
	instance := NewConsoleProgram(command, args...)
	if len(filename) > 0 {
		instance.filename = qb_utils.Paths.Absolute(filename)
	}
	return instance
}

func NewConsoleProgramWithDir(command string, dir string, args ...string) *ConsoleProgram {
	instance := NewConsoleProgram(command, args...)
	if len(dir) > 0 {
		instance.dir = qb_utils.Paths.Absolute(dir)
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ConsoleProgram) SetFileName(filename string) {
	instance.filename = filename
}

func (instance *ConsoleProgram) SetDir(dir string) {
	instance.dir = dir
}

func (instance *ConsoleProgram) InputAppend(args ...string) {
	instance.inputs = append(instance.inputs, args...)
}

func (instance *ConsoleProgram) OutWriterAppend(w io.Writer) {
	instance.outWriters = append(instance.outWriters, w)
}

func (instance *ConsoleProgram) ErrorWriterAppend(w io.Writer) {
	instance.errWriters = append(instance.errWriters, w)
}

func (instance *ConsoleProgram) Prepare(args ...string) (session *ConsoleProgramSession, err error) {
	if nil != instance {
		// args are in right order
		runArgs := make([]string, 0)
		runArgs = append(runArgs, instance.args...)
		runArgs = append(runArgs, args...)

		session, err = instance.createSession(runArgs...)
	}
	return
}

func (instance *ConsoleProgram) Run(args ...string) (*ConsoleProgramSession, error) {
	if nil != instance {
		session, err := instance.Prepare(args...)
		if nil == err {
			return session.Run()
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func (instance *ConsoleProgram) RunAsync(args ...string) (*ConsoleProgramSession, error) {
	if nil != instance {
		session, err := instance.Prepare(args...)
		if nil == err {
			return session.RunAsync(args...)
		} else {
			return nil, err
		}
	}
	return nil, nil
}

// RunWrapped run program wrapped into runnable file launcher.
func (instance *ConsoleProgram) RunWrapped(args ...string) (*ConsoleProgramSession, error) {
	if nil != instance {
		filename, err := CreateRunnableFile(instance.dir, instance.command, args...)
		if nil != err {
			return nil, err
		}
		uid := qb_rnd.Rnd.Uuid() // session
		session := NewProgramSession(filename, uid, instance.inputs, instance.outWriters, instance.errWriters)
		session.SetDir(qb_utils.Paths.Dir(filename))
		_, err = session.RunAsync()
		return session, err
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ConsoleProgram) createSession(args ...string) (*ConsoleProgramSession, error) {
	if nil != instance {
		if len(instance.filename) > 0 {
			if b, err := qb_utils.Paths.Exists(instance.filename); !b {
				if nil != err {
					return nil, err
				}
				return nil, NotExistsFileError
			}
		}
		uid := qb_rnd.Rnd.Uuid() // session
		session := NewProgramSession(instance.command, uid, instance.inputs, instance.outWriters, instance.errWriters)
		if len(instance.filename) > 0 {
			session.SetFileName(qb_utils.Paths.Absolute(instance.filename))
		}
		if len(instance.dir) > 0 {
			session.SetDir(qb_utils.Paths.Absolute(instance.dir))
		}

		session.initArgs = args

		return session, nil
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

// CreateRunnableFile creates executable .sh or .bat file to run a command
func CreateRunnableFile(dir string, cmd string, params ...string) (string, error) {
	if len(dir) == 0 {
		dir = "./"
	}
	dir = qb_utils.Paths.Absolute(dir)
	ext := ".sh"
	if qb_sys.Sys.IsWindows() {
		ext = ".bat"
	}
	filename := qb_utils.Paths.Concat(dir, "_runnable"+ext)
	content := fmt.Sprintf("#!/bin/sh\n%v", cmd)
	if len(params) > 0 {
		content += " "
		for i, param := range params {
			if i > 0 {
				content += " "
			}
			content += param
		}
	}
	_, err := qb_utils.IO.WriteTextToFile(content, filename)
	if nil != err {
		return "", err
	}
	err = os.Chmod(filename, 0755)
	return filename, err
}

//----------------------------------------------------------------------------------------------------------------------
//	ConsoleProgramSession
//----------------------------------------------------------------------------------------------------------------------

type ConsoleProgramSession struct {
	command    string
	uid        string
	filename   string
	execDir    string
	stdout     bytes.Buffer // all outputs
	stderr     bytes.Buffer // only errors
	outWriters []io.Writer
	errWriters []io.Writer
	inputs     []string

	initArgs  []string
	execArgs  []string
	executor  *Executor
	pidLatest int
}

func NewProgramSession(command, uid string, inputs []string, outWriters, errWriters []io.Writer) *ConsoleProgramSession {
	instance := new(ConsoleProgramSession)
	instance.command = command
	instance.uid = uid
	instance.filename = ""
	instance.execDir = ""
	instance.inputs = inputs
	instance.outWriters = append(instance.outWriters, outWriters...)
	instance.errWriters = append(instance.errWriters, errWriters...)

	instance.outWriters = append(instance.outWriters, &instance.stdout)
	instance.errWriters = append(instance.errWriters, &instance.stdout)
	instance.errWriters = append(instance.errWriters, &instance.stderr)

	instance.initArgs = make([]string, 0)
	instance.execArgs = make([]string, 0)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ConsoleProgramSession) String() string {
	if nil != instance {
		return instance.GoString()
	}
	return ""
}

func (instance *ConsoleProgramSession) GoString() string {
	if nil != instance {
		info := map[string]interface{}{
			"command":  instance.command,
			"filename": instance.filename,
			"out":      instance.StdOut(),
			"args":     instance.execArgs,
		}
		return qb_utils.JSON.Stringify(info)
	}
	return ""
}

func (instance *ConsoleProgramSession) StdOut() string {
	if nil != instance {
		return string(instance.stdout.Bytes())
	}
	return ""
}

func (instance *ConsoleProgramSession) StdErr() string {
	if nil != instance {
		return string(instance.stderr.Bytes())
	}
	return ""
}

func (instance *ConsoleProgramSession) StdOutJson() interface{} {
	if nil != instance {
		value := instance.StdOut()
		if a, b := qb_utils.JSON.StringToArray(value); b {
			return a
		} else if o, b := qb_utils.JSON.StringToMap(value); b {
			return o
		}
		return value
	}
	return ""
}

func (instance *ConsoleProgramSession) SetFileName(filename string) {
	instance.filename = filename
}

func (instance *ConsoleProgramSession) SetDir(dir string) {
	instance.execDir = dir
}

func (instance *ConsoleProgramSession) Run(args ...string) (*ConsoleProgramSession, error) {
	if nil != instance {
		return instance.run(args...)
	}
	return nil, nil
}

func (instance *ConsoleProgramSession) RunAsync(args ...string) (*ConsoleProgramSession, error) {
	if nil != instance {
		return instance.runAsync(args...)
	}
	return nil, nil
}

// Close try to close gracefully sending CTRL+C command
func (instance *ConsoleProgramSession) Close() {
	if nil != instance {
		if nil != instance.executor {
			instance.executor.Close()
		}
	}
}

func (instance *ConsoleProgramSession) Kill() (err error) {
	if nil != instance {
		if nil != instance.executor {
			err = instance.executor.Kill()
		} else if instance.pidLatest > 0 {
			err = qb_sys.Sys.KillProcessByPid(instance.pidLatest)
		}
	}
	return
}

func (instance *ConsoleProgramSession) Wait() (err error) {
	if nil != instance && nil != instance.executor {
		err = instance.executor.Wait()
	}
	return
}

func (instance *ConsoleProgramSession) PidCurrent() int {
	if nil != instance && nil != instance.executor {
		return instance.executor.PidCurrent()
	}
	return 0
}

func (instance *ConsoleProgramSession) PidLatest() int {
	if nil != instance {
		if nil != instance.executor {
			return instance.executor.PidLatest()
		}
		return instance.pidLatest
	}
	return 0
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ConsoleProgramSession) prepare(args ...string) (*Executor, []string) {

	// creates executor
	executor := NewExecutorWithDir(instance.command, instance.execDir)
	for _, input := range instance.inputs {
		executor.InputsAppend(input)
	}
	for _, w := range instance.outWriters {
		executor.OutWriterAppend(w)
	}
	for _, w := range instance.errWriters {
		executor.ErrorWriterAppend(w)
	}

	// create args for execution
	params := make([]string, 0)
	if len(instance.filename) > 0 {
		if b, _ := qb_utils.Paths.IsFile(instance.filename); b {
			params = append(params, instance.filename)
		}
	}
	if len(instance.initArgs) > 0 {
		params = append(params, instance.initArgs...)
		instance.initArgs = make([]string, 0)
	}
	params = append(params, args...)

	// set args
	instance.execArgs = params

	return executor, params
}

func (instance *ConsoleProgramSession) run(args ...string) (*ConsoleProgramSession, error) {
	defer instance.Close()

	instance.pidLatest = 0
	executor, params := instance.prepare(args...)
	err := executor.Run(params...)
	if nil != err {
		return nil, err
	}
	// assign executor
	instance.pidLatest = executor.PidLatest()
	instance.executor = executor
	err = wait(executor)
	instance.executor = nil

	return instance, err
}

func (instance *ConsoleProgramSession) runAsync(args ...string) (*ConsoleProgramSession, error) {
	instance.pidLatest = 0
	executor, params := instance.prepare(args...)
	err := executor.Run(params...)
	if nil != err {
		return nil, err
	}
	// assign executor
	instance.pidLatest = executor.PidLatest()
	instance.executor = executor

	go func() {
		err = wait(executor)
		instance.executor = nil
		instance.Close()
	}()

	return instance, nil
}

func wait(executor *Executor) (err error) {
	err = executor.Wait()
	if nil != err {
		if err.Error() == "exit status 1" {
			outErr := executor.StdErr()
			if len(outErr) > 0 {
				err = qb_utils.Errors.Prefix(err, outErr)
			}
		}
	}
	return
}

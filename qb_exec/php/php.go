package php

import (
	"errors"
	"io"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_utils"
)

var phpCommand = "php"

type PHPExec struct {
	command   string
	dirWork   string
	filename  string
	outWriter []io.Writer
	errWriter []io.Writer
}

// SetCommand Replace php command.
func SetCommand(cmd string) {
	if len(cmd) > 0 {
		phpCommand = cmd
	}
}

func NewPHPExec(args ...interface{}) *PHPExec {
	instance := new(PHPExec)
	instance.dirWork = qb_utils.Paths.Absolute("./")
	instance.command = phpCommand
	switch len(args) {
	case 1:
		instance.filename = qb_utils.Convert.ToString(args[0])
		instance.dirWork = qb_utils.Paths.Dir(instance.filename)
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *PHPExec) SetCommand(cmd string) {
	if len(cmd) > 0 {
		instance.command = cmd
	}
}

func (instance *PHPExec) SetDir(dir string) {
	instance.dirWork = dir
}

func (instance *PHPExec) GetDir() string {
	return instance.dirWork
}

func (instance *PHPExec) OutWriterAppend(w io.Writer) {
	instance.outWriter = append(instance.outWriter, w)
}

func (instance *PHPExec) ErrWriterAppend(w io.Writer) {
	instance.errWriter = append(instance.errWriter, w)
}

func (instance *PHPExec) Version() (string, error) {
	exec, err := instance.program().Run("--version")
	if nil != err {
		return "", err
	}
	response := strings.ToLower(exec.StdOut())

	if strings.Index(response, "php") == 0 {
		version := strings.TrimSpace(strings.ReplaceAll(response, "php", ""))
		return version, nil
	} else if strings.Index(response, "php") > 0 {
		nums := qb_utils.Regex.Numbers(response)
		if len(nums) > 0 {
			return nums[0], nil
		}
	}
	return "", errors.New(exec.StdOut())
}

func (instance *PHPExec) Run(args ...string) (*executor.ConsoleProgramSession, error) {
	return instance.program().Run(args...)
}

func (instance *PHPExec) RunAsync(args ...string) (*executor.ConsoleProgramSession, error) {
	return instance.program().RunAsync(args...)
}

func (instance *PHPExec) RunWrapped(args ...string) (*executor.ConsoleProgramSession, error) {
	return instance.program().RunWrapped(args...)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *PHPExec) program() *executor.ConsoleProgram {
	var program *executor.ConsoleProgram
	if len(instance.filename) > 0 {
		program = executor.NewConsoleProgramWithFile(instance.command, instance.filename)
		program.SetDir(instance.dirWork)
	} else {
		program = executor.NewConsoleProgramWithDir(instance.command, instance.dirWork)
	}
	for _, w := range instance.outWriter {
		program.OutWriterAppend(w)
	}
	for _, w := range instance.errWriter {
		program.ErrorWriterAppend(w)
	}
	return program
}

func init() {

}

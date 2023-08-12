package node_ts

import (
	"errors"
	"io"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_utils"
)

type TsNodeExec struct {
	execPath      string
	dirController *qb_utils.DirCentral
	logger        qb_log.ILogger
	initialized   bool

	filename  string
	outWriter []io.Writer
	errWriter []io.Writer
}

func NewTsNodeExec(execPath string, logger qb_log.ILogger, filename string) *TsNodeExec {
	instance := new(TsNodeExec)
	instance.execPath = execPath
	instance.logger = logger
	instance.filename = filename
	instance.dirController = qb_utils.Dir.NewCentral(fsName, ".tmp", true)

	instance.dirController.SetRoot(".")

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *TsNodeExec) SetLogger(logger qb_log.ILogger) *TsNodeExec {
	instance.logger = logger
	return instance
}

func (instance *TsNodeExec) SetRoot(dir string) *TsNodeExec {
	instance.dirController.SetRoot(dir)
	return instance
}

func (instance *TsNodeExec) SetTemp(dir string) *TsNodeExec {
	instance.dirController.SetTemp(dir)
	return instance
}

func (instance *TsNodeExec) SetSubTemp(enabled bool) *TsNodeExec {
	instance.dirController.SetSubTemp(enabled)
	return instance
}

func (instance *TsNodeExec) Root() string {
	return instance.dirController.DirRoot()
}

func (instance *TsNodeExec) Temp() string {
	return instance.dirController.DirTemp()
}

func (instance *TsNodeExec) Work() string {
	return instance.dirController.DirWork()
}

func (instance *TsNodeExec) GetPath(path string) (response string) {
	response = instance.dirController.GetPath(path)
	return
}

func (instance *TsNodeExec) GetWorkPath(subPath string) (response string) {
	response = instance.dirController.GetWorkPath(subPath)
	return
}

func (instance *TsNodeExec) GetTempPath(subPath string) (response string) {
	response = instance.dirController.GetTempPath(subPath)
	return
}

func (instance *TsNodeExec) LogFlush() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

func (instance *TsNodeExec) LogDisableRotation() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.RotateEnable(false)
	}
}

func (instance *TsNodeExec) OutWriterAppend(w io.Writer) {
	instance.outWriter = append(instance.outWriter, w)
}

func (instance *TsNodeExec) ErrWriterAppend(w io.Writer) {
	instance.errWriter = append(instance.errWriter, w)
}

func (instance *TsNodeExec) IsInstalled() bool {
	version, err := instance.Version()
	if nil != err {
		return false
	}
	return len(version) > 0
}

func (instance *TsNodeExec) Version() (string, error) {
	exec, err := instance.program().Run("--version")
	if nil != err {
		return "", err
	}
	response := strings.ToLower(exec.StdOut())
	if strings.Index(response, "v") == 0 {
		version := strings.TrimSpace(strings.ReplaceAll(response, "v", ""))
		return version, nil
	}
	return "", errors.New(exec.StdOut())
}

func (instance *TsNodeExec) Run(args ...string) (*executor.ConsoleProgramSession, error) {
	return instance.program().Run(args...)
}

func (instance *TsNodeExec) RunAsync(args ...string) (*executor.ConsoleProgramSession, error) {
	return instance.program().RunAsync(args...)
}

func (instance *TsNodeExec) RunWrapped(args ...string) (*executor.ConsoleProgramSession, error) {
	return instance.program().RunWrapped(args...)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *TsNodeExec) program() *executor.ConsoleProgram {
	instance.init()
	dir := instance.dirController.DirRoot()

	var program *executor.ConsoleProgram
	if len(instance.filename) > 0 {
		program = executor.NewConsoleProgramWithFile(tsnodeCommand, instance.filename)
		program.SetDir(dir)
	} else {
		program = executor.NewConsoleProgramWithDir(tsnodeCommand, dir)
	}
	for _, w := range instance.outWriter {
		program.OutWriterAppend(w)
	}
	for _, w := range instance.errWriter {
		program.ErrorWriterAppend(w)
	}
	return program
}

func (instance *TsNodeExec) init() {
	if nil != instance && !instance.initialized {
		instance.initialized = true

		instance.dirController.Refresh()

		logPath := instance.dirController.PathLog()
		if nil != instance.logger {
			instance.logger.(*qb_log.Logger).SetFilename(logPath)
		} else {
			instance.logger = qb_log.Log.New("info", logPath)
			if l, ok := instance.logger.(*qb_log.Logger); ok {
				l.SetMessageFormat("* " + l.GetMessageFormat())
			}
		}
	}
}

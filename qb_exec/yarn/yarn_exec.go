package yarn

import (
	"errors"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_utils"
)

type YarnExec struct {
	execPath      string
	dirController *qb_utils.DirCentral
	logger        qb_log.ILogger

	session     *executor.ConsoleProgramSession // current session
	initialized bool
}

func NewExec(execPath string, logger qb_log.ILogger) *YarnExec {
	instance := new(YarnExec)
	instance.execPath = execPath
	instance.logger = logger
	instance.dirController = qb_utils.Dir.NewCentral(fsName, ".tmp", true)

	instance.SetRoot(".")

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *YarnExec) SetRoot(dir string) *YarnExec {
	instance.dirController.SetRoot(dir)
	return instance
}

func (instance *YarnExec) SetTemp(dir string) *YarnExec {
	instance.dirController.SetTemp(dir)
	return instance
}

func (instance *YarnExec) SetSubTemp(enabled bool) *YarnExec {
	instance.dirController.SetSubTemp(enabled)
	return instance
}

func (instance *YarnExec) Root() string {
	return instance.dirController.DirRoot()
}

func (instance *YarnExec) Temp() string {
	return instance.dirController.DirTemp()
}

func (instance *YarnExec) Work() string {
	return instance.dirController.DirWork()
}

func (instance *YarnExec) GetPath(path string) (response string) {
	response = instance.dirController.GetPath(path)
	return
}

func (instance *YarnExec) GetWorkPath(subPath string) (response string) {
	response = instance.dirController.GetWorkPath(subPath)
	return
}

func (instance *YarnExec) GetTempPath(subPath string) (response string) {
	response = instance.dirController.GetTempPath(subPath)
	return
}

func (instance *YarnExec) LogFlush() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

func (instance *YarnExec) LogDisableRotation() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.RotateEnable(false)
	}
}

func (instance *YarnExec) SetLogger(logger qb_log.ILogger) *YarnExec {
	instance.logger = logger
	return instance
}

func (instance *YarnExec) TryKill() error {
	if nil != instance && nil != instance.session {
		return instance.session.Kill()
	}
	return nil
}

func (instance *YarnExec) Pid() int {
	if nil != instance && nil != instance.session {
		return instance.session.PidLatest()
	}
	return 0
}

func (instance *YarnExec) IsInstalled() bool {
	version, err := instance.Version()
	if nil != err {
		return false
	}
	return len(version) > 0
}

func (instance *YarnExec) Version() (response string, err error) {
	args := []string{"--version"}
	response, err = instance.ExecuteCommand(args...)
	return
}

func (instance *YarnExec) ExecuteCommand(arguments ...string) (out string, err error) {
	program := instance.program()
	session, execErr := program.Run(arguments...)
	if nil != execErr {
		err = execErr
		return
	}
	defer session.Close()
	instance.session = session

	stdErr := session.StdErr()
	if len(stdErr) > 0 {
		err = errors.New(stdErr)
	} else {
		out = strings.TrimSpace(session.StdOut())
	}
	return
}

func (instance *YarnExec) FlushLog() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *YarnExec) program(args ...string) *executor.ConsoleProgram {
	instance.init()
	dir := instance.dirController.DirRoot()
	return executor.NewConsoleProgramWithDir(instance.execPath, dir, args...)
}

func (instance *YarnExec) init() {
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

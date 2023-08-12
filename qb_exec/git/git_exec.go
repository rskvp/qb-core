package git

import (
	"errors"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_utils"
)

type GitExec struct {
	execPath    string
	dirRoot     string
	dirWork     string
	dirTemp     string
	logger      qb_log.ILogger
	initialized bool
}

func NewExec(execPath string, logger qb_log.ILogger) *GitExec {
	instance := new(GitExec)
	instance.execPath = execPath
	instance.logger = logger

	instance.SetRoot(".")

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *GitExec) SetLogger(logger qb_log.ILogger) *GitExec {
	instance.logger = logger
	return instance
}

func (instance *GitExec) SetRoot(dir string) *GitExec {
	instance.dirRoot = qb_utils.Paths.Absolute(dir)
	instance.dirWork = qb_utils.Paths.Concat(instance.dirRoot, fsName)
	instance.dirTemp = qb_utils.Paths.Concat(instance.dirWork, ".tmp")

	if instance.initialized {
		instance.initialized = false
		instance.init()
	}

	return instance
}

func (instance *GitExec) SetTemp(dir string) *GitExec {
	instance.dirTemp = dir
	return instance
}

func (instance *GitExec) Root() string {
	return instance.dirRoot
}

func (instance *GitExec) Temp() string {
	return instance.dirTemp
}

func (instance *GitExec) GetPath(path string) (response string) {
	response = qb_utils.Paths.Absolutize(path, instance.dirRoot)
	_ = qb_utils.Paths.Mkdir(response)
	return
}

func (instance *GitExec) GetWorkPath(subPath string) (response string) {
	response = qb_utils.Paths.Absolutize(subPath, instance.dirWork)
	_ = qb_utils.Paths.Mkdir(response)
	return
}

func (instance *GitExec) GetTempPath(subPath string) (response string) {
	response = qb_utils.Paths.Absolutize(subPath, instance.dirTemp)
	_ = qb_utils.Paths.Mkdir(response)
	return
}

func (instance *GitExec) LogFlush() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

func (instance *GitExec) LogDisableRotation() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.RotateEnable(false)
	}
}

func (instance *GitExec) IsInstalled() bool {
	version, err := instance.Version()
	if nil != err {
		return false
	}
	return len(version) > 0
}

func (instance *GitExec) Version() (response string, err error) {
	args := []string{"--version"}
	response, err = instance.ExecuteCommand(args...)
	return
}

func (instance *GitExec) ExecuteCommand(arguments ...string) (out string, err error) {
	program := instance.program()
	session, execErr := program.Run(arguments...)
	if nil != execErr {
		err = execErr
		return
	}
	defer session.Close()

	stdErr := session.StdErr()
	if len(stdErr) > 0 {
		err = errors.New(stdErr)
	} else {
		out = strings.TrimSpace(session.StdOut())
	}
	return
}

func (instance *GitExec) FlushLog() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *GitExec) program(args ...string) *executor.ConsoleProgram {
	instance.init()
	return executor.NewConsoleProgramWithDir(instance.execPath, instance.dirRoot, args...)
}

func (instance *GitExec) init() {
	if nil != instance && !instance.initialized {
		instance.initialized = true

		// creates paths
		_ = qb_utils.Paths.Mkdir(instance.dirWork + qb_utils.OS_PATH_SEPARATOR)
		_ = qb_utils.Paths.Mkdir(instance.dirTemp + qb_utils.OS_PATH_SEPARATOR)

		logPath := qb_utils.Paths.Concat(instance.dirWork, "logging.log")
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

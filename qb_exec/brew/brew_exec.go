package brew

import (
	"errors"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_utils"
)

type BrewExec struct {
	execPath string
	dirWork  string
	logger   qb_log.ILogger

	session *executor.ConsoleProgramSession // current session
}

func NewExec(execPath string, logger qb_log.ILogger) *BrewExec {
	instance := new(BrewExec)
	instance.execPath = execPath
	instance.dirWork = qb_utils.Paths.GetWorkspace(wpName).GetPath()
	instance.logger = logger

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *BrewExec) SetDir(dir string) *BrewExec {
	instance.dirWork = qb_utils.Paths.Absolute(dir)
	return instance
}

func (instance *BrewExec) GetDir() string {
	return instance.dirWork
}

func (instance *BrewExec) TryKill() error {
	if nil != instance && nil != instance.session {
		return instance.session.Kill()
	}
	return nil
}

func (instance *BrewExec) Pid() int {
	if nil != instance && nil != instance.session {
		return instance.session.PidLatest()
	}
	return 0
}

func (instance *BrewExec) IsInstalled() bool {
	version, err := instance.Version()
	if nil != err {
		return false
	}
	return len(version) > 0
}

func (instance *BrewExec) Version() (response string, err error) {
	args := []string{"--version"}
	response, err = instance.ExecuteCommand(args...)
	return
}

func (instance *BrewExec) ExecuteCommand(arguments ...string) (out string, err error) {
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

func (instance *BrewExec) FlushLog() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//	c o m m a n d s
// ---------------------------------------------------------------------------------------------------------------------

// Install run brew install ....
func (instance *BrewExec) Install(arguments ...string) (out string, err error) {
	program := instance.program()
	args := []string{"install"}
	args = append(args, arguments...)
	session, execErr := program.Run(args...)
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

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *BrewExec) program(args ...string) *executor.ConsoleProgram {
	return executor.NewConsoleProgramWithDir(instance.execPath, instance.dirWork, args...)
}

func (instance *BrewExec) run(args ...string) (out string, err error) {
	session, e := instance.program().Run(args...)
	if nil != e {
		err = e
		return
	}
	out = strings.TrimSpace(session.StdOut())
	session.Close()
	return
}

package npm

import (
	"errors"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_utils"
)

type NpmExec struct {
	execPath      string
	dirController *qb_utils.DirCentral
	logger        qb_log.ILogger
	initialized   bool
}

func NewExec(execPath string, logger qb_log.ILogger) *NpmExec {
	instance := new(NpmExec)
	instance.execPath = execPath
	instance.logger = logger
	instance.dirController = qb_utils.Dir.NewCentral(fsName, ".tmp", true)

	instance.SetRoot(".")

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NpmExec) SetLogger(logger qb_log.ILogger) *NpmExec {
	instance.logger = logger
	return instance
}

func (instance *NpmExec) SetRoot(dir string) *NpmExec {
	instance.dirController.SetRoot(dir)
	return instance
}

func (instance *NpmExec) SetTemp(dir string) *NpmExec {
	instance.dirController.SetTemp(dir)
	return instance
}

func (instance *NpmExec) SetSubTemp(enabled bool) *NpmExec {
	instance.dirController.SetSubTemp(enabled)
	return instance
}

func (instance *NpmExec) Root() string {
	return instance.dirController.DirRoot()
}

func (instance *NpmExec) Temp() string {
	return instance.dirController.DirTemp()
}

func (instance *NpmExec) Work() string {
	return instance.dirController.DirWork()
}

func (instance *NpmExec) GetPath(path string) (response string) {
	response = instance.dirController.GetPath(path)
	return
}

func (instance *NpmExec) GetWorkPath(subPath string) (response string) {
	response = instance.dirController.GetWorkPath(subPath)
	return
}

func (instance *NpmExec) GetTempPath(subPath string) (response string) {
	response = instance.dirController.GetTempPath(subPath)
	return
}

func (instance *NpmExec) LogFlush() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

func (instance *NpmExec) LogDisableRotation() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.RotateEnable(false)
	}
}

func (instance *NpmExec) IsInstalled() bool {
	version, err := instance.Version()
	if nil != err {
		return false
	}
	return len(version) > 0
}

func (instance *NpmExec) Version() (string, error) {
	exec, err := instance.program().Run("--version")
	if nil != err {
		return "", err
	}
	response := strings.ToLower(exec.StdOut())
	if strings.Index(response, ".") > -1 {
		version := strings.TrimSpace(strings.ReplaceAll(response, "v", ""))
		return version, nil
	}
	return "", errors.New(exec.StdOut())
}

func (instance *NpmExec) Help() (string, error) {
	if instance.IsInstalled() {
		exec, err := instance.program().Run("--help")
		if nil != err && err.Error() != "exit status 1" {
			return "", err
		}
		if nil != exec {
			return exec.StdOut(), nil
		}
		return "", err
	}
	return "", executor.NotInstalledError
}

func (instance *NpmExec) HelpOn(command string) (string, error) {
	if instance.IsInstalled() {
		exec, err := instance.program().Run(command, "-h")
		if nil != err {
			return "", err
		}
		return exec.StdOut(), nil
	}
	return "", executor.NotInstalledError
}

func (instance *NpmExec) Init(data *DataPackage) (string, error) {
	return instance.InitFromTemplate(tplPackage, data)
}

func (instance *NpmExec) InitFromTemplate(tpl string, data *DataPackage) (string, error) {
	if len(data.Name) == 0 {
		data.Name = "program"
	}
	if len(data.License) == 0 {
		data.License = "MIT"
	}
	if len(data.Version) == 0 {
		data.Version = "1.0.0"
	}
	if len(data.Main) == 0 {
		data.Main = "index.js"
	}
	dir := instance.dirController.DirWork()
	filename := qb_utils.Paths.Absolute(qb_utils.Paths.Concat(dir, "package.json"))
	if b, _ := qb_utils.Paths.Exists(filename); b {
		return filename, qb_utils.Errors.Prefix(errors.New("file_already_exists"), filename)
	}
	text, err := MergeTpl(tpl, data)
	if nil != err {
		return "", err
	}
	_, err = qb_utils.IO.WriteTextToFile(text, filename)
	return filename, err
}

func (instance *NpmExec) Install() (string, error) {
	p := instance.program()
	return getResponse(p.Run("install"))
}

func (instance *NpmExec) Run(command string) (string, error) {
	p := instance.program()
	return getResponse(p.Run("run", command))
}

func (instance *NpmExec) RunAsync(command string) (*executor.ConsoleProgramSession, error) {
	p := instance.program()
	return p.RunAsync("run", command)
}

func (instance *NpmExec) ExecuteCommand(arguments ...string) (out string, err error) {
	program := instance.program()
	session, execErr := program.Run(arguments...)
	if nil != execErr {
		err = execErr
		return
	}

	stdErr := session.StdErr()
	if len(stdErr) > 0 {
		err = errors.New(stdErr)
	} else {
		out = strings.TrimSpace(session.StdOut())
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NpmExec) program() *executor.ConsoleProgram {
	instance.init()
	dir := instance.dirController.DirRoot()

	return executor.NewConsoleProgramWithDir(npmCommand, dir)
}

func (instance *NpmExec) init() {
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

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func getResponse(exec *executor.ConsoleProgramSession, err error) (string, error) {
	if nil != err {
		message := err.Error()
		if strings.Index(message, "exit status") == -1 {
			return "", err
		}
	}
	if nil != exec {
		stderr := exec.StdErr()
		if len(stderr) > 0 {
			return "", errors.New(stderr)
		}
		return exec.StdOut(), nil
	}
	return "", err
}

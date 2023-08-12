package nvm

import (
	"fmt"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/shell"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_utils"
)

type NvmExec struct {
	dirController *qb_utils.DirCentral
	logger        qb_log.ILogger
	initialized   bool
}

func NewExec(logger qb_log.ILogger) *NvmExec {
	instance := new(NvmExec)
	instance.logger = logger
	instance.dirController = qb_utils.Dir.NewCentral(root, ".tmp", false)

	instance.SetRoot(root)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NvmExec) SetLogger(logger qb_log.ILogger) *NvmExec {
	instance.logger = logger
	return instance
}

func (instance *NvmExec) SetRoot(dir string) *NvmExec {
	instance.dirController.SetRoot(dir)
	return instance
}

func (instance *NvmExec) SetTemp(dir string) *NvmExec {
	instance.dirController.SetTemp(dir)
	return instance
}

func (instance *NvmExec) SetSubTemp(enabled bool) *NvmExec {
	instance.dirController.SetSubTemp(enabled)
	return instance
}

func (instance *NvmExec) Root() string {
	return instance.dirController.DirRoot()
}

func (instance *NvmExec) Temp() string {
	return instance.dirController.DirTemp()
}

func (instance *NvmExec) Work() string {
	return instance.dirController.DirWork()
}

func (instance *NvmExec) GetPath(path string) (response string) {
	response = instance.dirController.GetPath(path)
	return
}

func (instance *NvmExec) GetWorkPath(subPath string) (response string) {
	response = instance.dirController.GetWorkPath(subPath)
	return
}

func (instance *NvmExec) GetTempPath(subPath string) (response string) {
	response = instance.dirController.GetTempPath(subPath)
	return
}

func (instance *NvmExec) LogFlush() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

func (instance *NvmExec) LogDisableRotation() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.RotateEnable(false)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	i n s t a l l a t i o n
//----------------------------------------------------------------------------------------------------------------------

// Install
// https://github.com/nvm-sh/nvm#installing-and-updating
func (instance *NvmExec) Install(version string) (response string, err error) {
	if len(version) == 0 {
		version = "v0.39.2"
	}
	command := fmt.Sprintf("curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/%s/install.sh | bash", version)
	var filename string
	filename, err = shell.CreateFile(command)
	if nil == err {
		defer qb_utils.IO.Remove(filename)

		response, err = shell.RunFile(filename)
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	c o m m a n d s
//----------------------------------------------------------------------------------------------------------------------

func (instance *NvmExec) IsInstalled() bool {
	version, err := instance.Version()
	if nil != err {
		return false
	}
	return len(version) > 0
}

func (instance *NvmExec) Version() (response string, err error) {
	script := instance.getNvmScript("--version", "")
	response, err = shell.Run(script)
	return
}

func (instance *NvmExec) Help() (response string, err error) {
	script := instance.getNvmScript("--help", "")
	response, err = shell.Run(script)
	return
}

func (instance *NvmExec) Use(nodeVersion string) (response string, err error) {
	script := instance.getNvmScript("use ", nodeVersion)
	response, err = shell.Run(script)
	return
}

func (instance *NvmExec) NodeRun(nodeVersion string, command string) (response string, err error) {
	script := instance.getNodeScript("use ", nodeVersion, command)
	response, err = shell.Run(script)
	return
}

func (instance *NvmExec) NpmRun(nodeVersion string, command string) (response string, err error) {
	script := instance.getNpmScript("use ", nodeVersion, command)
	response, err = shell.Run(script)
	return
}

func (instance *NvmExec) NpxRun(nodeVersion string, command string) (response string, err error) {
	response, err = instance.Run(nodeVersion, "npx "+command)
	return
}

func (instance *NvmExec) Run(nodeVersion string, commands ...string) (response string, err error) {
	script := instance.getNvmsScript("use ", nodeVersion, commands...)
	response, err = shell.Run(script)
	return
}

func (instance *NvmExec) RunBackground(nodeVersion string, commands ...string) (c *shell.ShellCmd, err error) {
	script := instance.getNvmsScript("use ", nodeVersion, commands...)
	c, err = shell.RunBackground(script)
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NvmExec) getNvmScript(command, param string) string {
	instance.init()
	response, _ := qb_utils.Formatter.Merge(tplNvm, map[string]interface{}{
		"command": command,
		"param":   param,
	})
	return response
}

func (instance *NvmExec) getNvmsScript(command, param string, ncommands ...string) string {
	instance.init()
	model := map[string]interface{}{
		"command":   command,
		"param":     param,
		"ncommands": strings.Join(ncommands, "\n"),
	}
	response, _ := qb_utils.Formatter.Merge(tplNvms, model)
	return response
}

func (instance *NvmExec) getNodeScript(command, param, ncommand string) string {
	instance.init()
	response, _ := qb_utils.Formatter.Merge(tplNode, map[string]interface{}{
		"command":  command,
		"param":    param,
		"ncommand": ncommand,
	})
	return response
}

func (instance *NvmExec) getNpmScript(command, param, ncommand string) string {
	instance.init()
	response, _ := qb_utils.Formatter.Merge(tplNpm, map[string]interface{}{
		"command":  command,
		"param":    param,
		"ncommand": ncommand,
	})
	return response
}

func (instance *NvmExec) getNpxScript(command, param, ncommand string) string {
	instance.init()
	response, _ := qb_utils.Formatter.Merge(tplNpx, map[string]interface{}{
		"command":  command,
		"param":    param,
		"ncommand": ncommand,
	})
	return response
}

func (instance *NvmExec) init() {
	if nil != instance && !instance.initialized {
		instance.initialized = true

		instance.dirController.Refresh()

		logPath := instance.dirController.PathLog()
		if nil == instance.logger {
			instance.logger = qb_log.Log.New("info", logPath)
			if l, ok := instance.logger.(*qb_log.Logger); ok {
				l.SetMessageFormat("* " + l.GetMessageFormat())
			}
		}
	}
}

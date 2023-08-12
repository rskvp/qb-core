package python

import (
	"errors"
	"io"
	"os/exec"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_utils"
)

var pythonCommand = "python"

type PythonExec struct {
	command   string
	dirWork   string
	filename  string
	outWriter []io.Writer
	errWriter []io.Writer
}

// SetCommand Replace python command. ex: "python3"
func SetCommand(cmd string) {
	if len(cmd) > 0 {
		pythonCommand = cmd
	}
}

func NewPythonExec(args ...interface{}) *PythonExec {
	instance := new(PythonExec)
	instance.command = pythonCommand
	instance.dirWork = qb_utils.Paths.Absolute("./")
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

func (instance *PythonExec) SetCommand(cmd string) {
	if len(cmd) > 0 {
		instance.command = cmd
	}
}

func (instance *PythonExec) SetDir(dir string) {
	instance.dirWork = dir
}

func (instance *PythonExec) GetDir() string {
	return instance.dirWork
}

func (instance *PythonExec) OutWriterAppend(w io.Writer) {
	instance.outWriter = append(instance.outWriter, w)
}

func (instance *PythonExec) ErrWriterAppend(w io.Writer) {
	instance.errWriter = append(instance.errWriter, w)
}

func (instance *PythonExec) Version() (string, error) {
	exec, err := instance.program(true).Run("--version")
	if nil != err {
		return "", err
	}
	response := strings.ToLower(exec.StdOut())

	if strings.Index(response, "python") > -1 {
		version := strings.TrimSpace(strings.ReplaceAll(response, "python", ""))
		return version, nil
	}
	return "", errors.New(exec.StdOut())
}

func (instance *PythonExec) Run(args ...string) (*executor.ConsoleProgramSession, error) {
	return instance.program(false).Run(args...)
}

func (instance *PythonExec) RunAsync(args ...string) (*executor.ConsoleProgramSession, error) {
	return instance.program(false).RunAsync(args...)
}

func (instance *PythonExec) RunWrapped(args ...string) (*executor.ConsoleProgramSession, error) {
	return instance.program(false).RunWrapped(args...)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *PythonExec) program(ignoreFile bool) *executor.ConsoleProgram {
	var program *executor.ConsoleProgram
	if !ignoreFile && len(instance.filename) > 0 {
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
	pythonCommand = findExecPath()
}

func findExecPath() string {
	for _, path := range [...]string{
		// Unix-like
		"python",
		"/usr/bin/python",
		"/usr/local/bin/python",

		// Windows
		"python.exe", // in case PATHEXT is misconfigured
		// `C:\Program Files (x86)\Microsoft VS Code\chrome.exe`,
		//filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\Programs\Microsoft VS Code\code.exe`),

		// Mac
		//"/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code",
	} {
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	// Fall back
	return "python"
}

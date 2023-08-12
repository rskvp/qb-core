package vscode

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_utils"
)

var vscodeCommand = "code"

type VSCodeExec struct {
	dirWork string
}

// SetCommand Replace command.
func SetCommand(cmd string) {
	if len(cmd) > 0 {
		vscodeCommand = cmd
	}
}

func NewVSCodeExec() *VSCodeExec {
	instance := new(VSCodeExec)
	instance.dirWork = qb_utils.Paths.Absolute("./")

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *VSCodeExec) SetDir(dir string) {
	instance.dirWork = dir
}

func (instance *VSCodeExec) GetDir() string {
	return instance.dirWork
}

func (instance *VSCodeExec) Open(path string) error {
	if nil != instance {
		_, err := instance.program().Run(path)
		if nil != err {
			return err
		}
	}
	return nil
}

func (instance *VSCodeExec) Version() (version, commitId, architecture string) {
	if nil != instance {
		session, _ := instance.program().Run("--version")
		if nil != session {
			out := session.StdOut()
			tokens := strings.Split(out, "\n")
			return qb_utils.Arrays.GetAt(tokens, 0, "").(string),
				qb_utils.Arrays.GetAt(tokens, 1, "").(string),
				qb_utils.Arrays.GetAt(tokens, 2, "").(string)
		}
	}
	return "", "", ""
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *VSCodeExec) program() *executor.ConsoleProgram {
	return executor.NewConsoleProgramWithDir(vscodeCommand, instance.dirWork)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func init() {
	vscodeCommand = findExecPath()
}

func findExecPath() string {
	for _, path := range [...]string{
		// Unix-like
		"code",
		"/usr/bin/code",

		// Windows
		"code.exe", // in case PATHEXT is misconfigured
		`C:\Program Files (x86)\Microsoft VS Code\chrome.exe`,
		`C:\Program Files\Microsoft VS Code\chrome.exe`,
		filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\Programs\Microsoft VS Code\code.exe`),

		// Mac
		"/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code",
	} {
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	// Fall back to something simple and sensible, to give a useful error
	// message.
	return "code"
}

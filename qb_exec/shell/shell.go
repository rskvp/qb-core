package shell

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rskvp/qb-core/qb_rnd"
	"github.com/rskvp/qb-core/qb_sys"
	"github.com/rskvp/qb-core/qb_utils"
)

//go:embed tpl_sh.txt
var tplSh string

//go:embed tpl_bat.txt
var tplBat string

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// CreateFile creates an executable shell file
func CreateFile(shellScript string) (filename string, err error) {
	filename = getShellTempFile()
	content := fmt.Sprintf(getShellTemplate(), shellScript)
	_, err = qb_utils.IO.WriteTextToFile(content, filename)
	if nil != err {
		return
	}
	_, err = qb_utils.IO.Chmod(filename, 0755)
	if nil != err {
		return
	}

	return
}

// Run a shell command wrapping into .sh or .bat file
func Run(shellScript string, args ...string) (response string, err error) {
	var filename string
	filename, err = CreateFile(shellScript)
	defer qb_utils.IO.Remove(filename)
	if nil == err {
		response, err = RunFile(filename, args...)
	}
	return
}

func RunBackground(shellScript string, args ...string) (c *ShellCmd, err error) {
	var filename string
	filename, err = CreateFile(shellScript)
	defer qb_utils.IO.Remove(filename)
	if nil == err {
		c, err = RunFileBackground(filename, args...)
	}
	return
}

func RunFile(filename string, args ...string) (string, error) {
	cmd := getShellCmd()
	params := make([]string, 0)
	params = append(params, filename)
	params = append(params, args...)
	out, err := RunExec(cmd, params...)
	if nil != err {
		return "", err
	}
	return strings.Trim(string(out), " \n"), nil
}

func RunFileBackground(filename string, args ...string) (c *ShellCmd, err error) {
	cmd := getShellCmd()
	params := make([]string, 0)
	params = append(params, filename)
	params = append(params, args...)
	c = NewShellCmd(cmd, params...)
	c.SetBackground(true)
	err = c.Run()
	return
}

func RunExec(cmd string, args ...string) ([]byte, error) {
	c := exec.Command(cmd, args...)
	return toResponse(c.CombinedOutput())
}

// RunExecBackground start command in background a
func RunExecBackground(cmd string, args ...string) ([]byte, error) {
	ctx := context.Background()
	c := exec.CommandContext(ctx, cmd, args...)
	return toResponse(c.CombinedOutput())
}

func RunOutput(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)

	// output
	var out bytes.Buffer
	var err bytes.Buffer
	c.Stdout = &out
	c.Stderr = &err

	e := c.Run()
	if nil != e {
		return "", e
	}

	se := err.String()
	if len(se) > 0 {
		return "", errors.New(se)
	}
	return out.String(), nil
}

func StartExecBackground(cmd string, args ...string) (c *ShellCmd, err error) {
	c = NewShellCmd(cmd, args...)
	c.SetBackground(true)
	err = c.Run()
	if nil != err {
		return
	}
	return
}

func CommandStart(cmd string, args ...string) (*exec.Cmd, error) {
	c := exec.Command(cmd, args...)
	err := c.Start()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func CommandCombinedOut(cmd string, args ...string) (*exec.Cmd, *bytes.Buffer) {
	c := exec.Command(cmd, args...)

	// combined output
	var b bytes.Buffer
	c.Stdout = &b
	c.Stderr = &b

	return c, &b
}

func Command(cmd string, args ...string) (c *exec.Cmd, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	c = exec.Command(cmd, args...)

	// combined output
	var out bytes.Buffer
	var err bytes.Buffer
	c.Stdout = &out
	c.Stderr = &err

	stdout = &out
	stderr = &err

	return c, stdout, stderr
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func toResponse(out []byte, err error) ([]byte, error) {
	if err != nil {
		if len(out) > 0 {
			s := strings.Trim(string(out), " \n")
			return nil, qb_utils.Errors.Prefix(err, s+": ")
		}
		return nil, err
	}
	return out, nil
}

func getShellCmd() string {
	cmd := "/bin/sh"
	if qb_sys.Sys.IsWindows() {
		cmd = "cmd"
	}
	return cmd
}

func getShellTempFile() string {
	ext := ".sh"
	if qb_sys.Sys.IsWindows() {
		ext = ".bat"
	}

	return qb_utils.Paths.Absolute(qb_rnd.Rnd.Uuid() + ext)
}

func getShellTemplate() string {
	response := tplSh
	if qb_sys.Sys.IsWindows() {
		response = tplBat
	}

	return response
}

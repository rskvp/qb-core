package qb_exec

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/brew"
	"github.com/rskvp/qb-core/qb_exec/certbot"
	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_exec/getssl"
	"github.com/rskvp/qb-core/qb_exec/ghostscript"
	"github.com/rskvp/qb-core/qb_exec/git"
	"github.com/rskvp/qb-core/qb_exec/libreoffice"
	"github.com/rskvp/qb-core/qb_exec/node"
	"github.com/rskvp/qb-core/qb_exec/node_ts"
	"github.com/rskvp/qb-core/qb_exec/npm"
	"github.com/rskvp/qb-core/qb_exec/nvm"
	"github.com/rskvp/qb-core/qb_exec/php"
	"github.com/rskvp/qb-core/qb_exec/ping"
	"github.com/rskvp/qb-core/qb_exec/python"
	"github.com/rskvp/qb-core/qb_exec/shell"
	"github.com/rskvp/qb-core/qb_exec/vscode"
	"github.com/rskvp/qb-core/qb_exec/yarn"
	"github.com/rskvp/qb-core/qb_utils"
)

type ExecHelper struct {
}

var Exec *ExecHelper

func init() {
	Exec = new(ExecHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ExecHelper) RunShell(shellScript string, args ...string) (string, error) {
	return shell.RunFile(shellScript, args...)
}

func (instance *ExecHelper) Run(cmd string, args ...string) ([]byte, error) {
	return shell.RunExec(cmd, args...)
}

func (instance *ExecHelper) RunBackground(cmd string, args ...string) ([]byte, error) {
	return shell.RunExecBackground(cmd, args...)
}

func (instance *ExecHelper) RunOutput(cmd string, args ...string) (string, error) {
	return shell.RunOutput(cmd, args...)
}

func (instance *ExecHelper) Start(cmd string, args ...string) (*exec.Cmd, error) {
	return shell.CommandStart(cmd, args...)
}

// CommandCombinedOut creates a command and return command instance and output buffers (both Stdout and Stderr)
func (instance *ExecHelper) CommandCombinedOut(cmd string, args ...string) (*exec.Cmd, *bytes.Buffer) {
	return shell.CommandCombinedOut(cmd, args...)
}

func (instance *ExecHelper) Command(cmd string, args ...string) (c *exec.Cmd, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	return shell.Command(cmd, args...)
}

func (instance *ExecHelper) Open(args ...string) error {
	c := openFileCommand(args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return err
	}
	if len(out) > 0 {
		// may be an error
		return errors.New(string(out))
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	s p e c i a l i z e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *ExecHelper) NewExecutor(cmd string) *executor.Executor {
	return executor.NewExecutor(cmd)
}

func (instance *ExecHelper) NewExecutorWithDir(cmd, dir string) *executor.Executor {
	return executor.NewExecutorWithDir(cmd, dir)
}

func (instance *ExecHelper) NewPing() *ping.PingExec {
	return ping.NewPingExec()
}

func (instance *ExecHelper) NewLibreOffice() *libreoffice.LibreOfficeExec {
	return libreoffice.NewLibreOfficeExec()
}

func (instance *ExecHelper) NewVSCode() *vscode.VSCodeExec {
	return vscode.NewVSCodeExec()
}

func (instance *ExecHelper) NewPython() *python.PythonExec {
	return python.NewPythonExec()
}

func (instance *ExecHelper) NewPythonWithCommand(command string) *python.PythonExec {
	program := python.NewPythonExec()
	program.SetCommand(command)
	return program
}

func (instance *ExecHelper) NewPHP() *php.PHPExec {
	return php.NewPHPExec()
}

func (instance *ExecHelper) NewNode() *node.NodeExec {
	return node.Node.NewExec()
}

func (instance *ExecHelper) Node() *node.NodeHelper {
	return node.Node
}

func (instance *ExecHelper) TsNode() *node_ts.TsNodeHelper {
	return node_ts.TsNode
}

func (instance *ExecHelper) NewTsNode() *node_ts.TsNodeExec {
	return node_ts.TsNode.NewExec()
}

func (instance *ExecHelper) NewNpm() *npm.NpmExec {
	return npm.Npm.NewExec()
}

func (instance *ExecHelper) Npm() *npm.NpmHelper {
	return npm.Npm
}

func (instance *ExecHelper) Nvm() *nvm.NvmHelper {
	return nvm.Nvm
}

func (instance *ExecHelper) NewGhostScript() *ghostscript.GhostScriptExec {
	return ghostscript.NewGhostScriptExec()
}

func (instance *ExecHelper) NewGit() *git.GitExec {
	return git.Git.NewExec()
}

func (instance *ExecHelper) Git() *git.GitHelper {
	return git.Git
}

func (instance *ExecHelper) NewBrew() *brew.BrewExec {
	return brew.Brew.NewExec()
}

func (instance *ExecHelper) Brew() *brew.BrewHelper {
	return brew.Brew
}

func (instance *ExecHelper) NewYarn() *yarn.YarnExec {
	return yarn.Yarn.NewExec()
}

func (instance *ExecHelper) Yarn() *yarn.YarnHelper {
	return yarn.Yarn
}

func (instance *ExecHelper) NewCertbot() *certbot.CertbotExec {
	return certbot.Certbot.NewExec()
}

func (instance *ExecHelper) Certbot() *certbot.CertbotHelper {
	return certbot.Certbot
}

func (instance *ExecHelper) NewGetSSL() *getssl.GetSslExec {
	return getssl.GetSSL.NewExec()
}

func (instance *ExecHelper) GetSSL() *getssl.GetSslHelper {
	return getssl.GetSSL
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
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

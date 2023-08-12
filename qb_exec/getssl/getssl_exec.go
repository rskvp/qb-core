package getssl

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_utils"
)

const cmdInstall = "curl --silent https://raw.githubusercontent.com/srvrco/getssl/latest/getssl > getssl ; chmod 700 getssl"

type GetSslExec struct {
	execPath    string
	dirRoot     string // start dir
	dirWork     string // start dir - workspace
	dirBin      string // start dir - workspace - bin
	dirGetssl   string // start dir - workspace - .getssl
	logger      qb_log.ILogger
	initialized bool

	checkedInstallation bool
	version             string
}

func NewExec(execPath string, logger qb_log.ILogger) *GetSslExec {
	instance := new(GetSslExec)
	instance.execPath = execPath
	instance.logger = logger

	instance.SetRoot(".")

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *GetSslExec) SetLogger(logger qb_log.ILogger) *GetSslExec {
	instance.logger = logger
	return instance
}

func (instance *GetSslExec) SetRoot(dir string) *GetSslExec {
	instance.dirRoot = qb_utils.Paths.Absolute(dir)
	instance.dirWork = qb_utils.Paths.Concat(instance.dirRoot, fsName)
	instance.dirBin = qb_utils.Paths.Concat(instance.dirWork, "bin")
	instance.dirGetssl = qb_utils.Paths.UserHomePath(".getssl")

	if instance.initialized {
		instance.initialized = false
		instance.init()
	}

	return instance
}

func (instance *GetSslExec) Root() string {
	return instance.dirRoot
}

func (instance *GetSslExec) Work() string {
	return instance.dirWork
}

func (instance *GetSslExec) GetPath(path string) (response string) {
	response = qb_utils.Paths.Absolutize(path, instance.dirRoot)
	_ = qb_utils.Paths.Mkdir(response)
	return
}

func (instance *GetSslExec) GetWorkPath(subPath string) (response string) {
	response = qb_utils.Paths.Absolutize(subPath, instance.dirWork)
	_ = qb_utils.Paths.Mkdir(response)
	return
}

func (instance *GetSslExec) LogFlush() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

func (instance *GetSslExec) LogDisableRotation() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.RotateEnable(false)
	}
}

func (instance *GetSslExec) FlushLog() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//	c o m m a n d s
// ---------------------------------------------------------------------------------------------------------------------

func (instance *GetSslExec) IsInstalled() bool {
	version, err := instance.checkInstall()
	if nil != err {
		return false
	}
	return len(version) > 0
}

func (instance *GetSslExec) Version() (response string, err error) {
	response, err = instance.checkInstall()
	return
}

func (instance *GetSslExec) Upgrade() (response string, err error) {
	response, err = instance.ExecuteCommand("-u")
	return
}

func (instance *GetSslExec) ReadDefaultConfig() (response string, err error) {
	filename := qb_utils.Paths.Concat(instance.dirGetssl, "getssl.cfg")
	response, err = qb_utils.IO.ReadTextFromFile(filename)
	return
}

func (instance *GetSslExec) ReadConfig(domain string) (response string, err error) {
	filename := qb_utils.Paths.Concat(instance.dirGetssl, domainName(domain), "getssl.cfg")
	response, err = qb_utils.IO.ReadTextFromFile(filename)
	return
}

func (instance *GetSslExec) InitCertificate(domain string) (response string, err error) {
	response, err = instance.ExecuteCommand("-c", domainName(domain))
	return
}

func (instance *GetSslExec) GetCertificate(domain string) (response string, err error) {
	response, err = instance.ExecuteCommand(domainName(domain))
	return
}

func (instance *GetSslExec) Configure(domain string, acl []string) (response string, err error) {
	filename := qb_utils.Paths.Concat(instance.dirGetssl, domainName(domain), "getssl.cfg")
	config, e := qb_utils.IO.ReadTextFromFile(filename)
	if nil != e {
		err = e
		return
	}

	// start replacing text
	config = setACL(config, acl)

	// save settings
	_, err = qb_utils.IO.WriteTextToFile(config, filename)
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	g e n e r i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *GetSslExec) ExecuteCommand(arguments ...string) (out string, err error) {
	_, err = instance.checkInstall()
	if nil != err {
		return
	}
	return instance.exec(arguments...)
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *GetSslExec) params() []string {
	return []string{}
}

func (instance *GetSslExec) exec(arguments ...string) (out string, err error) {
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

func (instance *GetSslExec) program(args ...string) *executor.ConsoleProgram {
	instance.init()
	return executor.NewConsoleProgramWithDir(instance.execPath, instance.dirWork, args...)
}

func (instance *GetSslExec) init() {
	if nil != instance && !instance.initialized {
		instance.initialized = true

		// creates paths
		_ = qb_utils.Paths.Mkdir(instance.dirWork + qb_utils.OS_PATH_SEPARATOR)
		_ = qb_utils.Paths.Mkdir(instance.dirBin + qb_utils.OS_PATH_SEPARATOR)
		_ = qb_utils.Paths.Mkdir(instance.dirGetssl + qb_utils.OS_PATH_SEPARATOR)

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

func (instance *GetSslExec) checkInstall() (version string, err error) {
	if !instance.checkedInstallation {
		version, err = instance.exec("--version")
		if nil != err {
			err = instance.tryInstall()
			if nil != err {
				err = errors.New(fmt.Sprintf("Please, install GetSsl manually using this command: %s or check website %s",
					cmdInstall, "https://github.com/srvrco/getssl"))
			} else {
				version, err = instance.exec("--version")
			}
		}
		if nil == err && len(version) > 0 {
			instance.checkedInstallation = true
			tokens := strings.Split(version, " ")
			if len(tokens) == 2 {
				version = tokens[1]
			}
			instance.version = version
		}
	} else {
		version = instance.version
	}
	return
}

func (instance *GetSslExec) tryInstall() error {
	c := exec.Command("bash", "-c", "curl --silent https://raw.githubusercontent.com/srvrco/getssl/latest/getssl > getssl")
	c.Dir = instance.dirBin
	out, err := c.CombinedOutput()
	if nil != err {
		instance.logger.Error(err)
		instance.FlushLog()
		return err
	}
	c = exec.Command("bash", "-c", "chmod 700 getssl")
	c.Dir = instance.dirBin
	out, err = c.CombinedOutput()
	if nil != err {
		instance.logger.Error(err)
		instance.FlushLog()
		return err
	}
	if len(out) > 0 {
		instance.logger.Debug(string(out))
		instance.FlushLog()
	}
	instance.execPath = qb_utils.Paths.Concat(instance.dirBin, "getssl")
	return nil
}

func domainName(domain string) (response string) {
	response = strings.ToLower(domain)
	response = strings.ReplaceAll(response, "www.", "")
	response = strings.ReplaceAll(response, "https://", "")
	response = strings.ReplaceAll(response, "http://", "")
	return
}

func setACL(text string, acl []string) (response string) {
	response = text
	if len(text) > 0 && len(acl) > 0 {
		sub := qb_utils.Strings.SubBetween(text, "#ACL=(", ")")
		if len(sub) == 0 {
			sub = qb_utils.Strings.SubBetween(text, "ACL=(", ")")
		}
		if len(sub) == 0 {
			return
		}

		var repl strings.Builder
		repl.WriteString("ACL=(\n")
		for _, s := range acl {
			repl.WriteString(s)
			repl.WriteString("\n")
		}
		repl.WriteString(")")

		// replace
		response = strings.Replace(text, sub, repl.String(), 1)
	}
	return
}

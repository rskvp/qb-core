package certbot

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_utils"
)

type CertbotExec struct {
	execPath    string
	dirRoot     string
	dirWork     string
	dirCerts    string
	dirConfig   string
	dirLogs     string
	logger      qb_log.ILogger
	initialized bool

	checkedInstallation bool
	version             string
}

func NewExec(execPath string, logger qb_log.ILogger) *CertbotExec {
	instance := new(CertbotExec)
	instance.execPath = execPath
	instance.logger = logger

	instance.SetRoot(".")

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *CertbotExec) SetLogger(logger qb_log.ILogger) *CertbotExec {
	instance.logger = logger
	return instance
}

func (instance *CertbotExec) SetRoot(dir string) *CertbotExec {
	instance.dirRoot = qb_utils.Paths.Concat(qb_utils.Paths.Absolute(dir), fsName)
	instance.dirWork = qb_utils.Paths.Concat(instance.dirRoot, "work")
	instance.dirConfig = qb_utils.Paths.Concat(instance.dirRoot, "config")
	instance.dirLogs = qb_utils.Paths.Concat(instance.dirRoot, "logging")
	instance.dirCerts = qb_utils.Paths.Concat(instance.dirRoot, "certificates")

	if instance.initialized {
		instance.initialized = false
		instance.init()
	}

	return instance
}

func (instance *CertbotExec) Root() string {
	return instance.dirRoot
}

func (instance *CertbotExec) Work() string {
	return instance.dirWork
}

func (instance *CertbotExec) GetPath(path string) (response string) {
	response = qb_utils.Paths.Absolutize(path, instance.dirRoot)
	_ = qb_utils.Paths.Mkdir(response)
	return
}

func (instance *CertbotExec) GetWorkPath(subPath string) (response string) {
	response = qb_utils.Paths.Absolutize(subPath, instance.dirWork)
	_ = qb_utils.Paths.Mkdir(response)
	return
}

func (instance *CertbotExec) LogFlush() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

func (instance *CertbotExec) LogDisableRotation() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.RotateEnable(false)
	}
}

func (instance *CertbotExec) FlushLog() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//	c o m m a n d s
// ---------------------------------------------------------------------------------------------------------------------

func (instance *CertbotExec) IsInstalled() bool {
	version, err := instance.checkInstall()
	if nil != err {
		return false
	}
	return len(version) > 0
}

func (instance *CertbotExec) Version() (response string, err error) {
	response, err = instance.checkInstall()
	return
}

func (instance *CertbotExec) CertOnly(email, domain string) (response string, err error) {
	dir := qb_utils.Paths.Concat(instance.dirCerts, qb_utils.Strings.Slugify(domain))
	_ = qb_utils.Paths.Mkdir(dir + qb_utils.OS_PATH_SEPARATOR)

	params := make([]string, 0)
	params = append(params, "certonly")
	params = append(params, "--standalone")
	//params = append(params, "--webroot", "-w", dir)
	params = append(params, "-d", domain)
	params = append(params, "--email", email)
	// append system params
	params = append(params, instance.params()...)
	response, err = instance.ExecuteCommand(params...)
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	g e n e r i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *CertbotExec) ExecuteCommand(arguments ...string) (out string, err error) {
	_, err = instance.checkInstall()
	if nil != err {
		return
	}
	return instance.exec(arguments...)
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *CertbotExec) checkInstall() (version string, err error) {
	if !instance.checkedInstallation {
		version, err = instance.exec("--version")
		if nil != err {
			err = errors.New(fmt.Sprintf("Please, install Certbot Client following setup instructions: %s", "https://certbot.eff.org/instructions"))
		} else {
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

func (instance *CertbotExec) params() []string {
	return []string{
		"--config-dir", instance.dirConfig,
		"--work-dir", instance.dirWork,
		"--logs-dir", instance.dirLogs,
		"--non-interactive", "--agree-tos", "--force-renewal",
	}
}

func (instance *CertbotExec) exec(arguments ...string) (out string, err error) {
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

func (instance *CertbotExec) program(args ...string) *executor.ConsoleProgram {
	instance.init()
	return executor.NewConsoleProgramWithDir(instance.execPath, instance.dirRoot, args...)
}

func (instance *CertbotExec) init() {
	if nil != instance && !instance.initialized {
		instance.initialized = true

		// creates paths
		_ = qb_utils.Paths.Mkdir(instance.dirWork + qb_utils.OS_PATH_SEPARATOR)
		_ = qb_utils.Paths.Mkdir(instance.dirConfig + qb_utils.OS_PATH_SEPARATOR)
		_ = qb_utils.Paths.Mkdir(instance.dirCerts + qb_utils.OS_PATH_SEPARATOR)
		_ = qb_utils.Paths.Mkdir(instance.dirLogs + qb_utils.OS_PATH_SEPARATOR)

		logPath := qb_utils.Paths.Concat(instance.dirLogs, "logging.log")
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

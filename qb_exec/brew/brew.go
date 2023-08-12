package brew

import (
	"os/exec"
	"strings"

	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_utils"
)

var (
	Brew *BrewHelper
	// path on M1 ARM is /opt/homebrew/bin/brew
	brewCommand = "brew"
)

const wpName = "brew"
const fsName = "./fs-brew"

type BrewHelper struct {
	root        string
	rootTmp     string
	logger      qb_log.ILogger
	initialized bool
}

func init() {
	brewCommand = findExecPath()
	Brew = new(BrewHelper)
	Brew.SetRoot(fsName)
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *BrewHelper) SetRoot(root string) {
	qb_utils.Paths.GetWorkspace(wpName).SetPath(root)
	instance.root = qb_utils.Paths.GetWorkspace(wpName).GetPath()
	instance.rootTmp = qb_utils.Paths.Concat(qb_utils.Paths.GetWorkspace(wpName).GetPath(), ".tmp")
	logPath := qb_utils.Paths.Concat(qb_utils.Paths.GetWorkspace(wpName).GetPath(), "logging.log")
	if nil != instance.logger {
		instance.logger.(*qb_log.Logger).SetFilename(logPath)
	} else {
		instance.logger = qb_log.Log.New("info", logPath)
		if l, ok := instance.logger.(*qb_log.Logger); ok {
			l.SetMessageFormat("* " + l.GetMessageFormat())
		}
	}
}

func (instance *BrewHelper) Root() string {
	return instance.root
}

func (instance *BrewHelper) Temp() string {
	return instance.rootTmp
}

func (instance *BrewHelper) GetTempPath(subPath string) (response string) {
	response = qb_utils.Paths.Absolutize(subPath, instance.rootTmp)
	_ = qb_utils.Paths.Mkdir(response)
	return
}

func (instance *BrewHelper) GetPath(path string) (response string) {
	response = qb_utils.Paths.Absolutize(path, instance.root)
	_ = qb_utils.Paths.Mkdir(response)
	return
}

func (instance *BrewHelper) LogFlush() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.Flush()
	}
}

func (instance *BrewHelper) LogDisableRotation() {
	if l, ok := instance.logger.(*qb_log.Logger); ok {
		l.RotateEnable(false)
	}
}

func (instance *BrewHelper) IsInstalled() bool {
	if nil != instance {
		return instance.NewExec().IsInstalled()
	}
	return false
}

func (instance *BrewHelper) Version() (version string, err error) {
	program := instance.NewExec()
	version, err = program.Version()
	if nil == err {
		// Homebrew 3.4.6-60-ge1c1157
		tokens := strings.Split(version, " ")
		if len(tokens) > 1 {
			version = tokens[1]
		}
	}
	return
}

// NewExec
// Creates new exec command with default password
func (instance *BrewHelper) NewExec() *BrewExec {
	if nil != instance {
		instance.init()
		return NewExec(brewCommand, instance.logger)
	}
	return nil
}

func (instance *BrewHelper) init() {
	if nil != instance && !instance.initialized {
		instance.initialized = true
		// creates paths
		_ = qb_utils.Paths.Mkdir(qb_utils.Paths.GetWorkspace(wpName).GetPath() + qb_utils.OS_PATH_SEPARATOR)
		_ = qb_utils.Paths.Mkdir(Brew.rootTmp + qb_utils.OS_PATH_SEPARATOR)
	}
}

// findExecPath tries to find the Chrome browser somewhere in the current
// system. It performs a rather aggressive search, which is the same in all systems.
func findExecPath() string {
	for _, path := range [...]string{
		// Unix-like
		"/opt/homebrew/bin/brew",        // arm
		"/user/local/Homebrew/bin/brew", // intel
		"brew",
	} {
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	// Fall back to something simple and sensible, to give a useful error
	// message.
	return "brew"
}

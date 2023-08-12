package chrome

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_utils"
)

var chromeCommand = "google-chrome"

// https://peter.sh/experiments/chromium-command-line-switches/
type ChromeParams struct {
	Url         string
	AppMode     bool
	WindowSize  string // fmt.Sprintf("%d,%d", width, height)
	DisableGPU  bool
	UserAgent   string
	ProxyServer string
}

type ChromeExec struct {
	dirWork string
	params  *ChromeParams
	session *executor.ConsoleProgramSession
}

func NewChromeExec() *ChromeExec {
	instance := new(ChromeExec)
	instance.dirWork = qb_utils.Paths.Absolute("./")
	instance.params = new(ChromeParams)
	instance.params.AppMode = true
	instance.params.WindowSize = "800x600"

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ChromeExec) SetDir(dir string) {
	instance.dirWork = dir
}

func (instance *ChromeExec) GetDir() string {
	return instance.dirWork
}

func (instance *ChromeExec) Params() *ChromeParams {
	return instance.params
}

func (instance *ChromeExec) Open(url string) error {
	instance.params.Url = url
	session, err := instance.program().Run()
	if nil == err {
		instance.session = session
	}
	return err
}

func (instance *ChromeExec) OpenApp(url string) error {
	instance.params.AppMode = true
	return instance.Open(url)
}

func (instance *ChromeExec) OpenBrowser(url string) error {
	instance.params.AppMode = false
	return instance.Open(url)
}

func (instance *ChromeExec) Pid() int {
	if nil != instance && nil != instance.session {
		return instance.session.PidLatest()
	}
	return 0
}

func (instance *ChromeExec) TryKill() error {
	if nil != instance && nil != instance.session {
		return instance.session.Kill()
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ChromeExec) program() *executor.ConsoleProgram {
	args := parseParams(instance.params)
	return executor.NewConsoleProgramWithDir(chromeCommand, instance.dirWork, args...)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func NewChromeProgram(params *ChromeParams) *executor.ConsoleProgram {
	args := parseParams(params)
	return executor.NewConsoleProgram(chromeCommand, args...)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func init() {
	chromeCommand = findExecPath()
}

// findExecPath tries to find the Chrome browser somewhere in the current
// system. It performs a rather aggressive search, which is the same in all systems.
func findExecPath() string {
	for _, path := range [...]string{
		// Unix-like
		"headless_shell",
		"headless-shell",
		"chromium",
		"chromium-browser",
		"google-chrome",
		"google-chrome-stable",
		"google-chrome-beta",
		"google-chrome-unstable",
		"/usr/bin/google-chrome",

		// Windows
		"chrome",
		"chrome.exe", // in case PATHEXT is misconfigured
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\Google\Chrome\Application\chrome.exe`),

		// Mac
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
	} {
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	// Fall back to something simple and sensible, to give a useful error
	// message.
	return "google-chrome"
}

func parseParams(params *ChromeParams) []string {
	args := make([]string, 0)
	args = append(args, "--new-window")

	address := params.Url
	if params.AppMode {
		address = "--app=" + address
	}
	args = append(args, address)
	if len(params.WindowSize) > 0 {
		args = append(args, "--window-size="+params.WindowSize)
	}
	if params.DisableGPU {
		args = append(args, fmt.Sprintf("--disable-gpu=%v", true))
	}
	if len(params.UserAgent) > 0 {
		args = append(args, fmt.Sprintf("--user-agent=%v", params.UserAgent))
	}
	if len(params.ProxyServer) > 0 {
		args = append(args, fmt.Sprintf("--proxy-server=%v", params.ProxyServer))
	}
	return args
}

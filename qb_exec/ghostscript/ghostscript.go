package ghostscript

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_utils"
)

// DOWNLOAD: https://www.ghostscript.com/releases/gsdnld.html
// DOCS: https://ghostscript.com/doc/current/Use.htm
// SAMPLES: https://github.com/MrSaints/go-ghostscript/blob/master/examples/pdf-concat/main.go
//
//	https://gist.github.com/brenopolanski/2ae13095ed7e865b60f5
var (
	_ghostscriptCommand = "gs"

	optionsSilent = []string{
		"-dBATCH",
		"-dNOPROMPT",
		"-dNOPAUSE",
		"-dQUIET",
	}
)

type GhostScriptExec struct {
	command string
	dirWork string
	dirOut  string
	session *executor.ConsoleProgramSession
}

func NewGhostScriptExec() *GhostScriptExec {
	instance := new(GhostScriptExec)
	instance.command = _ghostscriptCommand
	instance.dirWork = qb_utils.Paths.Absolute("")

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *GhostScriptExec) SetDir(dir string) {
	instance.dirWork = qb_utils.Paths.Absolute(dir)
}

func (instance *GhostScriptExec) GetDir() string {
	return instance.dirWork
}

func (instance *GhostScriptExec) SetDirOut(dir string) {
	instance.dirOut = qb_utils.Paths.Absolute(dir)
}

func (instance *GhostScriptExec) GetDirOut() string {
	return instance.dirOut
}

func (instance *GhostScriptExec) TryKill() error {
	if nil != instance && nil != instance.session {
		return instance.session.Kill()
	}
	return nil
}

func (instance *GhostScriptExec) Pid() int {
	if nil != instance && nil != instance.session {
		return instance.session.PidLatest()
	}
	return 0
}

func (instance *GhostScriptExec) IsInstalled() bool {
	version, err := instance.CmdVersion()
	if nil != err {
		return false
	}
	return len(version) > 0
}

//----------------------------------------------------------------------------------------------------------------------
//	c o m m a n d   t o o l
//----------------------------------------------------------------------------------------------------------------------

func (instance *GhostScriptExec) CmdHelp() (string, error) {
	if nil != instance {
		session, err := instance.program().Run("--help")
		if nil != err {
			return "", err
		}
		instance.session = session
		return instance.session.StdOut(), nil
	}
	return "", nil
}

func (instance *GhostScriptExec) CmdVersion() (string, error) {
	if nil != instance {
		session, err := instance.program().Run("--version")
		if nil != err {
			return "", err
		}
		instance.session = session
		out := strings.TrimSpace(instance.session.StdOut())
		return out, nil
	}
	return "", nil
}

// CmdExecRaw execute a command. i.e. "-sDEVICE=pdfwrite -dBATCH -dNOPROMPT -dNOPAUSE -dQUIET -sOwnerPassword=mypassword -sUserPassword=mypassword -sOutputFile=MyOutputFile.pdf MyInputFile.pdf"
func (instance *GhostScriptExec) CmdExecRaw(arguments ...string) (out string, err error) {
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
		out = session.StdOut()
	}
	return
}

// CmdExecQuiet execute a command adding some argument to avoid prompt or logs
func (instance *GhostScriptExec) CmdExecQuiet(arguments ...string) (out string, err error) {
	options := append(make([]string, 0), optionsSilent...)
	options = append(options, arguments...)
	return instance.CmdExecRaw(options...)
}

// CmdPasswordProtect Protect a PDF with password
// gs -sDEVICE=pdfwrite -dBATCH -dNOPROMPT -dNOPAUSE -dQUIET -sOwnerPassword=mypassword -sUserPassword=mypassword -sOutputFile=MyOutputFile.pdf MyInputFile.pdf
func (instance *GhostScriptExec) CmdPasswordProtect(sourceFile, password string, optTargetFileName string) (targetFile string, err error) {
	if nil != instance {
		// normalize paths
		sourceFile = qb_utils.Paths.Absolutize(sourceFile, instance.dirWork)
		targetFile = optTargetFileName
		if len(targetFile) == 0 {
			targetFile = qb_utils.Paths.ChangeFileNameWithPrefix(sourceFile, "protected_")
		}
		// check file exists
		if ok, pathErr := qb_utils.Paths.Exists(sourceFile); ok {
			program := instance.program()
			options := append(make([]string, 0), optionsSilent...)
			options = append(options, []string{
				"-sDEVICE=pdfwrite",
				"-sOwnerPassword=" + password,
				"-sUserPassword=" + password,
				"-sOutputFile=" + targetFile,
				sourceFile,
			}...)
			session, e := program.Run(options...)
			if nil != e {
				err = e
			} else {
				stdErr := session.StdErr()
				if len(stdErr) > 0 {
					err = errors.New(stdErr)
				}
			}
		} else {
			err = pathErr
		}
	}
	return
}

func (instance *GhostScriptExec) CmdExportToImage(sourceFile, format string, optResolution, optTargetDir, optPassword string) (targetFiles []string, err error) {
	targetFiles = make([]string, 0)
	if nil != instance {
		// normalize paths
		sourceFile = strings.ToLower(qb_utils.Paths.Absolutize(sourceFile, instance.dirWork))
		if len(optTargetDir) == 0 {
			optTargetDir = qb_utils.Paths.Dir(sourceFile)
		}
		if len(format) == 0 {
			format = "jpg"
		}
		if len(optResolution) == 0 {
			optResolution = "100x100"
		}
		ext := qb_utils.Paths.ExtensionName(sourceFile)
		name := qb_utils.Paths.FileName(sourceFile, false)
		// check file exists
		if ok, pathErr := qb_utils.Paths.Exists(sourceFile); ok {
			options := append(make([]string, 0), optionsSilent...)
			if len(optPassword) > 0 {
				switch ext {
				case "pdf":
					options = append(options, "-sPDFPassword="+optPassword)
				}
			}
			options = append(options, []string{
				"-sDEVICE=png16m",
				fmt.Sprintf("-r%s", optResolution),
				fmt.Sprintf("-sOutputFile=%s-%%03d.%s", name, format),
				sourceFile,
			}...)
			program := instance.program()
			session, e := program.Run(options...)
			if nil != e {
				err = e
			} else {
				stdErr := session.StdErr()
				if len(stdErr) > 0 {
					err = errors.New(stdErr)
				} else {
					// lookup files
					files, _ := qb_utils.Paths.ListFiles(optTargetDir, fmt.Sprintf("%s-*.%s", name, format))
					if nil != files {
						targetFiles = append(targetFiles, files...)
					}
				}
			}
		} else {
			err = pathErr
		}
	}
	return
}

func (instance *GhostScriptExec) CmdExportToText(sourceFile string, optTargetDir, optPassword string) (out string, err error) {
	if nil != instance {
		// normalize paths
		sourceFile = strings.ToLower(qb_utils.Paths.Absolutize(sourceFile, instance.dirWork))
		if len(optTargetDir) == 0 {
			optTargetDir = qb_utils.Paths.Dir(sourceFile)
		}
		name := qb_utils.Paths.FileName(sourceFile, false)
		outputFile := qb_utils.Paths.Concat(optTargetDir, fmt.Sprintf("%s.txt", name))
		if ok, pathErr := qb_utils.Paths.Exists(sourceFile); ok {
			program := instance.program()
			options := append(make([]string, 0), optionsSilent...)
			if len(optPassword) > 0 {
				options = append(options, "-sPDFPassword="+optPassword)
			}
			options = append(options, []string{
				"-sDEVICE=txtwrite",
				fmt.Sprintf("-sOutputFile=%s", outputFile),
				sourceFile,
			}...)
			session, e := program.Run(options...)
			if nil != e {
				err = e
			} else {
				stdErr := session.StdErr()
				if len(stdErr) > 0 {
					err = errors.New(stdErr)
				} else {
					out = outputFile
				}
			}
		} else {
			err = pathErr
		}
	}
	return
}

func (instance *GhostScriptExec) CmdChainDocuments(sourceFiles []string, optTargetFileName, optPassword string) (targetFile string, err error) {
	if nil != instance && len(sourceFiles) > 0 {
		ext := qb_utils.Paths.ExtensionName(sourceFiles[0])
		targetFile = optTargetFileName
		if len(targetFile) == 0 {
			targetFile = qb_utils.Paths.Absolutize(fmt.Sprintf("combined.%s", ext), instance.dirWork)
		} else {
			targetFile = qb_utils.Paths.Absolutize(targetFile, instance.dirWork)
		}
		options := append(make([]string, 0), optionsSilent...)
		if len(optPassword) > 0 {
			switch ext {
			case "pdf":
				options = append(options, "-sPDFPassword="+optPassword)
			}
		}
		options = append(options, []string{
			"-sDEVICE=pdfwrite",
			"-sOutputFile=" + targetFile,
			"-dBATCH",
		}...)
		for _, sourceFile := range sourceFiles {
			options = append(options, sourceFile)
		}
		program := instance.program()
		session, e := program.Run(options...)
		if nil != e {
			err = e
		} else {
			stdErr := session.StdErr()
			if len(stdErr) > 0 {
				err = errors.New(stdErr)
			}
		}
	}
	return
}

// ----------------------------------------------------------------------------------------------------------------------
//
//	p r i v a t e
//
// ----------------------------------------------------------------------------------------------------------------------

func (instance *GhostScriptExec) getCommand() string {
	if nil != instance {
		return instance.command
	}
	return ""
}

func (instance *GhostScriptExec) program(args ...string) *executor.ConsoleProgram {
	return executor.NewConsoleProgramWithDir(instance.getCommand(), instance.dirWork, args...)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func init() {
	_ghostscriptCommand = findExecPath()
}

// findExecPath tries to find the Chrome browser somewhere in the current
// system. It performs a rather aggressive search, which is the same in all systems.
func findExecPath() string {
	for _, path := range [...]string{
		// Unix-like
		"gs",
		"/usr/bin/gs",

		// Windows
		"gs.exe", // in case PATHEXT is misconfigured
		`C:\Program Files (x86)\LibreOffice\program\gs.exe`,
		`C:\Program Files\LibreOffice\program\gs.exe`,
		filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\LibreOffice\Application\gs.exe`),

		// Mac
		"/Applications/GhostScript.app/Contents/MacOS/gs",
	} {
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	// Fall back to something simple and sensible, to give a useful error
	// message.
	return "gs"
}

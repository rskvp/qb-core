package libreoffice

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_utils"
)

// https://askubuntu.com/questions/519082/how-to-install-libre-office-without-gui
// https://help.libreoffice.org/Common/Starting_the_Software_With_Parameters

var libreOfficeCommand string = "soffice"

type LibreOfficeExec struct {
	dirWork string
	dirOut  string
	session *executor.ConsoleProgramSession
}

func NewLibreOfficeExec() *LibreOfficeExec {
	instance := new(LibreOfficeExec)
	instance.dirWork = qb_utils.Paths.Absolute("./")

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *LibreOfficeExec) SetDir(dir string) {
	instance.dirWork = qb_utils.Paths.Absolute(dir)
}

func (instance *LibreOfficeExec) GetDir() string {
	return instance.dirWork
}

func (instance *LibreOfficeExec) SetDirOut(dir string) {
	instance.dirOut = qb_utils.Paths.Absolute(dir)
}

func (instance *LibreOfficeExec) GetDirOut() string {
	return instance.dirOut
}

func (instance *LibreOfficeExec) TryKill() error {
	if nil != instance && nil != instance.session {
		return instance.session.Kill()
	}
	return nil
}

func (instance *LibreOfficeExec) Pid() int {
	if nil != instance && nil != instance.session {
		return instance.session.PidLatest()
	}
	return 0
}

func (instance *LibreOfficeExec) OpenDoc(filename string) error {
	session, err := instance.program(filename).Run()
	if nil != err {
		return err
	}
	instance.session = session
	return err
}

func (instance *LibreOfficeExec) IsInstalled() bool {
	version, err := instance.CmdVersion()
	if nil != err {
		return false
	}
	return len(version) > 0
}

//----------------------------------------------------------------------------------------------------------------------
//	c o m m a n d   t o o l
//----------------------------------------------------------------------------------------------------------------------

func (instance *LibreOfficeExec) CmdHelp() (string, error) {
	session, err := instance.program().Run("--headless", "--help")
	if nil != err {
		return "", err
	}
	instance.session = session
	return instance.session.StdOut(), nil
}

func (instance *LibreOfficeExec) CmdVersion() (string, error) {
	session, err := instance.program().Run("--headless", "--version")
	if nil != err {
		return "", err
	}
	instance.session = session
	return instance.session.StdOut(), nil
}

// CmdConvertTo run command "$ libreoffice --convert-to docx file.txt"
func (instance *LibreOfficeExec) CmdConvertTo(sourceFile, targetFormat string) (string, error) {
	outDir := ""
	ext := qb_utils.Paths.ExtensionName(targetFormat)
	if len(ext) == 0 {
		ext = targetFormat
	} else {
		outDir = qb_utils.Paths.Dir(targetFormat)
		if !qb_utils.Paths.IsAbs(outDir) {
			outDir = qb_utils.Paths.Concat(instance.dirWork, outDir)
		}
		_ = qb_utils.Paths.Mkdir(outDir + qb_utils.OS_PATH_SEPARATOR)
	}
	if len(outDir) > 0 {
		// outdir
		return instance.CmdConvertToDir(sourceFile, ext, outDir)
	} else {
		return instance.CmdConvertToDir(sourceFile, ext, instance.dirOut)
	}
}

func (instance *LibreOfficeExec) CmdConvertToDir(sourceFile, targetFormat, outDir string) (string, error) {
	program := instance.program("--headless", "--convert-to")
	var session *executor.ConsoleProgramSession
	var err error
	if len(outDir) > 0 {
		// outdir
		session, err = program.Run(targetFormat, sourceFile, "--outdir", outDir)
	} else {
		session, err = program.Run(targetFormat, sourceFile)
	}
	if nil != err {
		return "", err
	}
	instance.session = session
	return instance.session.StdOut(), nil
}

func (instance *LibreOfficeExec) CmdConvertToTxt(sourceFile string) (string, error) {
	return instance.CmdConvertTo(sourceFile, "txt")
}

func (instance *LibreOfficeExec) CmdConvertToDocx(sourceFile string) (string, error) {
	return instance.CmdConvertTo(sourceFile, "docx")
}

func (instance *LibreOfficeExec) CmdConvertToEpub(sourceFile string) (string, error) {
	return instance.CmdConvertTo(sourceFile, "epub")
}

func (instance *LibreOfficeExec) CmdConvertToPdf(sourceFile string) (string, error) {
	return instance.CmdConvertTo(sourceFile, "pdf")
}

// CmdPrint This option prints to the default printer without opening LibreOffice; it just sends the document to your printer.
func (instance *LibreOfficeExec) CmdPrint(sourceFile string) (string, error) {
	session, err := instance.program().Run("--headless", "-p", sourceFile)
	if nil != err {
		return "", err
	}
	instance.session = session
	return instance.session.StdOut(), nil
}

// CmdCat Dump text content of the following files to console
func (instance *LibreOfficeExec) CmdCat(sourceFile string) (string, error) {
	session, err := instance.program().Run("--headless", "--cat", sourceFile)
	if nil != err {
		return "", err
	}
	instance.session = session
	return instance.session.StdOut(), nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *LibreOfficeExec) program(args ...string) *executor.ConsoleProgram {
	return executor.NewConsoleProgramWithDir(libreOfficeCommand, instance.dirWork, args...)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func init() {
	libreOfficeCommand = findExecPath()
}

// findExecPath tries to find the Chrome browser somewhere in the current
// system. It performs a rather aggressive search, which is the same in all systems.
func findExecPath() string {
	for _, path := range [...]string{
		// Unix-like
		"libreoffice",
		"soffice",
		"/usr/bin/libreoffice",

		// Windows
		"libreoffice.exe", // in case PATHEXT is misconfigured
		"soffice.exe",
		`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
		`C:\Program Files\LibreOffice\program\soffice.exe`,
		filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\LibreOffice\Application\soffice.exe`),

		// Mac
		"/Applications/LibreOffice.app/Contents/MacOS/soffice",
	} {
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	// Fall back to something simple and sensible, to give a useful error
	// message.
	return "soffice"
}

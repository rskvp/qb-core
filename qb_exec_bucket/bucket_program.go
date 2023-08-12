package qb_exec_bucket

import (
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_exec/executor"
	"github.com/rskvp/qb-core/qb_rnd"
	"github.com/rskvp/qb-core/qb_utils"
)

type BucketProgram struct {
	root           string
	tmpDir         string
	execPath       string // i.e. "ts-node", "python", "node", "brew", ecc..
	programFile    string
	executableFile string

	handlerTimeout func()
	handlerSuccess func(response interface{})
	handlerError   func(err error)
	handlerFinish  func(response interface{}, err error)
	timeout        time.Duration

	session *executor.ConsoleProgramSession // current session
}

func newBucketProgram(exec *BucketExec, programFile string) (*BucketProgram, error) {
	dirTemp := exec.dirController.DirTemp()
	dirWork := exec.dirController.DirWork()
	programFile = qb_utils.Paths.Absolutize(programFile, dirWork)
	executableFile, err := copyFileToExecutable(programFile, dirTemp)
	return &BucketProgram{
		root:           dirWork,
		tmpDir:         dirTemp,
		execPath:       exec.execPath,
		programFile:    programFile,
		executableFile: executableFile,
		handlerTimeout: exec.handlerTimeout,
		handlerSuccess: exec.handlerSuccess,
		handlerError:   exec.handlerError,
		handlerFinish:  exec.handlerFinish,
		timeout:        exec.timeout,
	}, err
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Merge
// Format a parametrized script replacing tags with model data
func (instance *BucketProgram) Merge(model map[string]interface{}) error {
	if nil != instance {
		source := instance.executableFile
		target := instance.executableFile
		if len(source) == 0 {
			source = instance.programFile
			target = getTargetFile(source, instance.tmpDir)
			instance.executableFile = target
		}
		if len(source) > 0 {
			text, err := qb_utils.IO.ReadTextFromFile(source)
			if nil != err {
				return err
			}
			text, err = qb_utils.Formatter.MergeText(text, model)
			if nil != err {
				return err
			}

			_, err = qb_utils.IO.WriteTextToFile(text, target)
			if nil != err {
				return err
			}
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	run and session
//----------------------------------------------------------------------------------------------------------------------

func (instance *BucketProgram) TryKill() error {
	if nil != instance && nil != instance.session {
		return instance.session.Kill()
	}
	return nil
}

func (instance *BucketProgram) Pid() int {
	if nil != instance && nil != instance.session {
		return instance.session.PidLatest()
	}
	return 0
}

// RunUnboxed
// execute the command out of a bucket wrapper
func (instance *BucketProgram) RunUnboxed(args ...interface{}) (out string, err error) {
	return instance.run(args...)
}

// RunSync
// Run program and wait
func (instance *BucketProgram) RunSync(args ...interface{}) (elapsed int, response interface{}, err error) {
	task := instance.Run(args...)
	return task.Wait()
}

func (instance *BucketProgram) Run(args ...interface{}) (task *qb_utils.AsyncTask) {
	task = qb_utils.Async.NewAsyncTask()
	task.SetTimeout(instance.timeout)
	task.OnError(instance.handlerError)
	task.OnSuccess(instance.handlerSuccess)
	task.OnTimeout(instance.handlerTimeout)
	task.OnFinish(instance.handlerFinish)
	task.Run(func(ctx *qb_utils.AsyncContext, params ...interface{}) (interface{}, error) {
		return instance.run(params...)
	}, args...)
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *BucketProgram) program(args ...string) *executor.ConsoleProgram {
	return executor.NewConsoleProgramWithDir(instance.execPath, instance.root, args...)
}

func (instance *BucketProgram) run(args ...interface{}) (out string, err error) {
	// create params
	params := make([]string, 0)
	if len(instance.executableFile) > 0 {
		instance.executableFile = qb_utils.Paths.Absolutize(instance.executableFile, instance.root)
		params = append(params, instance.executableFile)
	}
	params = append(params, qb_utils.Convert.ToArrayOfString(args)...)
	// prepare the session of program
	session, e := instance.program().Prepare(params...)
	if nil != e {
		err = e
		return
	}

	// execute
	instance.session = session
	_, err = session.Run()
	instance.session = nil

	if nil == err {
		out = strings.TrimSpace(session.StdOut())
		session.Close()
	}

	if len(instance.executableFile) > 0 {
		// delete
		_ = qb_utils.IO.Remove(instance.executableFile)
	}

	return
}

func copyFileToExecutable(source, targetDir string) (string, error) {
	target := getTargetFile(source, targetDir)
	_ = qb_utils.Paths.Mkdir(target)
	_, err := qb_utils.IO.CopyFile(source, target)
	if nil == err {
		return target, nil
	}
	return "", err
}

func getTargetFile(source, targetDir string) string {
	targetFile := qb_utils.Paths.ChangeFileNameWithPrefix(qb_utils.Paths.FileName(source, true), qb_rnd.Rnd.Uuid()+"-")
	return qb_utils.Paths.Concat(targetDir, targetFile)
}

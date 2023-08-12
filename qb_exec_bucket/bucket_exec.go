package qb_exec_bucket

import (
	"time"

	"github.com/rskvp/qb-core/qb_exec"
	"github.com/rskvp/qb-core/qb_exec/git"
	"github.com/rskvp/qb-core/qb_exec/npm"
	"github.com/rskvp/qb-core/qb_exec/yarn"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_utils"
)

// BucketExec
// main executable container
type BucketExec struct {
	uid           string
	dirController *qb_utils.DirCentral
	// root      string
	// dirBucket string // "source program" home and download target
	// dirTemp   string
	execPath string // i.e. "ts-node", "python", "node", "brew", ecc..

	logger      qb_log.ILogger
	downloader  *BucketResourceDownloader
	initialized bool

	handlerTimeout func()
	handlerSuccess func(response interface{})
	handlerError   func(err error)
	handlerFinish  func(response interface{}, err error)
	timeout        time.Duration
}

func NewBucketExec(root, execPath string, global bool) (instance *BucketExec) {
	root = qb_utils.Paths.Absolute(root)
	var uid, dirBucket, dirTemp string
	if global {
		uid = ""
		dirBucket = qb_utils.Paths.Concat(root, ".bucket")
	} else {
		uid = qb_utils.Coding.MD5(root + execPath)
		dirBucket = qb_utils.Paths.Concat(root, ".bucket-"+uid)
	}
	dirTemp = qb_utils.Paths.Concat(dirBucket, ".tmp")

	instance = new(BucketExec)
	instance.uid = uid
	instance.dirController = qb_utils.Dir.NewCentral(dirBucket, dirTemp, true)
	instance.dirController.SetRoot(root)
	instance.execPath = execPath
	instance.downloader = NewBucketResourceDownloader(dirBucket)

	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *BucketExec) SetTimeout(timeout time.Duration) *BucketExec {
	if nil != instance {
		instance.timeout = timeout
	}
	return instance
}

func (instance *BucketExec) SetRoot(root string) *BucketExec {
	if nil != instance {
		instance.dirController.SetRoot(root)
	}
	return instance
}

func (instance *BucketExec) OnTimeout(callback func()) *BucketExec {
	if nil != instance {
		instance.handlerTimeout = callback
	}
	return instance
}

func (instance *BucketExec) OnSuccess(callback func(response interface{})) *BucketExec {
	if nil != instance {
		instance.handlerSuccess = callback
	}
	return instance
}

func (instance *BucketExec) OnError(callback func(err error)) *BucketExec {
	if nil != instance {
		instance.handlerError = callback
	}
	return instance
}

func (instance *BucketExec) OnFinish(callback func(response interface{}, err error)) *BucketExec {
	if nil != instance {
		instance.handlerFinish = callback
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	d o w n l o a d
//----------------------------------------------------------------------------------------------------------------------

func (instance *BucketExec) DownloaderSetRemoteRoot(remotePath string) *BucketExec {
	if nil != instance && nil != instance.downloader {
		instance.downloader.SetRemoteRoot(remotePath)
	}
	return instance
}

func (instance *BucketExec) DownloaderAddResource(remotePath, localRelativePath string) *BucketExec {
	if nil != instance && nil != instance.downloader {
		instance.downloader.AddResource(remotePath, localRelativePath)
	}
	return instance
}

func (instance *BucketExec) DownloaderRun(force bool) ([]string, []error) {
	if nil != instance && nil != instance.downloader {
		return instance.downloader.DownloadAll(force)
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	program
//----------------------------------------------------------------------------------------------------------------------

func (instance *BucketExec) NewProgram(programFile string) (*BucketProgram, error) {
	if nil != instance {
		instance.init()
		return newBucketProgram(instance, programFile)
	}
	return nil, nil
}

func (instance *BucketExec) Run(programFile string, model map[string]interface{}, args ...interface{}) (pid int, task *qb_utils.AsyncTask, err error) {
	if nil != instance {
		program, e := instance.NewProgram(programFile)
		if nil != e {
			err = e
			return
		}
		if nil != model {
			err = program.Merge(model)
			if nil != err {
				return
			}
		}
		task = program.Run(args...)
		time.Sleep(10 * time.Millisecond)
		pid = program.Pid()
	}
	return
}

func (instance *BucketExec) RunSync(programFile string, model map[string]interface{}, args ...interface{}) (pid, elapsed int, response interface{}, err error) {
	if nil != instance {
		p, task, e := instance.Run(programFile, model, args...)
		if nil != e {
			err = e
			return
		}
		pid = p
		elapsed, response, err = task.Wait()
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	u t i l i t y
//----------------------------------------------------------------------------------------------------------------------

func (instance *BucketExec) NewNpm() *npm.NpmExec {
	if nil != instance {
		instance.init()
		dirBucket := instance.dirController.DirWork()
		dirTemp := instance.dirController.DirTemp()
		return qb_exec.Exec.NewNpm().SetRoot(dirBucket).SetTemp(dirTemp).SetLogger(instance.logger)
	}
	return nil
}

func (instance *BucketExec) NewYarn() *yarn.YarnExec {
	if nil != instance {
		instance.init()
		dirBucket := instance.dirController.DirWork()
		dirTemp := instance.dirController.DirTemp()
		return qb_exec.Exec.NewYarn().SetRoot(dirBucket).SetTemp(dirTemp).SetLogger(instance.logger)
	}
	return nil
}

func (instance *BucketExec) NewGit() *git.GitExec {
	if nil != instance {
		instance.init()
		dirBucket := instance.dirController.DirWork()
		dirTemp := instance.dirController.DirTemp()
		return qb_exec.Exec.NewGit().SetRoot(dirBucket).SetTemp(dirTemp).SetLogger(instance.logger)
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *BucketExec) init() {
	if nil != instance {
		if !instance.initialized {
			instance.initialized = true
			instance.dirController.Refresh()

			logPath := instance.dirController.PathLog()
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
}

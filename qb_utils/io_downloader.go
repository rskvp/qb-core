package qb_utils

import (
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rskvp/qb-core/qb_rnd"
)

//----------------------------------------------------------------------------------------------------------------------
//	DownloadSession
//----------------------------------------------------------------------------------------------------------------------

type DownloadSession struct {
	pool    *ConcurrentPool
	mux     sync.Mutex
	actions map[string]*DownloaderAction
	files   []string
	errs    []error
}

func newDownloadSession(actions interface{}) *DownloadSession {
	instance := new(DownloadSession)
	instance.pool = Async.NewConcurrentPool(10)
	instance.files = make([]string, 0)
	instance.errs = make([]error, 0)
	instance.actions = make(map[string]*DownloaderAction)

	if m, b := actions.(map[string]*DownloaderAction); b {
		instance.actions = m
	} else if a, b := actions.([]*DownloaderAction); b {
		for _, v := range a {
			if len(v.Uid) == 0 {
				v.Uid = qb_rnd.Rnd.Uuid()
			}
			instance.actions[v.Uid] = v
		}
	}

	return instance
}

func (instance *DownloadSession) DownloadAll(force bool) ([]string, []error) {
	for _, v := range instance.actions {
		source := v.Source
		sourceVersion := v.SourceVersion
		target := v.Target
		// time.Sleep(time.Duration(lygo_rnd.Between(100, 2000)) * time.Millisecond)
		_ = instance.pool.RunArgs(func(args ...interface{}) error {
			so := args[0].(string)
			sv := args[1].(string)
			ta := args[2].(string)
			fo := args[3].(bool)
			f, e := download(so, sv, ta, fo)
			instance.mux.Lock()
			instance.files = append(instance.files, f...)
			instance.mux.Unlock()
			return e
		}, source, sourceVersion, target, force)
	}
	time.Sleep(300 * time.Millisecond)
	pe := instance.pool.Wait()
	if nil != pe && nil != pe.Errors {
		instance.errs = append(instance.errs, pe.Errors...)
	}
	return instance.files, instance.errs
}

//----------------------------------------------------------------------------------------------------------------------
//	Downloader
//----------------------------------------------------------------------------------------------------------------------

type DownloaderAction struct {
	Uid           string `json:"uid"`
	Source        string `json:"source"`
	SourceVersion string `json:"source-version"`
	Target        string `json:"target"`
}

type Downloader struct {
	Actions map[string]*DownloaderAction
}

func newDownloader() *Downloader {
	instance := new(Downloader)
	instance.Actions = make(map[string]*DownloaderAction)
	return instance
}

func newAction(uid, source, sourceversion, target string) *DownloaderAction {
	action := new(DownloaderAction)
	if len(uid) == 0 {
		uid = qb_rnd.Rnd.Uuid()
	}
	action.Uid = uid
	action.Source = source
	action.SourceVersion = sourceversion
	action.Target = target
	return action
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Downloader) Names() []string {
	response := make([]string, 0, len(instance.Actions))
	for k := range instance.Actions {
		response = append(response, k)
	}
	return response
}

func (instance *Downloader) PutAction(uid, source, sourceversion, target string) *Downloader {
	instance.Put(newAction(uid, source, sourceversion, target))
	return instance
}

func (instance *Downloader) Put(action *DownloaderAction) *Downloader {
	if nil != action {
		if len(action.Uid) == 0 {
			action.Uid = qb_rnd.Rnd.Uuid()
		}
		instance.Actions[action.Uid] = action
	}
	return instance
}

func (instance *Downloader) Download(uid string) (files []string, err error) {
	if v, b := instance.Actions[uid]; b {
		return download(v.Source, v.SourceVersion, v.Target, false)
	}
	return []string{}, nil
}

func (instance *Downloader) ForceDownload(uid string) (files []string, err error) {
	if v, b := instance.Actions[uid]; b {
		return download(v.Source, v.SourceVersion, v.Target, true)
	}
	return
}

func (instance *Downloader) DownloadAll() ([]string, []error) {
	names := instance.Names()
	if len(names) == 1 {
		files, err := instance.Download(names[0])
		if nil != err {
			return files, []error{err}
		}
		return files, []error{}
	} else {
		session := newDownloadSession(instance.Actions)
		return session.DownloadAll(false)
	}
}

func (instance *Downloader) ForceDownloadAll() ([]string, []error) {
	names := instance.Names()
	if len(names) == 1 {
		files, err := instance.ForceDownload(names[0])
		if nil != err {
			return files, []error{err}
		}
		return files, []error{}
	} else {
		session := newDownloadSession(instance.Actions)
		return session.DownloadAll(true)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func isZip(source string) bool {
	return strings.ToLower(Paths.ExtensionName(source)) == "zip"
}

func isGreaterThan(rv, lv string) bool {
	rv = strings.Trim(rv, " \n")
	lv = strings.Trim(lv, " \n")
	r := Convert.ToInt(strings.ReplaceAll(rv, ".", ""))
	l := Convert.ToInt(strings.ReplaceAll(lv, ".", ""))
	return r > l
}

func isDownloadable(source, sourceVersion, target string) bool {
	if len(sourceVersion) > 0 {
		// check VERSION
		dir := target
		if b, _ := Paths.IsFile(dir); b {
			dir = Paths.Dir(dir)
		}
		versionFile := Paths.Concat(dir, Paths.FileName(sourceVersion, true))
		data, err := IO.Download(sourceVersion)
		if nil == err {
			remoteVersion := string(data)
			localVersion, err := IO.ReadTextFromFile(versionFile)
			if nil == err {
				return isGreaterThan(remoteVersion, localVersion)
			}
			// write version file
			_, _ = IO.WriteBytesToFile(data, versionFile)
		}
	} else {
		// check EXISTS
		if isZip(source) {
			content, _ := Paths.ReadDir(target)
			return len(content) == 0
		} else {
			targetFile := target
			if b, _ := Paths.IsFile(targetFile); !b {
				targetFile = Paths.Concat(targetFile, Paths.FileName(source, true))
			}
			if b, _ := Paths.Exists(targetFile); b {
				return false
			}
		}
	}
	return true
}

func download(source, sourceversion, target string, force bool) (files []string, err error) {
	target = Paths.Absolute(target)
	dirTarget := mkdir(target)

	if force || isDownloadable(source, sourceversion, target) {
		var bytes []byte
		bytes, err = IO.Download(source)
		if nil != err {
			return nil, err
		}
		if !isValidContent(bytes) {
			return nil, Errors.Prefix(errors.New("invalid_content"), "Download error: ")
		}

		if isZip(source) {
			tmp := Paths.Concat(dirTarget, qb_rnd.Rnd.Uuid()+".tmp")
			defer remove(tmp)
			_, err = IO.WriteBytesToFile(bytes, tmp)
			if nil == err {
				files, err = Zip.Unzip(tmp, dirTarget)
			}
		} else {
			targetFile := target
			if b, _ := Paths.IsFile(targetFile); !b {
				targetFile = Paths.Concat(targetFile, Paths.FileName(source, true))
			}

			// write a file
			_, err = IO.WriteBytesToFile(bytes, targetFile)
			files = append(files, target)
		}
		return
	}
	return
}

func remove(path string) {
	err := IO.Remove(path)
	if nil == err {
		parent := Paths.Dir(path)
		files, _ := Paths.ReadDir(parent)
		if len(files) == 0 {
			_ = IO.RemoveAll(parent)
		}
	}
}

func mkdir(path string) string {
	if len(filepath.Ext(path)) > 0 {
		_ = Paths.Mkdir(path)
		return Paths.Dir(path)
	} else {
		_ = Paths.Mkdir(Paths.ConcatDir(path, OS_PATH_SEPARATOR))
		return path
	}
}

func isValidContent(bytes []byte) bool {
	if len(bytes) > 10 {
		start := strings.TrimSpace(string(bytes[:10]))
		if strings.HasPrefix(start, "<!") {
			// is html. check html is not a bitbucket error
			content := string(bytes)
			if strings.Index(content, "That link has no power here") > -1 {
				return false
			}
		}
	}
	return true
}

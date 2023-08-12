package qb_watchfiles

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rskvp/qb-core/qb_events"
	"github.com/rskvp/qb-core/qb_ticker"
	"github.com/rskvp/qb-core/qb_utils"
)

var ErrorSystemPanic = errors.New("panic_system_error")

type MultipleFileWatcher struct {
	settings        *MultipleFileWatcherSettings
	events          *qb_events.Emitter
	fileMux         sync.Mutex
	stopTicker      *qb_ticker.Ticker
	validationRules *MultipleFileWatcherRules

	foundFiles *FoundFiles
}

// ---------------------------------------------------------------------------------------------------------------------
//		c o n s t r u c t o r
// ---------------------------------------------------------------------------------------------------------------------

func NewMultipleFileWatcher(settings *MultipleFileWatcherSettings, optEvents *qb_events.Emitter) (instance *MultipleFileWatcher) {
	instance = new(MultipleFileWatcher)

	if nil == settings {
		settings = &MultipleFileWatcherSettings{
			Dirs:             make([]string, 0),
			FilePatterns:     make([]string, 0),
			ExtraConstraints: make([]string, 0),
			CheckIntervalMs:  1000,
		}
	}
	instance.settings = settings
	instance.events = optEvents
	if nil == instance.events {
		instance.events = qb_events.Events.NewEmitter()
	}
	instance.validationRules = NewMultipleFileWatcherRules()
	instance.foundFiles = NewFoundFiles()

	instance.init()

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	s e t t i n g s
// ---------------------------------------------------------------------------------------------------------------------

func (instance *MultipleFileWatcher) GetEventNameOnFileWatch() string {
	if nil != instance {
		return instance.settings.TriggerEventName
	}
	return ""
}

func (instance *MultipleFileWatcher) SetEventNameOnFileWatch(value string) *MultipleFileWatcher {
	if nil != instance {
		instance.settings.TriggerEventName = value
	}
	return instance
}

func (instance *MultipleFileWatcher) GetFilePatterns() (response []string) {
	if nil != instance {
		response = instance.settings.FilePatterns
	}
	return
}

func (instance *MultipleFileWatcher) SetFilePatterns(value []string) *MultipleFileWatcher {
	if nil != instance {
		instance.settings.FilePatterns = value
	}
	return instance
}

func (instance *MultipleFileWatcher) GetDirs() (response []string) {
	if nil != instance {
		response = instance.settings.Dirs
	}
	return
}

func (instance *MultipleFileWatcher) SetDirs(value []string) *MultipleFileWatcher {
	if nil != instance {
		instance.settings.Dirs = value
	}
	return instance
}

func (instance *MultipleFileWatcher) GetIncludeSubFolders() (response bool) {
	if nil != instance {
		response = instance.settings.IncludeSubFolders
	}
	return
}

func (instance *MultipleFileWatcher) SetIncludeSubFolders(value bool) *MultipleFileWatcher {
	if nil != instance {
		instance.settings.IncludeSubFolders = value
	}
	return instance
}

func (instance *MultipleFileWatcher) GetConstraintsMode() (response bool) {
	if nil != instance {
		response = instance.settings.HasConstraint
	}
	return
}

func (instance *MultipleFileWatcher) SetConstraintsMode(value bool) *MultipleFileWatcher {
	if nil != instance {
		instance.settings.HasConstraint = value
	}
	return instance
}

func (instance *MultipleFileWatcher) GetExtraConstraints() (response []string) {
	if nil != instance {
		response = instance.settings.ExtraConstraints
	}
	return
}

func (instance *MultipleFileWatcher) SetExtraConstraints(value []string) *MultipleFileWatcher {
	if nil != instance {
		instance.settings.ExtraConstraints = value
	}
	return instance
}

func (instance *MultipleFileWatcher) GetCheckIntervalMs() (response int) {
	if nil != instance {
		response = instance.settings.CheckIntervalMs
	}
	return
}

func (instance *MultipleFileWatcher) SetCheckIntervalMs(value int) *MultipleFileWatcher {
	if nil != instance {
		instance.settings.CheckIntervalMs = value
	}
	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//		p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *MultipleFileWatcher) String() string {
	if nil != instance && nil != instance.foundFiles {
		return instance.foundFiles.String()
	}
	return ""
}

func (instance *MultipleFileWatcher) OnFileWatch(callback func(event *qb_events.Event)) *MultipleFileWatcher {
	if nil != instance && nil != instance.events {
		instance.events.On(instance.settings.TriggerEventName, callback)
	}
	return instance
}

func (instance *MultipleFileWatcher) Start() {
	if nil != instance && nil != instance.settings && len(instance.settings.FilePatterns) > 0 && nil == instance.stopTicker {
		interval := instance.settings.CheckIntervalMs
		instance.stopTicker = qb_ticker.NewTicker(time.Duration(interval)*time.Millisecond, func(t *qb_ticker.Ticker) {
			t.Pause()
			defer t.Resume()
			instance.checkFiles()
		})
		instance.stopTicker.Start()
	}
}

func (instance *MultipleFileWatcher) Stop() {
	if nil != instance && nil != instance.stopTicker {
		instance.stopTicker.Stop()
		instance.stopTicker = nil
	}
}

func (instance *MultipleFileWatcher) Join() {
	if nil != instance && nil != instance.stopTicker {
		instance.stopTicker.Join()
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//		p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *MultipleFileWatcher) init() {
	settings := instance.settings
	if nil != settings {
		if len(settings.Dirs) == 0 {
			settings.Dirs = []string{qb_utils.Paths.GetWorkspacePath()}
		}
		if len(settings.FilePatterns) == 0 {
			settings.FilePatterns = []string{"*.*"}
		}
		if len(settings.TriggerEventName) == 0 {
			settings.TriggerEventName = "on_multiple_file_match"
		}
		if settings.CheckIntervalMs == 0 {
			settings.CheckIntervalMs = 1000
		}
	}
}

func (instance *MultipleFileWatcher) checkFiles() {
	if nil != instance && nil != instance.settings && nil != instance.stopTicker && nil != instance.events && len(instance.settings.FilePatterns) > 0 {
		// lock
		instance.fileMux.Lock()
		defer instance.fileMux.Unlock()

		roots := instance.settings.Dirs
		patterns := instance.settings.FilePatterns

		// check if file exists
		for _, root := range roots {
			root = qb_utils.Paths.Absolute(root)
			files, err := qb_utils.Paths.ListFiles(root, patterns)
			if nil != err {
				// [ASYNC] error event
				param := newMultipleFiles()
				param.Error = err
				param.ErrorContext = fmt.Sprintf("Listing files from '%s' with patterns '%v'", root, patterns)
				instance.events.Emit(instance.settings.TriggerEventName, param)
				return
			}

			for _, filename := range files {
				if ok, _ := qb_utils.Paths.Exists(filename); ok {
					uid := getUid(root, filename)
					if !instance.exists(uid, filename) {
						base := qb_utils.Paths.Dir(filename)
						if !instance.canContinue(root, base) {
							continue
						}

						// validate file before trigger
						if !instance.validationRules.IsValid(filename) {
							continue
						}

						instance.tryTrigger(uid, filename)
					}
				}
			}
		}
	}
}

func (instance *MultipleFileWatcher) canContinue(root, dir string) bool {
	if nil != instance {
		sameDir := root == dir
		if !instance.settings.IncludeSubFolders {
			return sameDir
		}
	}
	return true
}

func (instance *MultipleFileWatcher) exists(uid, filename string) bool {
	if nil != instance && nil != instance.foundFiles {
		return instance.foundFiles.ExistsFile(uid, filename)
	}
	return false
}

func (instance *MultipleFileWatcher) tryTrigger(uid, filename string) {
	if nil != instance && nil != instance.settings {
		hasConstraint := instance.settings.HasConstraint
		if !hasConstraint {
			// SINGLE FILE --------------------------------------------
			param := newMultipleFiles()
			w, e := wrapAndRemove(uid, filename)
			if nil == e {
				param.Files = append(param.Files, w)
				instance.events.Emit(instance.settings.TriggerEventName, param)
			} else {
				param.Error = e
				param.ErrorContext = fmt.Sprintf("Wrapping single file '%s'", filename)
				instance.events.Emit(instance.settings.TriggerEventName, param)
			}
		} else {
			// MULTIPLE FILES WITH CONSTRAINTS ------------------------
			patterns := instance.settings.FilePatterns
			// check pattern
			for _, pattern := range patterns {
				if !instance.foundFiles.Exists(uid, pattern) {
					if instance.patternMatch(uid, filename, pattern) {
						// wrap and add file
						w, e := wrapAndRemove(uid, filename)
						if nil == e {
							_ = instance.foundFiles.Add(uid, pattern, w)
							if instance.foundFiles.Len(uid) == len(patterns) {
								param := newMultipleFiles()
								param.Files = instance.foundFiles.GetFiles(uid)
								instance.foundFiles.Reset(uid) // RESET
								instance.events.Emit(instance.settings.TriggerEventName, param)
								return
							} else {
								break
							}
						} else {
							// error
							instance.foundFiles.Reset(uid) // RESET
							param := newMultipleFiles()
							param.Error = e
							param.ErrorContext = fmt.Sprintf("Wrapping single file '%s'", filename)
							instance.events.Emit(instance.settings.TriggerEventName, param)
							return
						}
					}
				}
			}
		}
	}
}

func (instance *MultipleFileWatcher) patternMatch(uid, filename, pattern string) bool {
	if nil != instance {
		if len(instance.settings.ExtraConstraints) > 0 {
			if qb_utils.Paths.PatternMatchBase(filename, pattern) {
				for _, constraint := range instance.settings.ExtraConstraints {
					switch constraint {
					case NameConstraint:
						// all files must have same name
						if instance.foundFiles.Len(uid) > 0 {
							name := qb_utils.Paths.FileName(filename, false)
							files := instance.foundFiles.GetFiles(uid)
							for _, file := range files {
								if name != file.Name {
									return false
								}
							}
						}
						return true
					default:
						return false
					}
				}
			}
		} else {
			return qb_utils.Paths.PatternMatchBase(filename, pattern)
		}
	}
	return false
}

// ---------------------------------------------------------------------------------------------------------------------
//	S T A T I C
// ---------------------------------------------------------------------------------------------------------------------

func newMultipleFiles() *MultipleFileEventParam {
	instance := new(MultipleFileEventParam)
	instance.Files = make([]*FileWrapper, 0)
	return instance
}

func getUid(root, filename string) string {
	dir := qb_utils.Paths.FileName(root, false)
	sub := strings.ReplaceAll(qb_utils.Paths.Dir(filename), root, "")
	return qb_utils.Paths.Concat(dir, sub)
}

func wrap(uid, filename string) (response *FileWrapper, err error) {
	var stats os.FileInfo
	stats, err = os.Stat(filename)
	if nil == err {
		var bytes []byte
		bytes, err = qb_utils.IO.ReadBytesFromFile(filename)
		if nil == err {
			response = &FileWrapper{
				Uid:      uid,
				Filename: filename,
				Name:     qb_utils.Paths.FileName(filename, false),
				Ext:      qb_utils.Paths.Extension(filename),
				Stats:    stats,
				Bytes:    bytes,
			}
		}
	}
	return
}

func wrapAndRemove(uid, filename string) (response *FileWrapper, err error) {
	response, err = wrap(uid, filename)
	if nil == err {
		err = qb_utils.IO.Remove(filename)
	}
	return
}

package qb_ticker

import (
	"sync"
	"time"

	"github.com/rskvp/qb-core/qb_events"
	"github.com/rskvp/qb-core/qb_utils"
)

// ---------------------------------------------------------------------------------------------------------------------
//		t y p e
// ---------------------------------------------------------------------------------------------------------------------

const FileMonitorEventName = "on_file_match"
const FileMonitorErrorEventName = "on_file_match_error"

type FileMonitorSettings struct {
	Dirs                  []string `json:"dirs"` // directories to monitor for file
	IncludeSubFolders     bool     `json:"include_sub_folders"`
	FilePattern           string   `json:"pattern"` // file name or file pattern i.e. "*.xml"
	TriggerEventName      string   `json:"trigger-event-name"`
	TriggerErrorEventName string   `json:"trigger-error-event-name"`
	CheckIntervalMs       int      `json:"check-interval-ms"`
	DeleteAfterMatch      bool     `json:"delete-after-match"`
	MoveAfterMatchDir     string   `json:"move-after-match-dir"`
}

type FileMonitor struct {
	roots             []string // where stop file is stored
	includeSubFolders bool
	filePattern       string
	eventName         string
	errorEventName    string
	deleteAfterMatch  bool
	moveAfterMatchDir string
	checkIntervalMs   int
	events            *qb_events.Emitter
	fileMux           sync.Mutex
	stopTicker        *Ticker
}

// ---------------------------------------------------------------------------------------------------------------------
//		c o n s t r u c t o r
// ---------------------------------------------------------------------------------------------------------------------

func NewFileMonitor(settings *FileMonitorSettings, events *qb_events.Emitter) *FileMonitor {
	instance := new(FileMonitor)

	instance.roots = []string{qb_utils.Paths.GetWorkspacePath()}
	instance.filePattern = "*.*"
	instance.eventName = FileMonitorEventName
	instance.errorEventName = FileMonitorErrorEventName
	instance.deleteAfterMatch = false
	instance.moveAfterMatchDir = "./matched"
	instance.checkIntervalMs = 1000

	if nil != settings {
		instance.includeSubFolders = settings.IncludeSubFolders
		instance.deleteAfterMatch = settings.DeleteAfterMatch
		if len(settings.Dirs) > 0 {
			instance.roots = settings.Dirs
		}
		if len(settings.FilePattern) > 0 {
			instance.filePattern = settings.FilePattern
		}
		if len(settings.TriggerEventName) > 0 {
			instance.eventName = settings.TriggerEventName
		}
		if len(settings.TriggerErrorEventName) > 0 {
			instance.errorEventName = settings.TriggerErrorEventName
		}

		instance.moveAfterMatchDir = settings.MoveAfterMatchDir

		if settings.CheckIntervalMs > 100 {
			instance.checkIntervalMs = settings.CheckIntervalMs
		}
	}

	instance.events = events
	if nil == instance.events {
		instance.events = qb_events.Events.NewEmitter()
	}
	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//		p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *FileMonitor) Events() *qb_events.Emitter {
	if nil != instance {
		return instance.events
	}
	return nil
}

func (instance *FileMonitor) Start() {
	if nil != instance && len(instance.filePattern) > 0 && nil == instance.stopTicker {
		instance.stopTicker = NewTicker(time.Duration(instance.checkIntervalMs)*time.Millisecond, func(t *Ticker) {
			instance.checkFile()
		})
		instance.stopTicker.Start()
	}
}

func (instance *FileMonitor) Stop() {
	if nil != instance && nil != instance.stopTicker {
		instance.stopTicker.Stop()
		instance.stopTicker = nil
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//		p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *FileMonitor) checkFile() {
	if nil != instance && nil != instance.stopTicker && len(instance.filePattern) > 0 {
		instance.stopTicker.Pause()
		defer instance.stopTicker.Resume()
		instance.fileMux.Lock()
		defer instance.fileMux.Unlock()

		// check if file exists
		for _, root := range instance.roots {
			root = qb_utils.Paths.Absolute(root)
			files, _ := qb_utils.Paths.ListFiles(root, instance.filePattern)
			for _, filename := range files {
				base := qb_utils.Paths.Dir(filename)
				if !instance.canContinue(root, base) {
					continue
				}

				// [ASYNC] file found event
				instance.events.Emit(instance.eventName, filename)

				// proceed moving or deleting the file
				moved, err := instance.move(filename)
				if nil == err {
					if !moved {
						err = instance.delete(filename)
					}
				}
				if nil != err {
					// [ASYNC] error event
					instance.events.Emit(instance.errorEventName, err, filename)
				}
			}
		}
	}
}

func (instance *FileMonitor) canContinue(root, dir string) bool {
	if !instance.includeSubFolders {
		return root == dir
	}
	return true
}

func (instance *FileMonitor) delete(filename string) (err error) {
	if nil != instance && instance.deleteAfterMatch {
		if ok, _ := qb_utils.Paths.Exists(filename); ok {
			err = qb_utils.IO.Remove(filename)
		}
	}
	return
}

func (instance *FileMonitor) move(filename string) (moved bool, err error) {
	if nil != instance && len(instance.moveAfterMatchDir) > 0 {
		path := qb_utils.Paths.Absolutize(instance.moveAfterMatchDir, qb_utils.Paths.GetWorkspacePath())
		name := qb_utils.Paths.FileName(filename, true)
		copyTo := qb_utils.Paths.Concat(path, name)
		_, err = qb_utils.IO.CopyFile(filename, copyTo)
		if nil == err {
			err = qb_utils.IO.Remove(filename)
		}
		moved = nil == err
	}
	return
}

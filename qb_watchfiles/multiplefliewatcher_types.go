package qb_watchfiles

import (
	"os"

	gg "github.com/rskvp/qb-core"
)

const (
	NameConstraint = "name"
)

type MultipleFileWatcherSettings struct {
	Dirs              []string `json:"dirs"` // directories to monitor for file
	IncludeSubFolders bool     `json:"include-sub-folders"`
	FilePatterns      []string `json:"patterns"`         // define the number of files and pattern required i.e. "*.xml"
	HasConstraint     bool     `json:"has-constraint"`   // "patterns" are used as constraint
	ExtraConstraints  []string `json:"extra-constraint"` // ["name"] a special constraint to be evaluated before considering a right constraint match
	TriggerEventName  string   `json:"trigger-event-name"`
	CheckIntervalMs   int      `json:"check-interval-ms"`
}

type MultipleFileEventParam struct {
	Files        []*FileWrapper `json:"files"`
	Error        error          `json:"error"`
	ErrorContext string         `json:"error-context"`
}

func (instance *MultipleFileEventParam) String() string {
	return gg.JSON.Stringify(instance)
}

type FileWrapper struct {
	Uid      string      `json:"sub-folder"`
	Filename string      `json:"filename"`
	Name     string      `json:"name"`
	Ext      string      `json:"ext"`
	Stats    os.FileInfo `json:"-"`
	Bytes    []byte      `json:"-"`
}

func (instance *FileWrapper) String() string {
	return gg.JSON.Stringify(instance)
}

type FoundFiles struct {
	files map[string]map[string]*FileWrapper // map of dirs - map of patterns
}

func NewFoundFiles() *FoundFiles {
	instance := new(FoundFiles)
	instance.ResetAll()
	return instance
}

func (instance *FoundFiles) String() string {
	if nil != instance && nil != instance.files {
		return gg.JSON.Stringify(instance.files)
	}
	return ""
}

func (instance *FoundFiles) ResetAll() {
	instance.files = make(map[string]map[string]*FileWrapper)
}

func (instance *FoundFiles) Reset(uid string) {
	if nil != instance && nil != instance.files {
		if nil != instance && nil != instance.files {
			if _, ok := instance.files[uid]; ok {
				instance.files[uid] = make(map[string]*FileWrapper)
			}
		}
	}
}

func (instance *FoundFiles) Exists(uid, pattern string) bool {
	if nil != instance && nil != instance.files {
		// ensure uid exists
		mapPatterns := instance.getMapPattern(uid)
		if _, ok := mapPatterns[pattern]; ok {
			return true
		}
	}
	return false
}

func (instance *FoundFiles) Add(uid, pattern string, w *FileWrapper) bool {
	if nil != instance && nil != instance.files {
		mapPatterns := instance.getMapPattern(uid)
		mapPatterns[pattern] = w
		return true
	}
	return false
}

func (instance *FoundFiles) Len(uid string) int {
	if nil != instance && nil != instance.files {
		mapPatterns := instance.getMapPattern(uid)
		return len(mapPatterns)
	}
	return 0
}

func (instance *FoundFiles) ExistsFile(uid, filename string) bool {
	if nil != instance && nil != instance.files {
		mapPatterns := instance.getMapPattern(uid)
		for _, f := range mapPatterns {
			if f.Filename == filename {
				return true
			}
		}
	}
	return false
}

func (instance *FoundFiles) GetFiles(uid string) (response []*FileWrapper) {
	response = make([]*FileWrapper, 0)
	if nil != instance && nil != instance.files {
		mapPatterns := instance.getMapPattern(uid)
		for _, w := range mapPatterns {
			response = append(response, w)
		}
	}
	return
}

func (instance *FoundFiles) getMapPattern(uid string) map[string]*FileWrapper {
	if nil != instance && nil != instance.files {
		if _, ok := instance.files[uid]; !ok {
			instance.files[uid] = make(map[string]*FileWrapper)
		}
		return instance.files[uid]
	}
	return make(map[string]*FileWrapper)
}

package qb_rnd

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ValuesRandomizerParams struct {
	Index  int           `json:"index,omitempty"`
	Values []interface{} `json:"values,omitempty"`
	Mode   string        `json:"mode,omitempty"` // one of "random" or "sequential"
}

func (instance *ValuesRandomizerParams) String() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

func (instance *ValuesRandomizerParams) Map() map[string]interface{} {
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(instance.String()), &m)
	return m
}

func (instance *ValuesRandomizerParams) ResetValues() *ValuesRandomizerParams {
	if nil != instance {
		instance.Values = make([]interface{}, 0)
	}
	return instance
}

func (instance *ValuesRandomizerParams) SetValues(values []interface{}) *ValuesRandomizerParams {
	if nil != instance {
		instance.Values = make([]interface{}, 0)
		instance.Values = append(instance.Values, values...)
	}
	return instance
}

func (instance *ValuesRandomizerParams) AddValue(value interface{}) *ValuesRandomizerParams {
	if nil != instance {
		instance.Values = append(instance.Values, value)
	}
	return instance
}

type ValuesRandomizer struct {
	Params   *ValuesRandomizerParams
	Autosave bool
	filename string // settings filename
	done     []int  //indexes already done
}

func NewValuesRandomizer(args ...interface{}) (instance *ValuesRandomizer, err error) {
	instance = new(ValuesRandomizer)
	instance.Params = new(ValuesRandomizerParams)
	instance.done = make([]int, 0)
	err = instance.init(args...)
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ValuesRandomizer) Map() map[string]interface{} {
	m := map[string]interface{}{
		"filename": instance.filename,
		"params":   instance.Params.Map(),
	}
	return m
}

func (instance *ValuesRandomizer) String() string {
	b, err := json.Marshal(instance.Map())
	if nil == err {
		return string(b)
	}
	return ""
}

func (instance *ValuesRandomizer) Save() (err error) {
	if nil != instance {
		_, err = instance.saveToFile()
	}
	return
}

func (instance *ValuesRandomizer) SaveTo(filename string) (err error) {
	if nil != instance {
		instance.filename = absolute(filename)
		_, err = instance.saveToFile()
	}
	return
}

func (instance *ValuesRandomizer) Next() interface{} {
	if nil != instance {
		switch instance.Params.Mode {
		case "random", "rnd":
			// random
			undone := instance.getUndoneIndexes()
			if len(undone) == len(instance.Params.Values) {
				instance.done = make([]int, 0)
			}
			i := int(Rnd.Between(0, int64(len(undone)-1)))
			instance.Params.Index = undone[i]
		default:
			// sequential
			instance.Params.Index += 1
			if instance.Params.Index > len(instance.Params.Values) {
				instance.Params.Index = 1
				instance.done = make([]int, 0)
			}
		}

		instance.done = append(instance.done, instance.Params.Index)

		if instance.Autosave {
			_ = instance.Save()
		}

		return instance.Params.Values[instance.Params.Index-1]
	}
	return nil
}

func (instance *ValuesRandomizer) GetDoneIndexes() []int {
	if nil != instance && nil != instance.done {
		return instance.done
	}
	return []int{}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ValuesRandomizer) init(args ...interface{}) (err error) {
	if nil != instance && len(args) > 0 {
		if len(args) == 1 {
			arg1 := args[0]
			if s, ok := arg1.(string); ok {
				// check if is a file
				if fileExists(s) {
					// load status from file
					err = instance.loadFromFile(s)
					return
				}
				// check if is JSON
				if isJSON(s) {
					// load status from json
					err = instance.loadFromJSON(s)
					return
				}
			} else if p, ok := arg1.(ValuesRandomizerParams); ok {
				instance.Params = &p
			} else if p, ok := arg1.(*ValuesRandomizerParams); ok {
				instance.Params = p
			} else if m, ok := arg1.(map[string]interface{}); ok {
				err = instance.loadFromMap(m)
				return
			}
		}
	}
	return
}

func (instance *ValuesRandomizer) loadFromFile(filename string) (err error) {
	// load status from file
	instance.filename = absolute(filename)
	var data []byte
	data, err = readBytesFromFile(filename)
	if nil != err {
		return
	}
	if len(data) > 0 {
		err = instance.loadFromJSON(string(data))
	}
	return
}

func (instance *ValuesRandomizer) loadFromMap(m map[string]interface{}) (err error) {
	b, err := json.Marshal(&m)
	if nil == err {
		err = instance.loadFromJSON(string(b))
	}
	return
}

func (instance *ValuesRandomizer) loadFromJSON(text string) (err error) {
	err = json.Unmarshal([]byte(text), &instance.Params)
	return
}

func (instance *ValuesRandomizer) getFilename() string {
	if len(instance.filename) == 0 {
		instance.filename = absolute("./" + Rnd.Uuid() + ".json")
	}
	return instance.filename
}

func (instance *ValuesRandomizer) saveToFile() (filename string, err error) {
	filename = instance.getFilename()
	text := instance.Params.String()
	if len(text) > 0 && len(filename) > 0 {
		var f *os.File
		f, err = os.Create(filename)

		if nil == err {
			defer f.Close()
			w := bufio.NewWriter(f)
			_, err = w.WriteString(text)
			_ = w.Flush()
		}
	}
	return
}

func (instance *ValuesRandomizer) getUndoneIndexes() []int {
	undone := make([]int, 0)
	all := make([]int, 0)
	for i, _ := range instance.Params.Values {
		all = append(all, i+1)
		if indexOf(instance.done, i+1) == -1 {
			undone = append(undone, i+1)
		}
	}
	if len(undone) > 0 {
		return undone
	}
	return all
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func isJSON(text string) bool {
	text = strings.TrimSpace(text)
	return strings.Index(text, "{") == 0 || strings.Index(text, "[") == 0
}

func fileExists(filename string) bool {
	if len(filename) > 0 && !isJSON(filename) {
		info, err := os.Stat(filename)
		if err == nil {
			return info.Size() > 0
		}
		if os.IsNotExist(err) {
			return false
		}
		return true
	}
	return false
}

func readBytesFromFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	return b, err
}

func absolute(path string) string {
	abs, err := filepath.Abs(path)
	if nil == err {
		return abs
	}
	return path
}

func indexOf(list []int, value int) int {
	for i, v := range list {
		if v == value {
			return i
		}
	}
	return -1
}

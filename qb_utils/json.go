package qb_utils

import (
	"encoding/json"
	"strings"
)

type JSONHelper struct {
}

var JSON *JSONHelper

//----------------------------------------------------------------------------------------------------------------------
//	i n i t
//----------------------------------------------------------------------------------------------------------------------

func init() {
	JSON = new(JSONHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *JSONHelper) StringToMap(text string) (map[string]interface{}, bool) {
	return instance.BytesToMap([]byte(text))
}

func (instance *JSONHelper) StringToArray(text string) ([]interface{}, bool) {
	return instance.BytesToArray([]byte(text))
}

func (instance *JSONHelper) BytesToMap(data []byte) (map[string]interface{}, bool) {
	var js map[string]interface{}
	err := json.Unmarshal(data, &js)
	if nil != err {
		return nil, false
	}
	return js, true
}

func (instance *JSONHelper) BytesToArray(data []byte) ([]interface{}, bool) {
	var js []interface{}
	err := json.Unmarshal(data, &js)
	if nil != err {
		return nil, false
	}
	return js, true
}

func (instance *JSONHelper) Bytes(entity interface{}) []byte {
	if nil != entity {
		if s, b := entity.(string); b {
			return []byte(s)
		}
		b, err := json.Marshal(&entity)
		if nil == err {
			return b
		}
	}

	return []byte{}
}

func (instance *JSONHelper) Stringify(entity interface{}) string {
	if s, b := entity.(string); b {
		if strings.Index(s, "\"") != 0 {
			// quote to not quoted string
			// return string(instance.Bytes(Strings.Quote(s)))
			return Strings.Quote(s)
		}
	}
	return string(instance.Bytes(entity))
}

func (instance *JSONHelper) Parse(input interface{}) interface{} {
	if v, b := input.(string); b {
		if o, b := instance.StringToMap(v); b {
			return o // map
		} else if a, b := instance.StringToArray(v); b {
			return a // array
		}
		return v // simple string
	} else if v, b := input.([]byte); b {
		if o, b := instance.BytesToMap(v); b {
			return o // map
		} else if a, b := instance.BytesToArray(v); b {
			return a // array
		}
		return v // simple string
	}
	return input
}

func (instance *JSONHelper) Read(input interface{}, entity interface{}) (err error) {
	if s, b := input.(string); b {
		err = json.Unmarshal([]byte(s), &entity)
	} else if s, b := input.([]byte); b {
		err = json.Unmarshal(s, &entity)
	}
	if nil != err {
		return err
	}
	return nil
}

func (instance *JSONHelper) ReadFromFile(fileName string, entity interface{}) error {
	b, err := IO.ReadBytesFromFile(fileName)
	if nil != err {
		return err
	}
	err = json.Unmarshal(b, &entity)
	if nil != err {
		return err
	}
	return nil
}

func (instance *JSONHelper) ReadMapFromFile(fileName string) (map[string]interface{}, error) {
	b, err := IO.ReadBytesFromFile(fileName)
	if nil != err {
		return nil, err
	}
	var response map[string]interface{}
	err = json.Unmarshal(b, &response)
	if nil != err {
		return nil, err
	}
	return response, nil
}

func (instance *JSONHelper) ReadArrayFromFile(fileName string) ([]map[string]interface{}, error) {
	b, err := IO.ReadBytesFromFile(fileName)
	if nil != err {
		return nil, err
	}
	var response []map[string]interface{}
	err = json.Unmarshal(b, &response)
	if nil != err {
		return nil, err
	}
	return response, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

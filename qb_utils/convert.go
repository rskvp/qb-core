package qb_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_num2word"
)

type ConversionHelper struct {
}

var Convert *ConversionHelper

//----------------------------------------------------------------------------------------------------------------------
//	i n i t
//----------------------------------------------------------------------------------------------------------------------

func init() {
	Convert = new(ConversionHelper)
}

var (
	Kb = uint64(1024)
	Mb = Kb * 1024
	Gb = Mb * 1024
	Tb = Gb * 1024
	Pb = Tb * 1024
	Eb = Pb * 1024
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ConversionHelper) Num2Word(args ...interface{}) string {
	var value int
	var langCode string
	var opts *qb_num2word.Num2WordOptions

	switch len(args) {
	case 1:
		value = instance.ToInt(args[0])
	case 2:
		value = instance.ToInt(args[0])
		langCode = instance.ToString(args[1])
	case 3:
		value = instance.ToInt(args[0])
		langCode = instance.ToString(args[1])
		if v, b := args[2].(*qb_num2word.Num2WordOptions); b {
			opts = v
		} else if v, b := args[2].(qb_num2word.Num2WordOptions); b {
			opts = &v
		}
	default:
		return ""
	}
	return qb_num2word.Num2Word2.ConvertOpts(value, langCode, opts)
}

func (instance *ConversionHelper) ToKiloBytes(val int64) float64 {
	return instance.ToFloat64(val) / instance.ToFloat64(Kb)
}

func (instance *ConversionHelper) ToMegaBytes(val int64) float64 {
	return instance.ToFloat64(val) / instance.ToFloat64(Mb)
}

func (instance *ConversionHelper) ToGigaBytes(val int64) float64 {
	return instance.ToFloat64(val) / instance.ToFloat64(Gb)
}

func (instance *ConversionHelper) ToTeraBytes(val int64) float64 {
	return instance.ToFloat64(val) / instance.ToFloat64(Tb)
}

func (instance *ConversionHelper) ToPetaBytes(val int64) float64 {
	return instance.ToFloat64(val) / instance.ToFloat64(Pb)
}

func (instance *ConversionHelper) ToEsaBytes(val int64) float64 {
	return instance.ToFloat64(val) / instance.ToFloat64(Eb)
}

func (instance *ConversionHelper) ToDuration(val interface{}) time.Duration {
	value := int64(1)
	if nil != val {
		if s, b := val.(string); b {
			// TIMELINE
			tl := strings.Split(s, ":") // hour:12
			um := tl[0]
			if len(tl) == 2 {
				value = instance.ToInt64(tl[1])
				if value == 0 {
					value = 1
				}
			}
			switch um {
			case "millisecond":
				return time.Duration(value) * time.Millisecond
			case "second":
				return time.Duration(value) * time.Second
			case "minute":
				return time.Duration(value) * time.Minute
			case "hour":
				return time.Duration(value) * time.Hour
			case "day":
				return time.Duration(value) * time.Hour * 24
			default:
				return 12 * time.Hour
			}
		} else if i, b := val.(int); b {
			value = int64(i)
		} else if i, b := val.(int8); b {
			value = int64(i)
		} else if i, b := val.(int16); b {
			value = int64(i)
		} else if i, b := val.(int32); b {
			value = int64(i)
		} else if i, b := val.(int64); b {
			value = i
		} else if i, b := val.(uint); b {
			value = int64(i)
		} else if i, b := val.(uint8); b {
			value = int64(i)
		} else if i, b := val.(uint16); b {
			value = int64(i)
		} else if i, b := val.(uint32); b {
			value = int64(i)
		} else if i, b := val.(uint64); b {
			value = int64(i)
		}
	}
	return time.Duration(value) * time.Millisecond
}

func (instance *ConversionHelper) ToArray(val ...interface{}) []interface{} {
	if nil == val {
		return nil
	}
	aa := instance.toArray(val...)
	return aa
}

func (instance *ConversionHelper) ToArrayOfString(val ...interface{}) []string {
	if nil == val {
		return nil
	}
	aa := instance.toArrayOfString(val...)
	return aa
}

func (instance *ConversionHelper) ToArrayOfUrlEncodedString(val ...interface{}) (response []string) {
	response = make([]string, 0)
	a := instance.ToArrayOfString(val...)
	for _, s := range a {
		response = append(response, url.QueryEscape(s))
	}
	return
}

func (instance *ConversionHelper) ToArrayOfByte(val interface{}) []byte {
	if nil == val {
		return nil
	}
	return instance.toArrayOfByte(val)
}

func (instance *ConversionHelper) ToArrayOfInt(val interface{}) []int {
	if b, _ := Compare.IsArray(val); b {
		if v, b := val.([]interface{}); b {
			return instance.toArrayOfInt(v)
		} else if v, b := val.([]int); b {
			return v
		}
	}
	return nil
}

func (instance *ConversionHelper) ToArrayOfInt8(val interface{}) []int8 {
	if b, _ := Compare.IsArray(val); b {
		if v, b := val.([]interface{}); b {
			return instance.toArrayOfInt8(v)
		} else if v, b := val.([]int8); b {
			return v
		}
	}
	return nil
}

func (instance *ConversionHelper) ToArrayOfInt16(val interface{}) []int16 {
	if b, _ := Compare.IsArray(val); b {
		if v, b := val.([]interface{}); b {
			return instance.toArrayOfInt16(v)
		} else if v, b := val.([]int16); b {
			return v
		}
	}
	return nil
}

func (instance *ConversionHelper) ToArrayOfInt32(val interface{}) []int32 {
	if b, _ := Compare.IsArray(val); b {
		if v, b := val.([]interface{}); b {
			return instance.toArrayOfInt32(v)
		} else if v, b := val.([]int32); b {
			return v
		}
	}
	return nil
}

func (instance *ConversionHelper) ToArrayOfInt64(val interface{}) []int64 {
	if b, _ := Compare.IsArray(val); b {
		if v, b := val.([]interface{}); b {
			return instance.toArrayOfInt64(v)
		} else if v, b := val.([]int64); b {
			return v
		}
	}
	return nil
}

func (instance *ConversionHelper) ToArrayOfFloat32(val interface{}) []float32 {
	if b, _ := Compare.IsArray(val); b {
		if v, b := val.([]interface{}); b {
			return instance.toArrayOfFloat32(v)
		} else if v, b := val.([]float32); b {
			return v
		}
	}
	return nil
}

func (instance *ConversionHelper) ToArrayOfFloat64(val interface{}) []float64 {
	if b, _ := Compare.IsArray(val); b {
		if v, b := val.([]interface{}); b {
			return instance.toArrayOfFloat64(v)
		} else if v, b := val.([]float64); b {
			return v
		}
	}
	return nil
}

func (instance *ConversionHelper) Int8ToStr(arr []int8) string {
	b := make([]byte, 0, len(arr))
	for _, v := range arr {
		if v == 0x00 {
			break
		}
		b = append(b, byte(v))
	}
	return string(b)
}

func (instance *ConversionHelper) ToInt(val interface{}) int {
	return instance.ToIntDef(val, -1)
}

func (instance *ConversionHelper) ToIntDef(val interface{}, def int) int {
	if b, s := Compare.IsString(val); b {
		v, err := strconv.Atoi(strings.TrimSpace(s))
		if nil == err {
			return v
		}
	}
	switch i := val.(type) {
	case float32:
		return int(i)
	case float64:
		return int(i)
	case int:
		return i
	case int8:
		return int(i)
	case int16:
		return int(i)
	case int32:
		return int(i)
	case int64:
		return int(i)
	case uint:
		return int(i)
	case uint8:
		return int(i)
	case uint16:
		return int(i)
	case uint32:
		return int(i)
	case uint64:
		return int(i)
	}

	return def
}

func (instance *ConversionHelper) ToInt64(val interface{}) int64 {
	return instance.ToInt64Def(val, -1.0)
}

func (instance *ConversionHelper) ToInt64Def(val interface{}, defVal int64) int64 {
	switch i := val.(type) {
	case float32:
		return int64(i)
	case float64:
		return int64(i)
	case int:
		return int64(i)
	case int8:
		return int64(i)
	case int16:
		return int64(i)
	case int32:
		return int64(i)
	case int64:
		return i
	case uint8:
		return int64(i)
	case uint16:
		return int64(i)
	case uint32:
		return int64(i)
	case uint64:
		return int64(i)
	case string:
		v, err := strconv.ParseInt(i, 10, 64)
		if nil == err {
			return v
		}
	}
	return defVal
}

func (instance *ConversionHelper) ToInt32(val interface{}) int32 {
	return instance.ToInt32Def(val, -1)
}

func (instance *ConversionHelper) ToInt32Def(val interface{}, defVal int32) int32 {
	switch i := val.(type) {
	case float32:
		return int32(i)
	case float64:
		return int32(i)
	case int:
		return int32(i)
	case int8:
		return int32(i)
	case int16:
		return int32(i)
	case int32:
		return i
	case int64:
		return int32(i)
	case uint8:
		return int32(i)
	case uint16:
		return int32(i)
	case uint32:
		return int32(i)
	case uint64:
		return int32(i)
	case string:
		v, err := strconv.ParseInt(i, 10, 32)
		if nil == err {
			return int32(v)
		}
	}
	return defVal
}

func (instance *ConversionHelper) ToInt8(val interface{}) int8 {
	return instance.ToInt8Def(val, -1)
}

func (instance *ConversionHelper) ToInt8Def(val interface{}, defVal int8) int8 {
	switch i := val.(type) {
	case float32:
		return int8(i)
	case float64:
		return int8(i)
	case int:
		return int8(i)
	case int8:
		return i
	case int16:
		return int8(i)
	case int32:
		return int8(i)
	case int64:
		return int8(i)
	case uint8:
		return int8(i)
	case uint16:
		return int8(i)
	case uint32:
		return int8(i)
	case uint64:
		return int8(i)
	case string:
		v, err := strconv.ParseInt(i, 10, 8)
		if nil == err {
			return int8(v)
		}
	}
	return defVal
}

func (instance *ConversionHelper) ToInt16(val interface{}) int16 {
	return instance.ToInt16Def(val, -1)
}

func (instance *ConversionHelper) ToInt16Def(val interface{}, defVal int16) int16 {
	switch i := val.(type) {
	case float32:
		return int16(i)
	case float64:
		return int16(i)
	case int:
		return int16(i)
	case int8:
		return int16(i)
	case int16:
		return i
	case int32:
		return int16(i)
	case int64:
		return int16(i)
	case uint8:
		return int16(i)
	case uint16:
		return int16(i)
	case uint32:
		return int16(i)
	case uint64:
		return int16(i)
	case string:
		v, err := strconv.ParseInt(i, 10, 16)
		if nil == err {
			return int16(v)
		}
	}
	return defVal
}

func (instance *ConversionHelper) ToFloat32(val interface{}) float32 {
	return instance.ToFloat32Def(val, -1)
}

func (instance *ConversionHelper) ToFloat32Def(val interface{}, defVal float32) float32 {
	if b, s := Compare.IsString(val); b {
		v, err := strconv.ParseFloat(s, 32)
		if nil == err {
			return float32(v)
		}
	}
	switch i := val.(type) {
	case float32:
		return i
	case float64:
		return float32(i)
	case int:
		return float32(i)
	case int8:
		return float32(i)
	case int16:
		return float32(i)
	case int32:
		return float32(i)
	case int64:
		return float32(i)
	case uint8:
		return float32(i)
	case uint16:
		return float32(i)
	case uint32:
		return float32(i)
	case uint64:
		return float32(i)
	}
	return defVal
}

func (instance *ConversionHelper) ToFloat64(val interface{}) float64 {
	return instance.ToFloat64Def(val, -1.0)
}

func (instance *ConversionHelper) ToFloat64Def(val interface{}, defVal float64) float64 {
	switch i := val.(type) {
	case float32:
		return float64(i)
	case float64:
		return i
	case int:
		return float64(i)
	case int8:
		return float64(i)
	case int16:
		return float64(i)
	case int32:
		return float64(i)
	case int64:
		return float64(i)
	case uint8:
		return float64(i)
	case uint16:
		return float64(i)
	case uint32:
		return float64(i)
	case uint64:
		return float64(i)
	case string:
		v, err := strconv.ParseFloat(i, 64)
		if nil == err {
			return v
		}
	}
	return defVal
}

func (instance *ConversionHelper) ToBool(val interface{}) bool {
	if v, b := val.(bool); b {
		return v
	}
	if b, s := Compare.IsString(val); b {
		v, err := strconv.ParseBool(s)
		if nil == err {
			return v
		}
	}
	if a, b := val.([]byte); b {
		if len(a) > 1 {
			v, err := strconv.ParseBool(string(a))
			if nil == err {
				return v
			}
		} else {
			if a[0] == 0 {
				return false
			} else {
				return true
			}
		}
	}
	return false
}

func (instance *ConversionHelper) ToMap(val interface{}) map[string]interface{} {
	if b, s := Compare.IsString(val); b {
		m, err := instance.stringToMap(s)
		if nil == err {
			return m
		}
	}
	if b, _ := Compare.IsMap(val); b {
		return instance.toMap(val)
	}

	if m, err := instance.stringToMap(instance.ToString(val)); nil == err {
		return m
	}

	return nil
}

func (instance *ConversionHelper) ToMapOfString(val interface{}) map[string]string {
	if b, s := Compare.IsString(val); b {
		m, err := instance.stringToMap(s)
		if nil == err {
			return instance.toMapOfString(m)
		}
	}
	if b, _ := Compare.IsMap(val); b {
		return instance.toMapOfString(val)
	}
	return nil
}

func (instance *ConversionHelper) ToMapOfStringArray(val interface{}) map[string][]string {
	if m, b := val.(map[string][]string); b {
		return m
	}

	// warning: this change the pointer to original object
	data, err := json.Marshal(val)
	if nil == err {
		var m map[string][]string
		err = json.Unmarshal(data, &m)
		if nil == err {
			return m
		}
	}
	return nil
}

func (instance *ConversionHelper) ForceMap(val interface{}) map[string]interface{} {
	m := instance.ToMap(val)
	if nil == m {
		m = instance.toMap(val)
	}
	if nil == m {
		text := strings.ReplaceAll(instance.ToString(val), "'", "\"")
		m, _ = instance.stringToMap(text)
	}
	if nil == m {
		m = make(map[string]interface{})
	}
	return m
}

func (instance *ConversionHelper) ForceMapOfString(val interface{}) map[string]string {
	m := instance.ToMapOfString(val)
	if nil == m {
		return instance.toMapOfString(val)
	}
	return m
}

func (instance *ConversionHelper) ToString(val interface{}) string {
	if nil == val {
		return ""
	}
	// string
	s, ss := val.(string)
	if ss {
		return s
	}
	// integer
	i, ii := val.(int)
	if ii {
		return strconv.Itoa(i)
	}
	// float32
	f, ff := val.(float32)
	if ff {
		return fmt.Sprintf("%g", f) // Exponent as needed, necessary digits only
	}
	// float 64
	F, FF := val.(float64)
	if FF {
		return fmt.Sprintf("%g", F) // Exponent as needed, necessary digits only
		// return strconv.FormatFloat(F, 'E', -1, 64)
	}

	// boolean
	b, bb := val.(bool)
	if bb {
		return strconv.FormatBool(b)
	}

	// array
	if aa, _ := Compare.IsArray(val); aa {
		// byte array??
		if ba, b := val.([]byte); b {
			return string(ba)
		} else {
			data, err := json.Marshal(val)
			if nil == err {
				return string(data)
			}
			/*
				response := []string{}
				// array := make([]interface{}, tt.Len())
				for i := 0; i < tt.Len(); i++ {
					v := tt.Index(i).Interface()
					s := ToString(v)
					response = append(response, s)
				}
				return "[" + strings.Wait(response, ",") + "]"
			*/
		}
	}

	// map
	if b, _ := Compare.IsMap(val); b {
		data, err := json.Marshal(val)
		if nil == err {
			return string(data)
		}
	}

	// struct
	if b, _ := Compare.IsStruct(val); b {
		data, err := json.Marshal(val)
		if nil == err {
			return string(data)
		}
	}

	// undefined value
	return fmt.Sprintf("%v", val)
}

func (instance *ConversionHelper) ToStringQuoted(val interface{}) string {
	response := instance.ToString(val)
	if _, b := val.(string); b {
		if strings.Index(response, "\"") == -1 {
			response = strconv.Quote(response)
		}
	}
	return response
}

func (instance *ConversionHelper) ToUrlQuery(params map[string]interface{}) string {
	if len(params) > 0 {
		buff := bytes.Buffer{}
		for k, v := range params {
			if buff.Len() > 0 {
				buff.WriteString("&")
			}
			var value string
			if ok, _ := Compare.IsArray(v); ok {
				value = strings.Join(instance.ToArrayOfUrlEncodedString(v), ",")
			} else {
				value = url.QueryEscape(instance.ToString(v))
			}
			if len(value) > 0 {
				buff.WriteString(fmt.Sprintf("%s=%s", k, value))
			}
		}
		return buff.String()
	}
	return ""
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ConversionHelper) toArray(args ...interface{}) []interface{} {
	response := make([]interface{}, 0)
	for _, val := range args {
		aa, tt := Compare.IsArray(val)
		if aa {
			for i := 0; i < tt.Len(); i++ {
				v := tt.Index(i).Interface()
				response = append(response, v)
			}
		} else {
			response = append(response, val)
		}
	}
	return response
}

func (instance *ConversionHelper) toArrayOfString(args ...interface{}) []string {
	response := make([]string, 0)
	for _, val := range args {
		b, tt := Compare.IsArray(val)
		if b {
			for i := 0; i < tt.Len(); i++ {
				v := tt.Index(i).Interface()
				response = append(response, instance.ToString(v))
			}
		} else {
			response = append(response, instance.ToString(val))
		}
	}

	return response
}

func (instance *ConversionHelper) toArrayOfInt(val []interface{}) []int {
	result := make([]int, 0)
	for _, v := range val {
		result = append(result, instance.ToInt(v))
	}
	return result
}

func (instance *ConversionHelper) toArrayOfInt8(val []interface{}) []int8 {
	result := make([]int8, 0)
	for _, v := range val {
		result = append(result, instance.ToInt8(v))
	}
	return result
}

func (instance *ConversionHelper) toArrayOfInt16(val []interface{}) []int16 {
	result := make([]int16, 0)
	for _, v := range val {
		result = append(result, instance.ToInt16(v))
	}
	return result
}

func (instance *ConversionHelper) toArrayOfInt32(val []interface{}) []int32 {
	result := make([]int32, 0)
	for _, v := range val {
		result = append(result, instance.ToInt32(v))
	}
	return result
}

func (instance *ConversionHelper) toArrayOfInt64(val []interface{}) []int64 {
	result := make([]int64, 0)
	for _, v := range val {
		result = append(result, instance.ToInt64(v))
	}
	return result
}

func (instance *ConversionHelper) toArrayOfFloat32(val []interface{}) []float32 {
	result := make([]float32, 0)
	for _, v := range val {
		result = append(result, instance.ToFloat32(v))
	}
	return result
}

func (instance *ConversionHelper) toArrayOfFloat64(val []interface{}) []float64 {
	result := make([]float64, 0)
	for _, v := range val {
		result = append(result, instance.ToFloat64(v))
	}
	return result
}

func (instance *ConversionHelper) tryToArrayOfByte(val interface{}) []byte {
	if v, b := val.([]uint8); b {
		return v
	} else if v, b := val.([]byte); b {
		return v
	} else if v, b := val.(string); b {
		return []byte(v)
	} else if v, b := val.(bool); b {
		if v {
			return []byte{1}
		} else {
			return []byte{0}
		}
	} else {
		return []byte(instance.ToString(val))
	}
}

func (instance *ConversionHelper) toArrayOfByte(args interface{}) []byte {
	if nil != args {
		t := instance.tryToArrayOfByte(args)
		if nil != t {
			return t
		}
		refVal := reflect.ValueOf(args)
		refKind := refVal.Kind()
		switch refKind {
		case reflect.Array, reflect.Slice:
			response := make([]byte, 0)
			for i := 0; i < refVal.Len(); i++ {
				t := instance.tryToArrayOfByte(args)
				if nil != t {
					response = append(response, t...)
				}
			}
			return response
		default:
		}
	}
	return make([]byte, 0)
}

func (instance *ConversionHelper) toMap(val interface{}) map[string]interface{} {
	if m, b := val.(map[string]interface{}); b {
		return m
	}

	// warning: this change the pointer to original object
	data, err := json.Marshal(val)
	if nil == err {
		var m map[string]interface{}
		err = json.Unmarshal(data, &m)
		if nil == err {
			return m
		}
	}

	if b, _ := Compare.IsArray(val); b {
		arr := instance.ToArray(val)
		m := make(map[string]interface{})
		for i, item := range arr {
			m[fmt.Sprintf("%v", i)] = item
		}
		return m
	}

	return nil
}

func (instance *ConversionHelper) toMapOfString(val interface{}) map[string]string {
	if m, b := val.(map[string]string); b {
		return m
	}

	// warning: this change the pointer to original object
	data, err := json.Marshal(val)
	if nil == err {
		var m map[string]string
		err = json.Unmarshal(data, &m)
		if nil == err {
			return m
		}
	}

	return nil
}

func (instance *ConversionHelper) stringToMap(text string) (map[string]interface{}, error) {
	var response map[string]interface{}
	if len(text) > 0 {
		if instance.isValidJsonObject(text) {
			err := json.Unmarshal([]byte(text), &response)
			if nil != err {
				return nil, err
			}
			return response, nil
		} else if array, b := instance.isValidJsonArray(text); b {
			return instance.toMap(array), nil
		} else if strings.Index(text, "=") > 0 {
			uri, err := url.Parse("?" + text)
			if nil == err {
				query := uri.Query()
				if nil != query && len(query) > 0 {
					response = make(map[string]interface{})
					for k, v := range query {
						if len(v) == 1 {
							value := v[0]
							if instance.isValidJsonObject(value) {
								item := map[string]interface{}{}
								err := json.Unmarshal([]byte(value), &item)
								if nil == err {
									response[k] = item
								}
							} else {
								response[k] = value
							}
						}
					}
				}
			} else {
				return nil, err
			}
		}
	}
	return response, nil
}

func (instance *ConversionHelper) isValidJsonObject(text string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(text), &js) == nil
}

func (instance *ConversionHelper) isValidJsonArray(text string) ([]interface{}, bool) {
	var js []interface{}
	err := json.Unmarshal([]byte(text), &js)
	return js, nil == err
}

package qb_utils

import (
	"reflect"
	"strings"
)

type ReflectHelper struct {
}

var Reflect *ReflectHelper

func init() {
	Reflect = new(ReflectHelper)
}

func (instance *ReflectHelper) ValueOf(item interface{}) reflect.Value {
	v := reflect.ValueOf(item)
	k := v.Kind()
	switch k {
	case reflect.Ptr:
		return instance.ValueOf(v.Elem().Interface())
	}
	return v
}

func (instance *ReflectHelper) InterfaceOf(item interface{}) interface{} {
	v := instance.ValueOf(item)
	return v.Interface()
}

func (instance *ReflectHelper) Get(object interface{}, name string) interface{} {
	if m, b := object.(map[string]interface{}); b {
		if nil != m {
			return Maps.Get(m, name)
		}
	} else if b, _ := Compare.IsMap(object); b {
		m := Convert.ToMap(object)
		if nil != m {
			return Maps.Get(m, name)
		}
	} else if b, _ := Compare.IsArray(object); b {
		i := Convert.ToInt(name)
		if i > -1 {
			return Arrays.GetAt(object, i, "")
		}
	} else {
		v := reflect.ValueOf(object)
		if v.IsValid() {
			return getFieldValue(v, name)
		}
	}
	return nil
}

func (instance *ReflectHelper) GetString(object interface{}, name string) string {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToString(v)
	}
	return ""
}

func (instance *ReflectHelper) GetBytes(object interface{}, name string) []byte {
	v := instance.Get(object, name)
	if nil != v {
		return []byte(Convert.ToString(v))
	}
	return make([]byte, 0)
}

func (instance *ReflectHelper) GetInt(object interface{}, name string) int {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToInt(v)
	}
	return 0
}

func (instance *ReflectHelper) GetInt32(object interface{}, name string) int32 {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToInt32(v)
	}
	return 0
}

func (instance *ReflectHelper) GetInt64(object interface{}, name string) int64 {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToInt64(v)
	}
	return 0
}

func (instance *ReflectHelper) GetInt8(object interface{}, name string) int8 {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToInt8(v)
	}
	return 0
}

func (instance *ReflectHelper) GetUint(object interface{}, name string) uint {
	v := instance.Get(object, name)
	if nil != v {
		return uint(Convert.ToInt(v))
	}
	return 0
}

func (instance *ReflectHelper) GetUint8(object interface{}, name string) uint8 {
	v := instance.Get(object, name)
	if nil != v {
		return uint8(Convert.ToInt8(v))
	}
	return 0
}

func (instance *ReflectHelper) GetUint16(object interface{}, name string) uint16 {
	v := instance.Get(object, name)
	if nil != v {
		return uint16(Convert.ToInt16(v))
	}
	return 0
}

func (instance *ReflectHelper) GetUint32(object interface{}, name string) uint32 {
	v := instance.Get(object, name)
	if nil != v {
		return uint32(Convert.ToInt32(v))
	}
	return 0
}

func (instance *ReflectHelper) GetUint64(object interface{}, name string) uint64 {
	v := instance.Get(object, name)
	if nil != v {
		return uint64(Convert.ToInt64(v))
	}
	return 0
}

func (instance *ReflectHelper) GetFloat32(object interface{}, name string) float32 {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToFloat32(v)
	}
	return 0
}

func (instance *ReflectHelper) GetFloat64(object interface{}, name string) float64 {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToFloat64(v)
	}
	return 0
}

func (instance *ReflectHelper) GetBool(object interface{}, name string) bool {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToBool(v)
	}
	return false
}

func (instance *ReflectHelper) GetArray(object interface{}, name string) []interface{} {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToArray(v)
	}
	return []interface{}{}
}

func (instance *ReflectHelper) GetArrayOfString(object interface{}, name string) []string {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToArrayOfString(v)
	}
	return []string{}
}

func (instance *ReflectHelper) GetArrayOfInt(object interface{}, name string) []int {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToArrayOfInt(v)
	}
	return make([]int, 0)
}

func (instance *ReflectHelper) GetArrayOfInt8(object interface{}, name string) []int8 {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToArrayOfInt8(v)
	}
	return make([]int8, 0)
}

func (instance *ReflectHelper) GetArrayOfInt16(object interface{}, name string) []int16 {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToArrayOfInt16(v)
	}
	return make([]int16, 0)
}

func (instance *ReflectHelper) GetArrayOfInt32(object interface{}, name string) []int32 {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToArrayOfInt32(v)
	}
	return make([]int32, 0)
}

func (instance *ReflectHelper) GetArrayOfInt64(object interface{}, name string) []int64 {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToArrayOfInt64(v)
	}
	return make([]int64, 0)
}

func (instance *ReflectHelper) GetArrayOfByte(object interface{}, name string) []byte {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToArrayOfByte(v)
	}
	return make([]byte, 0)
}

func (instance *ReflectHelper) GetMap(object interface{}, name string) map[string]interface{} {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToMap(v)
	}
	return make(map[string]interface{})
}

func (instance *ReflectHelper) GetMapOfString(object interface{}, name string) map[string]string {
	v := instance.Get(object, name)
	if nil != v {
		return Convert.ToMapOfString(v)
	}
	return make(map[string]string)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func getFieldValue(e reflect.Value, name string) interface{} {
	switch e.Kind() {
	case reflect.Ptr:
		elem := e.Elem()
		return getFieldValue(elem, name)
	case reflect.Struct:
		f := e.FieldByName(strings.Title(name))
		if f.IsValid() {
			return f.Interface()
		}
	case reflect.Map:
		m, _ := e.Interface().(map[string]interface{})
		return m[name]
	}
	return nil
}

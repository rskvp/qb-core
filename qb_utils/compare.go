package qb_utils

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode"
)

type CompareHelper struct {
}

var Compare *CompareHelper

//----------------------------------------------------------------------------------------------------------------------
//	i n i t
//----------------------------------------------------------------------------------------------------------------------

func init() {
	Compare = new(CompareHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *CompareHelper) Compare(item1, item2 interface{}) int {
	v1 := Reflect.ValueOf(item1)
	v2 := Reflect.ValueOf(item2)
	if v1.Kind() == v2.Kind() {

		// native
		switch v1.Kind() {
		case reflect.Struct, reflect.Map, reflect.Slice:
			c1 := v1.Interface()
			c2 := v2.Interface()
			if reflect.DeepEqual(c1, c2) {
				return 0
			} else {
				// string check
				s1 := fmt.Sprintf("%v", c1)
				s2 := fmt.Sprintf("%v", c2)
				if s1 > s2 {
					return 1
				} else {
					return -1
				}
			}
		case reflect.Bool:
			c1 := v1.Interface().(bool)
			c2 := v2.Interface().(bool)
			if c1 == c2 {
				return 0
			} else if c1 {
				return 1
			} else {
				return -1
			}
		case reflect.Int:
			c1 := v1.Interface().(int)
			c2 := v2.Interface().(int)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		case reflect.Int8:
			c1 := v1.Interface().(int8)
			c2 := v2.Interface().(int8)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		case reflect.Int16:
			c1 := v1.Interface().(int16)
			c2 := v2.Interface().(int16)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		case reflect.Int32:
			c1 := v1.Interface().(int32)
			c2 := v2.Interface().(int32)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		case reflect.Int64:
			c1 := v1.Interface().(int64)
			c2 := v2.Interface().(int64)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		case reflect.Uint:
			c1 := v1.Interface().(uint)
			c2 := v2.Interface().(uint)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		case reflect.Uint8:
			c1 := v1.Interface().(uint8)
			c2 := v2.Interface().(uint8)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		case reflect.Uint16:
			c1 := v1.Interface().(uint16)
			c2 := v2.Interface().(uint16)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		case reflect.Uint32:
			c1 := v1.Interface().(uint32)
			c2 := v2.Interface().(uint32)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		case reflect.Uint64:
			c1 := v1.Interface().(uint64)
			c2 := v2.Interface().(uint64)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		case reflect.Float32:
			c1 := v1.Interface().(float32)
			c2 := v2.Interface().(float32)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		case reflect.Float64:
			c1 := v1.Interface().(float64)
			c2 := v2.Interface().(float64)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		case reflect.String:
			c1 := v1.Interface().(string)
			c2 := v2.Interface().(string)
			if c1 == c2 {
				return 0
			} else if c1 > c2 {
				return 1
			} else {
				return -1
			}
		}
	}
	return -1
}

func (instance *CompareHelper) Equals(item1, item2 interface{}) bool {
	return instance.Compare(item1, item2) == 0
}

func (instance *CompareHelper) NotEquals(val1, val2 interface{}) bool {
	return !instance.Equals(val1, val2)
}

func (instance *CompareHelper) IsZero(item interface{}) bool {
	if nil != item {
		i := Convert.ToIntDef(item, -1)
		if i == 0 {
			return true
		} else {
			return len(Convert.ToString(item)) == 0
		}
	}
	return true
}

func (instance *CompareHelper) IsGreater(item1, item2 interface{}) bool {
	return instance.Compare(item1, item2) > 0
}

func (instance *CompareHelper) IsLower(item1, item2 interface{}) bool {
	return instance.Compare(item1, item2) < 0
}

func (instance *CompareHelper) IsString(val interface{}) (bool, string) {
	v, vv := val.(string)
	if vv {
		return true, v
	}
	return false, ""
}

func (instance *CompareHelper) IsStringASCII(val string) bool {
	for i := 0; i < len(val); i++ {
		if val[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func (instance *CompareHelper) IsInt(val interface{}) (bool, int) {
	v, vv := val.(int)
	if vv {
		return true, v
	}
	return false, 0
}

func (instance *CompareHelper) IsBool(val interface{}) (bool, bool) {
	v, vv := val.(bool)
	if vv {
		return true, v
	}
	return false, false
}

func (instance *CompareHelper) IsFloat32(val interface{}) (bool, float32) {
	v, vv := val.(float32)
	if vv {
		return true, v
	}
	return false, 0
}

func (instance *CompareHelper) IsFloat64(val interface{}) (bool, float64) {
	v, vv := val.(float64)
	if vv {
		return true, v
	}
	return false, 0
}

func (instance *CompareHelper) IsNumeric(val interface{}) (bool, float64) {
	switch i := val.(type) {
	case float32:
		return true, float64(i)
	case float64:
		return true, i
	case int:
		return true, float64(i)
	case int8:
		return true, float64(i)
	case int16:
		return true, float64(i)
	case int32:
		return true, float64(i)
	case int64:
		return true, float64(i)
	case string:
		v, err := strconv.ParseFloat(i, 64)
		if nil == err {
			return true, v
		}
	}
	return false, 0
}

func (instance *CompareHelper) IsArray(val interface{}) (bool, reflect.Value) {
	rt := Reflect.ValueOf(val)
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		return true, rt
	default:
		return false, rt
	}
}

func (instance *CompareHelper) IsArrayNotEmpty(array interface{}) (bool, reflect.Value) {
	s := Reflect.ValueOf(array)
	b := (s.Kind() == reflect.Slice || s.Kind() == reflect.Array) && s.Len() > 0
	if b {
		return true, s
	}
	return false, s
}

func (instance *CompareHelper) IsMap(val interface{}) (bool, reflect.Value) {
	rt := Reflect.ValueOf(val)
	switch rt.Kind() {
	case reflect.Map:
		return true, rt
	case reflect.Ptr:
		return instance.IsMap(rt.Elem().Interface())
	default:
		return false, rt
	}
}

func (instance *CompareHelper) IsStruct(val interface{}) (bool, reflect.Value) {
	rt := Reflect.ValueOf(val)
	switch rt.Kind() {
	case reflect.Struct:
		return true, rt
	case reflect.Ptr:
		return instance.IsStruct(rt.Elem().Interface())
	default:
		return false, rt
	}
}

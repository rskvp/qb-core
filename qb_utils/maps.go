package qb_utils

import (
	"strings"
)

type MapsHelper struct {
}

var Maps *MapsHelper

//----------------------------------------------------------------------------------------------------------------------
//	i n i t
//----------------------------------------------------------------------------------------------------------------------

func init() {
	Maps = new(MapsHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *MapsHelper) Keys(m map[string]interface{}) []string {
	response := make([]string, 0)
	for k, _ := range m {
		return	append(response, k)
	}
	return response
}

func (instance *MapsHelper) Values(m map[string]interface{}) []interface{} {
	response := make([]interface{}, 0)
	for _, v := range m {
		response = append(response, v)
	}
	return response
}

func (instance *MapsHelper) ValuesOfKeys(m map[string]interface{}, keys []string) []interface{} {
	response := make([]interface{}, 0)
	for _, k := range keys {
		if v, b := m[k]; b {
			response = append(response, v)
		}
	}
	return response
}

func (instance *MapsHelper) KeyValuePairs(m map[string]interface{}) ([]string, []interface{}) {
	kr := make([]string, 0)
	vr := make([]interface{}, 0)
	for k, v := range m {
		kr = append(kr, k)
		vr = append(vr, v)
	}
	return kr, vr
}

func (instance *MapsHelper) Clone(m map[string]interface{}) map[string]interface{} {
	response := make(map[string]interface{})
	for k, v := range m {
		tm, tmb := instance.isMap(v)
		if tmb {
			response[k] = instance.Clone(tm)
		} else {
			response[k] = v
		}
	}
	return response
}

func (instance *MapsHelper) Merge(overwrite bool, target map[string]interface{}, sources ...map[string]interface{}) map[string]interface{} {
	if nil != target {
		for _, source := range sources {
			_, _ = instance.merge(target, source, overwrite, nil)
		}
	}
	return target
}

func (instance *MapsHelper) MergeCount(overwrite bool, target map[string]interface{}, sources ...map[string]interface{}) int {
	count := 0
	if nil != target {
		for _, source := range sources {
			c, _ := instance.merge(target, source, overwrite, nil)
			count += c
		}
	}
	return count
}

func (instance *MapsHelper) MergeFields(overwrite bool, target map[string]interface{}, sources ...map[string]interface{}) (fields []string) {
	if nil != target {
		for _, source := range sources {
			_, f := instance.merge(target, source, overwrite, nil)
			fields = append(fields, f...)
		}
	}
	return
}

func (instance *MapsHelper) MergeExclusion(overwrite bool, exclusions []string, target map[string]interface{}, sources ...map[string]interface{}) map[string]interface{} {
	if nil != target {
		for _, source := range sources {
			instance.merge(target, source, overwrite, exclusions)
		}
	}
	return target
}

func (instance *MapsHelper) Get(m map[string]interface{}, path string) interface{} {
	if nil != m && len(path) > 0 {
		itemMap := m
		tokens := strings.Split(path, ".")
		length := len(tokens)
		for i := 0; i < length; i++ {
			token := tokens[i]
			tv, tb := itemMap[token]
			if tb {
				if i == length-1 {
					return tv
				} else {
					if mm, b := instance.isMap(tv); b {
						itemMap = mm
					}
				}
			}
		}
	}
	return nil
}

func (instance *MapsHelper) Set(m map[string]interface{}, path string, value interface{}) {
	if nil != m && len(path) > 0 {
		itemMap := m
		tokens := strings.Split(path, ".")
		length := len(tokens)
		for i := 0; i < length; i++ {
			token := tokens[i]
			if i == length-1 {
				itemMap[token] = value
			} else {
				if _, tb := itemMap[token]; !tb {
					itemMap[token] = map[string]interface{}{}
				}
				itemMap = itemMap[token].(map[string]interface{})
			}
		}
	}
}

func (instance *MapsHelper) GetString(m map[string]interface{}, path string) string {
	return Convert.ToString(instance.Get(m, path))
}

func (instance *MapsHelper) GetBytes(m map[string]interface{}, path string) []byte {
	return []byte(Convert.ToString(instance.Get(m, path)))
}

func (instance *MapsHelper) GetBool(m map[string]interface{}, path string) bool {
	return Convert.ToBool(instance.Get(m, path))
}

func (instance *MapsHelper) GetInt(m map[string]interface{}, path string) int {
	return Convert.ToInt(instance.Get(m, path))
}

func (instance *MapsHelper) GetInt8(m map[string]interface{}, path string) int8 {
	return Convert.ToInt8(instance.Get(m, path))
}

func (instance *MapsHelper) GetInt16(m map[string]interface{}, path string) int16 {
	return Convert.ToInt16(instance.Get(m, path))
}

func (instance *MapsHelper) GetInt32(m map[string]interface{}, path string) int32 {
	return Convert.ToInt32(instance.Get(m, path))
}

func (instance *MapsHelper) GetInt64(m map[string]interface{}, path string) int64 {
	return Convert.ToInt64(instance.Get(m, path))
}

func (instance *MapsHelper) GetUint(m map[string]interface{}, path string) uint {
	return uint(Convert.ToInt(instance.Get(m, path)))
}

func (instance *MapsHelper) GetUint8(m map[string]interface{}, path string) uint8 {
	return uint8(Convert.ToInt8(instance.Get(m, path)))
}

func (instance *MapsHelper) GetUint16(m map[string]interface{}, path string) uint16 {
	return uint16(Convert.ToInt16(instance.Get(m, path)))
}

func (instance *MapsHelper) GetUint32(m map[string]interface{}, path string) uint32 {
	return uint32(Convert.ToInt32(instance.Get(m, path)))
}

func (instance *MapsHelper) GetUint64(m map[string]interface{}, path string) uint64 {
	return uint64(Convert.ToInt64(instance.Get(m, path)))
}

func (instance *MapsHelper) GetFloat32(m map[string]interface{}, path string) float32 {
	return Convert.ToFloat32(instance.Get(m, path))
}

func (instance *MapsHelper) GetFloat64(m map[string]interface{}, path string) float64 {
	return Convert.ToFloat64(instance.Get(m, path))
}

func (instance *MapsHelper) GetArray(m map[string]interface{}, path string) []interface{} {
	return Convert.ToArray(instance.Get(m, path))
}

func (instance *MapsHelper) GetArrayOfString(m map[string]interface{}, path string) []string {
	return Convert.ToArrayOfString(instance.Get(m, path))
}

func (instance *MapsHelper) GetArrayOfByte(m map[string]interface{}, path string) []byte {
	return Convert.ToArrayOfByte(instance.Get(m, path))
}

func (instance *MapsHelper) GetMap(m map[string]interface{}, path string) map[string]interface{} {
	return Convert.ToMap(instance.Get(m, path))
}

func (instance *MapsHelper) GetMapOfString(m map[string]interface{}, path string) map[string]string {
	return Convert.ToMapOfString(instance.Get(m, path))
}

func (instance *MapsHelper) GetMapOfStringArray(m map[string]interface{}, path string) map[string][]string {
	return Convert.ToMapOfStringArray(instance.Get(m, path))
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *MapsHelper) isMap(v interface{}) (map[string]interface{}, bool) {
	if m, b := v.(map[string]interface{}); b {
		return m, true
	}
	return nil, false
}

func (instance *MapsHelper) merge(target map[string]interface{}, source map[string]interface{}, overwrite bool, exclusion []string) (count int, fields []string) {
	count = 0
	fields = make([]string, 0)
	if nil != target && nil != source {
		for sk, sv := range source {
			if len(exclusion) == 0 || Arrays.IndexOf(sk, exclusion) < 0 {
				tv, tb := target[sk]
				tm, tmb := instance.isMap(tv)
				sm, smb := instance.isMap(sv)
				if smb && tmb {
					c, f := instance.merge(tm, sm, overwrite, exclusion)
					count += c
					fields = append(fields, f...)
					continue
				} else if b, _ := Compare.IsArray(tv); b && !overwrite {
					target[sk] = Arrays.AppendUnique(tv, sv)
					count++
					fields = append(fields, sk)
					continue
				}
				if !tb || overwrite {
					target[sk] = sv
					count++
					fields = append(fields, sk)
					continue
				}
			}
		}
	}
	return
}

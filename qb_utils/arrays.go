package qb_utils

import (
	"math/rand"
	"reflect"
	"sort"
)

type ArraysHelper struct {
}

var Arrays *ArraysHelper

func init() {
	Arrays = new(ArraysHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ArraysHelper) GetLast(array interface{}, defValue interface{}) interface{} {
	v := Reflect.ValueOf(array)
	k := v.Kind()
	switch k {
	case reflect.Slice, reflect.Array:
		return v.Index(v.Len() - 1).Interface()
	}
	return defValue
}

func (instance *ArraysHelper) GetAt(array interface{}, index int, defValue interface{}) interface{} {
	v := Reflect.ValueOf(array)
	k := v.Kind()
	switch k {
	case reflect.Slice, reflect.Array:
		if v.Len() > index {
			return v.Index(index).Interface()
		}
	}
	return defValue
}

func (instance *ArraysHelper) Sort(array interface{}) {
	if a, b := array.([]interface{}); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]string); b {
		sort.Strings(a)
	} else if a, b := array.([]byte); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]int); b {
		sort.Ints(a)
	} else if a, b := array.([]int8); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]int16); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]int32); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]int64); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]uint); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]uint8); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]uint16); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]uint32); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]uint64); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]uintptr); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]float32); b {
		sort.Slice(array, func(i, j int) bool {
			return Compare.IsLower(a[i], a[j])
		})
	} else if a, b := array.([]float64); b {
		sort.Float64s(a)
	}
}

func (instance *ArraysHelper) SortDesc(array interface{}) {
	instance.Sort(array)
	instance.Reverse(array)
}

func (instance *ArraysHelper) Reverse(array interface{}) {
	if a, b := array.([]interface{}); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]string); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]byte); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]int); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]int8); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]int16); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]int32); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]int64); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]uint); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]uint8); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]uint16); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]uint32); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]uint64); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]uintptr); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]float32); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]float64); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	}
}

func (instance *ArraysHelper) AppendUnique(target interface{}, source interface{}) interface{} {
	vt := Reflect.ValueOf(target)
	vs := Reflect.ValueOf(source)
	if vt.Kind() == reflect.Slice && vs.Kind() == reflect.Slice {
		for i := 0; i < vs.Len(); i++ {
			vsv := vs.Index(i)
			if instance.IndexOf(vsv.Interface(), target) == -1 {
				vt = reflect.Append(vt, vsv)
			}
		}
	} else {
		if instance.IndexOf(source, target) == -1 {
			vt = reflect.Append(vt, reflect.ValueOf(source))
		}
	}
	return vt.Interface()
}

func (instance *ArraysHelper) AppendUniqueFunc(target interface{}, source interface{}, callback func(t interface{}, s interface{}) bool) interface{} {
	if nil == callback {
		return target
	}
	vt := Reflect.ValueOf(target) // value of target
	vs := Reflect.ValueOf(source) // value of source
	if vt.Kind() == reflect.Slice && vs.Kind() == reflect.Slice {
		for i := 0; i < vs.Len(); i++ {
			sourceItem := vs.Index(i) // source item
			if vt.Len() > 0 {
				for ii := 0; ii < vt.Len(); ii++ {
					targetItem := vt.Index(ii) // target item
					addThis := callback(targetItem.Interface(), sourceItem.Interface())
					if addThis {
						vt = reflect.Append(vt, sourceItem)
					}
				}
			} else {
				vt = reflect.Append(vt, sourceItem)
			}
		}
	} else if vt.Kind() == reflect.Slice && (vs.Kind() == reflect.Struct || vs.Kind() == reflect.Map) {
		if vt.Len() > 0 {
			addThis := true
			for ii := 0; ii < vt.Len(); ii++ {
				targetItem := vt.Index(ii) // target item
				addThis = addThis && callback(targetItem.Interface(), source)
				if !addThis {
					break
				}
			}
			if addThis {
				vt = reflect.Append(vt, reflect.ValueOf(source))
			}
		} else {
			vt = reflect.Append(vt, reflect.ValueOf(source))
		}
	}
	return vt.Interface()
}

func (instance *ArraysHelper) ForEach(array interface{}, callback func(item interface{}) bool) {
	if nil == callback {
		return
	}
	s := Reflect.ValueOf(array)
	if s.Kind() == reflect.Slice || s.Kind() == reflect.Array {
		for i := 0; i < s.Len(); i++ {
			v := s.Index(i)
			if callback(v.Interface()) {
				return // exit loop
			}
		}
	}
}

func (instance *ArraysHelper) Filter(array interface{}, callback func(item interface{}) bool) interface{} {
	if nil == callback {
		return array
	}
	s := Reflect.ValueOf(array)
	if s.Kind() == reflect.Slice || s.Kind() == reflect.Array {
		if s.Len() > 0 {
			retType := reflect.TypeOf(array)
			r := reflect.MakeSlice(retType, 0, 0)
			for i := 0; i < s.Len(); i++ {
				v := s.Index(i)
				if v.IsValid() && callback(v.Interface()) {
					r = reflect.Append(r, v)
				}
			}
			return r.Interface()
		}
	}
	return array
}

func (instance *ArraysHelper) Map(array interface{}, callback func(item interface{}) interface{}) interface{} {
	if nil == callback {
		return array
	}
	s := Reflect.ValueOf(array)
	if s.Kind() == reflect.Slice || s.Kind() == reflect.Array {
		if s.Len() > 0 {
			retType := reflect.TypeOf(array)
			r := reflect.MakeSlice(retType, 0, 0)
			for i := 0; i < s.Len(); i++ {
				v := s.Index(i)
				if v.IsValid() {
					newValue := callback(v.Interface())
					if nil != newValue {
						r = reflect.Append(r, reflect.ValueOf(newValue))
					} else {
						r = reflect.Append(r, v)
					}
				}
			}
			return r.Interface()
		}
	}
	return array
}

func (instance *ArraysHelper) IndexOfFunc(array interface{}, callback func(item interface{}) bool) int {
	if nil == callback {
		return -1
	}
	s := Reflect.ValueOf(array)
	if s.Kind() == reflect.Slice || s.Kind() == reflect.Array {
		if s.Len() > 0 {
			for i := 0; i < s.Len(); i++ {
				v := s.Index(i)
				if v.IsValid() && callback(v.Interface()) {
					return i
				}
			}
		}
	}
	return -1
}

func (instance *ArraysHelper) IndexOf(item interface{}, array interface{}) int {
	s := Reflect.ValueOf(array)
	if s.Kind() == reflect.Slice || s.Kind() == reflect.Array {
		for i := 0; i < s.Len(); i++ {
			v := s.Index(i)
			if v.IsValid() && Compare.Equals(v.Interface(), item) {
				return i
			}
		}
	}
	return -1
}

func (instance *ArraysHelper) Count(item interface{}, array interface{}) int {
	count := 0
	s := Reflect.ValueOf(array)
	if s.Kind() == reflect.Slice || s.Kind() == reflect.Array {
		for i := 0; i < s.Len(); i++ {
			v := s.Index(i)
			if v.IsValid() && Compare.Equals(v.Interface(), item) {
				count++
			}
		}
	}
	return count
}

func (instance *ArraysHelper) Remove(item interface{}, array interface{}) interface{} {
	if b, s := Compare.IsArrayNotEmpty(array); b {
		retType := reflect.TypeOf(array)
		r := reflect.MakeSlice(retType, 0, 0)
		for i := 0; i < s.Len(); i++ {
			v := s.Index(i)
			if v.IsValid() {
				if Compare.Equals(v.Interface(), item) {
					continue
				}
				r = reflect.Append(r, v)
			}
		}
		return r.Interface()
	}
	return array
}

func (instance *ArraysHelper) RemoveIndex(idx int, array interface{}) interface{} {
	s := Reflect.ValueOf(array)
	if s.Kind() == reflect.Slice || s.Kind() == reflect.Array {
		if s.Len() > 0 {
			retType := reflect.TypeOf(array)
			r := reflect.MakeSlice(retType, 0, 0)
			for i := 0; i < s.Len(); i++ {
				if i == idx {
					continue
				}
				v := s.Index(i)
				if v.IsValid() {
					r = reflect.Append(r, v)
				}
			}
			return r.Interface()
		}
	}
	return array
}

func (instance *ArraysHelper) ReplaceAll(source, target interface{}, array interface{}) interface{} {
	return instance.Replace(source, target, array, -1, false)
}

func (instance *ArraysHelper) Replace(source interface{}, target interface{}, array interface{}, n int, shrink bool) interface{} {
	if b, s := Compare.IsArrayNotEmpty(array); b {
		retType := reflect.TypeOf(array)
		r := reflect.MakeSlice(retType, 0, 0)
		count := 0
		replaced := make([]interface{}, 0)
		for i := 0; i < s.Len(); i++ {
			v := s.Index(i)
			if v.IsValid() {
				if equals(source, v.Interface()) && (n < 1 || count < n) {
					count++
					if !shrink {
						r = reflect.Append(r, reflect.ValueOf(target))
					} else {
						if instance.IndexOf(target, r.Interface()) == -1 {
							r = reflect.Append(r, reflect.ValueOf(target))
							replaced = append(replaced, v.Interface())
						}
					}
				} else {
					r = reflect.Append(r, v)
				}
			}
		}
		return r.Interface()
	}
	return array
}

// Copy a slice and return new slice with same items
func (instance *ArraysHelper) Copy(array interface{}) interface{} {
	if a, b := array.([]interface{}); b {
		response := make([]interface{}, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]string); b {
		response := make([]string, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]byte); b {
		response := make([]byte, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]int); b {
		response := make([]int, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]int8); b {
		response := make([]int8, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]int16); b {
		response := make([]int16, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]int32); b {
		response := make([]int32, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]int64); b {
		response := make([]int64, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]uint); b {
		response := make([]uint, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]uint8); b {
		response := make([]uint8, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]uint16); b {
		response := make([]uint16, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]uint32); b {
		response := make([]uint32, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]uint64); b {
		response := make([]uint64, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]uintptr); b {
		response := make([]uintptr, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]float32); b {
		response := make([]float32, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]float64); b {
		response := make([]float64, len(a))
		copy(response, a)
		return response
	}
	return nil
}

// Group a slice in batch.
// Returns a slice of slice.
// usage:
// response := Group([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
// fmt.Println(response) // [[0 1 2] [3 4 5] [6 7 8] [9]]
func (instance *ArraysHelper) Group(groupSize int, array interface{}) interface{} {
	if a, b := array.([]interface{}); b {
		groups := make([][]interface{}, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]string); b {
		groups := make([][]string, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]byte); b {
		groups := make([][]byte, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]int); b {
		groups := make([][]int, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]int8); b {
		groups := make([][]int8, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]int16); b {
		groups := make([][]int16, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]int32); b {
		groups := make([][]int32, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]int64); b {
		groups := make([][]int64, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]uint); b {
		groups := make([][]uint, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]uint8); b {
		groups := make([][]uint8, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]uint16); b {
		groups := make([][]uint16, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]uint32); b {
		groups := make([][]uint32, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]uint64); b {
		groups := make([][]uint64, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]uintptr); b {
		groups := make([][]uintptr, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]float32); b {
		groups := make([][]float32, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]float64); b {
		groups := make([][]float64, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	}
	return nil
}

func (instance *ArraysHelper) Sub(array interface{}, start, end int) interface{} {
	if start >= end {
		start = 0
	}
	if a, b := array.([]interface{}); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]interface{}, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]string); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]string, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]byte); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]byte, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]int); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]int, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]int8); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]int8, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]int16); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]int16, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]int32); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]int32, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]int64); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]int64, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]uint); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]uint, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]uint8); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]uint8, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]uint16); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]uint16, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]uint32); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]uint32, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]uint64); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]uint64, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]uintptr); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]uintptr, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]float32); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]float32, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]float64); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]float64, 0)
		result = append(result, a[start:end+1]...)
		return result
	}
	return nil
}

func (instance *ArraysHelper) Shuffle(array interface{}) {
	if a, b := array.([]interface{}); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]string); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]byte); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]int); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]int8); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]int16); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]int32); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]int64); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]uint); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]uint8); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]uint16); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]uint32); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]uint64); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]uintptr); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]float32); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]float64); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func equals(source, target interface{}) bool {
	if ab, sourceArray := Compare.IsArray(source); ab {
		for ii := 0; ii < sourceArray.Len(); ii++ {
			if sourceValue := sourceArray.Index(ii); sourceValue.IsValid() {
				if Compare.Equals(target, sourceValue.Interface()) {
					return true
				}
			}
		}
		return false
	}
	return Compare.Equals(target, source)
}

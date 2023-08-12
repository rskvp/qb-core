package qb_structs

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_utils"
)

type StructHelper struct {
}

var Structs *StructHelper

func init() {
	Structs = new(StructHelper)
}

func (instance *StructHelper) NewBuilder() *StructBuilder {
	return NewBuilder()
}

//----------------------------------------------------------------------------------------------------------------------
//	structs
//----------------------------------------------------------------------------------------------------------------------

type StructBuilder struct {
	autoTagJson bool
	fields      []reflect.StructField
	values      map[string]interface{}
}

func NewBuilder() *StructBuilder {
	instance := new(StructBuilder)
	instance.values = make(map[string]interface{})
	return instance
}

func (instance *StructBuilder) EnableAutoTagJson(value bool) *StructBuilder {
	instance.autoTagJson = value
	return instance
}

func (instance *StructBuilder) Set(name string, value interface{}) *StructBuilder {
	return instance.AddFieldByValue(name, "", value)
}

func (instance *StructBuilder) AddFieldByValue(name, tag string, value interface{}) *StructBuilder {
	name = qb_utils.Strings.CapitalizeFirst(name)
	instance.addField(name, tag, reflect.TypeOf(value))
	instance.values[name] = value
	return instance
}

// AddFieldByStringType add a field using its string type name with a default value
func (instance *StructBuilder) AddFieldByStringType(name, tag string, sType string) *StructBuilder {
	name = qb_utils.Strings.CapitalizeFirst(name)
	t, v := instance.typeFromText(sType)
	instance.addField(name, tag, t)
	instance.values[name] = v
	return instance
}

// Elem return element you can use to add fields: ex. elem.Field(0).SetInt(1234)
func (instance *StructBuilder) Elem() reflect.Value {
	t := reflect.StructOf(instance.fields)
	return reflect.New(t).Elem()
}

func (instance *StructBuilder) Interface() interface{} {
	elem := instance.Elem()
	for k, v := range instance.values {
		value := reflect.ValueOf(v)
		field := elem.FieldByName(k)
		field.Set(value)
	}
	return elem.Addr().Interface()
}

func (instance *StructBuilder) Json() string {
	return qb_utils.JSON.Stringify(instance.Interface())
}

func (instance *StructBuilder) addField(name, tag string, typeOf reflect.Type) *StructBuilder {
	if len(tag) == 0 && instance.autoTagJson {
		tag = fmt.Sprintf(`json:"%s"`, qb_utils.Strings.Underscore(name))
	}
	instance.fields = append(instance.fields, reflect.StructField{
		Name: name,
		Tag:  reflect.StructTag(tag),
		Type: typeOf,
	})
	return instance
}

func (instance *StructBuilder) typeFromText(sType string) (t reflect.Type, v interface{}) {
	sType = strings.ToLower(sType)
	switch sType {
	case "int", "integer", "number":
		v = 0
	case "float", "float32":
		v = float32(0.0)
	case "float64", "long":
		v = float64(0.0)
	case "time", "date", "datetime":
		v = time.Now()
	case "bool", "boolean":
		v = true
	case "byte":
		v = byte(0)
	case "rune":
		v = rune(0)
	default:
		v = ""
	}
	t = reflect.TypeOf(v)
	return
}

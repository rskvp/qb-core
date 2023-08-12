package qb_i18n_bundle

import (
	"strings"

	"github.com/rskvp/qb-core/qb_utils"
)

var (
	PLURAL_FIELDS = []string{"value", "count", "values", "qty", "length", "len", "size", "num"}
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type Deserializer func(data []byte) (map[string]*Element, error)

type Element struct {
	Zero  string `json:"zero"`
	One   string `json:"one"`
	Two   string `json:"two"`
	Other string `json:"other"`
}

func NewElementFromString(s string) *Element {
	return &Element{
		Zero:  s,
		One:   s,
		Two:   s,
		Other: s,
	}
}

func NewElementFromMap(m map[string]interface{}) (*Element, error) {
	var instance Element
	err := qb_utils.JSON.Read(qb_utils.Convert.ToString(m), &instance)
	if nil != err {
		return nil, err
	}
	if len(instance.Two) == 0 {
		instance.Two = instance.Other
	}
	if len(instance.Zero) == 0 {
		instance.Zero = instance.Other
	}

	return &instance, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Element) Value() string {
	if len(instance.One) > 0 {
		return instance.One
	} else if len(instance.Two) > 0 {
		return instance.Two
	} else if len(instance.Zero) > 0 {
		return instance.Zero
	}
	return instance.Other
}

func (instance *Element) Plural(value float64) string {
	switch value {
	case -1:
		return instance.Value()
	case 0:
		return instance.Zero
	case 1:
		return instance.One
	case 2:
		return instance.Two
	default:
		return instance.Other
	}
}

func (instance *Element) PluralValueWithFields(pluralFields []string, model ...interface{}) float64 {
	if nil == pluralFields {
		pluralFields = PLURAL_FIELDS
	}
	return lookUpPluralValue(pluralFields, model...)
}

/*
func (instance *Element) Get(model ...interface{}) string {
	text := instance.Value()
	// out, _ := Render(text, model...)
	return text//out
}

func (instance *Element) GetWithPluralValue(value float64, model ...interface{}) string {
	text := instance.Plural(value)
	// out, _ := Render(text, model...)
	return text //  out
}

func (instance *Element) GetWithPlural(pluralFields []string, model ...interface{}) string {
	if nil == pluralFields {
		pluralFields = PLURAL_FIELDS
	}
	value := lookUpPluralValue(pluralFields, model...)
	return instance.GetWithPluralValue(value, model...)
}*/

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func lookUpPluralValue(pluralFields []string, model ...interface{}) float64 {
	result := float64(-1)
	for _, m := range model {
		if mm, b := m.(map[string]interface{}); b {
			for k, v := range mm {
				if isPluralField(k, pluralFields) {
					result = qb_utils.Convert.ToFloat64(v)
					goto exit
				}
				if b, v := qb_utils.Compare.IsNumeric(v); b {
					result = v
				}
			}
		}
	}
exit:
	return result
}

func isPluralField(name string, pluralFields []string) bool {
	if qb_utils.Arrays.IndexOf(name, pluralFields) > -1 {
		return true
	}
	if strings.Index(name, SUFFIX_SEP) > -1 {
		name := qb_utils.Convert.ToString(qb_utils.Arrays.GetAt(strings.Split(name, SUFFIX_SEP), 1, ""))
		return isPluralField(name, pluralFields)
	}
	return false
}

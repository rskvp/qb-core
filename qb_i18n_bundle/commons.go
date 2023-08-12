package qb_i18n_bundle

import (
	"errors"

	"github.com/rskvp/qb-core/qb_utils"
	"golang.org/x/text/language"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const (
	SUFFIX_SEP    = "_"
	SUFFIX_LENGHT = "_length"
	SUFFIX_INDEX  = "_index"
)

//----------------------------------------------------------------------------------------------------------------------
//	e r r o r s
//----------------------------------------------------------------------------------------------------------------------

var (
	LanguageNotSupportedError      = errors.New("language_not_supported")
	DeserializerNotRegisteredError = errors.New("deserializer_not_registered")
	ResourceNotFoundError          = errors.New("resource_not_found")
)


var deserializers map[string]Deserializer

func init() {
	deserializers = make(map[string]Deserializer)
	deserializers["json"] = deserializeJSON
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r s
//----------------------------------------------------------------------------------------------------------------------

func NewCustomBundleFromTag(defTag language.Tag, deserializers *map[string]Deserializer) *Bundle {
	instance := new(Bundle)
	instance.defTag = defTag
	instance.defLang = langOf(defTag.String())
	instance.deserializers = deserializers
	instance.cache = make(map[string]*bundleCache)

	instance.Strictly = false
	instance.IgnoreUnknownDeserializers = true
	instance.renderer = new(Renderer)

	return instance
}

func NewCustomBundle(defLang string, deserializers *map[string]Deserializer) (*Bundle, error) {
	tag, err := language.Parse(defLang)
	if nil != err {
		return nil, err
	}
	return NewCustomBundleFromTag(tag, deserializers), nil
}

func NewBundleFromLang(defLang interface{}) (*Bundle, error) {
	if t, b := defLang.(language.Tag); b {
		return NewCustomBundleFromTag(t, &deserializers), nil
	} else if s, b := defLang.(string); b {
		return NewCustomBundle(s, &deserializers)
	}
	return nil, LanguageNotSupportedError
}

func NewBundleFromDir(defLang interface{}, path string) (*Bundle, error) {
	bundle, err := NewBundleFromLang(defLang)
	if nil != err {
		return nil, err
	}
	err = bundle.LoadAll(path)
	if nil != err {
		return nil, err
	}
	return bundle, nil
}

func NewBundleFromFiles(defLang interface{}, files []string) (*Bundle, error) {
	bundle, err := NewBundleFromLang(defLang)
	if nil != err {
		return nil, err
	}
	err = bundle.LoadFiles(files)
	if nil != err {
		return nil, err
	}
	return bundle, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func RegisterDeserializer(extension string, deserializer Deserializer) {
	deserializers[extension] = deserializer
}

func DeserializeFile(filename string) (map[string]*Element, error) {
	if b, err := qb_utils.Paths.Exists(filename); b {
		ext := qb_utils.Paths.ExtensionName(filename)
		if deserializer, b := deserializers[ext]; b {
			data, err := qb_utils.IO.ReadBytesFromFile(filename)
			if nil != err {
				return nil, err
			}
			return deserializer(data)
		} else {
			return nil, DeserializerNotRegisteredError
		}
	} else {
		return nil, err
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func deserializeJSON(data []byte) (map[string]*Element, error) {
	response := make(map[string]*Element)
	var m map[string]interface{}
	err := qb_utils.JSON.Read(data, &m)
	if nil == err {
		parse("", m, response)
	}
	return response, err
}

func parse(parentKey string, m map[string]interface{}, data map[string]*Element) {
	for k, v := range m {
		key := parentKey
		if len(key) > 0 {
			key += "."
		}
		key += k
		if mm, b := v.(map[string]interface{}); b {
			// test if is an element
			elem, err := NewElementFromMap(mm)
			if nil == err && len(elem.One) > 0 {
				data[key] = elem
			} else {
				parse(key, mm, data)
			}
		} else if s, b := v.(string); b {
			// from simple string
			elem := NewElementFromString(s)
			data[key] = elem
		}
	}
}

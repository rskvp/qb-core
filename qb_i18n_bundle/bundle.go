package qb_i18n_bundle

import (
	"github.com/rskvp/qb-core/qb_utils"
	"golang.org/x/text/language"
)

//----------------------------------------------------------------------------------------------------------------------
//	bundleCache
//----------------------------------------------------------------------------------------------------------------------

type bundleCache struct {
	filename     string
	deserializer Deserializer
	lang         string
	data         map[string]*Element // late initialized
}

func (instance *bundleCache) getData() (map[string]*Element, error) {
	if nil == instance.data {
		if b, err := qb_utils.Paths.Exists(instance.filename); b {
			data, err := qb_utils.IO.ReadBytesFromFile(instance.filename)
			if nil != err {
				return nil, err
			}
			d, err := instance.deserializer(data)
			if nil != err {
				return nil, err
			}
			instance.data = d
		} else {
			return nil, err
		}
	}
	return instance.data, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	Bundle
//----------------------------------------------------------------------------------------------------------------------

type Bundle struct {
	Strictly                   bool // must load data or go panic
	IgnoreUnknownDeserializers bool

	defTag        language.Tag
	defLang       string
	deserializers *map[string]Deserializer
	cache         map[string]*bundleCache // deserializers or array of Element
	renderer      I18nRenderer
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Bundle) SetRenderer(renderer I18nRenderer) {
	instance.renderer = renderer
}

func (instance *Bundle) Renderer() I18nRenderer {
	return instance.renderer
}

func (instance *Bundle) String() string {
	m := map[string]interface{}{
		"def-lang":      instance.defLang,
		"strictly":      instance.Strictly,
		"deserializers": instance.Deserializers(),
	}
	return qb_utils.JSON.Stringify(m)
}
func (instance *Bundle) Deserializers() []string {
	response := make([]string, 0)
	for k, _ := range *instance.deserializers {
		response = append(response, k)
	}
	return response
}
func (instance *Bundle) LoadAll(dirName string) error {
	if b, err := qb_utils.Paths.Exists(dirName); !b {
		return err
	}
	files, err := qb_utils.Paths.ListFiles(dirName, "*.*")
	if nil != err {
		if instance.Strictly {
			panic(err)
		}
		return err
	}
	return instance.LoadFiles(files)
}

func (instance *Bundle) LoadFiles(files []string) error {
	for _, file := range files {
		err := instance.LoadFile(file)
		if nil != err {
			if instance.Strictly {
				panic(err)
			}
			return err
		}
	}
	return nil
}

func (instance *Bundle) LoadFile(file string) error {
	err := instance.addFileToCache(file)
	if nil != err {
		if instance.Strictly {
			panic(err)
		}
		return err
	}
	return nil
}

func (instance *Bundle) GetDictionary(rawLang string) (map[string]*Element, error) {
	lang := langOf(rawLang)
	var cache *bundleCache
	if v, b := instance.cache[lang]; b {
		cache = v
	} else {
		cache = instance.cache[instance.defLang]
	}
	return cache.getData()
}

func (instance *Bundle) GetDictionaryByTag(tag language.Tag) (map[string]*Element, error) {
	return instance.GetDictionary(tag.String())
}

func (instance *Bundle) AddElement(rawLang string, path string, elem *Element) {
	lang := langOf(rawLang)
	if _, b := instance.cache[lang]; !b {
		instance.cache[lang] = &bundleCache{
			filename:     "",
			deserializer: nil,
			lang:         lang,
			data:         make(map[string]*Element),
		}
	}
	cache := instance.cache[lang]
	cache.data[path] = elem
}

func (instance *Bundle) GetElement(lang string, path string) (*Element, error) {
	dictionary, err := instance.GetDictionary(lang)
	if nil != err {
		return nil, err
	}
	var elem *Element
	if v, b := dictionary[path]; !b {
		return nil, ResourceNotFoundError
	} else {
		elem = v
	}
	return elem, nil
}

func (instance *Bundle) Get(lang string, path string, data ...interface{}) (string, error) {
	elem, err := instance.GetElement(lang, path)
	if nil != err {
		return "", err
	}
	elemValue := elem.Value()
	return instance.renderer.Render(elemValue, data...) // elem.Get(data...), nil
}

func (instance *Bundle) GetWithPluralValue(lang string, path string, value float64,
	data ...interface{}) (string, error) {
	elem, err := instance.GetElement(lang, path)
	if nil != err {
		return "", err
	}
	elemPlural := elem.Plural(value)
	return instance.renderer.Render(elemPlural, data...) // elem.GetWithPluralValue(value, data...), nil
}

func (instance *Bundle) GetWithPlural(lang string, path string, pluralFields []string,
	data ...interface{}) (string, error) {
	elem, err := instance.GetElement(lang, path)
	if nil != err {
		return "", err
	}
	value := elem.PluralValueWithFields(pluralFields, data...)
	elemPlural := elem.Plural(value)
	return instance.renderer.Render(elemPlural, data...) //elem.GetWithPlural(pluralFields, data...), nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Bundle) addFileToCache(filename string) error {
	ext := qb_utils.Paths.ExtensionName(filename)
	deserializers := *instance.deserializers
	if v, b := deserializers[ext]; !b {
		if instance.IgnoreUnknownDeserializers {
			return nil
		}
		return DeserializerNotRegisteredError
	} else {
		name := qb_utils.Paths.FileName(filename, false)
		lang, err := validate(name)
		if nil == err {
			return instance.addToCache(lang, filename, v)
		} else {
			// not an i18n file name
		}
	}
	return nil
}

func (instance *Bundle) addToCache(lang string, filename string, deserializer Deserializer) error {
	// add default
	if _, b := instance.cache[instance.defLang]; !b || isLangFile(instance.defTag, filename) {
		instance.cache[instance.defLang] = &bundleCache{
			lang:         instance.defLang,
			filename:     filename,
			deserializer: deserializer,
		}
	}
	// add custom lang
	instance.cache[lang] = &bundleCache{
		lang:         lang,
		filename:     filename,
		deserializer: deserializer,
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func validate(lang string) (string, error) {
	tag, err := language.Parse(lang)
	if nil != err {
		return "", err
	}
	tokens := qb_utils.Strings.Split(tag.String(), "_-")
	if len(tokens) > 0 {
		return tokens[0], nil
	}
	return tag.String(), nil
}

func langOf(rawLang string) string {
	lang, _ := validate(rawLang)
	return lang
}

func isLangFile(tag language.Tag, filename string) bool {
	name := qb_utils.Paths.FileName(filename, false)
	return langOf(name) == langOf(tag.String())
}

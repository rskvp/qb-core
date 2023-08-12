package qb_i18n_bundle

import (
	"github.com/rskvp/qb-core/qb_utils"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type Localizer struct {
	bundle *Bundle
}

func NewLocalizer(bundle *Bundle) *Localizer {
	instance := new(Localizer)
	instance.bundle = bundle
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Localizer) SetRenderer(renderer I18nRenderer) {
	if nil != instance.bundle {
		instance.bundle.renderer = renderer
	}
}

func (instance *Localizer) Renderer() I18nRenderer {
	if nil == instance.bundle {
		return nil
	}
	return instance.bundle.renderer
}

func (instance *Localizer) String() string {
	if nil != instance.bundle {
		return instance.bundle.String()
	}
	return ""
}

func (instance *Localizer) Localize(rawLang string, text string, context ...interface{}) (string, error) {
	// consolidate the context into a single model
	model := instance.toMap(context...)

	// localize the model
	instance.localizeModel(&model, rawLang, instance.bundle)

	// get all template tags
	// all missing tags should be added to model and retrieved from
	// an i18n bundle
	missing, err := instance.getMissingModelKeys(text, model)
	if nil != err {
		return "", err
	}
	for _, key := range missing {
		value, err := instance.bundle.GetWithPlural(rawLang, key, nil, model)
		if nil == err {
			// deep assign
			qb_utils.Maps.Set(model, key, value)
			//model[key] = value
		}
	}

	// ready to Render the template
	return instance.bundle.Renderer().Render(text, model)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Localizer) toMap(context ...interface{}) map[string]interface{} {
	response := make(map[string]interface{})
	for _, c := range context {
		if mm, b := c.(map[string]interface{}); b {
			for k, v := range mm {
				response[k] = v
			}
		}
	}
	return response
}

func (instance *Localizer) getMissingModelKeys(text string, context ...interface{}) ([]string, error) {
	response := make([]string, 0)
	tags, err := instance.bundle.Renderer().TagsFrom(text)
	if nil != err {
		return nil, err
	}
	keys := instance.getKeys(true, context...)
	for _, tag := range tags {
		// check keywords
		if !isGoTemplateKeyword(tag) && qb_utils.Arrays.IndexOf(tag, keys) == -1 {
			response = append(response, tag)
		}
	}
	return response, nil
}

func (instance *Localizer) getKeys(allowDuplicates bool, context ...interface{}) []string {
	response := make([]string, 0)
	for _, c := range context {
		if mm, b := c.(map[string]interface{}); b {
			instance.deepKeys(mm, &response, allowDuplicates)
		}
	}
	return response
}

func (instance *Localizer) deepKeys(m map[string]interface{}, response *[]string, allowDuplicates bool) {
	for k, v := range m {
		if mm, b := v.(map[string]interface{}); b {
			instance.deepKeys(mm, response, allowDuplicates)
		} else if a, b := v.([]map[string]interface{}); b && len(a) > 0 {
			instance.deepKeys(a[0], response, allowDuplicates)
		} else {
			if !allowDuplicates && qb_utils.Arrays.IndexOf(k, *response) == -1 {
				*response = append(*response, k)
			} else {
				*response = append(*response, k)
			}
		}
	}
}

// localizeModel loop on all model data and check for tags to resolve.
// localization data are retrieved from i18n bundle
func (instance *Localizer) localizeModel(pmodel *map[string]interface{}, lang string, bundle *Bundle) {
	model := *pmodel
	for k, v := range model {
		if s, b := v.(string); b {
			value, err := instance.localize(s, lang, bundle, model)
			if nil == err && s != value {
				model[k] = value
			}
		} else if m, b := v.(map[string]interface{}); b {
			instance.localizeModel(&m, lang, bundle)
		} else if a, b := v.([]map[string]interface{}); b {
			// extend model with array counter
			model[k+SUFFIX_LENGHT] = len(a)
			for mi, mm := range a {
				mm[k+SUFFIX_INDEX] = mi
				instance.localizeModel(&mm, lang, bundle)
			}
		} else if a, b := v.([]interface{}); b {
			// extend model with array counter
			model[k+SUFFIX_LENGHT] = len(a)
			for iii, ii := range a {
				if mm, b := ii.(map[string]interface{}); b {
					mm[k+SUFFIX_INDEX] = iii
					instance.localizeModel(&mm, lang, bundle)
				}
			}
		}
	}
}

func (instance *Localizer) localize(text string, lang string, bundle *Bundle, i18nModel map[string]interface{}) (string, error) {
	tags, err := instance.bundle.Renderer().TagsFrom(text)
	if nil != err {
		return "", err
	}
	if len(tags) > 0 {
		model := make(map[string]interface{})
		for _, tag := range tags {
			value, err := bundle.GetWithPlural(lang, tag, nil, i18nModel)
			if nil == err && len(value) > 0 {
				model[tag] = value
			}
		}
		return instance.bundle.Renderer().Render(text, model)
	}
	return text, nil
}
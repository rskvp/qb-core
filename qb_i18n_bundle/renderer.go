package qb_i18n_bundle

import (
	"errors"
	"strings"

	"github.com/rskvp/qb-core/qb_utils"
)

type I18nRenderer interface {
	Render(text string, model ...interface{}) (string, error)
	TagsFrom(text string) ([]string, error)
}

type Renderer struct {
}

func (instance *Renderer) Render(text string, model ...interface{}) (string, error) {
	if isGoTemplate(text) {
		return renderGo(text, model...)
	} else {
		return text, errors.New("bad_template_syntax")
	}
}

func (instance *Renderer) TagsFrom(text string) ([]string, error) {
	if isGoTemplate(text) {
		return tagsFromGo(text)
	} else {
		return nil, errors.New("bad_template_syntax")
	}
}

func isGoTemplate(text string) bool {
	return strings.Index(text, "{{") == -1 || (strings.Index(text, "{{.") > -1 || strings.Index(text, "{{ .") > -1)
}

func isGoTemplateKeyword(tag string)bool{
	if strings.Index(tag, "with ")==0 || tag=="end" || tag=="range" {
		return true
	}
	return false
}

func renderGo(text string, model ...interface{}) (string, error) {
	res, err := qb_utils.Formatter.MergeText(text, qb_utils.Arrays.GetAt(model, 0, nil)) // mustache.Render(text, model...)
	if nil != err {
		return text, nil
	}
	return res, nil
}

func tagsFromGo(text string) ([]string, error) {
	response := make([]string, 0)
	tags := qb_utils.Regex.TextBetweenBraces(text)
	for _, tag := range tags {
		response = append(response, strings.TrimSpace(strings.Replace(tag, ".", "", 1)))
	}
	return response, nil
}


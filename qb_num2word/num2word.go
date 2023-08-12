package qb_num2word

import (
	"strings"

	"github.com/rskvp/qb-core/qb_num2word/qb_num2word_languages"
)

var Num2Word2 *Num2WordHelper

func init() {
	Num2Word2 = newNum2Word()
}

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type Num2WordHelper struct {
	Options *Num2WordOptions
}

type Num2WordOptions struct {
	WordSeparator string
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func newNum2Word() *Num2WordHelper {
	instance := new(Num2WordHelper)
	instance.Options = new(Num2WordOptions)
	instance.Options.WordSeparator = " "

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Num2WordHelper) ConvertDefault(input int) string {
	lang := qb_num2word_languages.Languages.Default()
	return num2Word(input, &lang, instance.Options)
}

func (instance *Num2WordHelper) Convert(input int, langCode string) string {
	return instance.ConvertOpts(input, langCode, instance.Options)
}

func (instance *Num2WordHelper) ConvertOpts(input int, langCode string, opts *Num2WordOptions) string {
	lang := qb_num2word_languages.Languages.Lookup(langCode)
	if nil == lang {
		lang = qb_num2word_languages.Languages.Lookup("en-us")
	}
	if nil == opts {
		opts = instance.Options
	}
	return num2Word(input, lang, opts)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func num2Word(input int, lang *qb_num2word_languages.Language, options *Num2WordOptions) string {
	response := ""
	if len(lang.Name) > 0 && nil != lang.IntegerToWords {
		response = lang.IntegerToWords(input)
	}

	if nil != options {
		response = strings.Replace(response, " ", options.WordSeparator, -1)
	}

	return response
}

package qb_i18n_bundle

import (
	"errors"

	"github.com/rskvp/qb-core/qb_utils"
)

type I18NHelper struct {
}

var I18N *I18NHelper

func init() {
	I18N = new(I18NHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *I18NHelper) NewLocalizer(args ...interface{}) (localizer *Localizer, err error) {
	if len(args) > 0 {
		var bundle *Bundle
		switch len(args) {
		case 1:
			lang := qb_utils.Convert.ToString(args[0])
			bundle, err = NewBundleFromLang(lang)
			if nil != err {
				return
			}
			localizer = NewLocalizer(bundle)
			return
		case 2:
			lang := qb_utils.Convert.ToString(args[0])
			fs := args[1]
			if dir, b := fs.(string); b {
				bundle, err = NewBundleFromDir(lang, dir)
				if nil != err {
					return
				}
				localizer = NewLocalizer(bundle)
				return
			}
			if array, b := fs.([]string); b {
				bundle, err = NewBundleFromFiles(lang, array)
				if nil != err {
					return
				}
				localizer = NewLocalizer(bundle)
				return
			}
		}
		return
	}
	return nil, errors.New("missing parameters")
}

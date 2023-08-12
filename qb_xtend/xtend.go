package qb_xtend

import "github.com/rskvp/qb-core/qb_xtend/pdf"

type XtendHelper struct {
	xpdf pdf.IPdfExtension
}

var Xtend *XtendHelper

func init() {
	Xtend = new(XtendHelper)
}

func (instance *XtendHelper) SetPdf(value pdf.IPdfExtension) *XtendHelper {
	if nil != instance {
		instance.xpdf = value
	}
	return instance
}

func (instance *XtendHelper) Pdf() pdf.IPdfExtension {
	if nil != instance {
		return instance.xpdf
	}
	return nil
}

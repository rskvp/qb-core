package qb_vcal

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/rskvp/qb-core/qb_utils"
)

type VCalHelper struct {
}

var VCal *VCalHelper

func init() {
	VCal = new(VCalHelper)
}

func (instance *VCalHelper) New() *Calendar {
	return NewCalendarFor("gg")
}

func (instance *VCalHelper) NewForService(service string) *Calendar {
	return NewCalendarFor(service)
}

func (instance *VCalHelper) Parse(data interface{}) (*Calendar, error) {
	if r, b := data.(io.Reader); b {
		return ParseCalendar(r)
	}
	if r, b := data.([]byte); b {
		return ParseCalendar(bytes.NewReader(r))
	}
	if s, b := data.(string); b {
		if qb_utils.Paths.IsFilePath(s) {
			r, err := os.Open(s)
			if nil != err {
				return nil, err
			}
			return ParseCalendar(r)
		} else {
			return ParseCalendar(bytes.NewReader([]byte(s)))
		}
	}
	return nil, errors.New("unsupported_data_error")
}

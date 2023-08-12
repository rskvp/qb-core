package qb_utils

import (
	"bytes"
	"encoding/xml"
	"errors"
)

type XMLHelper struct {
}

var XML *XMLHelper

//----------------------------------------------------------------------------------------------------------------------
//	i n i t
//----------------------------------------------------------------------------------------------------------------------

func init() {
	XML = new(XMLHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *XMLHelper) Stringify(entity interface{}) string {
	data, _ := xml.MarshalIndent(entity, "", "  ")
	return string(data)
}

func (instance *XMLHelper) Read(input interface{}, entity interface{}) (err error) {
	var decoder *xml.Decoder
	if s, b := input.(string); b {
		decoder = xml.NewDecoder(bytes.NewReader([]byte(s)))
	} else if s, b := input.([]byte); b {
		decoder = xml.NewDecoder(bytes.NewReader(s))
	}
	if nil != decoder {
		err = decoder.Decode(entity)
	} else {
		err = errors.New("unsupported input type")
	}
	return
}

func (instance *XMLHelper) ReadFromFile(fileName string, entity interface{}) error {
	b, err := IO.ReadBytesFromFile(fileName)
	if nil != err {
		return err
	}
	return instance.Read(b, entity)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

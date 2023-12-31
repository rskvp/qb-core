package qb_utils

import (
	"bytes"
	"io"
	"io/ioutil"
)

type BOMHelper struct {
}

var BOM *BOMHelper

func init() {
	BOM = new(BOMHelper)
}

const (
	bom0 = 0xef
	bom1 = 0xbb
	bom2 = 0xbf
)

// CleanBom returns b with the 3 byte BOM stripped off the front if it is present.
// If the BOM is not present, then b is returned.
func (instance *BOMHelper) CleanBom(b []byte) []byte {
	if len(b) >= 3 &&
		b[0] == bom0 &&
		b[1] == bom1 &&
		b[2] == bom2 {
		return b[3:]
	}
	return b
}

// NewReaderWithoutBom returns an io.Reader that will skip over initial UTF-8 byte order marks.
func (instance *BOMHelper) NewReaderWithoutBom(r io.Reader) (io.Reader, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(instance.CleanBom(bs)), nil
}


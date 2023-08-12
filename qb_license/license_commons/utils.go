package license_commons

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_utils"
)

//----------------------------------------------------------------------------------------------------------------------
//	const
//----------------------------------------------------------------------------------------------------------------------

var KEY = "G&G-license-key-fake______________"

var (
	LicenseNotFoundError                  = errors.New("license_not_found_error")
	LicenseConfigurationFileNotFoundError = errors.New("license_config_not_found_error")
	LicenseExpiredError                   = errors.New("license_expired_error")
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func IsValidPacket(bytes []byte) bool {
	if len(bytes) > 0 {
		text := string(bytes)
		return strings.Index(text, "not found") == -1
	}
	return false
}

func IsJSON(bytes []byte) bool {
	return qb_utils.Regex.IsValidJsonObject(string(bytes))
}

func Encode(data []byte) (response []byte, err error) {
	key := qb_utils.Strings.FillLeftBytes([]byte(KEY), 32, '0')
	response, err = qb_utils.Coding.EncryptBytesAES(data, key)
	return
}

func EncodeText(text string) (response string, err error) {
	data, e := Encode([]byte(text))
	if nil != e {
		err = e
	} else {
		response = string(data)
	}
	return
}

func Decode(data []byte) (response []byte, err error) {
	if !IsJSON(data) {
		key := qb_utils.Strings.FillLeftBytes([]byte(KEY), 32, '0')
		response, err = qb_utils.Coding.DecryptBytesAES(data, key)
	} else {
		response = data
	}
	return
}

func DecodeText(text string) (response string, err error) {
	data, e := Decode([]byte(text))
	if nil != e {
		err = e
	} else {
		response = string(data)
	}
	return
}

func Download(url string) ([]byte, error) {
	if len(url) > 0 {

		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    15 * time.Second,
			DisableCompression: true,
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(url)
		if nil == err {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if nil == err {
				return body, nil
			} else {
				return []byte{}, err
			}
		} else {
			return []byte{}, err
		}
	}
	return []byte{}, errors.New("missing_path: path parameter is empty string")
}

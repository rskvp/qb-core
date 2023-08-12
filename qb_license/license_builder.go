package qb_license

import (
	"fmt"

	"github.com/rskvp/qb-core/qb_license/license_commons"
	"github.com/rskvp/qb-core/qb_utils"
)

type LicenseBuilder struct {
}

func NewLicenseBuilder() (instance *LicenseBuilder) {
	instance = new(LicenseBuilder)
	return
}

func (instance *LicenseBuilder) NewLicense(uid string) (license *license_commons.License) {
	license = license_commons.NewLicense(uid)
	return license
}

func (instance *LicenseBuilder) ToText(uid, name, email, lang string, durationDays int64, params map[string]interface{}) (text string) {
	license := instance.licenseFrom(uid, name, email, lang, durationDays, params)
	return license.String()
}

func (instance *LicenseBuilder) EncodeToText(uid, name, email, lang string, durationDays int64, params map[string]interface{}) (text string, err error) {
	license := instance.licenseFrom(uid, name, email, lang, durationDays, params)
	return instance.EncodeLicenseToText(license)
}

func (instance *LicenseBuilder) EncodeToBytes(uid, name, email, lang string, durationDays int64, params map[string]interface{}) (bytes []byte, err error) {
	license := instance.licenseFrom(uid, name, email, lang, durationDays, params)
	return instance.EncodeLicenseToBytes(license)
}

func (instance *LicenseBuilder) SaveToFile(uid, name, email, lang string, durationDays int64, params map[string]interface{}, filename string) error {
	license := instance.licenseFrom(uid, name, email, lang, durationDays, params)
	return instance.SaveLicenseToFileName(license, filename)
}

func (instance *LicenseBuilder) SaveToTempFile(uid, name, email, lang string, durationDays int64, params map[string]interface{}) (string, error) {
	license := instance.licenseFrom(uid, name, email, lang, durationDays, params)
	return instance.SaveLicenseToTempFile(license)
}

func (instance *LicenseBuilder) EncodeLicenseToText(license *license_commons.License) (text string, err error) {
	if nil != license {
		text, err = license_commons.EncodeText(license.String())
	}
	return
}

func (instance *LicenseBuilder) EncodeLicenseToBytes(license *license_commons.License) (bytes []byte, err error) {
	if nil != license {
		bytes, err = license_commons.Encode([]byte(license.String()))
	}
	return
}

func (instance *LicenseBuilder) SaveLicenseToFileName(license *license_commons.License, filename string) (err error) {
	if nil != license {
		err = license.SaveToFile(filename)
	}
	return
}

func (instance *LicenseBuilder) SaveLicenseToFile(license *license_commons.License, dir string) (filename string, err error) {
	if nil != license {
		dir = qb_utils.Paths.Absolutize(dir, qb_utils.Paths.GetWorkspacePath())
		filename = qb_utils.Paths.Concat(dir, fmt.Sprintf(fmt.Sprintf("%s.lic", license.Uid)))
		err = license.SaveToFile(filename)
	}
	return
}

func (instance *LicenseBuilder) SaveLicenseToTempFile(license *license_commons.License) (filename string, err error) {
	if nil != license {
		filename = qb_utils.Paths.TempPath(fmt.Sprintf("%s.lic", license.Uid))
		err = license.SaveToFile(filename)
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *LicenseBuilder) licenseFrom(uid, name, email, lang string, durationDays int64, params map[string]interface{}) (license *license_commons.License) {
	license = instance.NewLicense(uid)
	license.Name = name
	license.Email = email
	license.Lang = lang
	license.DurationDays = durationDays
	license.Params = params

	return
}

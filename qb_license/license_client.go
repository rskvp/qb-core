package qb_license

import (
	"sync"

	"github.com/rskvp/qb-core/qb_license/license_commons"
	"github.com/rskvp/qb-core/qb_utils"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type LicenseClient struct {
	Config *license_commons.LicenseConfig

	mux sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewLicenseClient(config *license_commons.LicenseConfig) *LicenseClient {
	instance := new(LicenseClient)
	instance.Config = config

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *LicenseClient) GetUrl() string {
	protocol := "http"
	if instance.Config.UseSSL {
		protocol = "https"
	}
	host := instance.Config.Host
	port := instance.Config.Port
	return qb_utils.Strings.Format("%s://%s:%s/", protocol, host, port)
}

func (instance *LicenseClient) RequestLicense(path string) (license *license_commons.License, err error) {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	license = new(license_commons.License)
	bytes := make([]byte, 0)
	if len(path) > 0 {
		bytes, err = license_commons.Download(instance.GetUrl() + path)
	} else if len(instance.Config.Path) > 0 {
		bytes, err = license_commons.Download(instance.GetUrl() + instance.Config.Path)
	}

	if nil == err {
		if license_commons.IsValidPacket(bytes) {
			if string(bytes[0]) != "{" {
				bytes, err = license_commons.Decode(bytes)
			}
			data := string(bytes)
			if len(data) > 0 {
				err = license.Parse(data)
			}
		} else {
			err = license_commons.LicenseNotFoundError
		}
	}

	return license, err
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

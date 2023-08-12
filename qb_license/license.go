package qb_license

import (
	"fmt"
	"time"

	"github.com/rskvp/qb-core/qb_license/license_commons"
	"github.com/rskvp/qb-core/qb_utils"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

var License *LicenseHelper

type LicenseErrorCallback func(license *license_commons.License, err error)

type LicenseHelper struct {
	OnFail LicenseErrorCallback

	root          string
	configFile    string
	lastFail      time.Time
	failTolerance time.Duration
	listeners     []LicenseErrorCallback

	_config  *license_commons.LicenseConfig
	_client  *LicenseClient
	_ticker  *LicenseTicker
	_builder *LicenseBuilder
}

func init() {
	License = new(LicenseHelper)
	License.failTolerance = 1 * 24 * time.Hour
	License.lastFail = time.Time{}
	License.listeners = make([]LicenseErrorCallback, 0)
}

//----------------------------------------------------------------------------------------------------------------------
//	s e t t i n g s
//----------------------------------------------------------------------------------------------------------------------

func (instance *LicenseHelper) SetSecret(key string) {
	license_commons.KEY = key
}

func (instance *LicenseHelper) SetTolerance(duration time.Duration) {
	instance.failTolerance = duration
}

func (instance *LicenseHelper) SetRoot(root string) {
	instance.root = qb_utils.Paths.Absolutize(root, qb_utils.Paths.GetWorkspacePath())
	instance.configFile = qb_utils.Paths.Concat(instance.root, "license.config")
}

func (instance *LicenseHelper) SetConfigFile(filename string) {
	instance.configFile = qb_utils.Paths.Absolutize(filename, qb_utils.Paths.GetWorkspacePath())
	if len(instance.root) == 0 {
		instance.root = qb_utils.Paths.Dir(instance.configFile)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	b u i l d e r
//----------------------------------------------------------------------------------------------------------------------

func (instance *LicenseHelper) Builder() (*LicenseBuilder, error) {
	if nil != instance {
		return instance.builder()
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	e v e n t s
//----------------------------------------------------------------------------------------------------------------------

func (instance *LicenseHelper) AddErrorListener(callback LicenseErrorCallback) *LicenseHelper {
	if nil != instance {
		instance.listeners = append(instance.listeners, callback)
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *LicenseHelper) Check() (err error) {
	if nil != instance {
		var cli *LicenseClient
		cli, err = instance.client()
		if nil != err {
			return
		}

		var license *license_commons.License
		license, err = cli.RequestLicense("")
		if nil != err {
			return
		}

		if !license.IsValid() {
			err = expiredError(license)
		}
	}
	return
}

func (instance *LicenseHelper) Start(interval time.Duration) (err error) {
	if nil != instance {
		var ticker *LicenseTicker
		ticker, err = instance.ticker()
		if nil != err {
			return
		}

		if ticker.IsRunning() {
			ticker.Stop()
		}
		ticker.Interval = interval
		ticker.Start()
	}
	return
}

func (instance *LicenseHelper) StartWithCallback(interval time.Duration, onFail LicenseErrorCallback) (err error) {
	if nil != instance {
		instance.OnFail = onFail
		err = instance.Start(interval)
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *LicenseHelper) config() (*license_commons.LicenseConfig, error) {
	if nil != instance {
		if nil == instance._config {
			if len(instance.configFile) == 0 {
				instance.SetRoot("")
			}

			if ok, _ := qb_utils.Paths.Exists(instance.configFile); ok {
				text, err := qb_utils.IO.ReadTextFromFile(instance.configFile)
				if nil != err {
					return nil, err
				}
				config := new(license_commons.LicenseConfig)
				err = config.Parse(text)
				if nil != err {
					return nil, err
				}
				instance._config = config
			} else {
				// missing config file
				return nil,
					qb_utils.Errors.Prefix(license_commons.LicenseConfigurationFileNotFoundError, fmt.Sprintf("Missing config file '%s': ", instance.configFile))
			}
		}
	}
	return instance._config, nil
}

func (instance *LicenseHelper) client() (*LicenseClient, error) {
	if nil != instance {
		if nil == instance._client {
			config, err := instance.config()
			if nil != err {
				return nil, err
			}
			instance._client = NewLicenseClient(config)
		}
	}
	return instance._client, nil
}

func (instance *LicenseHelper) ticker() (*LicenseTicker, error) {
	if nil != instance {
		if nil == instance._ticker {
			config, err := instance.config()
			if nil != err {
				return nil, err
			}
			instance._ticker = NewLicenseTicker(config)
			instance._ticker._client, err = instance.client()
			if nil == err {
				instance._ticker.RequestLicenseHook = instance.onTick
				instance._ticker.ErrorLicenseHook = instance.onTick
				instance._ticker.ExpiredLicenseHook = instance.onTick
			}
		}
	}
	return instance._ticker, nil
}

func (instance *LicenseHelper) builder() (*LicenseBuilder, error) {
	if nil != instance {
		if nil == instance._builder {
			instance._builder = NewLicenseBuilder()
		}
	}
	return instance._builder, nil
}

func (instance *LicenseHelper) onTick(ctx *LicenseTickerContext) {
	if nil != instance && nil != instance._ticker && instance._ticker.IsRunning() {
		if nil != ctx.Error {
			if instance.lastFail.IsZero() {
				instance.lastFail = time.Now()
			}
		} else {
			if nil != ctx.License && !ctx.License.IsValid() {
				ctx.Error = expiredError(ctx.License)
				if instance.lastFail.IsZero() {
					instance.lastFail = time.Now()
				}
			} else {
				// reset last fail
				instance.lastFail = time.Time{}
			}

		}

		if !instance.lastFail.IsZero() {
			// we have an error
			limit := qb_utils.Dates.Add(instance.lastFail, instance.failTolerance)
			if time.Now().Unix() > limit.Unix() {
				// reset last fail
				instance.lastFail = time.Time{}
				// fail event
				instance.triggerError(ctx.License, ctx.Error)
			}
		}
	}
}

func (instance *LicenseHelper) triggerError(license *license_commons.License, err error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	if nil != instance.listeners {
		for _, callback := range instance.listeners {
			if nil != callback {
				callback(license, err)
			}
		}
	}

	if nil != instance.OnFail {
		instance.OnFail(license, err)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func expiredError(license *license_commons.License) error {
	return qb_utils.Errors.Prefix(license_commons.LicenseExpiredError,
		fmt.Sprintf("License expired on '%v': ", license.GetExpireDate()))
}

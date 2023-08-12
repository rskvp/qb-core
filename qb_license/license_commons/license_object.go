package license_commons

import (
	"encoding/json"
	"math"
	"time"

	"github.com/rskvp/qb-core/qb_rnd"
	"github.com/rskvp/qb-core/qb_utils"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type License struct {
	Uid          string                 `json:"uid"`
	CreationTime time.Time              `json:"creation_time"`
	DurationDays int64                  `json:"duration_days"`
	Name         string                 `json:"name"`
	Email        string                 `json:"email"`
	Lang         string                 `json:"lang"`
	Enabled      bool                   `json:"enabled"`
	Params       map[string]interface{} `json:"params"`
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewLicense(uid string) *License {
	instance := new(License)
	instance.Uid = uid

	instance.init()
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	License
//----------------------------------------------------------------------------------------------------------------------

func (instance *License) Parse(text string) error {
	return instance.ParseBytes([]byte(text))
}

func (instance *License) ParseBytes(bytes []byte) error {
	return json.Unmarshal(bytes, &instance)
}

func (instance *License) String() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

func (instance *License) GetDataAsString() string {
	return qb_utils.Strings.Format("\tId: %s\n\tOwner: %s\n\tCreation Date: %s\n\tExpire Date: %s\n\tExpired from days: %s",
		instance.Uid,
		instance.Name,
		instance.CreationTime,
		instance.GetExpireDate(),
		instance.RemainingDays()*-1,
	)
}

func (instance *License) IsValid() bool {
	if instance.Enabled && len(instance.Uid) > 0 {
		created := instance.CreationTime
		duration := instance.DurationDays
		now := time.Now()
		days := int64(now.Sub(created).Hours() / 24)

		return days <= duration
	}
	return false
}

func (instance *License) RemainingDays() int64 {
	if instance.Enabled && len(instance.Uid) > 0 {
		created := instance.CreationTime
		duration := instance.DurationDays
		now := time.Now()
		days := int64(now.Sub(created).Hours() / 24)

		return duration - days
	}
	return 0
}

func (instance *License) GetExpireDate() time.Time {
	created := instance.CreationTime
	duration := instance.DurationDays
	durationHours := duration * 24
	return created.Add(time.Duration(durationHours) * time.Hour)
}

func (instance *License) SetExpireDate(date time.Time) {
	created := instance.CreationTime
	days := int64(math.Round(date.Sub(created).Hours() / 24))
	instance.DurationDays = days
}

func (instance *License) ParseExpireDate(layout string, value string) error {
	date, err := time.Parse(layout, value)
	if nil == err {
		instance.SetExpireDate(date)
	}
	return err
}

func (instance *License) Add(days int64) {
	instance.Enabled = true
	instance.DurationDays = instance.DurationDays + days
}

func (instance *License) Encode() (text string, err error) {
	text, err = EncodeText(instance.String())
	return
}

func (instance *License) SaveToFile(filename string) (err error) {
	encoded, e := EncodeText(instance.String())
	if nil != e {
		err = e
	} else {
		_, err = qb_utils.IO.WriteTextToFile(encoded, filename)
	}
	return
}

func (instance *License) ReadFromFile(filename string) (err error) {
	encoded, e := qb_utils.IO.ReadBytesFromFile(filename)
	if nil != e {
		err = e
	} else {
		decoded, e := Decode(encoded)
		if nil != e {
			err = e
		} else {
			err = instance.ParseBytes(decoded)
		}
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	private
//----------------------------------------------------------------------------------------------------------------------

func (instance *License) init() {
	if nil != instance {
		instance.Enabled = true
		if nil == instance.Params {
			instance.Params = make(map[string]interface{})
		}
		if len(instance.Uid) == 0 {
			instance.Uid = qb_rnd.Rnd.Uuid()
		}
		if instance.CreationTime.IsZero() {
			instance.CreationTime = time.Now()
		}
		if instance.DurationDays == 0 {
			instance.DurationDays = 1
		}
		if len(instance.Lang) == 0 {
			instance.Lang = "en"
		}
		if len(instance.Name) == 0 {
			instance.Name = "Anonymous"
		}
	}
}

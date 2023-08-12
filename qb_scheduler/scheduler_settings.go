package qb_scheduler

import "github.com/rskvp/qb-core/qb_utils"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type SchedulerSettings struct {
	Uid       string      `json:"uid"`
	Sync      bool        `json:"sync"`
	Schedules []*Schedule `json:"schedules"`
}

func (instance *SchedulerSettings) String() string {
	return qb_utils.JSON.Stringify(instance)
}

type Schedule struct {
	Uid       string                 `json:"uid,omitempty"`
	StartAt   string                 `json:"start_at"` // hh:mm ss (optional)
	Timeline  string                 `json:"timeline"` // minutes:1, hour:24, second:10
	Payload   map[string]interface{} `json:"payload,omitempty"`
	Arguments []interface{}          `json:"-"` // custom attachments
}

func (instance *Schedule) String() string {
	return qb_utils.JSON.Stringify(instance)
}

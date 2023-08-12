package qb_scheduler

import (
	"fmt"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_utils"
)

const (
	defaultTimeline = "hour:12"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type SchedulerTask struct {
	Uid       string
	Arguments []interface{}
	Payload   map[string]interface{}

	schedulerUid string
	startAt      time.Time // fixed setting start value
	nextStartAt  time.Time // next tick
	timeline     time.Duration
	err          error
	settings     *Schedule
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewSchedulerTask(schedulerUid string, settings *Schedule) *SchedulerTask {
	instance := new(SchedulerTask)
	instance.schedulerUid = schedulerUid
	instance.settings = settings
	instance.init()

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *SchedulerTask) String() string {
	if nil != instance {
		return instance.GoString()
	}
	return ""
}

func (instance *SchedulerTask) GoString() string {
	if nil != instance {
		data := map[string]interface{}{
			"scheduler": instance.schedulerUid,
			"uid":       instance.settings.Uid,
			"error":     instance.Error(),
			"start_at":  instance.startAt,
			"timeline":  instance.settings.Timeline,
		}
		return qb_utils.JSON.Stringify(data)
	}
	return ""
}

func (instance *SchedulerTask) Settings() *Schedule {
	return instance.settings
}

func (instance *SchedulerTask) Error() string {
	if nil != instance.err {
		return instance.err.Error()
	}
	return ""
}

func (instance *SchedulerTask) IsReady() bool {
	now := time.Now()
	diff := instance.nextStartAt.Sub(now)
	if diff <= 0 {
		// move to next tick
		instance.nextStartAt = instance.nextStartAt.Add(instance.timeline)

		return true
	}
	return false
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *SchedulerTask) init() {
	now := time.Now()
	settings := instance.settings

	instance.Uid = settings.Uid
	instance.Arguments = append(instance.Arguments, settings.Arguments...)
	instance.Payload = settings.Payload

	// START-AT
	if len(settings.StartAt) > 0 {
		t, err := qb_utils.Formatter.ParseDate(settings.StartAt, "HH:mm:ss")
		if nil == err {
			year := now.Year()
			month := qb_utils.Strings.FillLeft(qb_utils.Convert.ToString(int(now.Month())), 2, '0')
			day := qb_utils.Strings.FillLeft(qb_utils.Convert.ToString(now.Day()), 2, '0')
			hour := qb_utils.Strings.FillLeft(qb_utils.Convert.ToString(t.Hour()), 2, '0')
			min := qb_utils.Strings.FillLeft(qb_utils.Convert.ToString(t.Minute()), 2, '0')
			sec := qb_utils.Strings.FillLeft(qb_utils.Convert.ToString(t.Second()), 2, '0')
			_, zOffset := now.Zone()
			z := qb_utils.Strings.FillLeft(qb_utils.Convert.ToString(zOffset/36), 4, '0')
			if zOffset > 0 {
				z = "+" + z
			}
			sdate := fmt.Sprintf("%v-%v-%v %v:%v:%v %v", year, month, day, hour, min, sec, z)
			d, err := qb_utils.Formatter.ParseDate(sdate, "yyyy-MM-dd HH:mm:ss Z")
			if nil == err {
				instance.startAt = d
			} else {
				instance.err = err
				instance.startAt = now
			}
		} else {
			instance.err = err
			instance.startAt = now
		}
	} else {
		instance.startAt = now
	}

	// TIMELINE
	tl := strings.Split(settings.Timeline, ":") // hour:12
	if len(tl) == 2 {
		value := qb_utils.Convert.ToInt(tl[1])
		if value == 0 {
			value = 1
		}
		switch tl[0] {
		case "millisecond":
			instance.timeline = time.Duration(value) * time.Millisecond
		case "second":
			instance.timeline = time.Duration(value) * time.Second
		case "minute":
			instance.timeline = time.Duration(value) * time.Minute
		case "hour":
			instance.timeline = time.Duration(value) * time.Hour
		default:
			instance.timeline = 12 * time.Hour
		}
	} else {
		// invalid timeline add defaults
		settings.Timeline = defaultTimeline
		instance.timeline = 12 * time.Hour
	}

	// NEXT
	instance.nextStartAt = instance.startAt
}

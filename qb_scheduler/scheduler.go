package qb_scheduler

import (
	"fmt"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_events"
	"github.com/rskvp/qb-core/qb_utils"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const (
	onSchedule = "on_schedule"
	onError    = "on_error"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type SchedulerTaskHandler func(schedule *SchedulerTask)
type SchedulerErrorHandler func(err string)

type Scheduler struct {
	settings       *SchedulerSettings
	internalEvents *qb_events.Emitter
	taskHandlers   []SchedulerTaskHandler
	errorHandlers  []SchedulerErrorHandler
	timer          *time.Ticker
	stopChan       chan bool
	closed         bool
	paused         bool
	tasks          []*SchedulerTask
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewScheduler() *Scheduler {
	instance := new(Scheduler)
	instance.settings = new(SchedulerSettings)
	instance.internalEvents = qb_events.Events.NewEmitter()
	instance.taskHandlers = make([]SchedulerTaskHandler, 0)
	instance.errorHandlers = make([]SchedulerErrorHandler, 0)
	instance.stopChan = make(chan bool, 1)
	instance.closed = true
	instance.paused = false
	instance.tasks = make([]*SchedulerTask, 0)
	instance.initEvents()

	return instance
}

func NewSchedulerWithSettings(settings *SchedulerSettings) *Scheduler {
	instance := NewScheduler()
	instance.settings = settings

	return instance
}

func NewSchedulerFromFile(configFile string) *Scheduler {
	instance := NewScheduler()
	instance.load(configFile)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Scheduler) Uid() string {
	if nil != instance && nil != instance.settings {
		return instance.settings.Uid
	}
	return ""
}

func (instance *Scheduler) String() string {
	if nil != instance {
		return instance.GoString()
	}
	return ""
}

func (instance *Scheduler) GoString() string {
	if nil != instance && nil != instance.settings {
		return qb_utils.JSON.Stringify(instance.settings)
	}
	return ""
}

func (instance *Scheduler) IsAsync() bool {
	if nil != instance && nil != instance.settings {
		return !instance.settings.Sync
	}
	return false
}

func (instance *Scheduler) IsStarted() bool {
	if nil != instance && nil != instance.settings {
		return nil != instance.timer
	}
	return false
}

func (instance *Scheduler) CountHandlers() int {
	if nil != instance && nil != instance.settings && nil != instance.taskHandlers {
		return len(instance.taskHandlers)
	}
	return 0
}

func (instance *Scheduler) SetAsync(value bool) {
	if nil != instance && nil != instance.settings {
		instance.settings.Sync = !value
	}
}

func (instance *Scheduler) HasErrors() bool {
	if nil != instance {
		if len(instance.tasks) > 0 {
			for _, t := range instance.tasks {
				if nil != t.err {
					return true
				}
			}
		}
	}
	return false
}

func (instance *Scheduler) GetErrors() string {
	if nil != instance {
		builder := new(strings.Builder)
		if len(instance.tasks) > 0 {
			for _, t := range instance.tasks {
				if nil != t.err {
					builder.WriteString(t.err.Error() + "\n")
				}
			}
		}
		return builder.String()
	}
	return ""
}

func (instance *Scheduler) GetTimeout() time.Duration {
	if nil != instance {
		return instance.calculateTimeout()
	}
	return 0 * time.Second
}

func (instance *Scheduler) AddSchedule(item *Schedule, args ...interface{}) {
	if nil != instance {
		if len(args) > 0 {
			item.Arguments = append(item.Arguments, args...)
		}
		instance.settings.Schedules = append(instance.settings.Schedules, item)
	}
}

func (instance *Scheduler) AddScheduleByJson(json string, args ...interface{}) {
	if nil != instance {
		if qb_utils.Regex.IsValidJsonObject(json) {
			var item Schedule
			err := qb_utils.JSON.Read(json, &item)
			if nil == err {
				instance.AddSchedule(&item, args...)
			}
		}
	}
}

func (instance *Scheduler) OnSchedule(handler SchedulerTaskHandler) {
	if nil != instance && nil != handler {
		instance.taskHandlers = append(instance.taskHandlers, handler)
	}
}

func (instance *Scheduler) OnError(handler SchedulerErrorHandler) {
	if nil != instance && nil != handler {
		instance.errorHandlers = append(instance.errorHandlers, handler)
	}
}

func (instance *Scheduler) Start() {
	if nil != instance && instance.closed {
		instance.closed = false
		if nil != instance.timer {
			instance.timer.Stop()
			instance.timer = nil
		}
		instance.initTasks()
		go instance.run()
	}
}

func (instance *Scheduler) Stop() {
	if nil != instance {
		if nil != instance.timer {
			instance.timer.Stop()
		} else {
			instance.timer = nil
		}
		instance.closed = true
		instance.stopChan <- true
		instance.stopChan = make(chan bool, 1) // reset channel
		instance.tasks = make([]*SchedulerTask, 0)
	}
}

func (instance *Scheduler) Pause() {
	if nil != instance {
		instance.paused = true
	}
}

func (instance *Scheduler) Resume() {
	if nil != instance {
		instance.paused = false
	}
}

func (instance *Scheduler) IsPaused() bool {
	if nil != instance {
		return instance.paused
	}
	return false
}

func (instance *Scheduler) TogglePause() {
	if nil != instance {
		instance.paused = !instance.paused
	}
}

func (instance *Scheduler) Reload() {
	if nil != instance {
		instance.Stop()
		instance.Start()
	}
}

func (instance *Scheduler) Join() {
	if nil != instance && nil != instance.stopChan && !instance.closed {
		<-instance.stopChan
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Scheduler) load(filename string) {
	if len(filename) > 0 {
		txt, err := qb_utils.IO.ReadTextFromFile(filename)
		if nil == err && len(txt) > 0 && qb_utils.Regex.IsValidJsonObject(txt) {
			_ = qb_utils.JSON.Read(txt, instance.settings)
		}
	}
}

func (instance *Scheduler) initEvents() {
	instance.internalEvents.On(onSchedule, func(event *qb_events.Event) {
		if nil != instance && len(instance.taskHandlers) > 0 && !instance.closed && !instance.paused {
			item := event.Argument(0)
			if v, b := item.(*SchedulerTask); b {
				for _, handler := range instance.taskHandlers {
					if nil != handler && !instance.closed && !instance.paused {
						// external handler
						if instance.IsAsync() {
							go handler(v)
						} else {
							handler(v)
						}
					}
				}
			}
		}
	})
	instance.internalEvents.On(onError, func(event *qb_events.Event) {
		if nil != instance && len(instance.errorHandlers) > 0 && !instance.closed {
			var err string
			item := event.Argument(0)
			if v, b := item.(string); b {
				err = v
			} else if v, b := item.(error); b {
				err = v.Error()
			}
			if len(err) > 0 {
				for _, handler := range instance.errorHandlers {
					if nil != handler && !instance.closed {
						// handler(err)
						// external handler
						if instance.IsAsync() {
							go handler(err)
						} else {
							handler(err)
						}
					}
				}
			}
		}
	})
}

func (instance *Scheduler) initTasks() {
	// read configuration and creates task array
	schedules := instance.settings.Schedules
	for _, schedule := range schedules {
		task := NewSchedulerTask(instance.settings.Uid, schedule)
		instance.tasks = append(instance.tasks, task)
		if e := task.Error(); len(e) > 0 {
			instance.internalEmit(onError, e)
		}
	}
}

func (instance *Scheduler) newTicker(fixed bool) *time.Ticker {
	if nil != instance {
		if fixed {
			return time.NewTicker(1 * time.Second)
		}
		return time.NewTicker(instance.calculateTimeout())
	}
	return nil
}

func (instance *Scheduler) calculateTimeout() time.Duration {
	if nil != instance {
		min := 1 * time.Minute
		for _, t := range instance.tasks {
			if t.timeline < min {
				min = t.timeline
			}
		}

		return min
	}
	return 1 * time.Second
}

func (instance *Scheduler) run() {
	instance.timer = instance.newTicker(true)
	if nil == instance.timer {
		return
	}
	for {
		select {
		case <-instance.stopChan:
			return // exit
		case <-instance.timer.C:
			// timer tick
			if nil != instance.timer {
				instance.timer.Stop()
				instance.timer = nil
				instance.checkSchedule()
				instance.timer = instance.newTicker(false)
				if nil == instance.timer {
					return
				}
			}
		}
	}
}

func (instance *Scheduler) checkSchedule() {
	if nil != instance && !instance.closed {
		// PANIC RECOVERY
		defer func() {
			if r := recover(); r != nil {
				// recovered from panic
				message := qb_utils.Strings.Format("[panic] Scheduler '%s' ERROR: %s", instance.Uid(), r)
				if nil != instance {
					instance.internalEmit(onError, message)
				} else {
					fmt.Println(message)
				}
			}
		}()

		for _, task := range instance.tasks {
			if nil != task && task.IsReady() {
				// trigger: is sync because of internal use only
				instance.internalEmit(onSchedule, task)
			}
		}
	}
}

func (instance *Scheduler) internalEmit(event string, args ...interface{}) {
	if nil != instance && nil != instance.internalEvents {
		instance.internalEvents.Emit(event, args...)
	}
}

package qb_events

import "time"

//----------------------------------------------------------------------------------------------------------------------
//	Debouncer
//----------------------------------------------------------------------------------------------------------------------

type Debouncer struct {
	key       string
	eventName string
	args      []interface{}
	async     bool
	emitter   *Emitter
	timer     *time.Timer
}

func (instance *Debouncer) Init(emitter *Emitter, wait time.Duration, eventName string, async bool, args ...interface{}) {
	if nil != instance {
		instance.key = buildKey(eventName, async)
		instance.eventName = eventName
		instance.async = async
		instance.args = args
		instance.emitter = emitter
		if nil != instance.timer {
			instance.timer.Stop()
			instance.timer = nil
		}
		instance.timer = time.NewTimer(wait)
		go instance.start()
	}
}

func (instance *Debouncer) start() {
	if nil != instance && nil != instance.timer {
		<-instance.timer.C
		instance.timer.Stop()
		instance.timer = nil
		if nil != instance.emitter {
			instance.emitter.removeDebounce(instance)
			instance.emitter.emit(instance.eventName, instance.async, instance.args...)
		}
	}
}

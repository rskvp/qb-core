package qb_events

import (
	"time"
)

type EventsHelper struct {
}

var Events *EventsHelper

func init() {
	Events = new(EventsHelper)
}

func (instance *EventsHelper) NewEmitter(payload ...interface{}) *Emitter {
	emitter := NewEmitterInstance(0, payload...)
	return emitter
}

func (instance *EventsHelper) NewDebounceEmitter(waitTime time.Duration, payload ...interface{}) *Emitter {
	emitter := NewEmitterInstance(waitTime, payload...)
	return emitter
}

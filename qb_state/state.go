package qb_state

import (
	"sync"

	"github.com/rskvp/qb-core/qb_events"
	"github.com/rskvp/qb-core/qb_utils"
)

const (
	EventOnChangeState = "_on_change_state"
)

type StateHelper struct {
}

var StateH *StateHelper

func init() {
	StateH = new(StateHelper)
}

func (instance *StateHelper) New(name string) *State {
	return NewState(name)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// State simple state object that store data and emit events
type State struct {
	name   string
	data   map[string]interface{}
	events *qb_events.Emitter
	mux    sync.Mutex
}

func NewState(name string) (instance *State) {
	instance = new(State)
	instance.name = name
	instance.events = qb_events.Events.NewEmitter(name)
	instance.data = make(map[string]interface{})

	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *State) Name() string {
	if nil != instance {
		return instance.name
	}
	return ""
}

func (instance *State) Events() *qb_events.Emitter {
	if nil != instance {
		return instance.events
	}
	return nil
}

func (instance *State) OnStateChanged(callback func(event *qb_events.Event)) *State {
	if nil != instance {
		instance.events.On(EventOnChangeState, callback)
	}
	return instance
}

func (instance *State) OffStateChanged(callback ...func(event *qb_events.Event)) *State {
	if nil != instance {
		instance.events.Off(EventOnChangeState, callback...)
	}
	return instance
}

func (instance *State) On(eventName string, callback func(event *qb_events.Event)) *State {
	if nil != instance {
		instance.events.On(eventName, callback)
	}
	return instance
}

func (instance *State) Off(eventName string, callback ...func(event *qb_events.Event)) *State {
	if nil != instance {
		instance.events.Off(eventName, callback...)
	}
	return instance
}

func (instance *State) SetState(m map[string]interface{}) *State {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		fields := qb_utils.Maps.MergeFields(true, instance.data, m)
		if len(fields) > 0 {
			instance.events.EmitAsync(EventOnChangeState, fields)
		}
	}
	return instance
}

func (instance *State) Put(key string, value interface{}) *State {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		fields := qb_utils.Maps.MergeFields(true, instance.data, map[string]interface{}{key: value})
		if len(fields) > 0 {
			instance.events.EmitAsync(EventOnChangeState, fields)
		}
	}
	return instance
}

func (instance *State) GetState() map[string]interface{} {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		return qb_utils.Maps.Merge(true, map[string]interface{}{}, instance.data)
	}
	return map[string]interface{}{}
}

func (instance *State) Get(key string) interface{} {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		return qb_utils.Maps.Get(map[string]interface{}{}, key)
	}
	return map[string]interface{}{}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

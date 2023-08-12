package qb_events

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/rskvp/qb-core/qb_utils"
)

type EventCallback func(event *Event)

type stackItem struct {
	event    *Event
	callback EventCallback
}

//----------------------------------------------------------------------------------------------------------------------
//	Emitter
//----------------------------------------------------------------------------------------------------------------------

type Emitter struct {
	waitTime  time.Duration
	debounces map[string]*Debouncer
	listeners map[string][]EventCallback
	mux       sync.Mutex
	payload   interface{}
}

func NewEmitterInstance(waitTime time.Duration, payload ...interface{}) (instance *Emitter) {
	instance = new(Emitter)
	instance.waitTime = waitTime
	instance.listeners = make(map[string][]EventCallback)
	if waitTime > 0 {
		instance.debounces = make(map[string]*Debouncer)
	}
	if len(payload) == 1 {
		instance.payload = payload[0]
	} else {
		instance.payload = payload
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Debounce transform a standard event emitterinto a debounced one
func (instance *Emitter) Debounce(waitTime time.Duration) *Emitter {
	if nil != instance {
		instance.waitTime = waitTime
		if nil == instance.debounces {
			instance.debounces = make(map[string]*Debouncer)
		}
	}
	return instance
}

func (instance *Emitter) Has(eventName string) bool {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		if _, b := instance.listeners[eventName]; b {
			return len(instance.listeners[eventName]) > 0
		}
	}
	return false
}

func (instance *Emitter) On(eventName string, callback func(event *Event)) *Emitter {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		instance.listeners[eventName] = append(instance.listeners[eventName], callback)
	}
	return instance
}

func (instance *Emitter) Off(eventName string, callback ...func(event *Event)) *Emitter {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		if _, ok := instance.listeners[eventName]; ok {
			if len(callback) == 0 {
				instance.listeners[eventName] = make([]EventCallback, 0)
			} else {
				handlers := instance.listeners[eventName]
				// loop starting from end
				for i := len(handlers) - 1; i > -1; i-- {
					f := handlers[i]
					for _, h := range callback {
						v1 := reflect.ValueOf(f)
						v2 := reflect.ValueOf(h)
						if v1 == v2 {
							handlers = removeIndex(handlers, i)
							break
						}
					}

				}
				instance.listeners[eventName] = handlers
			}
		}
	}
	return instance
}

func (instance *Emitter) Clear() {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		// reset listeners
		instance.listeners = make(map[string][]EventCallback, 0)
		// remove debouncers
		for _, d := range instance.debounces {
			if nil != d.timer {
				d.timer.Stop()
			}
		}
		instance.debounces = make(map[string]*Debouncer)
	}
}

func (instance *Emitter) Emit(eventName string, args ...interface{}) *Emitter {
	if nil != instance {
		if nil != instance.debounces {
			instance.debounce(eventName, false, args...)
		} else {
			instance.emit(eventName, false, args...)
		}
	}
	return instance
}

func (instance *Emitter) EmitAsync(eventName string, args ...interface{}) *Emitter {
	if nil != instance {
		if nil != instance.debounces {
			instance.debounce(eventName, true, args...)
		} else {
			instance.emit(eventName, true, args...)
		}
	}
	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Emitter) removeDebounce(item *Debouncer) {
	if nil != instance && nil != instance.debounces && nil != item {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		delete(instance.debounces, item.key)
	}
}

func (instance *Emitter) debounce(eventName string, async bool, args ...interface{}) *Emitter {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		key := buildKey(eventName, async)

		if _, ok := instance.debounces[key]; !ok {
			instance.debounces[key] = new(Debouncer)
		}
		instance.debounces[key].Init(instance, instance.waitTime, eventName, async, args...)
	}
	return instance
}

func (instance *Emitter) emit(eventName string, async bool, args ...interface{}) {
	if nil != instance {
		defer func() {
			if r := recover(); r != nil {
				// recovered from panic
				message := qb_utils.Strings.Format("Emit '%s' ERROR: %s", eventName, r)
				fmt.Println(message)
			}
		}()

		instance.mux.Lock()
		defer instance.mux.Unlock()

		// creates internal execution stack
		stack := make([]*stackItem, 0)
		for k, handlers := range instance.listeners {
			if k == eventName {
				for _, handler := range handlers {
					if nil != handler {
						event := NewEvent(async, eventName, instance.payload, args...)
						item := &stackItem{
							event:    event,
							callback: handler,
						}
						stack = append(stack, item)
					}
				}
			}
		}

		go rawEmit(stack)
	}
}

func removeIndex(a []EventCallback, index int) []EventCallback {
	return append(a[:index], a[index+1:]...)
}

func buildKey(eventName string, async bool) string {
	return fmt.Sprintf("%s-%v", eventName, async)
}

func rawEmit(stack []*stackItem) {
	if nil != stack {
		for _, item := range stack {
			if nil != item && nil != item.event && nil != item.callback {
				if item.event.Async {
					go item.callback(item.event)
				} else {
					item.callback(item.event)
				}
			}
		}
	}
}

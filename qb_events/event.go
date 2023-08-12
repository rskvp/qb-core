package qb_events

import "github.com/rskvp/qb-core/qb_utils"

type Event struct {
	Name      string
	Arguments []interface{}
	Payload   interface{}
	Async     bool
}

func NewEvent(async bool, eventName string, payload interface{}, args ...interface{}) (event *Event) {
	event = new(Event)
	event.Name = eventName
	event.Async = async
	event.Payload = payload
	event.Arguments = make([]interface{}, 0)
	event.Arguments = append(event.Arguments, args...)
	return
}

func (instance *Event) ArgumentsInterface() interface{} {
	return interface{}(instance.Arguments)
}

func (instance *Event) Argument(index int) interface{} {
	if len(instance.Arguments) > index {
		return instance.Arguments[index]
	}
	return nil
}

func (instance *Event) ArgumentAsError(index int) error {
	v := instance.Argument(index)
	if nil != v {
		if e, b := v.(error); b {
			return e
		}
	}
	return nil
}

func (instance *Event) ArgumentAsString(index int) string {
	v := instance.Argument(index)
	return qb_utils.Convert.ToString(v)
}

func (instance *Event) ArgumentAsInt(index int) int {
	v := instance.Argument(index)
	if nil != v {
		return qb_utils.Convert.ToInt(v)
	}
	return -1
}

func (instance *Event) ArgumentAsBytes(index int) []byte {
	v := instance.Argument(index)
	if nil != v {
		if e, b := v.([]byte); b {
			return e
		}
	}
	return nil
}

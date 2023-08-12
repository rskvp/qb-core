package qb_watchfiles

import "github.com/rskvp/qb-core/qb_events"

type WatchFilesHelper struct {
}

var WatchFiles *WatchFilesHelper

func init() {
	WatchFiles = new(WatchFilesHelper)
}

func (instance *WatchFilesHelper) GetEventParam(event *qb_events.Event) (response *MultipleFileEventParam) {
	if nil != instance && nil != event {
		arg := event.Argument(0)
		if nil != arg {
			if v, ok := arg.(*MultipleFileEventParam); ok {
				response = v
			}
		}
	}
	return
}

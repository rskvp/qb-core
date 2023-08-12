package qb_vcal

import (
	"fmt"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_utils"
)

//----------------------------------------------------------------------------------------------------------------------
//	constants
//----------------------------------------------------------------------------------------------------------------------

type PropFrequency string

const (
	PropFreqSecondly PropFrequency = "SECONDLY"
	PropFreqMinutely PropFrequency = "MINUTELY"
	PropFreqHourly   PropFrequency = "HOURLY"
	PropFreqDaily    PropFrequency = "DAILY"
	PropFreqWeekly   PropFrequency = "WEEKLY"
	PropFreqMonthly  PropFrequency = "MONTHLY"
	PropFreqYearly   PropFrequency = "YEARLY"
)

var (
	PropFrequencyAll = []PropFrequency{PropFreqSecondly, PropFreqMinutely, PropFreqHourly, PropFreqDaily, PropFreqWeekly, PropFreqMonthly, PropFreqYearly}
)

func (instance PropFrequency) String() string {
	return string(instance)
}

func ParsePropFrequency(value string) PropFrequency {
	for _, prop := range PropFrequencyAll {
		if prop.String() == value {
			return prop
		}
	}
	return PropFreqYearly
}

type PropByDay string

const (
	PropByDaySunday    PropByDay = "SU"
	PropByDayMonday    PropByDay = "MO"
	PropByDayTuesday   PropByDay = "TU"
	PropByDayWednesday PropByDay = "WE"
	PropByDayThursday  PropByDay = "TH"
	PropByDayFriday    PropByDay = "FR"
	PropByDaySaturday  PropByDay = "SA"
	PropByDayFormula   PropByDay = ""
)

var (
	PropByDayAll = []PropByDay{PropByDaySunday, PropByDayMonday, PropByDayTuesday, PropByDayWednesday, PropByDayThursday, PropByDayFriday, PropByDaySaturday}
)

func (instance PropByDay) String() string {
	return string(instance)
}

func (instance PropByDay) Weekday() time.Weekday {
	switch instance {
	case PropByDaySunday:
		return time.Sunday
	case PropByDayMonday:
		return time.Monday
	case PropByDayTuesday:
		return time.Tuesday
	case PropByDayWednesday:
		return time.Wednesday
	case PropByDayThursday:
		return time.Thursday
	case PropByDayFriday:
		return time.Friday
	}
	return time.Saturday
}

func ParsePropByDay(value string) PropByDay {
	for _, prop := range PropByDayAll {
		if prop.String() == value {
			return prop
		}
	}
	return PropByDayFormula
}

type PropByDayValue struct {
	ByDay    []PropByDay
	Formula  string
	RawValue string
}

func (instance PropByDayValue) GoString() string {
	return instance.String()
}

func (instance PropByDayValue) String() string {
	if len(instance.Formula) == 0 {
		return qb_utils.Strings.ConcatSep(",", instance.ByDay)
	}
	return instance.RawValue
}

func (instance *PropByDayValue) MarshalJSON() ([]byte, error) {
	if len(instance.Formula) == 0 {
		return []byte(qb_utils.JSON.Stringify(instance.ByDay)), nil
	}
	return []byte(qb_utils.JSON.Stringify(instance.RawValue)), nil
}

//----------------------------------------------------------------------------------------------------------------------
//	RRuleSlot
//----------------------------------------------------------------------------------------------------------------------

type RRuleSlot struct {
	EventId  string        `json:"event_id"`
	StartAt  time.Time     `json:"start_at"`
	EndAt    time.Time     `json:"end_at"`
	Duration time.Duration `json:"duration"`
}

func (instance *RRuleSlot) String() string {
	return qb_utils.JSON.Stringify(instance)
}

func NewRRuleSlot(eventId string, start, end time.Time) *RRuleSlot {
	instance := new(RRuleSlot)
	instance.EventId = eventId
	instance.StartAt = start
	instance.EndAt = end
	instance.Duration = end.Sub(start)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	RRule
//----------------------------------------------------------------------------------------------------------------------

// RRule repeat rule
// https://icalendar.org/iCalendar-RFC-5545/3-3-10-recurrence-rule.html
// FREQ=WEEKLY;BYDAY=FR,MO,TH,TU,WE
// FREQ=WEEKLY;UNTIL=20220329T170000Z;INTERVAL=1;BYDAY=TU,WE,FR;WKST=SU
// FREQ=YEARLY;BYMONTH=3;BYDAY=-1SU
type RRule struct {
	EventId       string          `json:"event_id"`
	Frequency     PropFrequency   `json:"freq"`     // "SECONDLY" / "MINUTELY" / "HOURLY" / "DAILY" / "WEEKLY" / "MONTHLY" / "YEARLY"
	ByDay         *PropByDayValue `json:"byday"`    // "SU" / "MO" / "TU" / "WE" / "TH" / "FR" / "SA". Each BYDAY value can also be preceded by a positive (+n) or negative (-n) integer. If present, this indicates the nth occurrence of a specific day within the MONTHLY or YEARLY "RRULE". For example, within a MONTHLY rule, +1MO (or simply 1MO) represents the first Monday within the month, whereas -1MO represents the last Monday of the month.
	ByMonth       int             `json:"bymonth"`  // The BYMONTH rule part specifies a COMMA-separated list of months of the year. Valid values are 1 to 12.
	Until         time.Time       `json:"until"`    // The UNTIL rule part defines a DATE or DATE-TIME value that bounds the recurrence rule in an inclusive manner.
	Interval      int             `json:"interval"` // default is 1. The INTERVAL rule part contains a positive integer representing at which intervals the recurrence rule repeats. The default value is "1", meaning every second for a SECONDLY rule, every minute for a MINUTELY rule
	Count         int             `json:"count"`    // number or repeat. The COUNT rule part defines the number of occurrences at which to range-bound the recurrence. The "DTSTART" property value always counts as the first occurrence.
	WorkWeekStart PropByDay       `json:"wkst"`     // The WKST rule part specifies the day on which the workweek starts. Valid values are MO, TU, WE, TH, FR, SA, and SU. This is significant when a WEEKLY "RRULE" has an interval greater than 1, and a BYDAY rule part is specified.
}

func ParseRRule(text string) (*RRule, error) {
	instance := new(RRule)
	// default
	instance.Interval = 1
	instance.WorkWeekStart = PropByDayMonday

	props := strings.Split(text, ";")
	for _, prop := range props {
		tokens := strings.Split(prop, "=")
		if len(tokens) == 2 {
			propName := tokens[0]
			propValues := strings.Split(tokens[1], ",")
			if len(propValues) > 0 {
				switch propName {
				case "FREQ":
					if len(propValues[0]) > 0 {
						instance.Frequency = ParsePropFrequency(propValues[0])
					} else {
						instance.Frequency = PropFreqMonthly
					}
				case "UNTIL":
					t, err := ToTime(icalTimeFormat, propValues[0])
					if nil == err {
						instance.Until = t
					}
				case "BYMONTH":
					if len(propValues[0]) > 0 {
						instance.ByMonth = qb_utils.Convert.ToInt(propValues[0])
					}
				case "BYDAY":
					if len(propValues) > 0 {
						instance.ByDay = new(PropByDayValue)
						instance.ByDay.RawValue = tokens[1]
						for i, value := range propValues {
							day := ParsePropByDay(value)
							if i == 0 && len(day.String()) == 0 {
								instance.ByDay.Formula = value[:2]
								instance.ByDay.ByDay = append(instance.ByDay.ByDay, ParsePropByDay(value[2:]))
							} else {
								instance.ByDay.ByDay = append(instance.ByDay.ByDay, day)
							}
						}
					}
				case "WKST":
					if len(propValues[0]) > 0 {
						instance.WorkWeekStart = ParsePropByDay(propValues[0])
					}
				default:
					// not supported
				}
			}
		}
	}
	return instance, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *RRule) Serialize() string {
	if nil != instance {
		var sb strings.Builder
		if len(instance.Frequency) > 0 {
			if sb.Len() > 0 {
				sb.WriteString(";")
			}
			sb.WriteString("FREQ=")
			sb.WriteString(instance.Frequency.String())
		}
		if !qb_utils.Dates.IsZero(instance.Until) {
			if sb.Len() > 0 {
				sb.WriteString(";")
			}
			sb.WriteString("UNTIL=")
			sb.WriteString(instance.Until.Format(icalTimeFormat))
		}
		if instance.ByMonth > 0 {
			if sb.Len() > 0 {
				sb.WriteString(";")
			}
			sb.WriteString("BYMONTH=")
			sb.WriteString(fmt.Sprintf("%v", instance.ByMonth))
		}
		if nil != instance.ByDay {
			if sb.Len() > 0 {
				sb.WriteString(";")
			}
			sb.WriteString("BYDAY=")
			sb.WriteString(instance.ByDay.RawValue)
		}
		if len(instance.WorkWeekStart) > 0 {
			if sb.Len() > 0 {
				sb.WriteString(";")
			}
			sb.WriteString("WKST=")
			sb.WriteString(instance.WorkWeekStart.String())
		}
		return sb.String()
	}
	return ""
}

func (instance *RRule) Json() string {
	if nil != instance {
		return qb_utils.JSON.Stringify(instance)
	}
	return ""
}

func (instance *RRule) String() string {
	if nil != instance {
		return instance.Serialize()
	}
	return ""
}

func (instance *RRule) GetSlots(eventStart, eventEnd, until time.Time) []*RRuleSlot {
	if nil != instance {
		if qb_utils.Dates.IsZero(until) {
			until = instance.Until
		}
		if qb_utils.Dates.IsZero(until) {
			until = qb_utils.Dates.AddWeeks(time.Now(), 2) // add 2 weeks defaults
		}
		return instance.calculateSlots(eventStart, eventEnd, until)
	}
	return make([]*RRuleSlot, 0)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *RRule) calculateSlots(eventStart, eventEnd, until time.Time) []*RRuleSlot {
	response := make([]*RRuleSlot, 0)
	var done bool
	start := eventStart
	end := eventEnd
	// loop until end of "UNTIL" time
	for i := 0; !done; i++ {
		switch instance.Frequency {
		case PropFreqDaily:
			// every day FREQ=DAILY
			if i == 0 {
				response = append(response, NewRRuleSlot(instance.EventId, eventStart, eventEnd)) // add now
				done = start.Equal(until) || start.After(until) || end.Equal(until) || end.After(until)
			} else {
				start = qb_utils.Dates.AddDays(start, instance.Interval)
				end = qb_utils.Dates.AddDays(end, instance.Interval)
				done = start.Equal(until) || start.After(until) || end.Equal(until) || end.After(until)
				if !done {
					response = append(response, NewRRuleSlot(instance.EventId, start, end))
				}
			}
		case PropFreqWeekly:
			// every week FREQ=WEEKLY;BYDAY=WE
			if len(instance.ByDay.Formula) > 0 {
				// TODO: handle formulas
				done = true
			} else {
				days := instance.ByDay.ByDay
				for _, day := range days {
					start = qb_utils.Dates.NextWeekday(start, day.Weekday(), true)
					end = qb_utils.Dates.NextWeekday(end, day.Weekday(), true)
					done = start.Equal(until) || start.After(until) || end.Equal(until) || end.After(until)
					if !done {
						response = append(response, NewRRuleSlot(instance.EventId, start, end))
						start = qb_utils.Dates.AddDays(start, 1)
						end = qb_utils.Dates.AddDays(end, 1)
					}
				}
			}
		case PropFreqMonthly:
			if len(instance.ByDay.Formula) > 0 {
				// FREQ=MONTHLY;BYDAY=-1WE last we of month
				// FREQ=MONTHLY;BYDAY=+1WE first we of month
				day := instance.ByDay.ByDay[0]
				if strings.Index(instance.ByDay.Formula, "-") == 0 {
					// last day of month
					start = qb_utils.Dates.LastDayOfMonth(start, day.Weekday(), true, true)
					end = qb_utils.Dates.LastDayOfMonth(end, day.Weekday(), true, true)
					done = start.Equal(until) || start.After(until) || end.Equal(until) || end.After(until)
					if !done {
						response = append(response, NewRRuleSlot(instance.EventId, start, end))
						start = qb_utils.Dates.AddMonths(start, instance.Interval)
						end = qb_utils.Dates.AddMonths(end, instance.Interval)
					}
				} else {
					// first day of month
					start = qb_utils.Dates.FirstDayOfMonth(start, day.Weekday(), true, true)
					end = qb_utils.Dates.FirstDayOfMonth(end, day.Weekday(), true, true)
					done = start.Equal(until) || start.After(until) || end.Equal(until) || end.After(until)
					if !done {
						response = append(response, NewRRuleSlot(instance.EventId, start, end))
						start = qb_utils.Dates.AddMonths(start, instance.Interval)
						end = qb_utils.Dates.AddMonths(end, instance.Interval)
					}
				}
			} else {
				done = true
			}
		default:
			done = true
		}
	}
	return response
}

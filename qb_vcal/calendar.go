package qb_vcal

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_utils"
)

type CalendarProperty struct {
	BaseProperty
}

type Calendar struct {
	Components         []Component
	CalendarProperties []CalendarProperty
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r s
//----------------------------------------------------------------------------------------------------------------------

func NewCalendarFor(service string) *Calendar {
	c := &Calendar{
		Components:         []Component{},
		CalendarProperties: []CalendarProperty{},
	}
	c.SetVersion("1.0")
	c.SetProductId("-//" + service + "//GG library for VCALENDAR")
	return c
}

func ParseCalendar(r io.Reader) (*Calendar, error) {
	state := "begin"
	c := &Calendar{}
	cs := NewCalendarStream(r)
	cont := true
	for i := 0; cont; i++ {
		l, err := cs.ReadLine()
		if err != nil {
			switch err {
			case io.EOF:
				cont = false
			default:
				return c, err
			}
		}
		if l == nil || len(*l) == 0 {
			continue
		}
		line := ParseProperty(*l)
		if line == nil {
			return nil, qb_utils.Errors.Prefix(MalformedCalendarError, "Error parsing line")
		}
		switch state {
		case "begin":
			switch line.IANAToken {
			case "BEGIN":
				switch line.Value {
				case "VCALENDAR":
					state = "properties"
				default:
					return nil, qb_utils.Errors.Prefix(MalformedCalendarError, "Missing VCALENDAR property at the BEGIN")
				}
			default:
				return nil, qb_utils.Errors.Prefix(MalformedCalendarError, "Missing BEGIN statement.")
			}
		case "properties":
			switch line.IANAToken {
			case "END":
				switch line.Value {
				case "VCALENDAR":
					state = "end"
				default:
					return nil, qb_utils.Errors.Prefix(MalformedCalendarError, "Missing VCALENDAR property at the END")
				}
			case "BEGIN":
				state = "components"
			default:
				c.CalendarProperties = append(c.CalendarProperties, CalendarProperty{*line})
			}
			if state != "components" {
				break
			}
			fallthrough
		case "components":
			switch line.IANAToken {
			case "END":
				switch line.Value {
				case "VCALENDAR":
					state = "end"
				default:
					return nil, qb_utils.Errors.Prefix(MalformedCalendarError, "Missing VCALENDAR property at the END of components")
				}
			case "BEGIN":
				co, err := GeneralParseComponent(cs, line)
				if err != nil {
					return nil, err
				}
				if co != nil {
					c.Components = append(c.Components, co)
				}
			default:
				return nil, MalformedCalendarError
			}
		case "end":
			return nil, MalformedCalendarError
		default:
			return nil, CalendarStateError
		}
	}
	return c, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Calendar) String() string {
	b := bytes.NewBufferString("")
	_ = instance.SerializeTo(b)
	return b.String()
}

func (instance *Calendar) SerializeTo(w io.Writer) error {
	_, _ = fmt.Fprint(w, "BEGIN:VCALENDAR", "\r\n")
	for _, p := range instance.CalendarProperties {
		p.serialize(w)
	}
	for _, c := range instance.Components {
		c.serialize(w)
	}
	_, _ = fmt.Fprint(w, "END:VCALENDAR", "\r\n")
	return nil
}

func (instance *Calendar) Method() string {
	return instance.getPropertyAsString(PropertyMethod)
}

func (instance *Calendar) ProdId() string {
	return instance.getPropertyAsString(PropertyProductId)
}

func (instance *Calendar) ProdVendor() string {
	prodId := strings.ToLower(instance.ProdId())
	if strings.Index(prodId, "microsoft") > -1 {
		return "microsoft"
	} else if strings.Index(prodId, "google") > -1 {
		return "google"
	} else if strings.Index(prodId, "icloud") > -1 {
		return "apple"
	}
	return "caldav"
}

func (instance *Calendar) Name() string {
	return instance.getPropertyAsString(PropertyName)
}

func (instance *Calendar) Priority() string {
	return instance.getPropertyAsString(PropertyPriority)
}

func (instance *Calendar) Version() string {
	return instance.getPropertyAsString(PropertyVersion)
}

func (instance *Calendar) CalScale() string {
	return instance.getPropertyAsString(PropertyCalscale)
}

func (instance *Calendar) SetMethod(method Method, props ...PropertyParameter) {
	instance.setProperty(PropertyMethod, ToText(string(method)), props...)
}

func (instance *Calendar) SetXPublishedTTL(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyXPublishedTTL, string(s), props...)
}

func (instance *Calendar) SetVersion(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyVersion, ToText(s), props...)
}

func (instance *Calendar) SetProductId(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyProductId, ToText(s), props...)
}

func (instance *Calendar) SetName(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyName, string(s), props...)
}

func (instance *Calendar) SetColor(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyColor, string(s), props...)
}

func (instance *Calendar) SetXWRCalName(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyXWRCalName, string(s), props...)
}

func (instance *Calendar) SetXWRCalDesc(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyXWRCalDesc, string(s), props...)
}

func (instance *Calendar) SetXWRTimezone(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyXWRTimezone, string(s), props...)
}

func (instance *Calendar) SetXWRCalID(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyXWRCalID, string(s), props...)
}

func (instance *Calendar) SetDescription(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyDescription, ToText(s), props...)
}

func (instance *Calendar) SetLastModified(t time.Time, props ...PropertyParameter) {
	instance.setProperty(PropertyLastModified, t.UTC().Format(icalTimeFormat), props...)
}

func (instance *Calendar) SetRefreshInterval(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyRefreshInterval, string(s), props...)
}

func (instance *Calendar) SetCalscale(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyCalscale, string(s), props...)
}

func (instance *Calendar) SetTzid(s string, props ...PropertyParameter) {
	instance.setProperty(PropertyTzid, string(s), props...)
}

func (instance *Calendar) AddEvent(id string) *VEvent {
	e := NewEvent(id)
	instance.Components = append(instance.Components, e)
	return e
}

func (instance *Calendar) AddVEvent(e *VEvent) {
	instance.Components = append(instance.Components, e)
}

func (instance *Calendar) Events() (events []*VEvent) {
	events = []*VEvent{}
	for i := range instance.Components {
		switch event := instance.Components[i].(type) {
		case *VEvent:
			events = append(events, event)
		}
	}
	return
}

func (instance *Calendar) TimeZone() *VTimezone {
	for i := range instance.Components {
		switch comp := instance.Components[i].(type) {
		case *VTimezone:
			return comp
		}
	}
	return nil
}

func (instance *Calendar) IsTimeZoneGreenwichStandard() bool {
	timezone := instance.TimeZone()
	if nil != timezone {
		value := strings.ToLower(timezone.Value())
		return strings.Index(value, "standard") > -1
	}
	return true
}

func (instance *Calendar) EventGetStartAt() (t time.Time, err error) {
	events := instance.Events()
	if len(events) > 0 {
		event := events[0]
		t, err = event.GetStartAt()
	}
	return
}

func (instance *Calendar) EventGetEndAt() (t time.Time, err error) {
	events := instance.Events()
	if len(events) > 0 {
		event := events[0]
		t, err = event.GetEndAt()
	}
	return
}

func (instance *Calendar) EventGetLocation() string {
	events := instance.Events()
	if len(events) > 0 {
		return events[0].Location()
	}
	return ""
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Calendar) setProperty(property Property, value string, props ...PropertyParameter) {
	for i := range instance.CalendarProperties {
		if instance.CalendarProperties[i].IANAToken == string(property) {
			instance.CalendarProperties[i].Value = value
			instance.CalendarProperties[i].ICalParameters = map[string][]string{}
			for _, p := range props {
				k, v := p.KeyValue()
				instance.CalendarProperties[i].ICalParameters[k] = v
			}
			return
		}
	}
	r := CalendarProperty{
		BaseProperty{
			IANAToken:      string(property),
			Value:          value,
			ICalParameters: map[string][]string{},
		},
	}
	for _, p := range props {
		k, v := p.KeyValue()
		r.ICalParameters[k] = v
	}
	instance.CalendarProperties = append(instance.CalendarProperties, r)
}

func (instance *Calendar) getPropertyAsString(p Property) string {
	for _, prop := range instance.CalendarProperties {
		if prop.IANAToken == string(p) {
			return prop.Value
		}
	}
	return ""
}

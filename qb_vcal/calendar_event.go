package qb_vcal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_utils"
)

type VEvent struct {
	ComponentBase
}

func NewEvent(uniqueId string) *VEvent {
	e := &VEvent{
		ComponentBase{
			Properties: []IANAProperty{
				{BaseProperty{IANAToken: ToText(string(ComponentPropertyUniqueId)), Value: uniqueId}},
			},
		},
	}
	return e
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *VEvent) String() string {
	b := &bytes.Buffer{}
	instance.ComponentBase.serializeThis(b, "VEVENT")
	return b.String()
}

func (instance *VEvent) GetProperty(componentProperty ComponentProperty) *IANAProperty {
	return GetProperty(componentProperty, instance.Properties)
}

func (instance *VEvent) SetProperty(property ComponentProperty, value string, props ...PropertyParameter) {
	for i := range instance.Properties {
		if instance.Properties[i].IANAToken == string(property) {
			instance.Properties[i].Value = value
			instance.Properties[i].ICalParameters = map[string][]string{}
			for _, p := range props {
				k, v := p.KeyValue()
				instance.Properties[i].ICalParameters[k] = v
			}
			return
		}
	}
	instance.AddProperty(property, value, props...)
}

func (instance *VEvent) AddProperty(property ComponentProperty, value string, props ...PropertyParameter) {
	r := IANAProperty{
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
	instance.Properties = append(instance.Properties, r)
}

func (instance *VEvent) SetCreatedTime(t time.Time, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyCreated, t.UTC().Format(icalTimeFormat), props...)
}

func (instance *VEvent) SetDtStampTime(t time.Time, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyDtstamp, t.UTC().Format(icalTimeFormat), props...)
}

func (instance *VEvent) SetModifiedAt(t time.Time, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyLastModified, t.UTC().Format(icalTimeFormat), props...)
}

func (instance *VEvent) SetSequence(seq int, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertySequence, strconv.Itoa(seq), props...)
}

func (instance *VEvent) SetStartAt(t time.Time, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyDtStart, t.UTC().Format(icalTimeFormat), props...)
}

func (instance *VEvent) SetAllDayStartAt(t time.Time, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyDtStart, t.UTC().Format(icalAllDayTimeFormat), props...)
}

func (instance *VEvent) SetEndAt(t time.Time, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyDtEnd, t.UTC().Format(icalTimeFormat), props...)
}

func (instance *VEvent) SetAllDayEndAt(t time.Time, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyDtEnd, t.UTC().Format(icalAllDayTimeFormat), props...)
}

func (instance *VEvent) SetTimeTransparency(v TimeTransparency, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyTransp, string(v), props...)
}

func (instance *VEvent) SetSummary(s string, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertySummary, ToText(s), props...)
}

func (instance *VEvent) SetStatus(s ObjectStatus, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyStatus, ToText(string(s)), props...)
}

func (instance *VEvent) SetDescription(s string, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyDescription, ToText(s), props...)
}

func (instance *VEvent) SetLocation(s string, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyLocation, ToText(s), props...)
}

func (instance *VEvent) SetGeo(lat interface{}, lng interface{}, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyGeo, fmt.Sprintf("%v;%v", lat, lng), props...)
}

func (instance *VEvent) SetURL(s string, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyUrl, s, props...)
}

func (instance *VEvent) SetOrganizer(s string, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyOrganizer, s, props...)
}

func (instance *VEvent) SetColor(s string, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyColor, s, props...)
}

func (instance *VEvent) SetClass(c Classification, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyClass, string(c), props...)
}

func (instance *VEvent) AddAttendee(s string, props ...PropertyParameter) {
	instance.AddProperty(ComponentPropertyAttendee, "mailto:"+s, props...)
}

func (instance *VEvent) AddExdate(s string, props ...PropertyParameter) {
	instance.AddProperty(ComponentPropertyExdate, s, props...)
}

func (instance *VEvent) AddExrule(s string, props ...PropertyParameter) {
	instance.AddProperty(ComponentPropertyExrule, s, props...)
}

func (instance *VEvent) AddRdate(s string, props ...PropertyParameter) {
	instance.AddProperty(ComponentPropertyRdate, s, props...)
}

func (instance *VEvent) AddRrule(s string, props ...PropertyParameter) {
	instance.AddProperty(ComponentPropertyRrule, s, props...)
}

func (instance *VEvent) AddAttachment(s string, props ...PropertyParameter) {
	instance.AddProperty(ComponentPropertyAttach, s, props...)
}

func (instance *VEvent) AddAttachmentURL(uri string, contentType string) {
	instance.AddAttachment(uri, WithFmtType(contentType))
}

func (instance *VEvent) AddAttachmentBinary(binary []byte, contentType string) {
	instance.AddAttachment(base64.StdEncoding.EncodeToString(binary),
		WithFmtType(contentType), WithEncoding("base64"), WithValue("binary"),
	)
}

func (instance *VEvent) AddAlarm() *VAlarm {
	a := &VAlarm{
		ComponentBase: ComponentBase{},
	}
	instance.Components = append(instance.Components, a)
	return a
}

func (instance *VEvent) Id() string {
	response, _ := GetPropertyAsString(ComponentPropertyUniqueId, instance.Properties)
	return response
}

func (instance *VEvent) Summary() string {
	response, _ := GetPropertyAsString(ComponentPropertySummary, instance.Properties)
	return response
}

func (instance *VEvent) Description() string {
	response, _ := GetPropertyAsString(ComponentPropertyDescription, instance.Properties)
	return response
}

func (instance *VEvent) Links() []string {
	description := instance.Description()
	return qb_utils.Regex.Links(description)
}

func (instance *VEvent) LinkMeeting() string {
	links := instance.Links()
	for _, link := range links {
		for _, check := range ConferencePlatforms {
			if strings.Index(strings.ToLower(link), check) > -1 {
				return link
			}
		}
	}
	return ""
}

func (instance *VEvent) Status() string {
	response, _ := GetPropertyAsString(ComponentPropertyStatus, instance.Properties)
	return response
}

func (instance *VEvent) IsStatusCancelled() bool {
	status := instance.Status()
	return status == "CANCELLED"
}

func (instance *VEvent) Alarms() (r []*VAlarm) {
	r = []*VAlarm{}
	for i := range instance.Components {
		switch alarm := instance.Components[i].(type) {
		case *VAlarm:
			r = append(r, alarm)
		}
	}
	return
}

func (instance *VEvent) Attendees() (r []*Attendee) {
	r = []*Attendee{}
	for i := range instance.Properties {
		switch instance.Properties[i].IANAToken {
		case string(ComponentPropertyAttendee):
			a := &Attendee{
				instance.Properties[i],
			}
			r = append(r, a)
		}
	}
	return
}

func (instance *VEvent) GetStartAt() (t time.Time, err error) {
	t, err = GetPropertyAsTime(ComponentPropertyDtStart, icalTimeFormat, instance.Properties)
	if nil != err {
		return
	}
	return
}

func (instance *VEvent) GetEndAt() (t time.Time, err error) {
	t, err = GetPropertyAsTime(ComponentPropertyDtEnd, icalTimeFormat, instance.Properties)
	if nil != err {
		return
	}
	return
}

func (instance *VEvent) GetAllDayStartAt() (t time.Time, err error) {
	t, err = GetPropertyAsTime(ComponentPropertyDtStart, icalAllDayTimeFormat, instance.Properties)
	if nil != err {
		return
	}
	return
}

func (instance *VEvent) GetAllDayEndAt() (t time.Time, err error) {
	t, err = GetPropertyAsTime(ComponentPropertyDtEnd, icalAllDayTimeFormat, instance.Properties)
	if nil != err {
		return
	}
	return
}

func (instance *VEvent) Location() string {
	value, _ := GetPropertyAsString(ComponentPropertyLocation, instance.Properties)
	return value
}

func (instance *VEvent) Organizer() string {
	prop := GetProperty(ComponentPropertyOrganizer, instance.Properties)
	if nil != prop {
		return strings.TrimSpace(strings.Replace(prop.Value, "mailto:", "", 1))
	}
	return ""
}

func (instance *VEvent) AttendeesEmails() (response []string) {
	attendees := instance.Attendees()
	if len(attendees) > 0 {
		for _, attendee := range attendees {
			response = append(response, attendee.Email())
		}
	}
	return
}

func (instance *VEvent) Duration() time.Duration {
	startAt, es := instance.GetStartAt()
	endAt, ee := instance.GetEndAt()
	if nil == ee && nil == es && !qb_utils.Dates.IsZero(startAt) && !qb_utils.Dates.IsZero(endAt) {
		return endAt.Sub(startAt)
	}
	return 0 * time.Second
}

// RRule Repeat Rule
func (instance *VEvent) RRule() *RRule {
	if nil != instance {
		prop := GetProperty(ComponentPropertyRrule, instance.Properties)
		if nil != prop {
			rule, _ := ParseRRule(prop.Value)
			if nil != rule {
				rule.EventId = instance.Id()
			}
			return rule
		}
	}
	return nil
}

func (instance *VEvent) GetSlots(until time.Time) []*RRuleSlot {
	rrule := instance.RRule()
	if nil != rrule {
		startAt, es := instance.GetStartAt()
		endAt, ee := instance.GetEndAt()
		if nil == ee && nil == es && !qb_utils.Dates.IsZero(startAt) && !qb_utils.Dates.IsZero(endAt) {
			return rrule.GetSlots(startAt, endAt, until)
		}
	}
	return make([]*RRuleSlot, 0)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *VEvent) serialize(w io.Writer) {
	instance.ComponentBase.serializeThis(w, "VEVENT")
}

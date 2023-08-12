package qb_vcal

import (
	"bytes"
	"io"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	Component
//----------------------------------------------------------------------------------------------------------------------

type Component interface {
	UnknownPropertiesIANAProperties() []IANAProperty
	SubComponents() []Component
	serialize(b io.Writer)
}

type ComponentBase struct {
	Properties []IANAProperty
	Components []Component
}

func (instance *ComponentBase) UnknownPropertiesIANAProperties() []IANAProperty {
	return instance.Properties
}

func (instance *ComponentBase) SubComponents() []Component {
	return instance.Components
}
func (instance ComponentBase) serializeThis(writer io.Writer, componentType string) {
	SerializeThis(instance, writer, componentType)
}

//----------------------------------------------------------------------------------------------------------------------
//	Attendee
//----------------------------------------------------------------------------------------------------------------------

type Attendee struct {
	IANAProperty
}

func (instance *Attendee) Email() string {
	if strings.HasPrefix(instance.Value, "mailto:") {
		return instance.Value[len("mailto:"):]
	}
	return instance.Value
}

func (instance *Attendee) ParticipationStatus() ParticipationStatus {
	return ParticipationStatus(instance.getPropertyFirst(ParameterParticipationStatus))
}

func (instance *Attendee) getPropertyFirst(parameter Parameter) string {
	vs := instance.getProperty(parameter)
	if len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (instance *Attendee) getProperty(parameter Parameter) []string {
	if vs, ok := instance.ICalParameters[string(parameter)]; ok {
		return vs
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	VTodo
//----------------------------------------------------------------------------------------------------------------------

type VTodo struct {
	ComponentBase
}

func (instance *VTodo) serialize(w io.Writer) {
	instance.ComponentBase.serializeThis(w, "VTODO")
}

func (instance *VTodo) Serialize() string {
	b := &bytes.Buffer{}
	instance.ComponentBase.serializeThis(b, "VTODO")
	return b.String()
}

//----------------------------------------------------------------------------------------------------------------------
//	VJournal
//----------------------------------------------------------------------------------------------------------------------

type VJournal struct {
	ComponentBase
}

func (instance *VJournal) serialize(w io.Writer) {
	instance.ComponentBase.serializeThis(w, "VJOURNAL")
}

func (instance *VJournal) Serialize() string {
	b := &bytes.Buffer{}
	instance.ComponentBase.serializeThis(b, "VJOURNAL")
	return b.String()
}

//----------------------------------------------------------------------------------------------------------------------
//	VBusy
//----------------------------------------------------------------------------------------------------------------------

type VBusy struct {
	ComponentBase
}

func (instance *VBusy) Serialize() string {
	b := &bytes.Buffer{}
	instance.ComponentBase.serializeThis(b, "VBUSY")
	return b.String()
}

func (instance *VBusy) serialize(w io.Writer) {
	instance.ComponentBase.serializeThis(w, "VBUSY")
}

//----------------------------------------------------------------------------------------------------------------------
//	VTimezone
//----------------------------------------------------------------------------------------------------------------------

type VTimezone struct {
	ComponentBase
}

func (instance *VTimezone) Serialize() string {
	b := &bytes.Buffer{}
	instance.ComponentBase.serializeThis(b, "VTIMEZONE")
	return b.String()
}

func (instance *VTimezone) serialize(w io.Writer) {
	instance.ComponentBase.serializeThis(w, "VTIMEZONE")
}

func (instance *VTimezone) Value() string {
	if nil != instance && nil != instance.ComponentBase.Properties {
		for _, prop := range instance.ComponentBase.Properties {
			if prop.IANAToken == string(PropertyTzid) {
				return prop.Value
			}
		}
	}
	return ""
}

//----------------------------------------------------------------------------------------------------------------------
//	VAlarm
//----------------------------------------------------------------------------------------------------------------------

type VAlarm struct {
	ComponentBase
}

func (instance *VAlarm) Serialize() string {
	b := &bytes.Buffer{}
	instance.ComponentBase.serializeThis(b, "VALARM")
	return b.String()
}

func (instance *VAlarm) serialize(w io.Writer) {
	instance.ComponentBase.serializeThis(w, "VALARM")
}

func (instance *VAlarm) SetProperty(property ComponentProperty, value string, props ...PropertyParameter) {
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

func (instance *VAlarm) AddProperty(property ComponentProperty, value string, props ...PropertyParameter) {
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

func (instance *VAlarm) SetAction(a Action, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyAction, string(a), props...)
}

func (instance *VAlarm) SetTrigger(s string, props ...PropertyParameter) {
	instance.SetProperty(ComponentPropertyTrigger, s, props...)
}

//----------------------------------------------------------------------------------------------------------------------
//	Standard
//----------------------------------------------------------------------------------------------------------------------

type Standard struct {
	ComponentBase
}

func (instance *Standard) Serialize() string {
	b := &bytes.Buffer{}
	instance.ComponentBase.serializeThis(b, "STANDARD")
	return b.String()
}

func (instance *Standard) serialize(w io.Writer) {
	instance.ComponentBase.serializeThis(w, "STANDARD")
}

//----------------------------------------------------------------------------------------------------------------------
//	Daylight
//----------------------------------------------------------------------------------------------------------------------

type Daylight struct {
	ComponentBase
}

func (instance *Daylight) Serialize() string {
	b := &bytes.Buffer{}
	instance.ComponentBase.serializeThis(b, "DAYLIGHT")
	return b.String()
}

func (instance *Daylight) serialize(w io.Writer) {
	instance.ComponentBase.serializeThis(w, "DAYLIGHT")
}

//----------------------------------------------------------------------------------------------------------------------
//	GeneralComponent
//----------------------------------------------------------------------------------------------------------------------

type GeneralComponent struct {
	ComponentBase
	Token string
}

func (instance *GeneralComponent) Serialize() string {
	b := &bytes.Buffer{}
	instance.ComponentBase.serializeThis(b, instance.Token)
	return b.String()
}

func (instance *GeneralComponent) serialize(w io.Writer) {
	instance.ComponentBase.serializeThis(w, instance.Token)
}

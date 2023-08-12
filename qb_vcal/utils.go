package qb_vcal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_utils"
)

var textEscaper = strings.NewReplacer(
	`\`, `\\`,
	"\n", `\n`,
	`;`, `\;`,
	`,`, `\,`,
)

var textUnescaper = strings.NewReplacer(
	`\\`, `\`,
	`\n`, "\n",
	`\N`, "\n",
	`\;`, `;`,
	`\,`, `,`,
)

func ToText(s string) string {
	// Some special characters for iCalendar format should be escaped while
	// setting a value of a property with a TEXT type.
	return textEscaper.Replace(s)
}

func FromText(s string) string {
	// Some special characters for iCalendar format should be escaped while
	// setting a value of a property with a TEXT type.
	return textUnescaper.Replace(s)
}

func ToTime(format, value string) (t time.Time, err error) {
	if len(value) == 15 {
		if strings.HasSuffix(value, "Z") {
			value = strings.ReplaceAll(value, "Z", "0Z")
		} else {
			value = value + "Z"
		}
	}
	t, err = time.Parse(format, value)
	if nil != err {
		t, err = qb_utils.Dates.ParseAny(value)
	}
	return
}

func GetProperty(componentProperty ComponentProperty, properties []IANAProperty) *IANAProperty {
	for i := range properties {
		if properties[i].IANAToken == string(componentProperty) {
			return &properties[i]
		}
	}
	return nil
}

func GetPropertyAsTime(componentProperty ComponentProperty, tFormat string, properties []IANAProperty) (t time.Time, err error) {
	timeProp := GetProperty(componentProperty, properties)
	if timeProp == nil {
		return time.Time{}, errors.New("property not found")
	}
	value := timeProp.BaseProperty.Value
	t, err = ToTime(tFormat, value)
	// correct time adding
	return
}

func GetPropertyAsString(componentProperty ComponentProperty, properties []IANAProperty) (string, error) {
	timeProp := GetProperty(componentProperty, properties)
	if timeProp == nil {
		return "", errors.New("property not found")
	}
	return FromText(timeProp.BaseProperty.Value), nil
}

func Serialize(component GeneralComponent) string {
	b := &bytes.Buffer{}
	SerializeThis(component.ComponentBase, b, component.Token)
	return b.String()
}

func SerializeThis(component ComponentBase, writer io.Writer, componentType string) {
	_, _ = fmt.Fprint(writer, "BEGIN:"+componentType, "\r\n")
	for _, p := range component.Properties {
		p.serialize(writer)
	}
	for _, c := range component.Components {
		c.serialize(writer)
	}
	_, _ = fmt.Fprint(writer, "END:"+componentType, "\r\n")
}

package qb_vcal

import "errors"

const (
	icalTimeFormat       = "20060102T150405Z"
	icalAllDayTimeFormat = "20060102"
)

var (
	CalendarStateError     = errors.New("calendar_state_error")
	MalformedCalendarError = errors.New("malformed_calendar_error")

	ConferencePlatforms = []string{"skype", "zoom", "meet", "gotomeeting", "pexip",
		"facebook", "messenger", "telegram", "whatsapp"}
)

type ComponentType string

const (
	ComponentVCalendar ComponentType = "VCALENDAR"
	ComponentVEvent    ComponentType = "VEVENT"
	ComponentVTodo     ComponentType = "VTODO"
	ComponentVJournal  ComponentType = "VJOURNAL"
	ComponentVFreeBusy ComponentType = "VFREEBUSY"
	ComponentVTimezone ComponentType = "VTIMEZONE"
	ComponentVAlarm    ComponentType = "VALARM"
	ComponentStandard  ComponentType = "STANDARD"
	ComponentDaylight  ComponentType = "DAYLIGHT"
)

type ComponentProperty Property

const (
	ComponentPropertyUniqueId            = ComponentProperty(PropertyUid) // TEXT
	ComponentPropertyDtstamp             = ComponentProperty(PropertyDtstamp)
	ComponentPropertyOrganizer           = ComponentProperty(PropertyOrganizer)
	ComponentPropertyAttendee            = ComponentProperty(PropertyAttendee)
	ComponentPropertyAttach              = ComponentProperty(PropertyAttach)
	ComponentPropertyDescription         = ComponentProperty(PropertyDescription) // TEXT
	ComponentPropertyCategories          = ComponentProperty(PropertyCategories)  // TEXT
	ComponentPropertyClass               = ComponentProperty(PropertyClass)       // TEXT
	ComponentPropertyColor               = ComponentProperty(PropertyColor)       // TEXT
	ComponentPropertyCreated             = ComponentProperty(PropertyCreated)
	ComponentPropertySummary             = ComponentProperty(PropertySummary) // TEXT
	ComponentPropertyDtStart             = ComponentProperty(PropertyDtstart)
	ComponentPropertyDtEnd               = ComponentProperty(PropertyDtend)
	ComponentPropertyLocation            = ComponentProperty(PropertyLocation)            // TEXT
	ComponentPropertyXMicrosoftLocations = ComponentProperty(PropertyXMicrosoftLocations) // array of map
	ComponentPropertyStatus              = ComponentProperty(PropertyStatus)              // TEXT
	ComponentPropertyFreebusy            = ComponentProperty(PropertyFreebusy)
	ComponentPropertyLastModified        = ComponentProperty(PropertyLastModified)
	ComponentPropertyUrl                 = ComponentProperty(PropertyUrl)
	ComponentPropertyGeo                 = ComponentProperty(PropertyGeo)
	ComponentPropertyTransp              = ComponentProperty(PropertyTransp)
	ComponentPropertySequence            = ComponentProperty(PropertySequence)
	ComponentPropertyExdate              = ComponentProperty(PropertyExdate)
	ComponentPropertyExrule              = ComponentProperty(PropertyExrule)
	ComponentPropertyRdate               = ComponentProperty(PropertyRdate)
	ComponentPropertyRrule               = ComponentProperty(PropertyRrule)
	ComponentPropertyAction              = ComponentProperty(PropertyAction)
	ComponentPropertyTrigger             = ComponentProperty(PropertyTrigger)
)

type Property string

const (
	PropertyCalscale            Property = "CALSCALE" // TEXT
	PropertyMethod              Property = "METHOD"   // TEXT
	PropertyProductId           Property = "PRODID"   // TEXT
	PropertyVersion             Property = "VERSION"  // TEXT
	PropertyXPublishedTTL       Property = "X-PUBLISHED-TTL"
	PropertyRefreshInterval     Property = "REFRESH-INTERVAL;VALUE=DURATION"
	PropertyAttach              Property = "ATTACH"
	PropertyCategories          Property = "CATEGORIES"  // TEXT
	PropertyClass               Property = "CLASS"       // TEXT
	PropertyColor               Property = "COLOR"       // TEXT
	PropertyComment             Property = "COMMENT"     // TEXT
	PropertyDescription         Property = "DESCRIPTION" // TEXT
	PropertyXWRCalDesc          Property = "X-WR-CALDESC"
	PropertyGeo                 Property = "GEO"
	PropertyLocation            Property = "LOCATION" // TEXT
	PropertyXMicrosoftLocations Property = "X-MICROSOFT-LOCATIONS"
	PropertyPercentComplete     Property = "PERCENT-COMPLETE"
	PropertyPriority            Property = "PRIORITY"
	PropertyResources           Property = "RESOURCES" // TEXT
	PropertyStatus              Property = "STATUS"    // TEXT
	PropertySummary             Property = "SUMMARY"   // TEXT
	PropertyCompleted           Property = "COMPLETED"
	PropertyDtend               Property = "DTEND"
	PropertyDue                 Property = "DUE"
	PropertyDtstart             Property = "DTSTART"
	PropertyDuration            Property = "DURATION"
	PropertyFreebusy            Property = "FREEBUSY"
	PropertyTransp              Property = "TRANSP" // TEXT
	PropertyTzid                Property = "TZID"   // TEXT
	PropertyTzname              Property = "TZNAME" // TEXT
	PropertyTzoffsetfrom        Property = "TZOFFSETFROM"
	PropertyTzoffsetto          Property = "TZOFFSETTO"
	PropertyTzurl               Property = "TZURL"
	PropertyAttendee            Property = "ATTENDEE"
	PropertyContact             Property = "CONTACT" // TEXT
	PropertyOrganizer           Property = "ORGANIZER"
	PropertyRecurrenceId        Property = "RECURRENCE-ID"
	PropertyRelatedTo           Property = "RELATED-TO" // TEXT
	PropertyUrl                 Property = "URL"
	PropertyUid                 Property = "UID" // TEXT
	PropertyExdate              Property = "EXDATE"
	PropertyExrule              Property = "EXRULE"
	PropertyRdate               Property = "RDATE"
	PropertyRrule               Property = "RRULE"
	PropertyAction              Property = "ACTION" // TEXT
	PropertyRepeat              Property = "REPEAT"
	PropertyTrigger             Property = "TRIGGER"
	PropertyCreated             Property = "CREATED"
	PropertyDtstamp             Property = "DTSTAMP"
	PropertyLastModified        Property = "LAST-MODIFIED"
	PropertyRequestStatus       Property = "REQUEST-STATUS" // TEXT
	PropertyName                Property = "NAME"
	PropertyXWRCalName          Property = "X-WR-CALNAME"
	PropertyXWRTimezone         Property = "X-WR-TIMEZONE"
	PropertySequence            Property = "SEQUENCE"
	PropertyXWRCalID            Property = "X-WR-RELCALID"
)

type Parameter string

const (
	ParameterAltrep              Parameter = "ALTREP"
	ParameterCn                  Parameter = "CN"
	ParameterCutype              Parameter = "CUTYPE"
	ParameterDelegatedFrom       Parameter = "DELEGATED-FROM"
	ParameterDelegatedTo         Parameter = "DELEGATED-TO"
	ParameterDir                 Parameter = "DIR"
	ParameterEncoding            Parameter = "ENCODING"
	ParameterFmttype             Parameter = "FMTTYPE"
	ParameterFbtype              Parameter = "FBTYPE"
	ParameterLanguage            Parameter = "LANGUAGE"
	ParameterMember              Parameter = "MEMBER"
	ParameterParticipationStatus Parameter = "PARTSTAT"
	ParameterRange               Parameter = "RANGE"
	ParameterRelated             Parameter = "RELATED"
	ParameterReltype             Parameter = "RELTYPE"
	ParameterRole                Parameter = "ROLE"
	ParameterRsvp                Parameter = "RSVP"
	ParameterSentBy              Parameter = "SENT-BY"
	ParameterTzid                Parameter = "TZID"
	ParameterValue               Parameter = "VALUE"
)

type ValueDataType string

const (
	ValueDataTypeBinary     ValueDataType = "BINARY"
	ValueDataTypeBoolean    ValueDataType = "BOOLEAN"
	ValueDataTypeCalAddress ValueDataType = "CAL-ADDRESS"
	ValueDataTypeDate       ValueDataType = "DATE"
	ValueDataTypeDateTime   ValueDataType = "DATE-TIME"
	ValueDataTypeDuration   ValueDataType = "DURATION"
	ValueDataTypeFloat      ValueDataType = "FLOAT"
	ValueDataTypeInteger    ValueDataType = "INTEGER"
	ValueDataTypePeriod     ValueDataType = "PERIOD"
	ValueDataTypeRecur      ValueDataType = "RECUR"
	ValueDataTypeText       ValueDataType = "TEXT"
	ValueDataTypeTime       ValueDataType = "TIME"
	ValueDataTypeUri        ValueDataType = "URI"
	ValueDataTypeUtcOffset  ValueDataType = "UTC-OFFSET"
)

type CalendarUserType string

const (
	CalendarUserTypeIndividual CalendarUserType = "INDIVIDUAL"
	CalendarUserTypeGroup      CalendarUserType = "GROUP"
	CalendarUserTypeResource   CalendarUserType = "RESOURCE"
	CalendarUserTypeRoom       CalendarUserType = "ROOM"
	CalendarUserTypeUnknown    CalendarUserType = "UNKNOWN"
)

func (cut CalendarUserType) KeyValue(s ...interface{}) (string, []string) {
	return string(ParameterCutype), []string{string(cut)}
}

type FreeBusyTimeType string

const (
	FreeBusyTimeTypeFree            FreeBusyTimeType = "FREE"
	FreeBusyTimeTypeBusy            FreeBusyTimeType = "BUSY"
	FreeBusyTimeTypeBusyUnavailable FreeBusyTimeType = "BUSY-UNAVAILABLE"
	FreeBusyTimeTypeBusyTentative   FreeBusyTimeType = "BUSY-TENTATIVE"
)

type ParticipationStatus string

const (
	ParticipationStatusNeedsAction ParticipationStatus = "NEEDS-ACTION"
	ParticipationStatusAccepted    ParticipationStatus = "ACCEPTED"
	ParticipationStatusDeclined    ParticipationStatus = "DECLINED"
	ParticipationStatusTentative   ParticipationStatus = "TENTATIVE"
	ParticipationStatusDelegated   ParticipationStatus = "DELEGATED"
	ParticipationStatusCompleted   ParticipationStatus = "COMPLETED"
	ParticipationStatusInProcess   ParticipationStatus = "IN-PROCESS"
)

func (ps ParticipationStatus) KeyValue(s ...interface{}) (string, []string) {
	return string(ParameterParticipationStatus), []string{string(ps)}
}

type ObjectStatus string

const (
	ObjectStatusTentative   ObjectStatus = "TENTATIVE"
	ObjectStatusConfirmed   ObjectStatus = "CONFIRMED"
	ObjectStatusCancelled   ObjectStatus = "CANCELLED"
	ObjectStatusNeedsAction ObjectStatus = "NEEDS-ACTION"
	ObjectStatusCompleted   ObjectStatus = "COMPLETED"
	ObjectStatusInProcess   ObjectStatus = "IN-PROCESS"
	ObjectStatusDraft       ObjectStatus = "DRAFT"
	ObjectStatusFinal       ObjectStatus = "FINAL"
)

func (ps ObjectStatus) KeyValue(s ...interface{}) (string, []string) {
	return string(PropertyStatus), []string{ToText(string(ps))}
}

type RelationshipType string

const (
	RelationshipTypeChild   RelationshipType = "CHILD"
	RelationshipTypeParent  RelationshipType = "PARENT"
	RelationshipTypeSibling RelationshipType = "SIBLING"
)

type ParticipationRole string

const (
	ParticipationRoleChair          ParticipationRole = "CHAIR"
	ParticipationRoleReqParticipant ParticipationRole = "REQ-PARTICIPANT"
	ParticipationRoleOptParticipant ParticipationRole = "OPT-PARTICIPANT"
	ParticipationRoleNonParticipant ParticipationRole = "NON-PARTICIPANT"
)

func (pr ParticipationRole) KeyValue(s ...interface{}) (string, []string) {
	return string(ParameterRole), []string{string(pr)}
}

type Action string

const (
	ActionAudio     Action = "AUDIO"
	ActionDisplay   Action = "DISPLAY"
	ActionEmail     Action = "EMAIL"
	ActionProcedure Action = "PROCEDURE"
)

type Classification string

const (
	ClassificationPublic       Classification = "PUBLIC"
	ClassificationPrivate      Classification = "PRIVATE"
	ClassificationConfidential Classification = "CONFIDENTIAL"
)

type Method string

const (
	MethodPublish        Method = "PUBLISH"
	MethodRequest        Method = "REQUEST"
	MethodReply          Method = "REPLY"
	MethodAdd            Method = "ADD"
	MethodCancel         Method = "CANCEL"
	MethodRefresh        Method = "REFRESH"
	MethodCounter        Method = "COUNTER"
	MethodDeclinecounter Method = "DECLINECOUNTER"
)

type TimeTransparency string

const (
	TransparencyOpaque      TimeTransparency = "OPAQUE" // default
	TransparencyTransparent TimeTransparency = "TRANSPARENT"
)

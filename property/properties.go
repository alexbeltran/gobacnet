package property

import "fmt"

type PropertyID uint32

const (
	AckedTransitions               PropertyID = 0
	AckRequired                    PropertyID = 1
	Action                         PropertyID = 2
	ActionText                     PropertyID = 3
	ActiveText                     PropertyID = 4
	ActiveVTSessions               PropertyID = 5
	AlarmValue                     PropertyID = 6
	AlarmValues                    PropertyID = 7
	All                            PropertyID = 8
	AllWritesSuccessful            PropertyID = 9
	ApduSegmentTimeout             PropertyID = 10
	ApduTimeout                    PropertyID = 11
	ApplicationSoftwareVersion     PropertyID = 12
	Archive                        PropertyID = 13
	Bias                           PropertyID = 14
	ChangeOfStateCount             PropertyID = 15
	ChangeOfStateTime              PropertyID = 16
	NotificationClass              PropertyID = 17
	None                           PropertyID = 18
	ControlledVariableReference    PropertyID = 19
	ControlledVariableUnits        PropertyID = 20
	ControlledVariableValue        PropertyID = 21
	CovIncrement                   PropertyID = 22
	DateList                       PropertyID = 23
	DaylightSavingsStatus          PropertyID = 24
	Deadband                       PropertyID = 25
	DerivativeConstant             PropertyID = 26
	DerivativeConstantUnits        PropertyID = 27
	Description                    PropertyID = 28
	DescriptionOfHalt              PropertyID = 29
	DeviceAddressBinding           PropertyID = 30
	DeviceType                     PropertyID = 31
	EffectivePeriod                PropertyID = 32
	ElapsedActiveTime              PropertyID = 33
	ErrorLimit                     PropertyID = 34
	EventEnable                    PropertyID = 35
	EventState                     PropertyID = 36
	EventType                      PropertyID = 37
	ExceptionSchedule              PropertyID = 38
	FaultValues                    PropertyID = 39
	FeedbackValue                  PropertyID = 40
	FileAccessMethod               PropertyID = 41
	FileSize                       PropertyID = 42
	FileType                       PropertyID = 43
	FirmwareRevision               PropertyID = 44
	HighLimit                      PropertyID = 45
	InactiveText                   PropertyID = 46
	InProcess                      PropertyID = 47
	InstanceOf                     PropertyID = 48
	IntegralConstant               PropertyID = 49
	IntegralConstantUnits          PropertyID = 50
	IssueConfirmedNotifications    PropertyID = 51
	LimitEnable                    PropertyID = 52
	ListOfGroupMembers             PropertyID = 53
	ListOfObjectPropertyReferences PropertyID = 54
	ListOfSessionKeys              PropertyID = 55
	LocalDate                      PropertyID = 56
	LocalTime                      PropertyID = 57
	Location                       PropertyID = 58
	LowLimit                       PropertyID = 59
	ManipulatedVariableReference   PropertyID = 60
	MaximumOutput                  PropertyID = 61
	MaxApduLengthAccepted          PropertyID = 62
	MaxInfoFrames                  PropertyID = 63
	MaxMaster                      PropertyID = 64
	MaxPresValue                   PropertyID = 65
	MinimumOffTime                 PropertyID = 66
	MinimumOnTime                  PropertyID = 67
	MinimumOutput                  PropertyID = 68
	MinPresValue                   PropertyID = 69
	ModelName                      PropertyID = 70
	ModificationDate               PropertyID = 71
	NotifyType                     PropertyID = 72
	NumberOfAPDURetries            PropertyID = 73
	NumberOfStates                 PropertyID = 74
	ObjectIdentifier               PropertyID = 75
	ObjectList                     PropertyID = 76
	ObjectName                     PropertyID = 77
	ObjectPropertyReference        PropertyID = 78
	ObjectType                     PropertyID = 79
	Optional                       PropertyID = 80
	OutOfService                   PropertyID = 81
	OutputUnits                    PropertyID = 82
	EventParameters                PropertyID = 83
	Polarity                       PropertyID = 84
	PresentValue                   PropertyID = 85
	Priority                       PropertyID = 86
	PriorityArray                  PropertyID = 87
	PriorityForWriting             PropertyID = 88
	ProcessIdentifier              PropertyID = 89
	ProgramChange                  PropertyID = 90
	ProgramLocation                PropertyID = 91
	ProgramState                   PropertyID = 92
	ProportionalConstant           PropertyID = 93
	ProportionalConstantUnits      PropertyID = 94
	ProtocolConformanceClass       PropertyID = 95
	ProtocolObjectTypesSupported   PropertyID = 96
	ProtocolServicesSupported      PropertyID = 97
	ProtocolVersion                PropertyID = 98
	ReadOnly                       PropertyID = 99
	ReasonForHalt                  PropertyID = 100
	Recipient                      PropertyID = 101
	RecipientList                  PropertyID = 102
	Reliability                    PropertyID = 103
	RelinquishDefault              PropertyID = 104
	Required                       PropertyID = 105
	Resolution                     PropertyID = 106
	SegmentationSupported          PropertyID = 107
	Setpoint                       PropertyID = 108
	SetpointReference              PropertyID = 109
	StateText                      PropertyID = 110
	StatusFlags                    PropertyID = 111
	SystemStatus                   PropertyID = 112
	TimeDelay                      PropertyID = 113
	TimeOfActiveTimeReset          PropertyID = 114
	TimeOfStateCountReset          PropertyID = 115
	TimeSynchronizationRecipients  PropertyID = 116
	Units                          PropertyID = 117
	UpdateInterval                 PropertyID = 118
	UtcOffset                      PropertyID = 119
	VendorIdentifier               PropertyID = 120
	VendorName                     PropertyID = 121
	VTClassesSupported             PropertyID = 122
	WeeklySchedule                 PropertyID = 123
)

const (
	DescriptionStr = "Description"
	ObjectNameStr  = "ObjectName"
)

// enumMapping should be treated as read only.
var enumMapping = map[string]PropertyID{
	"AckedTransitions":               AckedTransitions,
	"AckRequired":                    AckRequired,
	"Action":                         Action,
	"ActionText":                     ActionText,
	"ActiveText":                     ActiveText,
	"ActiveVTSessions":               ActiveVTSessions,
	"AlarmValue":                     AlarmValue,
	"AlarmValues":                    AlarmValues,
	"All":                            All,
	"AllWritesSuccessful":            AllWritesSuccessful,
	"ApduSegmentTimeout":             ApduSegmentTimeout,
	"ApduTimeout":                    ApduTimeout,
	"ApplicationSoftwareVersion":     ApplicationSoftwareVersion,
	"Archive":                        Archive,
	"Bias":                           Bias,
	"ChangeOfStateCount":             ChangeOfStateCount,
	"ChangeOfStateTime":              ChangeOfStateTime,
	"NotificationClass":              NotificationClass,
	"None":                           None,
	"ControlledVariableReference":    ControlledVariableReference,
	"ControlledVariableUnits":        ControlledVariableUnits,
	"ControlledVariableValue":        ControlledVariableValue,
	"CovIncrement":                   CovIncrement,
	"DateList":                       DateList,
	"DaylightSavingsStatus":          DaylightSavingsStatus,
	"Deadband":                       Deadband,
	"DerivativeConstant":             DerivativeConstant,
	"DerivativeConstantUnits":        DerivativeConstantUnits,
	DescriptionStr:                   Description,
	"DescriptionOfHalt":              DescriptionOfHalt,
	"DeviceAddressBinding":           DeviceAddressBinding,
	"DeviceType":                     DeviceType,
	"EffectivePeriod":                EffectivePeriod,
	"ElapsedActiveTime":              ElapsedActiveTime,
	"ErrorLimit":                     ErrorLimit,
	"EventEnable":                    EventEnable,
	"EventState":                     EventState,
	"EventType":                      EventType,
	"ExceptionSchedule":              ExceptionSchedule,
	"FaultValues":                    FaultValues,
	"FeedbackValue":                  FeedbackValue,
	"FileAccessMethod":               FileAccessMethod,
	"FileSize":                       FileSize,
	"FileType":                       FileType,
	"FirmwareRevision":               FirmwareRevision,
	"HighLimit":                      HighLimit,
	"InactiveText":                   InactiveText,
	"InProcess":                      InProcess,
	"InstanceOf":                     InstanceOf,
	"IntegralConstant":               IntegralConstant,
	"IntegralConstantUnits":          IntegralConstantUnits,
	"IssueConfirmedNotifications":    IssueConfirmedNotifications,
	"LimitEnable":                    LimitEnable,
	"ListOfGroupMembers":             ListOfGroupMembers,
	"ListOfObjectPropertyReferences": ListOfObjectPropertyReferences,
	"ListOfSessionKeys":              ListOfSessionKeys,
	"LocalDate":                      LocalDate,
	"LocalTime":                      LocalTime,
	"Location":                       Location,
	"LowLimit":                       LowLimit,
	"ManipulatedVariableReference":   ManipulatedVariableReference,
	"MaximumOutput":                  MaximumOutput,
	"MaxApduLengthAccepted":          MaxApduLengthAccepted,
	"MaxInfoFrames":                  MaxInfoFrames,
	"MaxMaster":                      MaxMaster,
	"MaxPresValue":                   MaxPresValue,
	"MinimumOffTime":                 MinimumOffTime,
	"MinimumOnTime":                  MinimumOnTime,
	"MinimumOutput":                  MinimumOutput,
	"MinPresValue":                   MinPresValue,
	"ModelName":                      ModelName,
	"ModificationDate":               ModificationDate,
	"NotifyType":                     NotifyType,
	"NumberOfAPDURetries":            NumberOfAPDURetries,
	"NumberOfStates":                 NumberOfStates,
	"ObjectIdentifier":               ObjectIdentifier,
	"ObjectList":                     ObjectList,
	ObjectNameStr:                    ObjectName,
	"ObjectPropertyReference":        ObjectPropertyReference,
	"ObjectType":                     ObjectType,
	"Optional":                       Optional,
	"OutOfService":                   OutOfService,
	"OutputUnits":                    OutputUnits,
	"EventParameters":                EventParameters,
	"Polarity":                       Polarity,
	"PresentValue":                   PresentValue,
	"Priority":                       Priority,
	"PriorityArray":                  PriorityArray,
	"PriorityForWriting":             PriorityForWriting,
	"ProcessIdentifier":              ProcessIdentifier,
	"ProgramChange":                  ProgramChange,
	"ProgramLocation":                ProgramLocation,
	"ProgramState":                   ProgramState,
	"ProportionalConstant":           ProportionalConstant,
	"ProportionalConstantUnits":      ProportionalConstantUnits,
	"ProtocolConformanceClass":       ProtocolConformanceClass,
	"ProtocolObjectTypesSupported":   ProtocolObjectTypesSupported,
	"ProtocolServicesSupported":      ProtocolServicesSupported,
	"ProtocolVersion":                ProtocolVersion,
	"ReadOnly":                       ReadOnly,
	"ReasonForHalt":                  ReasonForHalt,
	"Recipient":                      Recipient,
	"RecipientList":                  RecipientList,
	"Reliability":                    Reliability,
	"RelinquishDefault":              RelinquishDefault,
	"Required":                       Required,
	"Resolution":                     Resolution,
	"SegmentationSupported":          SegmentationSupported,
	"Setpoint":                       Setpoint,
	"SetpointReference":              SetpointReference,
	"StateText":                      StateText,
	"StatusFlags":                    StatusFlags,
	"SystemStatus":                   SystemStatus,
	"TimeDelay":                      TimeDelay,
	"TimeOfActiveTimeReset":          TimeOfActiveTimeReset,
	"TimeOfStateCountReset":          TimeOfStateCountReset,
	"TimeSynchronizationRecipients":  TimeSynchronizationRecipients,
	"Units":                          Units,
	"UpdateInterval":                 UpdateInterval,
	"UtcOffset":                      UtcOffset,
	"VendorIdentifier":               VendorIdentifier,
	"VendorName":                     VendorName,
	"VTClassesSupported":             VTClassesSupported,
	"WeeklySchedule":                 WeeklySchedule,
}

// strMapping is a human readable printing of the priority
var strMapping = map[PropertyID]string{
	AckedTransitions:               "Acked Transitions",
	AckRequired:                    "Ack Required",
	Action:                         "Action",
	ActionText:                     "Action Text",
	ActiveText:                     "Active Text",
	ActiveVTSessions:               "Active VT Sessions",
	AlarmValue:                     "Alarm Value",
	AlarmValues:                    "Alarm Values",
	All:                            "All",
	AllWritesSuccessful:            "All Writes Successful",
	ApduSegmentTimeout:             "Apdu Segment Timeout",
	ApduTimeout:                    "Apdu Timeout",
	ApplicationSoftwareVersion:     "Application Software Version",
	Archive:                        "Archive",
	Bias:                           "Bias",
	ChangeOfStateCount:             "Change Of State Count",
	ChangeOfStateTime:              "Change Of State Time",
	NotificationClass:              "Notification Class",
	None:                           "None",
	ControlledVariableReference:    "Controlled Variable Reference",
	ControlledVariableUnits:        "Controlled Variable Units",
	ControlledVariableValue:        "Controlled Variable Value",
	CovIncrement:                   "Cov Increment",
	DateList:                       "Date List",
	DaylightSavingsStatus:          "Daylight Savings Status",
	Deadband:                       "Deadband",
	DerivativeConstant:             "Derivative Constant",
	DerivativeConstantUnits:        "Derivative Constant Units",
	Description:                    "Description",
	DescriptionOfHalt:              "Description Of Halt",
	DeviceAddressBinding:           "Device Address Binding",
	DeviceType:                     "Device Type",
	EffectivePeriod:                "Effective Period",
	ElapsedActiveTime:              "Elapsed Active Time",
	ErrorLimit:                     "Error Limit",
	EventEnable:                    "Event Enable",
	EventState:                     "Event State",
	EventType:                      "Event Type",
	ExceptionSchedule:              "Exception Schedule",
	FaultValues:                    "Fault Values",
	FeedbackValue:                  "Feedback Value",
	FileAccessMethod:               "File Access Method",
	FileSize:                       "File Size",
	FileType:                       "File Type",
	FirmwareRevision:               "Firmware Revision",
	HighLimit:                      "High Limit",
	InactiveText:                   "Inactive Text",
	InProcess:                      "In Process",
	InstanceOf:                     "Instance Of",
	IntegralConstant:               "Integral Constant",
	IntegralConstantUnits:          "Integral Constant Units",
	IssueConfirmedNotifications:    "Issue Confirmed Notifications",
	LimitEnable:                    "Limit Enable",
	ListOfGroupMembers:             "List Of Group Members",
	ListOfObjectPropertyReferences: "List Of Object Property References",
	ListOfSessionKeys:              "List Of Session Keys",
	LocalDate:                      "Local Date",
	LocalTime:                      "Local Time",
	Location:                       "Location",
	LowLimit:                       "Low Limit",
	ManipulatedVariableReference:   "Manipulated Variable Reference",
	MaximumOutput:                  "Maximum Output",
	MaxApduLengthAccepted:          "Max Apdu Length Accepted",
	MaxInfoFrames:                  "Max Info Frames",
	MaxMaster:                      "Max Master",
	MaxPresValue:                   "Max Pres Value",
	MinimumOffTime:                 "Minimum Off Time",
	MinimumOnTime:                  "Minimum On Time",
	MinimumOutput:                  "Minimum Output",
	MinPresValue:                   "Min Pres Value",
	ModelName:                      "Model Name",
	ModificationDate:               "Modification Date",
	NotifyType:                     "Notify Type",
	NumberOfAPDURetries:            "Number Of A P D U Retries",
	NumberOfStates:                 "Number Of States",
	ObjectIdentifier:               "Object Identifier",
	ObjectList:                     "Object List",
	ObjectName:                     "Object Name",
	ObjectPropertyReference:        "Object Property Reference",
	ObjectType:                     "Object Type",
	Optional:                       "Optional",
	OutOfService:                   "Out Of Service",
	OutputUnits:                    "Output Units",
	EventParameters:                "Event Parameters",
	Polarity:                       "Polarity",
	PresentValue:                   "Present Value",
	Priority:                       "Priority",
	PriorityArray:                  "Priority Array",
	PriorityForWriting:             "Priority For Writing",
	ProcessIdentifier:              "Process Identifier",
	ProgramChange:                  "Program Change",
	ProgramLocation:                "Program Location",
	ProgramState:                   "Program State",
	ProportionalConstant:           "Proportional Constant",
	ProportionalConstantUnits:      "Proportional Constant Units",
	ProtocolConformanceClass:       "Protocol Conformance Class",
	ProtocolObjectTypesSupported:   "Protocol Object Types Supported",
	ProtocolServicesSupported:      "Protocol Services Supported",
	ProtocolVersion:                "Protocol Version",
	ReadOnly:                       "Read Only",
	ReasonForHalt:                  "Reason For Halt",
	Recipient:                      "Recipient",
	RecipientList:                  "Recipient List",
	Reliability:                    "Reliability",
	RelinquishDefault:              "Relinquish Default",
	Required:                       "Required",
	Resolution:                     "Resolution",
	SegmentationSupported:          "Segmentation Supported",
	Setpoint:                       "Setpoint",
	SetpointReference:              "Setpoint Reference",
	StateText:                      "State Text",
	StatusFlags:                    "Status Flags",
	SystemStatus:                   "System Status",
	TimeDelay:                      "Time Delay",
	TimeOfActiveTimeReset:          "Time Of Active Time Reset",
	TimeOfStateCountReset:          "Time Of State Count Reset",
	TimeSynchronizationRecipients:  "Time Synchronization Recipients",
	Units:                          "Units",
	UpdateInterval:                 "Update Interval",
	UtcOffset:                      "Utc Offset",
	VendorIdentifier:               "Vendor Identifier",
	VendorName:                     "Vendor Name",
	VTClassesSupported:             "VT Classes Supported",
	WeeklySchedule:                 "Weekly Schedule",
}

// listOfKeys should be treated as read only after init
var listOfKeys []string

func init() {
	listOfKeys = make([]string, len(enumMapping))
	i := 0
	for k := range enumMapping {
		listOfKeys[i] = k
		i++
	}
}

func Keys() map[string]PropertyID {
	// A copy is made since we do not want outside packages editing our keys by
	// accident
	keys := make(map[string]PropertyID)
	for k, v := range enumMapping {
		keys[k] = v
	}
	return keys
}

func Get(s string) (PropertyID, error) {
	if v, ok := enumMapping[s]; ok {
		return v, nil
	}
	err := fmt.Errorf("%s is not a valid property", s)
	return 0, err
}

// String returns a human readible string of the given property
func String(prop PropertyID) string {
	s, ok := strMapping[prop]
	if !ok {
		return "Unknown"
	}
	return fmt.Sprintf("%s (%d)", s, prop)
}

// The bool in the map doesn't actually matter since it won't be used.
var deviceProperties = map[PropertyID]bool{
	ObjectList: true,
}

func IsDeviceProperty(id PropertyID) bool {
	_, ok := deviceProperties[id]
	return ok
}

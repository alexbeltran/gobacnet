package types

// DO NOT CHANGE THE ORDER OF THESE!
type ServicesSupported struct {
	AcknowledgeAlarm             bool
	ConfirmedCOVNotification     bool
	ConfirmedEventNotification   bool
	GetAlarmSummary              bool
	GetEnrollmentSummary         bool
	SubscribeCOV                 bool
	AtomicReadFile               bool
	AtomicWriteFile              bool
	AddListElement               bool
	RemoveListElement            bool
	CreateObject                 bool
	DeleteObject                 bool
	ReadProperty                 bool
	ReadPropertyConditional      bool
	ReadPropertyMultiple         bool
	WriteProperty                bool
	WritePropertyMultiple        bool
	DeviceCommunicationControl   bool
	ConfirmedPrivateTransfer     bool
	ConfirmedTextMessage         bool
	ReinitializeDevice           bool
	VtOpen                       bool
	VtClose                      bool
	VtData                       bool
	Authenticate                 bool
	RequestKey                   bool
	IAm                          bool
	IHave                        bool
	UnconfirmedCOVNotification   bool
	UnconfirmedEventNotification bool
	UnconfirmedPrivateTransfer   bool
	UnconfirmedTextMessage       bool
	TimeSynchronization          bool
	WhoHas                       bool
	WhoIs                        bool
}

// DO NOT CHANGE THE ORDER OF THESE
type TypesSupported struct {
	AnalogInput       bool
	AnalogOutput      bool
	AnalogValue       bool
	BinaryInput       bool
	BinaryOutput      bool
	BinaryValue       bool
	Calendar          bool
	Command           bool
	Device            bool
	EventEnrollment   bool
	File              bool
	Group             bool
	Loop              bool
	MultiStateInput   bool
	MultiStateOutput  bool
	NotificationClass bool
	Program           bool
	Schedule          bool
	Averaging         bool
	MultiStateValue   bool
	TrendLog          bool
}

// Generic BitString
type BitString struct {
	Bits []bool
}

func (a *BitString) ByteLengthOfBitsAndRemainder() int {
	length := len(a.Bits)
	if length == 0 {
		return 1
	}

	return (length-1)/8 + 2
}

func (a *BitString) String() string {
	s := ""
	for i := 0; i < len(a.Bits); i++ {
		if a.Bits[i] {
			s += "T"
		} else {
			s += "F"
		}
	}

	return s
}

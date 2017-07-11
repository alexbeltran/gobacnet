package types

type ServiceConfirmed uint8

const (
	/* Alarm and Event Services */
	ServiceConfirmedAcknowledgeAlarm     ServiceConfirmed = 0
	ServiceConfirmedCOVNotification      ServiceConfirmed = 1
	ServiceConfirmedEventNotification    ServiceConfirmed = 2
	ServiceConfirmedGetAlarmSummary      ServiceConfirmed = 3
	ServiceConfirmedGetEnrollmentSummary ServiceConfirmed = 4
	ServiceConfirmedGetEventInformation  ServiceConfirmed = 29
	ServiceConfirmedSubscribeCOV         ServiceConfirmed = 5
	ServiceConfirmedSubscribeCOVProperty ServiceConfirmed = 28
	ServiceConfirmedLifeSafetyOperation  ServiceConfirmed = 27
	/* File Access Services */
	ServiceConfirmedAtomicReadFile  ServiceConfirmed = 6
	ServiceConfirmedAtomicWriteFile ServiceConfirmed = 7
	/* Object Access Services */
	ServiceConfirmedAddListElement      ServiceConfirmed = 8
	ServiceConfirmedRemoveListElement   ServiceConfirmed = 9
	ServiceConfirmedCreateObject        ServiceConfirmed = 10
	ServiceConfirmedDeleteObject        ServiceConfirmed = 11
	ServiceConfirmedReadProperty        ServiceConfirmed = 12
	ServiceConfirmedReadPropConditional ServiceConfirmed = 13
	ServiceConfirmedReadPropMultiple    ServiceConfirmed = 14
	ServiceConfirmedReadRange           ServiceConfirmed = 26
	ServiceConfirmedWriteProperty       ServiceConfirmed = 15
	ServiceConfirmedWritePropMultiple   ServiceConfirmed = 16
	/* Remote Device Management Services */
	ServiceConfirmedDeviceCommunicationControl ServiceConfirmed = 17
	ServiceConfirmedPrivateTransfer            ServiceConfirmed = 18
	ServiceConfirmedTextMessage                ServiceConfirmed = 19
	ServiceConfirmedReinitializeDevice         ServiceConfirmed = 20
	/* Virtual Terminal Services */
	ServiceConfirmedVTOpen  ServiceConfirmed = 21
	ServiceConfirmedVTClose ServiceConfirmed = 22
	ServiceConfirmedVTData  ServiceConfirmed = 23
	/* Security Services */
	ServiceConfirmedAuthenticate ServiceConfirmed = 24
	ServiceConfirmedRequestKey   ServiceConfirmed = 25
	/* Services added after 1995 */
	/* readRange (26) see Object Access Services */
	/* lifeSafetyOperation (27) see Alarm and Event Services */
	/* subscribeCOVProperty (28) see Alarm and Event Services */
	/* getEventInformation (29) see Alarm and Event Services */
	maxBACnetConfirmedService ServiceConfirmed = 30
)

// APDU - Application Protocol Data Unit
type APDU struct {
	SegmentedMessage          bool
	MoreFollows               bool
	SegmentedResponseAccepted bool
	MaxSegs                   uint
	MaxApdu                   uint
	InvokeId                  uint8
	Sequence                  uint8
	WindowNumber              uint8
	Service                   ServiceConfirmed
}

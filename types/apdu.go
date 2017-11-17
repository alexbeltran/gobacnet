/*Copyright (C) 2017 Alex Beltran

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to:
The Free Software Foundation, Inc.
59 Temple Place - Suite 330
Boston, MA  02111-1307, USA.

As a special exception, if other files instantiate templates or
use macros or inline functions from this file, or you compile
this file and link it with other works to produce a work based
on this file, this file does not by itself cause the resulting
work to be covered by the GNU General Public License. However
the source code for this file must still be made available in
accordance with section (3) of the GNU General Public License.

This exception does not invalidate any other reasons why a work
based on this file might be covered by the GNU General Public
License.
*/

package types

import "fmt"

type ServiceConfirmed uint8
type ServiceUnconfirmed uint8

const MaxAPDUOverIP = 1476
const MaxAPDU = MaxAPDUOverIP

const (
	ServiceUnconfirmedIAm               ServiceUnconfirmed = 0
	ServiceUnconfirmedIHave             ServiceUnconfirmed = 1
	ServiceUnconfirmedCOVNotification   ServiceUnconfirmed = 2
	ServiceUnconfirmedEventNotification ServiceUnconfirmed = 3
	ServiceUnconfirmedPrivateTransfer   ServiceUnconfirmed = 4
	ServiceUnconfirmedTextMessage       ServiceUnconfirmed = 5
	ServiceUnconfirmedTimeSync          ServiceUnconfirmed = 6
	ServiceUnconfirmedWhoHas            ServiceUnconfirmed = 7
	ServiceUnconfirmedWhoIs             ServiceUnconfirmed = 8
	ServiceUnconfirmedUTCTimeSync       ServiceUnconfirmed = 9
	ServiceUnconfirmedWriteGroup        ServiceUnconfirmed = 10
	/* Other services to be added as they are defined. */
	/* All choice values in this production are reserved */
	/* for definition by ASHRAE. */
	/* Proprietary extensions are made by using the */
	/* UnconfirmedPrivateTransfer service. See Clause 23. */
	MaxServiceUnconfirmed ServiceUnconfirmed = 11
)

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
	DataType                  PDUType
	SegmentedMessage          bool
	MoreFollows               bool
	SegmentedResponseAccepted bool
	MaxSegs                   uint
	MaxApdu                   uint
	InvokeId                  uint8
	Sequence                  uint8
	WindowNumber              uint8
	Service                   ServiceConfirmed
	UnconfirmedService        ServiceUnconfirmed
	Error                     struct {
		Class uint32
		Code  uint32
	}

	// This is the raw data passed based on the service
	RawData []byte
}

// pduType encomposes all valid pdus.
type PDUType uint8

// pdu requests
const (
	ConfirmedServiceRequest   PDUType = 0
	UnconfirmedServiceRequest PDUType = 0x10
	ComplexAck                PDUType = 0x30
	SegmentAck                PDUType = 0x40
	Error                     PDUType = 0x50
	Reject                    PDUType = 0x60
	Abort                     PDUType = 0x70
)

// IsConfirmedServiceRequest checks to see if the APDU is in the list of known services
func (a *APDU) IsConfirmedServiceRequest() bool {
	return (0xF0 & a.DataType) == ConfirmedServiceRequest
}

func (s *ServiceConfirmed) String() string {
	switch *s {
	default:
		return fmt.Sprintf("Unknown %d", uint(*s))
	}
}

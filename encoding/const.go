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

package encoding

const MaxInstance = 0x3FFFFF
const InstanceBits = 22
const MaxPropertyID = 4194303

const initialTagPos = 0

const (
	size8  = 1
	size16 = 2
	size24 = 3
	size32 = 4
)

const (
	flag16bit uint8 = 254
	flag32bit uint8 = 255
)

// pduType encomposes all valid pdus.
type pduType uint8

// pdu requests
const (
	confirmedServiceRequest pduType = 0
	complexAck              pduType = 0x30
)

type serviceConfirmed uint8

const (
	/* Alarm and Event Services */
	serviceConfirmedAcknowledgeAlarm     serviceConfirmed = 0
	serviceConfirmedCOVNotification      serviceConfirmed = 1
	serviceConfirmedEventNotification    serviceConfirmed = 2
	serviceConfirmedGetAlarmSummary      serviceConfirmed = 3
	serviceConfirmedGetEnrollmentSummary serviceConfirmed = 4
	serviceConfirmedGetEventInformation  serviceConfirmed = 29
	serviceConfirmedSubscribeCOV         serviceConfirmed = 5
	serviceConfirmedSubscribeCOVProperty serviceConfirmed = 28
	serviceConfirmedLifeSafetyOperation  serviceConfirmed = 27
	/* File Access Services */
	serviceConfirmedAtomicReadFile  serviceConfirmed = 6
	serviceConfirmedAtomicWriteFile serviceConfirmed = 7
	/* Object Access Services */
	serviceConfirmedAddListElement      serviceConfirmed = 8
	serviceConfirmedRemoveListElement   serviceConfirmed = 9
	serviceConfirmedCreateObject        serviceConfirmed = 10
	serviceConfirmedDeleteObject        serviceConfirmed = 11
	serviceConfirmedReadProperty        serviceConfirmed = 12
	serviceConfirmedReadPropConditional serviceConfirmed = 13
	serviceConfirmedReadPropMultiple    serviceConfirmed = 14
	serviceConfirmedReadRange           serviceConfirmed = 26
	serviceConfirmedWriteProperty       serviceConfirmed = 15
	serviceConfirmedWritePropMultiple   serviceConfirmed = 16
	/* Remote Device Management Services */
	serviceConfirmedDeviceCommunicationControl serviceConfirmed = 17
	serviceConfirmedPrivateTransfer            serviceConfirmed = 18
	serviceConfirmedTextMessage                serviceConfirmed = 19
	serviceConfirmedReinitializeDevice         serviceConfirmed = 20
	/* Virtual Terminal Services */
	serviceConfirmedVTOpen  serviceConfirmed = 21
	serviceConfirmedVTClose serviceConfirmed = 22
	serviceConfirmedVTData  serviceConfirmed = 23
	/* Security Services */
	serviceConfirmedAuthenticate serviceConfirmed = 24
	serviceConfirmedRequestKey   serviceConfirmed = 25
	/* Services added after 1995 */
	/* readRange (26) see Object Access Services */
	/* lifeSafetyOperation (27) see Alarm and Event Services */
	/* subscribeCOVProperty (28) see Alarm and Event Services */
	/* getEventInformation (29) see Alarm and Event Services */
	maxBACnetConfirmedService serviceConfirmed = 30
)

const ArrayAll = 0xFFFFFFFF

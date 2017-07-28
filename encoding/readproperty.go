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

import (
	"fmt"

	bactype "github.com/alexbeltran/gobacnet/types"
)

func (e *Encoder) readPropertyHeader(tagPos uint8, data bactype.ReadPropertyData) uint8 {
	// Tag - Object Type and Instance
	if data.ObjectType < MaxObject {
		e.contextObjectID(tagPos, data.ObjectType, data.ObjectInstance)
	}
	tagPos++

	// Tag - Object Property
	if data.ObjectProperty < MaxPropertyID {
		e.contextEnumerated(tagPos, data.ObjectProperty)
	}
	tagPos++

	// Optional Tag - Array Index
	if data.ArrayIndex != ArrayAll {
		e.contextUnsigned(tagPos, data.ArrayIndex)
	}
	tagPos++
	return tagPos
}

// ReadProperty is a service request to read a property that is passed.
func (e *Encoder) ReadProperty(invokeID uint8, data bactype.ReadPropertyData) error {
	// PDU Type
	a := bactype.APDU{
		DataType:         bactype.ConfirmedServiceRequest,
		Service:          bactype.ServiceConfirmedReadProperty,
		MaxSegs:          0,
		MaxApdu:          MaxAPDU,
		InvokeId:         invokeID,
		SegmentedMessage: false,
	}
	e.APDU(a)
	e.readPropertyHeader(initialTagPos, data)
	return e.Error()
}

// ReadPropertyAck is the response made to a ReadProperty service request.
func (e *Encoder) ReadPropertyAck(invokeID uint8, data bactype.ReadPropertyData) {
	// PDU Type
	a := bactype.APDU{
		DataType: bactype.ComplexAck,
		Service:  bactype.ServiceConfirmedReadProperty,
		MaxSegs:  0,
		MaxApdu:  MaxAPDU,
		InvokeId: invokeID,
	}
	e.APDU(a)

	tagID := e.readPropertyHeader(initialTagPos, data)

	e.openingTag(tagID)
	tagID++
	for _, d := range data.ApplicationData {
		e.write(d)
	}
	e.closingTag(tagID)
}

func (d *Decoder) ReadProperty(data *bactype.ReadPropertyData) error {
	// Must have at least 7 bytes
	if d.buff.Len() < 7 {
		return fmt.Errorf("Missing parameters")
	}

	// Tag 0: Object ID
	tag, meta := d.tagNumber()

	var expectedTag uint8
	if tag != expectedTag {
		return &ErrorIncorrectTag{expectedTag, tag}
	}
	expectedTag++

	var objectType uint16
	var instance uint32
	var property uint32
	if !meta.isContextSpecific() {
		return fmt.Errorf("Tag %d should be context specific. %x", tag, meta)
	}
	objectType, instance = d.objectId()

	// Tag 1: Property ID
	tag, meta = d.tagNumber()
	if tag != expectedTag {
		return &ErrorIncorrectTag{expectedTag, tag}
	}
	expectedTag++

	lenValue := d.value(meta)
	property = d.enumerated(int(lenValue))
	if d.len() != 0 {
		tag, meta = d.tagNumber()
	}

	var arrIndex uint32
	// Check to see if we still have bytes to read.
	if d.buff.Len() != 0 || tag >= 2 {
		// If we do then that means we are reading the optional argument,
		// arra length

		// Tag 2: Array Length (OPTIONAL)
		var lenValue uint32
		lenValue = d.value(meta)

		var openTag uint8
		// I tried to not use magic numbers but it doesn't look like it can be avoid
		// If the attag we receive is a tag of 2 then set the value
		if tag == 2 {
			arrIndex = d.unsigned(int(lenValue))
			if d.len() > 0 {
				openTag, meta = d.tagNumber()
			}
		} else {
			openTag = tag
			arrIndex = ArrayAll
		}

		if openTag == 3 {
			// We subtract one to ignore the closing tag.
			data.ApplicationDataLen = d.buff.Len() - 1
			data.ApplicationData = make([]byte, d.buff.Len()-1)
			d.decode(data.ApplicationData)
		}
	} else {
		arrIndex = ArrayAll
	}

	// We now assemble all the values that we have read above
	data.ObjectInstance = instance
	data.ObjectType = objectType
	data.ObjectProperty = property
	data.ArrayIndex = arrIndex

	return d.Error()
}

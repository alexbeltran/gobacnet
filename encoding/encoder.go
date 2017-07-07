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
	"bytes"
	"encoding/binary"

	bactype "github.com/alexbeltran/gobacnet/types"
)

var EncodingEndian binary.ByteOrder = binary.BigEndian

type Encoder struct {
	buff *bytes.Buffer
	err  error
}

func NewEncoder() *Encoder {
	e := Encoder{
		buff: new(bytes.Buffer),
		err:  nil,
	}
	return &e
}

func (e *Encoder) Error() error {
	return e.err
}

func (e *Encoder) Bytes() []byte {
	return e.buff.Bytes()
}

func (e *Encoder) write(p interface{}) {
	if e.err != nil {
		return
	}
	e.err = binary.Write(e.buff, EncodingEndian, p)
}

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
func (e *Encoder) ReadProperty(invokeID uint8, data bactype.ReadPropertyData) {
	// PDU Type
	e.write(confirmedServiceRequest)
	e.write(encodeMaxSegsMaxApdu(0, maxApdu))
	e.write(invokeID)
	e.write(ReadPropertyService)
	e.readPropertyHeader(initialTagPos, data)
	return
}

func (e *Encoder) serviceConfirmed(in serviceConfirmed) {
	e.write(in)
}

// ReadPropertyAck is the response made to a ReadProperty service request.
func (e *Encoder) ReadPropertyAck(invokeID uint8, data bactype.ReadPropertyData) {
	// PDU Type
	e.write(complexAck)
	e.write(invokeID)
	e.serviceConfirmed(serviceConfirmedReadProperty)

	tagID := e.readPropertyHeader(initialTagPos, data)

	e.openingTag(tagID)
	tagID++
	for _, d := range data.ApplicationData {
		e.write(d)
	}
	e.closingTag(tagID)
}

func (e *Encoder) contextObjectID(tagNum uint8, objectType uint16, instance uint32) {
	/* length of object id is 4 octets, as per 20.2.14 */
	e.tag(tagNum, true, 4)
	e.objectId(objectType, instance)
}

// Write opening tag to the system
func (e *Encoder) openingTag(num uint8) {
	var meta tagMeta
	meta.setOpening()
	e.tagNum(meta, num)
}

func (e *Encoder) closingTag(num uint8) {
	var meta tagMeta
	meta.setClosing()
	e.tagNum(meta, num)
}

// pretags
func (e *Encoder) tagNum(meta tagMeta, num uint8) {
	t := uint8(meta)
	if num <= 14 {
		t |= (num << 4)
		e.write(t)

		// We don't have enough space so make it in a new byte
	} else {
		t |= 0xF0
		e.write(t)
		e.write(num)
	}
}

func (e *Encoder) tag(tagNum uint8, contextSpecific bool, lenValueType uint32) {
	var t uint8
	if contextSpecific {
		var meta tagMeta
		meta.setContextSpecific()
		t = uint8(meta)
	}

	if lenValueType <= 4 {
		t |= uint8(lenValueType)
	} else {
		t |= 5
	}

	// We have enough room to put it with the last value
	if tagNum <= 14 {
		t |= (tagNum << 4)
		e.write(t)

		// We don't have enough space so make it in a new byte
	} else {
		t |= 0xF0
		e.write(t)
		e.write(tagNum)
	}
	if lenValueType > 4 {
		// Depending on the length, we will either write it as an 8 bit, 32 bit, or 64 bit integer
		if lenValueType <= 253 {
			e.write(uint8(lenValueType))
		} else if lenValueType <= 65535 {
			e.write(flag16bit)
			e.write(uint16(lenValueType))
		} else {
			e.write(flag32bit)
			e.write(lenValueType)
		}
	}
}

/* from clause 20.2.14 Encoding of an Object Identifier Value
returns the number of apdu bytes consumed */
func (e *Encoder) objectId(objectType uint16, instance uint32) {
	var value uint32
	value = ((uint32(objectType) & MaxObject) << InstanceBits) | (instance & MaxInstance)
	e.write(value)
}

func (e *Encoder) contextEnumerated(tagNumber uint8, value uint32) {
	e.contextUnsigned(tagNumber, value)
}

func (e *Encoder) contextUnsigned(tagNumber uint8, value uint32) {
	len := valueLength(value)
	e.tag(tagNumber, true, uint32(len))
	e.unsigned(value)
}

func (e *Encoder) enumerated(value uint32) {
	e.unsigned(value)
}

// weird, huh?
func (e *Encoder) unsigned24(value uint32) {
	e.write(uint8((value & 0xFF0000) >> 16))
	e.write(uint8((value & 0x00FF00) >> 8))
	e.write(uint8(value & 0x0000FF))

}

func (e *Encoder) unsigned(value uint32) {
	if value < 0x100 {
		e.write(uint8(value))
	} else if value < 0x10000 {
		e.write(uint16(value))
	} else if value < 0x1000000 {
		// Really!? 24 bits?
		e.unsigned24(value)
	} else {
		e.write(value)
	}
}

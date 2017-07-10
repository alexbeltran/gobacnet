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
	"fmt"

	bactype "github.com/alexbeltran/gobacnet/types"
)

// Decoder used
type Decoder struct {
	buff       *bytes.Buffer
	err        error
	tagCounter int
}

func (d *Decoder) len() int {
	return d.buff.Len()
}
func NewDecoder(b []byte) *Decoder {
	return &Decoder{
		bytes.NewBuffer(b),
		nil,
		0,
	}
}

func (d *Decoder) Error() error {
	return d.err
}

func (d *Decoder) decode(data interface{}) {
	// Only decode if there have been no errors so far
	if d.err != nil {
		return
	}
	d.err = binary.Read(d.buff, EncodingEndian, data)
}

func (d *Decoder) tagCheck(inTag uint8) {
	if d.tagCounter != int(inTag) {
		d.err = fmt.Errorf("Mismatch in tag id. Tag ID should be %d but is %d", d.tagCounter, inTag)
	}
}

func (d *Decoder) tagIncr() {
	d.tagCounter++
}

func (d *Decoder) ReadProperty() (data bactype.ReadPropertyData, err error) {
	// Must have at least 7 bytes
	if d.buff.Len() < 7 {
		err = fmt.Errorf("Missing parameters")
		return
	}

	// Tag 0: Object ID
	tag, meta := d.tagNumber()
	d.tagCheck(tag)
	d.tagIncr()

	if !meta.isContextSpecific() {
		err = fmt.Errorf("Tag %d should be context specific. %x", tag, meta)
		return
	}

	objectType, instance := d.objectId()

	// Tag 1: Property ID
	tag, _, lenValue := d.tagNumberAndValue()
	d.tagCheck(tag)
	d.tagIncr()

	property := d.enumerated(int(lenValue))

	var arrIndex uint32
	// Check to see if we still have bytes to read.
	if d.buff.Len() != 0 {
		// If we do then that means we are reading the optional argument,
		// arra length

		// Tag 2: Array Length (OPTIONAL)
		tag, meta, lenValue = d.tagNumberAndValue()

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

		if openTag == 3 && meta.isOpening() {
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

	err = d.Error()
	return
}

// contexTag decoder

// Returns both a tag and additional meta data stored in this byte. If it is of
// extended type, then that means that the entire first byte is metadata, else
// the firrst 4 bytes store the tag
func (d *Decoder) tagNumber() (tag uint8, meta tagMeta) {
	// Read the first value
	d.decode(&meta)
	if meta.isExtendedTagNumber() {
		d.decode(&tag)
		return tag, meta
	}
	return uint8(meta) >> 4, meta
}

func (d *Decoder) tagNumberAndValue() (tag uint8, meta tagMeta, value uint32) {
	tag, meta = d.tagNumber()
	if meta.isExtendedValue() {
		var val uint8
		d.decode(&val)
		// Tagged as an uint32
		if val == flag32bit {
			var parse uint32
			d.decode(&parse)
			return tag, meta, parse

			// Tagged as a uint16
		} else if val == flag16bit {
			var parse uint16
			d.decode(&parse)
			return tag, meta, uint32(parse)

			// No tag, it must be a uint8
		} else {
			return tag, meta, uint32(val)
		}
	} else if meta.isOpening() || meta.isClosing() {
		return tag, meta, 0
	}
	// It must be a non extended/small value
	// Note this is a mask of the last 3 bits
	return tag, meta, uint32(meta & 0x07)
}

func (d *Decoder) objectId() (objectType uint16, instance uint32) {
	var value uint32
	d.decode(&value)
	objectType = uint16((value >> InstanceBits) & MaxObject)
	instance = value & MaxInstance
	return
}

func (d *Decoder) enumerated(len int) uint32 {
	return d.unsigned(len)
}

func (d *Decoder) unsigned24() uint32 {
	var a, b, c uint8
	d.decode(&a)
	d.decode(&b)
	d.decode(&c)

	var x uint32
	x = uint32((uint32(a) << 16) & 0x00ff0000)
	x |= uint32((uint32(b) << 8) & 0x0000ff00)
	x |= uint32(uint32(c) & 0x000000ff)
	return x
}

func (d *Decoder) unsigned(length int) uint32 {
	switch length {
	case size8:
		var val uint8
		d.decode(&val)
		return uint32(val)
	case size16:
		var val uint16
		d.decode(&val)
		return uint32(val)
	case size24:
		return d.unsigned24()
	case size32:
		var val uint32
		d.decode(&val)
		return val
	default:
		return 0
	}
}

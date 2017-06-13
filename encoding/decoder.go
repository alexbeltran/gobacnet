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
)

// Decoder used
type Decoder struct {
	buff *bytes.Buffer
	err  error
}

func NewDecoder(b []byte) *Decoder {
	return &Decoder{
		bytes.NewBuffer(b),
		nil,
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
	d.err = binary.Read(d.buff, encodingEndian, data)
}

func (d *Decoder) readProperty() (err error, data ReadPropertyData) {
	// tag counter should be incremented every time a tag is read. this is important
	// to make sure that the correct order is read.
	var tagCounter uint8 = 0

	// Tag checker checks the passed tag and increments the tag counter
	tagCheck := func(inTag uint8) error {
		if tagCounter != inTag {
			return fmt.Errorf("Mismatch in tag id. Tag ID should be %d but is %d", tagCounter, inTag)
		}
		tagCounter++
		return nil
	}

	// Must have at least 7 bytes
	if d.buff.Len() < 7 {
		err = fmt.Errorf("Missing parameters")
		return
	}

	// Tag 0: Object ID
	tag, meta := d.tagNumber()

	if !isContextSpecific(meta) {
		err = fmt.Errorf("Tag %d should be context specific. %x", tag, meta)
		return
	}

	if err = tagCheck(tag); err != nil {
		return
	}

	objectType, instance := d.objectId()

	// Tag 1: Property ID
	tag, lenValue := d.tagNumberAndValue()
	if err = tagCheck(tag); err != nil {
		return
	}
	property := d.enumerated(int(lenValue))

	var arrIndex uint32
	// Check to see if we still have bytes to read.
	if d.buff.Len() != 0 {
		// If we do then that means we are reading the optional argument,
		// arra length

		// Tag 2: Array Length (OPTIONAL)
		tag, lenValue = d.tagNumberAndValue()
		if err = tagCheck(tag); err != nil {
			return
		}
		arrIndex = d.unsigned(int(lenValue))
	} else {
		arrIndex = ArrayAll
	}

	if d.buff.Len() > 0 {
		err = fmt.Errorf("there are too man arguments included in passed buffer")
		return
	}

	// We now assemble all the values that we have read above
	data.ObjectInstance = instance
	data.ObjectType = objectType
	data.ObjectProperty = property
	data.ArrayIndex = arrIndex
	return
}

// contexTag decoder

// Returns both a tag and additional meta data stored in this byte. If it is of
// extended type, then that means that the entire first byte is metadata, else
// the firrst 4 bytes store the tag
func (d *Decoder) tagNumber() (tag uint8, meta uint8) {
	// Read the first value
	d.decode(&meta)
	if isExtendedTagNumber(meta) {

		d.decode(&tag)
		return tag, meta
	}
	return meta >> 4, meta
}

func (d *Decoder) tagNumberAndValue() (tag uint8, value uint32) {
	tag, meta := d.tagNumber()
	if isExtendedTagNumber(meta) {
		var val uint8
		d.decode(&val)
		// Tagged as an uint32
		if val == 255 {
			var parse uint32
			d.decode(&parse)
			return tag, parse

			// Tagged as a uint16
		} else if val == 254 {
			var parse uint16
			d.decode(&parse)
			return tag, uint32(parse)

			// No tag, it must be a uint8
		} else {
			return tag, uint32(val)
		}
	} else if isOpeningTag(meta) || isClosingTag(meta) {
		return tag, 0
	}
	// It must be a non extended/small value
	// Note this is a mask of the last 3 bits
	return tag, uint32(meta & 0x07)
}

func (d *Decoder) objectId() (objectType uint16, instance uint32) {
	var value uint32
	d.decode(&value)
	objectType = uint16((value >> InstanceBits) & MaxObject)
	instance = value & MaxInstance
	return
}

func isExtendedTagNumber(x uint8) bool {
	return ((x & 0xF0) == 0xF0)
}

/* from clause 20.2.1.3.2 Constructed Data */
/* true if the tag is an opening tag */
func isOpeningTag(x uint8) bool {
	return ((x & 0x07) == 6)
}

/* from clause 20.2.1.3.2 Constructed Data */
/* true if the tag is a closing tag */
func isClosingTag(x uint8) bool {
	return ((x & 0x07) == 7)
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

const contextSpecificBit = 0x08

// context specific flag is the third bit
func isContextSpecific(meta uint8) bool {
	return ((meta & contextSpecificBit) > 0)
}

func setContextSpecific(x uint8) uint8 {
	return (x | contextSpecificBit)
}

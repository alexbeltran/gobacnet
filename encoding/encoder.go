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

func (e *Encoder) contextObjectID(tagNum uint8, objectType bactype.ObjectType, instance bactype.ObjectInstance) {
	/* length of object id is 4 octets, as per 20.2.14 */
	e.tag(tagInfo{ID: tagNum, Context: true, Value: 4})
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

func (e *Encoder) tag(tg tagInfo) {
	var t uint8
	var meta tagMeta
	if tg.Context {
		meta.setContextSpecific()
	}
	if tg.Opening {
		meta.setOpening()
	}
	if tg.Closing {
		meta.setClosing()
	}

	t = uint8(meta)
	if tg.Value <= 4 {
		t |= uint8(tg.Value)
	} else {
		t |= 5
	}

	// We have enough room to put it with the last value
	if tg.ID <= 14 {
		t |= (tg.ID << 4)
		e.write(t)

		// We don't have enough space so make it in a new byte
	} else {
		t |= 0xF0
		e.write(t)
		e.write(tg.ID)
	}
	if tg.Value > 4 {
		// Depending on the length, we will either write it as an 8 bit, 32 bit, or 64 bit integer
		if tg.Value <= 253 {
			e.write(uint8(tg.Value))
		} else if tg.Value <= 65535 {
			e.write(flag16bit)
			e.write(uint16(tg.Value))
		} else {
			e.write(flag32bit)
			e.write(tg.Value)
		}
	}
}

/* from clause 20.2.14 Encoding of an Object Identifier Value
returns the number of apdu bytes consumed */
func (e *Encoder) objectId(objectType bactype.ObjectType, instance bactype.ObjectInstance) {
	var value uint32
	value = ((uint32(objectType) & MaxObject) << InstanceBits) | (uint32(instance) & MaxInstance)
	e.write(value)
}

func (e *Encoder) contextEnumerated(tagNumber uint8, value uint32) {
	e.contextUnsigned(tagNumber, value)
}

func (e *Encoder) contextUnsigned(tagNumber uint8, value uint32) {
	len := valueLength(value)
	e.tag(tagInfo{ID: tagNumber, Context: true, Value: uint32(len)})
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

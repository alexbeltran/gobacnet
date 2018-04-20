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

const (
	tagNull            uint8 = 0
	tagBool            uint8 = 1
	tagUint            uint8 = 2
	tagInt             uint8 = 3
	tagReal            uint8 = 4
	tagDouble          uint8 = 5
	tagOctetString     uint8 = 6
	tagCharacterString uint8 = 7
	tagBitString       uint8 = 8
	tagEnumerated      uint8 = 9
	tagDate            uint8 = 10
	tagTime            uint8 = 11
	tagObjectID        uint8 = 12
	tagReserve1        uint8 = 13
	tagReserve2        uint8 = 14
	tagReserve3        uint8 = 15
	maxTag             uint8 = 16
)

// Other values omitted here can have variable length
const (
	realLen     uint32 = 4
	doubleLen   uint32 = 8
	dateLen     uint32 = 4
	timeLen     uint32 = 4
	objectIDLen uint32 = 4
)

// epochYear is an increment to all non-stored values. This year is chosen in
// the standard. Why? No idea. God help us all if bacnet hits the 255 + 1990
// limit
const epochYear = 1990

// If the values == 0XFF, that means it is not specified. We will take that to
const notDefined = 0xff

// All app layer is non-context specific
const appLayerContext = false

func IsOddMonth(month int) bool {
	return month == 13
}

func IsEvenMonth(month int) bool {
	return month == 14
}

func IsLastDayOfMonth(day int) bool {
	return day == 32
}

func IsEvenDayOfMonth(day int) bool {
	return day == 33
}

func IsOddDayOfMonth(day int) bool {
	return day == 32
}
func (e *Encoder) string(s string) {
	e.write(stringUTF8)
	e.write([]byte(s))
}
func (d *Decoder) string(s *string, len int) error {
	var t stringType
	d.decode(&t)
	if t != stringUTF8 {
		return fmt.Errorf("unsupported string format %d", t)
	}

	b := make([]byte, len)
	d.decode(b)
	*s = string(b)
	return d.Error()
}
func (e *Encoder) octetstring(b []byte) {
	e.write([]byte(b))
}
func (d *Decoder) octetstring(b *[]byte, len int) {
	*b = make([]byte, len)
	d.decode(b)
}

func (e *Encoder) date(dt bactype.Date) {
	// We don't want to override an unspecified time date
	if dt.Year != bactype.UnspecifiedTime {
		e.write(uint8(dt.Year - epochYear))
	} else {
		e.write(uint8(dt.Year))
	}
	e.write(uint8(dt.Month))
	e.write(uint8(dt.Day))
	e.write(uint8(dt.DayOfWeek))
}

func (d *Decoder) date(dt *bactype.Date) {
	var year, month, day, dayOfWeek uint8

	if dt.Year != bactype.UnspecifiedTime {
		dt.Year = int(year) + epochYear
	} else {
		dt.Year = int(year)
	}

	dt.Month = int(month)
	dt.Day = int(day)
	dt.DayOfWeek = bactype.DayOfWeek(dayOfWeek)
}

func (e *Encoder) time(t bactype.Time) {
	e.write(uint8(t.Hour))
	e.write(uint8(t.Minute))
	e.write(uint8(t.Second))

	// Stored as 1/100 of a second
	e.write(uint8(t.Millisecond / 10))
}
func (d *Decoder) time(t *bactype.Time) {
	var hour, min, sec, centisec uint8
	d.decode(&hour)
	d.decode(&min)
	d.decode(&sec)
	// Yeah, they report centisecs instead of milliseconds.
	d.decode(&centisec)

	t.Hour = int(hour)
	t.Minute = int(min)
	t.Second = int(sec)
	t.Millisecond = int(centisec) * 10

}

func (e *Encoder) boolean(x bool) {
	// Boolean information is stored into the length field
	var length uint32
	if x {
		length = 1
	} else {
		length = 0
	}
	e.tag(tagInfo{ID: tagBool, Context: appLayerContext, Value: length})
}

func (e *Encoder) real(x float32) {
	e.write(x)
}

func (d *Decoder) real(x *float32) {
	d.decode(x)
}

func (e *Encoder) double(x float64) {
	e.write(x)
}

func (d *Decoder) double(x *float64) {
	d.decode(x)
}

func (e *Encoder) AppData(i interface{}) error {
	switch val := i.(type) {
	case float32:
		e.tag(tagInfo{ID: tagReal, Context: appLayerContext, Value: realLen})
		e.real(val)
	case float64:
		e.tag(tagInfo{ID: tagDouble, Context: appLayerContext, Value: realLen})
		e.double(val)
	case bool:
		e.boolean(val)
	case string:
		// Add 1 to length to account for the encoding byte
		e.tag(tagInfo{ID: tagCharacterString, Context: appLayerContext, Value: uint32(len(val) + 1)})
		e.string(val)
	case uint32:
		length := valueLength(val)
		e.tag(tagInfo{ID: tagUint, Context: appLayerContext, Value: uint32(length)})
		e.unsigned(val)

	// Enumerated is pretty much a wrapper for a uint32 with an enumerated associated with it.
	case bactype.Enumerated:
		v := uint32(val)
		length := valueLength(v)
		e.tag(tagInfo{ID: tagEnumerated, Context: appLayerContext, Value: uint32(length)})
		e.enumerated(v)
	case bactype.ObjectID:
		e.tag(tagInfo{ID: tagObjectID, Context: appLayerContext, Value: objectIDLen})
		e.objectId(val.Type, val.Instance)

	case bactype.Null:
		e.tag(tagInfo{ID: tagNull, Context: appLayerContext})

	default:
		err := fmt.Errorf("Unknown type %T", i)
		// Set global error
		e.err = err
		return err
	}
	return nil
}

func (d *Decoder) AppData() (interface{}, error) {
	tag, _, lenvalue := d.tagNumberAndValue()
	len := int(lenvalue)

	switch tag {
	case tagNull:
		return bactype.Null{}, nil
	case tagBool:
		// Originally this was in C so non 0 values are considered
		// true
		return len > 0, d.Error()
	case tagUint:
		return d.unsigned(len), d.Error()
	case tagInt:
		return d.signed(len), d.Error()
	case tagReal:
		var x float32
		d.real(&x)
		return x, d.Error()
	case tagDouble:
		var x float64
		d.double(&x)
		return x, d.Error()
	case tagOctetString:
		var b []byte
		d.octetstring(&b, len)
		return b, d.Error()

	case tagCharacterString:
		var s string
		// Subtract 1 to length to account for the encoding byte
		err := d.string(&s, len-1)
		return s, err
	case tagBitString:
		return nil, fmt.Errorf("decoding bit strings is currently unsupported")
	case tagEnumerated:
		return d.enumerated(len), d.Error()
	case tagDate:
		var date bactype.Date
		d.date(&date)
		return date, d.Error()
	case tagTime:
		var t bactype.Time
		d.time(&t)
		return t, d.Error()
	case tagObjectID:
		objType, objInstance := d.objectId()
		return bactype.ObjectID{
			Type:     bactype.ObjectType(objType),
			Instance: objInstance,
		}, d.Error()
	default:
		return nil, fmt.Errorf("Unsupported tag: %d", tag)
	}
}

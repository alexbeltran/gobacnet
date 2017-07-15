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
	"encoding/binary"
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

// If the values == 0XFF, that means it is not specified. We will take that to
const notDefined = 0xff

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

func (d *Decoder) date(dt *bactype.Date) {
	var year, month, day, dayOfWeek uint8
	// God help us all if bacnet hits the 255 + 1990 limit
	dt.Year = int(year) + 1990
	dt.Month = int(month)
	dt.Day = int(day)
	dt.DayOfWeek = bactype.DayOfWeek(dayOfWeek)
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

func (d *Decoder) AppData(tag uint8, len int) (interface{}, error) {
	switch tag {
	case tagNull:
		return nil, fmt.Errorf("Null tag")
	case tagBool:
		// Originally this was in C so non 0 values are considered
		// true
		return len > 0, nil
	case tagUint:
		d.unsigned(len)
	case tagInt:
		return d.signed(len), nil
	case tagReal:
		var x float32
		binary.Read(d.buff, binary.BigEndian, &x)
		return x, nil
	case tagDouble:
	case tagOctetString:
	case tagCharacterString:
	case tagBitString:
	case tagEnumerated:
		return d.enumerated(len), nil
	case tagDate:
		var date bactype.Date
		d.date(&date)
		return date, nil
	case tagTime:
		var t bactype.Time
		d.time(&t)
		return t, nil
	case tagObjectID:
		objType, objInstance := d.objectId()
		return bactype.ObjectID{
			Type:     objType,
			Instance: objInstance,
		}, nil
	default:
		return nil, fmt.Errorf("Unsupported tag")
	}

	return nil, fmt.Errorf("Unknown Tag")
}

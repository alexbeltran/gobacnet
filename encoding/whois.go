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
	bactype "github.com/alexbeltran/gobacnet/types"
)

func (e *Encoder) WhoIs(low, high int32) error {
	apdu := bactype.APDU{
		DataType:           bactype.UnconfirmedServiceRequest,
		UnconfirmedService: bactype.ServiceUnconfirmedWhoIs,
	}
	e.write(apdu.DataType)
	e.write(apdu.UnconfirmedService)

	// The range is optional. A scan for all objects is done when either low/high
	// are negative or when we are scanning above the max instance
	if low >= 0 && high >= 0 && low < bactype.MaxInstance && high <
		bactype.MaxInstance {
		// Tag 0
		e.contextUnsigned(0, uint32(low))

		// Tag 1
		e.contextUnsigned(1, uint32(high))
	}
	return e.Error()
}

func (d *Decoder) WhoIs(low, high *int32) error {
	// APDU read in a higher level
	if d.len() == 0 {
		*low = bactype.WhoIsAll
		*high = bactype.WhoIsAll
		return nil
	}
	// Tag 0 - Low Value
	var expectedTag uint8
	tag, _, value := d.tagNumberAndValue()
	if tag != expectedTag {
		return &ErrorIncorrectTag{Expected: expectedTag, Given: tag}
	}
	l := d.unsigned(int(value))
	*low = int32(l)

	// Tag 1 - High Value
	expectedTag = 1
	tag, _, value = d.tagNumberAndValue()
	if tag != expectedTag {
		return &ErrorIncorrectTag{Expected: expectedTag, Given: tag}
	}
	h := d.unsigned(int(value))
	*high = int32(h)

	return d.Error()
}

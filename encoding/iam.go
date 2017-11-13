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
	"github.com/alexbeltran/gobacnet/types"
)

func (enc *Encoder) IAm(id types.IAm) error {
	enc.AppData(id.ID)
	enc.AppData(id.MaxApdu)
	enc.AppData(id.Segmentation)
	enc.AppData(id.Vendor)
	return enc.Error()
}

func (d *Decoder) IAm(id *types.IAm) error {
	objID, err := d.AppData()
	if err != nil {
		return err
	}
	if i, ok := objID.(types.ObjectID); ok {
		id.ID = i
	}

	maxapdu, _ := d.AppData()
	if m, ok := maxapdu.(uint32); ok {
		id.MaxApdu = m
	}

	segmentation, _ := d.AppData()
	if m, ok := segmentation.(uint32); ok {
		id.Segmentation = types.Enumerated(m)
	}

	vendor, _ := d.AppData()
	if v, ok := vendor.(uint32); ok {
		id.Vendor = v
	}

	return d.Error()
}

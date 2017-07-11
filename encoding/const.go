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

const ArrayAll = 0xFFFFFFFF

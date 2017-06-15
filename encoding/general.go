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

// valueLength caclulates how large the necessary value needs to be to fit in the appropriate
// packet length
func valueLength(value uint32) int {
	/* length of enumerated is variable, as per 20.2.11 */
	if value < 0x100 {
		return size8
	} else if value < 0x10000 {
		return size16
	} else if value < 0x1000000 {
		return size24
	}
	return size32
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

const contextSpecificBit = 0x08

// context specific flag is the third bit
func isContextSpecific(meta uint8) bool {
	return ((meta & contextSpecificBit) > 0)
}

func setContextSpecific(x uint8) uint8 {
	return (x | contextSpecificBit)
}

func isExtendedValue(x uint8) bool {
	return (x & 0x07) == 5
}

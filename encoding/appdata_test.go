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
	"reflect"
	"testing"
)

func subTestSimpleData(t *testing.T, d *Decoder, x interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		y, err := d.AppData()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(x, y) {
			t.Errorf("Mismatch between decrypted values. Received %v, expected %v", y, x)
		}
	}
}

func TestSimpleDataTypes(t *testing.T) {
	t.Run("Specific function based encoding", generalSimpleDataTypes(t, false))
	t.Run("General Interface Encoding", generalSimpleDataTypes(t, true))
}

// generic states if we plan to use specific or generic encoding
func generalSimpleDataTypes(t *testing.T, generic bool) func(t *testing.T) {
	return func(t *testing.T) {
		enc := NewEncoder()
		var real float32 = 3.14

		// Default float 64
		double := 1234567.890
		boolean := false
		// Will be stored as an 8 bit uint
		var small uint32 = 12

		// Stored as 16
		var medium uint32 = 0xABF0

		// Stored as 24
		var wtf uint32 = 0x302010

		// Stored as 32
		var large uint32 = 0xFFFFFFF0

		str := "pizza pizza"
		if generic {
			values := []interface{}{real, double, boolean, !boolean, small, medium, wtf, large, str}
			for _, v := range values {
				enc.AppData(v)
			}

		} else {
			enc.tag(tagReal, appLayerContext, realLen)
			enc.real(real)

			enc.tag(tagDouble, appLayerContext, doubleLen)
			enc.double(double)

			enc.boolean(boolean)
			enc.boolean(!boolean)

			enc.tag(tagUint, appLayerContext, size8)
			enc.unsigned(small)

			enc.tag(tagUint, appLayerContext, size16)
			enc.unsigned(medium)

			enc.tag(tagUint, appLayerContext, size24)
			enc.unsigned(wtf)

			enc.tag(tagUint, appLayerContext, size32)
			enc.unsigned(large)

			enc.tag(tagCharacterString, appLayerContext, uint32(len(str)))
			enc.string(str)
		}

		if err := enc.Error(); err != nil {
			t.Fatal(err)
		}

		dec := NewDecoder(enc.Bytes())
		t.Run("Encoding Real", subTestSimpleData(t, dec, real))
		t.Run("Encoding Double", subTestSimpleData(t, dec, double))
		t.Run("Encoding Boolean", subTestSimpleData(t, dec, boolean))
		t.Run("Encoding Flipped Boolean", subTestSimpleData(t, dec, !boolean))
		t.Run("Encoding Uint8", subTestSimpleData(t, dec, small))
		t.Run("Encoding Uint16", subTestSimpleData(t, dec, medium))
		t.Run("Encoding uint24", subTestSimpleData(t, dec, wtf))
		t.Run("Encoding uint32", subTestSimpleData(t, dec, large))
		t.Run("Encoding string", subTestSimpleData(t, dec, str))

		if err := dec.Error(); err != nil {
			t.Fatal(err)
		}
	}
}

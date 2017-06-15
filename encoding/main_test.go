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
	"testing"

	bactype "github.com/alexbeltran/gobacnet/types"
)

func TestNPDU(t *testing.T) {
	n := encodeNPDU(false, Normal)
	_, err := EncodePDU(&n, &BacnetAddress{}, &BacnetAddress{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSegsApduEncode(t *testing.T) {
	// Test is structured as parameter 1, parameter 2, output
	tests := [][]int{
		[]int{0, 1, 0},
		[]int{64, 60, 0x61},
		[]int{80, 205, 0x72},
		[]int{80, 405, 0x73},
		[]int{80, 1005, 0x74},
		[]int{3, 1035, 0x15},
		[]int{9, 1035, 0x35},
	}

	for _, test := range tests {
		d := int(encodeMaxSegsMaxApdu(test[0], test[1]))
		if d != test[2] {
			t.Fatalf("Input was Segments %d and Apdu %d: Expected %x got %x", test[0], test[1], test[2], d)
		}
	}
}

func TestObject(t *testing.T) {
	e := NewEncoder()
	var inObjectType uint16 = 17
	var inInstance uint32 = 23
	e.objectId(inObjectType, inInstance)
	b := e.Bytes()
	t.Log(b)

	d := NewDecoder(b)
	outObject, outInstance := d.objectId()

	if inObjectType != outObject {
		t.Fatalf("There was an issue encoding/decoding objectType. Input value was %d and output value was %d", inObjectType, outObject)
	}

	if inInstance != outInstance {
		t.Fatalf("There was an issue encoding/decoding objectType. Input value was %d and output value was %d", inInstance, outInstance)
	}

	if err := d.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestEnumerated(t *testing.T) {
	lengths := []int{size8, size16, size24, size32, size32}
	tests := []uint32{1, 2 << 8, 3 << 17, 7 << 25, 8 << 26}
	e := NewEncoder()
	for _, val := range tests {
		e.enumerated(val)
	}
	b := e.Bytes()
	d := NewDecoder(b)
	for i, val := range tests {
		x := d.enumerated(lengths[i])
		if x != val {
			t.Fatalf("Test[%d]:Decoded value %d doesn't match encoded value %d", i+1, x, val)
		}
	}

	d = NewDecoder(b)
	// 1000 is not a valid length
	x := d.enumerated(1000)
	if x != 0 {
		t.Fatalf("For invalid lengths, the value 0 should be decoded. The value %d was decoded", x)
	}
}

const compareErrFmt = "Mismatch in %s when decoding values. Expected: %d, recieved: %d"

func compare(t *testing.T, name string, a uint, b uint) {
	// See if the initial read property data matches the output read property
	if a != b {
		t.Fatalf(compareErrFmt, name, a, b)
	}
}

func compareReadProperties(t *testing.T, rd bactype.ReadPropertyData, outRd bactype.ReadPropertyData) {
	// See if the initial read property data matches the output read property
	compare(t, "object instance", uint(rd.ObjectInstance), uint(outRd.ObjectInstance))
	compare(t, "boject type", uint(rd.ObjectType), uint(outRd.ObjectType))
	compare(t, "object property", uint(rd.ObjectProperty), uint(outRd.ObjectProperty))
	compare(t, "array index", uint(rd.ArrayIndex), uint(outRd.ArrayIndex))
	compare(t, "application data length", uint(len(rd.ApplicationData)), uint(len(outRd.ApplicationData)))
	if len(rd.ApplicationData) > 0 {
		for i, _ := range rd.ApplicationData {
			compare(t, "application data", uint(rd.ApplicationData[i]), uint(outRd.ApplicationData[i]))
		}
	}
}

func subTestReadProperty(t *testing.T, rd bactype.ReadPropertyData) {
	e := NewEncoder()
	e.ReadProperty(10, rd)
	if err := e.Error(); err != nil {
		t.Fatal(err)
	}

	b := e.Bytes()
	d := NewDecoder(b)

	// Read Property reads 4 extra fields that are not original encoded. Need to
	//find out where these 4 fields come from
	d.buff.Read(make([]uint8, 4))
	outRd, err := d.ReadProperty()
	if err != nil {
		t.Fatal(err)
	}
	compareReadProperties(t, rd, outRd)
}

func subTestReadPropertyAck(t *testing.T, rd bactype.ReadPropertyData) {
	e := NewEncoder()
	e.ReadPropertyAck(10, rd)
	if err := e.Error(); err != nil {
		t.Fatal(err)
	}

	b := e.Bytes()
	d := NewDecoder(b)

	// Read Property reads 4 extra fields that are not original encoded. Need to
	//find out where these 4 fields come from
	d.buff.Read(make([]uint8, 3))
	outRd, err := d.ReadProperty()
	if err != nil {
		t.Fatal(err)
	}
	compareReadProperties(t, rd, outRd)
}

func TestReadAckProperty(t *testing.T) {
	rd := bactype.ReadPropertyData{
		ObjectType:      37,
		ObjectInstance:  1000,
		ObjectProperty:  3921,
		ArrayIndex:      ArrayAll,
		ApplicationData: []byte{3, 7, 23, 5, 11},
	}
	subTestReadPropertyAck(t, rd)

	rd.ArrayIndex = 2
	subTestReadPropertyAck(t, rd)
}

func TestReadProperty(t *testing.T) {
	rd := bactype.ReadPropertyData{
		ObjectType:     37,
		ObjectInstance: 1000,
		ObjectProperty: 3921,
		ArrayIndex:     ArrayAll,
	}
	// Test a generic read property
	subTestReadProperty(t, rd)

	// Test with an array value given
	rd.ArrayIndex = 1
	subTestReadProperty(t, rd)
}

// Test for when the read property is too small and error handling
func TestReadPropertyTooSmall(t *testing.T) {
	e := NewEncoder()
	var garbage uint16 = 100
	e.write(garbage)
	d := NewDecoder(e.Bytes())
	_, err := d.ReadProperty()
	if err == nil {
		t.Fatal("Missed too small error")
	}
}

// Test for mismatch id error.
func TestReadPropertyMismatch(t *testing.T) {
	e := NewEncoder()
	var incorrectTag uint8 = 100
	var randomValue uint32 = 4

	// Has to be written 4 times at least since a minimum of 7 data is required
	// for read property
	for i := 0; i < 7; i++ {
		e.tag(incorrectTag, true, randomValue)
	}
	d := NewDecoder(e.Bytes())
	_, err := d.ReadProperty()
	if err == nil {
		t.Fatal("Incorrect tag number was allowed to pass")
	}
}

func TestTag(t *testing.T) {
	e := NewEncoder()
	// Respective to each other
	inTag := []uint8{4, 15, 30, 254, 1}
	inValue := []uint32{4, 20, 6000, 1, 70000}

	for i, tag := range inTag {
		e.tag(tag, true, inValue[i])
	}

	// Check for errors during the encoding processes
	if err := e.Error(); err != nil {
		t.Fatal(err)
	}

	b := e.Bytes()
	d := NewDecoder(b)
	for i, tag := range inTag {
		outTag, _, value := d.tagNumberAndValue()
		if tag != outTag {
			t.Fatalf("Test[%d]: Tag was not processed propertly. Expected %d, got %d", i, tag, outTag)
		}

		if value != inValue[i] {
			t.Fatalf("Test[%d]: Value was not processed propertly. Expected %d, got %d", i, inValue[i], value)
		}
	}

	// Check for errors during the decoding process
	if err := d.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestTagMetadata(t *testing.T) {
	var m tagMeta = 0
	m.setClosing()
	if !m.isClosing() {
		t.Fatal("Closing flag was not properly set.")
	}
	m.Clear()
	if m.isClosing() {
		t.Fatal("Flag was not cleared")
	}

	m.setOpening()
	if !m.isOpening() {
		t.Fatal("Opening flag was not properly set")
	}

	m.Clear()

	if m.isContextSpecific() {
		t.Fatal("Context specific was set when it shouldn't have been")
	}
	m.setContextSpecific()
	if !m.isContextSpecific() {
		t.Fatal("Context specific was not properly set.")
	}
}

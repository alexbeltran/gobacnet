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
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"testing"

	bactype "github.com/alexbeltran/gobacnet/types"
)

func TestReadPropertyService(t *testing.T) {
	// This value is based on a known sample
	expected := []byte{129, 10, 0, 22, 1, 36, 9, 124, 1, 29, 255, 0, 5, 1, 12,
		12, 0, 0, 0, 1, 25, 85}

	e := NewEncoder()
	//s := `{"ID":24289,"MaxAPDU":480,"Address":{"Mac":"ChQAzLrA","MacLen":6,"Net":2428,"Adr":"HQ==","AdrLen":1}}`
	var mac []uint8
	var adr []uint8
	json.Unmarshal([]byte("\"ChQAzLrA\""), &mac)
	json.Unmarshal([]byte("\"HQ==\""), &adr)
	readProp := bactype.ReadPropertyData{
		Object: bactype.Object{
			ID: bactype.ObjectID{
				Type:     0,
				Instance: 1,
			},
			Properties: []bactype.Property{
				bactype.Property{
					Type:       85,
					ArrayIndex: ArrayAll,
				},
			},
		},
	}

	dest := bactype.Address{
		Net:    2428,
		Mac:    mac,
		MacLen: 6,
		Len:    1,
		Adr:    adr,
	}
	e.NPDU(bactype.NPDU{
		Version:               bactype.ProtocolVersion,
		IsNetworkLayerMessage: false,
		ExpectingReply:        true,
		HopCount:              bactype.DefaultHopCount,
		Priority:              bactype.Normal,
		Destination:           &dest,
	})
	e.ReadProperty(1, readProp)
	data := e.Bytes()

	enc := NewEncoder()
	bv := bactype.BVLC{
		Type:     bactype.BVLCTypeBacnetIP,
		Function: bactype.BacFuncUnicast,
		Length:   4 + uint16(len(data)),
		Data:     data,
	}
	enc.BVLC(bv)

	raw := enc.Bytes()
	for i, b := range raw {
		if expected[i] != b {
			t.Errorf("Error during decoding: %x does not equal expected %x", b, expected[i])
		}
	}
	if len(raw) != len(expected) {
		t.Fatalf("There is a mismatch in sizes. Got: %d, Expected:%d", len(raw), len(expected))
	}
	t.Logf("Length: %d", len(raw))
}

func TestReadPropertyResponse(t *testing.T) {
	// This value is based on a known sample
	in := []byte{48, 1, 12, 12, 0, 0, 0, 1, 25, 85, 62, 68, 192, 160, 0, 0, 63}
	d := NewDecoder(in)
	apdu := bactype.APDU{}
	d.APDU(&apdu)

	rpd := bactype.ReadPropertyData{}
	err := d.ReadProperty(&rpd)
	if err != nil {
		t.Fatal(err)
	}

	x := rpd.Object.Properties[0].Data
	f := x.(float32)
	if f != -5.0 {
		t.Fatalf("Final value was not decrypted properly")
	}

}

func TestWhoIs(t *testing.T) {
	e := NewEncoder()
	var low int32 = 28
	var high int32 = 32
	err := e.WhoIs(low, high)
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(e.Bytes())
	a := bactype.APDU{}
	d.APDU(&a)

	d = NewDecoder(a.RawData)
	var lowOut, highOut int32
	d.WhoIs(&lowOut, &highOut)

	if err = d.Error(); err != nil {
		t.Fatal(err)
	}

	if low != lowOut || high != highOut {
		t.Fatalf("WhoIs was not decoded properly. Low was %d, given %d. High was %d, given %d", low, lowOut, high, highOut)
	}
}

func TestIAmRealData(t *testing.T) {
	b := []byte{196, 2, 3, 180, 113, 34, 1, 224, 145, 3, 33, 24}
	dec := NewDecoder(b)
	for dec.len() > 0 {
		x, err := dec.AppData()
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("app: %v", x)
	}
}

func TestIAm(t *testing.T) {
	iam := bactype.IAm{
		MaxApdu: 1234,
		ID: bactype.ObjectID{
			Instance: 401,
			Type:     17,
		},
		Segmentation: 100,
		Vendor:       413,
	}

	enc := NewEncoder()
	err := enc.IAm(iam)
	if err != nil {
		t.Fatal(err)
	}

	dec := NewDecoder(enc.Bytes())

	var after bactype.IAm
	err = dec.IAm(&after)
	if err != nil {
		t.Fatal(err)
	}

	equal := reflect.DeepEqual(iam, after)
	if !equal {
		t.Errorf("Encoding/Decoding Failed: %v does not equal %v", iam, after)
	}

}

func TestRealDataIAm(t *testing.T) {
	raw := []byte{196, 2, 0, 4, 210, 34, 5, 196, 145, 3, 34, 1, 4}
	dec := NewDecoder(raw)
	var out bactype.IAm
	err := dec.IAm(&out)
	t.Logf("%v", out)
	if err != nil {
		t.Fatal(err)
	}
}

/*
func TestIAm(t *testing.T) {
	ids := []bactype.ObjectID{
		bactype.ObjectID{Instance: 1, Type: 5},
		bactype.ObjectID{Instance: 99, Type: 6},
		bactype.ObjectID{Instance: 133, Type: 1},
	}
	enc := NewEncoder()
	err := enc.IAm(ids)
	if err != nil {
		t.Fatal(err)
	}

	dec := NewDecoder(enc.Bytes())

	decIds := make([]bactype.ObjectID, len(ids))
	err = dec.IAm(decIds[:])

	equal := reflect.DeepEqual(ids, decIds)
	if !equal {
		t.Errorf("Encoding/Decoding Failed: %v does not equal %v", ids, decIds)
	}

}
*/
func TestReadMultiple(t *testing.T) {
	// This is from a data trace of a sample server which returns the names of each object.
	raw := []byte{12, 2, 0, 4, 210, 30, 41, 77, 78, 117, 13, 0, 83, 105, 109,
		112, 108, 101, 83, 101, 114, 118, 101, 114, 79, 31, 12, 2, 128, 0, 0, 30,
		41, 77, 78, 117, 7, 0, 70, 73, 76, 69, 32, 48, 79, 31, 12, 2, 128, 0, 1, 30,
		41, 77, 78, 117, 7, 0, 70, 73, 76, 69, 32, 49, 79, 31, 12, 2, 128, 0, 2, 30,
		41, 77, 78, 117, 7, 0, 70, 73, 76, 69, 32, 50, 79, 31}
	names := []string{"SimpleServer", "FILE 0", "FILE 1", "FILE 2"}
	dec := NewDecoder(raw)
	rp := bactype.ReadMultipleProperty{}
	err := dec.ReadMultiplePropertyAck(&rp)
	if err != nil {
		t.Fatal(err)
	}

	counter := 0
	for _, obj := range rp.Objects {
		for _, prop := range obj.Properties {
			name, ok := prop.Data.(string)
			if !ok {
				t.Fatalf("Type mismatch. Type should be string, it is %T", prop.Data)
			}
			if strings.Compare(name, names[counter]) > 0 {
				t.Fatalf("Object name should be \"%s\" not \"%s\"", names[counter], name)
			}
			counter++
		}
	}
}

func TestReadMultipleTwo(t *testing.T) {
	raw := []byte{12, 2, 0, 4, 210, 30, 41, 77, 78, 117, 13, 0, 83, 105, 109, 112, 108, 101, 83,
		101, 114, 118, 101, 114, 79, 41, 75, 78, 196, 2, 0, 4, 210, 79, 31, 12, 2, 128, 0, 0, 30, 41, 77, 78,
		117, 7, 0, 70, 73, 76, 69, 32, 48, 79, 41, 75, 78, 196, 2, 128, 0, 0, 79, 31, 12, 2, 128, 0, 1, 30, 41,
		77, 78, 117, 7, 0, 70, 73, 76, 69, 32, 49, 79, 41, 75, 78, 196, 2, 128, 0, 1, 79, 31, 12, 2, 128, 0, 2,
		30, 41, 77, 78, 117, 7, 0, 70, 73, 76, 69, 32, 50, 79, 41, 75, 78, 196, 2, 128, 0, 2, 79, 31}
	names := []string{"SimpleServer", "FILE 0", "FILE 1", "FILE 2"}

	dec := NewDecoder(raw)
	rp := bactype.ReadMultipleProperty{}
	err := dec.ReadMultiplePropertyAck(&rp)
	if err != nil {
		t.Fatal(err)
	}

	counter := 0
	for _, obj := range rp.Objects {
		prop := obj.Properties[0]
		name, ok := prop.Data.(string)
		if !ok {
			t.Fatalf("Type mismatch. Type should be string, it is %T", prop.Data)
		}
		if strings.Compare(name, names[counter]) > 0 {
			t.Fatalf("Object name should be \"%s\" not \"%s\"", names[counter], name)
		}

		prop = obj.Properties[1]
		_, ok = prop.Data.(bactype.ObjectID)
		if !ok {
			t.Fatalf("Type mismatch. Type should be object id, it is %T", prop.Data)
		}
		counter++
	}
}

func TestReadMultipleThree(t *testing.T) {

	b := []byte{48, 54, 14, 12, 2, 0, 78, 233, 30, 41, 76, 57, 241, 78, 196, 0, 128,
		0, 62, 79, 41, 76, 57, 242, 78, 196, 5, 0, 0, 68, 79, 41, 76, 57, 243, 78,
		196, 0, 128, 0, 61, 79, 41, 76, 57, 244, 78, 196, 5, 0, 0, 67, 79, 41, 76,
		57, 245, 78, 196, 5, 0, 0, 65, 79, 41, 76, 57, 246, 78, 196, 0, 128, 0, 60,
		79, 41, 76, 57, 247, 78, 196, 5, 0, 0, 66, 79, 41, 76, 57, 248, 78, 196, 5,
		0, 0, 69, 79, 41, 76, 57, 249, 78, 196, 1, 64, 0, 40, 79, 41, 76, 57, 250,
		78, 196, 1, 64, 0, 42, 79, 41, 76, 57, 251, 78, 196, 4, 192, 0, 17, 79, 41,
		76, 57, 252, 78, 196, 1, 64, 0, 38, 79, 41, 76, 57, 253, 78, 196, 4, 192, 0,
		16, 79, 41, 76, 57, 254, 78, 196, 1, 64, 0, 37, 79, 41, 76, 57, 255, 78,
		196, 4, 192, 0, 18, 79, 41, 76, 58, 1, 0, 78, 196, 1, 64, 0, 39, 79, 41, 76,
		58, 1, 1, 78, 196, 1, 64, 0, 41, 79, 41, 76, 58, 1, 2, 78, 196, 5, 0, 0, 73,
		79, 41, 76, 58, 1, 3, 78, 196, 0, 128, 0, 75, 79, 41, 76, 58, 1, 4, 78, 196,
		5, 0, 0, 82, 79, 41, 76, 58, 1, 5, 78, 196, 0, 128, 0, 76, 79, 41, 76, 58,
		1, 6, 78, 196, 5, 0, 0, 83, 79, 41, 76, 58, 1, 7, 78, 196, 0, 128, 0, 77,
		79, 41, 76, 58, 1, 8, 78, 196, 5, 0, 0, 84, 79, 41, 76, 58, 1, 9, 78, 196,
		0, 128, 0, 67, 79, 41, 76, 58, 1, 10, 78, 196, 5, 0, 0, 74, 79, 41, 76, 58,
		1, 11, 78, 196, 0, 128, 0, 68, 79, 41, 76, 58, 1, 12, 78, 196, 5, 0, 0, 75,
		79, 41, 76, 58, 1, 13, 78, 196, 5, 0, 0, 76, 79, 41, 76, 58, 1, 14, 78, 196,
		0, 128, 0, 73, 79, 31}
	dec := NewDecoder(b)
	var apdu bactype.APDU
	dec.APDU(&apdu)

	rp := bactype.ReadMultipleProperty{}
	err := dec.ReadMultiplePropertyAck(&rp)
	if err != nil {
		t.Fatalf("failed to decode read multiple: %v", err)
	}
}

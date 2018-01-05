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

package types

import (
	"encoding/json"
	"fmt"
	"net"
)

type Enumerated uint32
type ObjectType uint16
type ObjectInstance uint32

type ObjectID struct {
	Type     ObjectType
	Instance ObjectInstance
}

type Object struct {
	Name        string
	Description string
	ID          ObjectID
	Properties  []Property `json:",omitempty"`
}

type Property struct {
	Type       uint32
	ArrayIndex uint32
	Data       interface{}
}

type ReadPropertyData struct {
	InvokeID   uint16
	Object     Object
	ErrorClass uint8
	ErrorCode  uint8
}

type ReadMultipleProperty struct {
	Objects    []Object
	ErrorClass uint8
	ErrorCode  uint8
}

type Address struct {
	Net    uint16
	Len    uint8
	MacLen uint8
	Mac    []uint8
	Adr    []uint8
}

type ObjectMap map[ObjectType]map[ObjectInstance]Object

type Device struct {
	ID           ObjectID
	MaxApdu      uint32
	Segmentation Enumerated
	Vendor       uint32
	Addr         Address
	Objects      ObjectMap
}

type IAm struct {
	ID           ObjectID
	MaxApdu      uint32
	Segmentation Enumerated
	Vendor       uint32
	Addr         Address
}

const broadcastNetwork uint16 = 0xFFFF

// IsBroadcast returns if the address is a broadcast address
func (a *Address) IsBroadcast() bool {
	if a.Net == broadcastNetwork || a.MacLen == 0 {
		return true
	}
	return false
}

func (a *Address) SetBroadcast(b bool) {
	if b {
		a.MacLen = 0
	} else {
		a.MacLen = uint8(len(a.Mac))
	}
}

// IsSubBroadcast checks to see if packet is meant to be a network
// specific broadcast
func (a *Address) IsSubBroadcast() bool {
	if a.Net > 0 && a.Len == 0 {
		return true
	}
	return false
}

// IsUnicast checks to see if packet is meant to be a unicast
func (a *Address) IsUnicast() bool {
	if a.MacLen == 6 {
		return true
	}
	return false
}

// UDPAddr parses the mac address and returns an proper net.UDPAddr
func (a *Address) UDPAddr() (net.UDPAddr, error) {
	if len(a.Mac) != 6 {
		return net.UDPAddr{}, fmt.Errorf("Mac is too short at %d", len(a.Mac))
	}
	port := uint(a.Mac[4])<<8 | uint(a.Mac[5])
	ip := net.IPv4(byte(a.Mac[0]), byte(a.Mac[1]), byte(a.Mac[2]), byte(a.Mac[3]))
	return net.UDPAddr{
		IP:   ip,
		Port: int(port),
	}, nil
}

// Address converts a given udp address into a bacnet address
func UDPToAddress(n *net.UDPAddr) Address {
	a := Address{}
	p := uint16(n.Port)

	// Length of IP plus the port
	length := net.IPv4len + 2
	a.Mac = make([]uint8, length)
	//Encode ip
	for i := 0; i < net.IPv4len; i++ {
		a.Mac[i] = n.IP[i]
	}

	// Encode port
	a.Mac[net.IPv4len+0] = uint8(p >> 8)
	a.Mac[net.IPv4len+1] = uint8(p & 0x00FF)

	a.MacLen = uint8(length)
	return a
}

// Len returns the total number of entries within the object map.
func (o ObjectMap) Len() int {
	counter := 0
	for _, t := range o {
		for _ = range t {
			counter++
		}

	}
	return counter
}

func (om ObjectMap) MarshalJSON() ([]byte, error) {
	m := make(map[string]map[ObjectInstance]Object)
	for typ, sub := range om {
		key := typ.String()
		if m[key] == nil {
			m[key] = make(map[ObjectInstance]Object)
		}
		for inst, obj := range sub {
			m[key][inst] = obj
		}
	}
	return json.Marshal(m)
}

func (om ObjectMap) UnmarshalJSON(data []byte) error {
	m := make(map[string]map[ObjectInstance]Object, 0)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	for t, sub := range m {
		key := GetType(t)
		if om[key] == nil {
			om[key] = make(map[ObjectInstance]Object)
		}
		for inst, obj := range sub {
			om[key][inst] = obj
		}
	}
	return nil
}

// ObjectSlice returns all the objects in the device as a slice (not thread-safe)
func (dev Device) ObjectSlice() []Object {
	objs := []Object{}
	for _, objMap := range dev.Objects {
		for _, o := range objMap {
			objs = append(objs, o)
		}
	}
	return objs
}

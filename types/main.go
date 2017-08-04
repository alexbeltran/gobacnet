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
	"fmt"
	"log"
	"net"
)

type ObjectID struct {
	Type     uint16
	Instance uint32
}

type Object struct {
	ID         ObjectID
	Properties []Property
}

type Property struct {
	Type       uint32
	ArrayIndex uint32
	Data       []uint8
	DataLen    int
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
	log.Println("ROCESSING")
	log.Printf("IP: %v ", n.IP)

	// Length of IP plus the port
	length := 4 + 2
	a.Mac = make([]uint8, length)
	//Encode ip
	for i := range n.IP {
		a.Mac[i] = n.IP[i]
	}

	// Encode port
	a.Mac[len(n.IP)+0] = uint8(p >> 8)
	a.Mac[len(n.IP)+1] = uint8(p & 0x00FF)

	a.MacLen = uint8(length)
	return a
}

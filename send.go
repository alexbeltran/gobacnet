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

package gobacnet

import (
	"fmt"
	"net"

	"github.com/alexbeltran/gobacnet/encoding"
	bactype "github.com/alexbeltran/gobacnet/types"
)

// address returns the address given
func (c *Client) address(addr bactype.Address) (net.UDPAddr, error) {
	if addr.IsBroadcast() {
		return net.UDPAddr{
			IP:   c.broadcastAddress,
			Port: c.port,
		}, nil
	} else if addr.IsSubBroadcast() {
		// Network specific
		if addr.IsUnicast() {
			return addr.UDPAddr()
		}

		// Broadcast
		return net.UDPAddr{
			IP:   c.broadcastAddress,
			Port: c.port,
		}, nil
	} else if addr.IsUnicast() {
		return addr.UDPAddr()
	}
	return net.UDPAddr{}, fmt.Errorf("Unable to parse bacnet address")
}

// Sets the udp version used to transfer data
// See https://golang.org/pkg/net/#DialUDP
const udpVersion = "udp"
const mtuHeaderLength = 4
const forwardHeaderLength = 10

// send transfers the raw apdu byte slice to the destination address.
func (c *Client) send(dest bactype.Address, data []byte) (int, error) {
	var header bactype.BVLC

	// Set packet type
	header.Type = bactype.BVLCTypeBacnetIP

	if dest.IsBroadcast() || dest.IsSubBroadcast() {
		// SET BROADCAST FLAG
		header.Function = bactype.BacFuncBroadcast
	} else {
		// SET UNICAST FLAG
		header.Function = bactype.BacFuncUnicast
	}
	header.Length = uint16(mtuHeaderLength + len(data))
	header.Data = data
	e := encoding.NewEncoder()
	err := e.BVLC(header)
	if err != nil {
		return 0, err
	}

	// Get IP Address
	d, err := dest.UDPAddr()
	if err != nil {
		return 0, err
	}

	// use default udp type, src = local address (nil)
	return c.listener.WriteTo(e.Bytes(), &d)
}

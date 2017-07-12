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

// NPDU encodes the network layer control message
func (e *Encoder) NPDU(n bactype.NPDU) {
	e.write(n.Version)

	// Prepare metadata into the second byte
	meta := NPDUMetadata(0)
	meta.SetNetworkLayerMessage(n.IsNetworkLayerMessage)
	meta.SetExpectingReply(n.ExpectingReply)
	meta.SetPriority(n.Priority)

	// Check to see if we have a net address. If so set destination true
	if n.Destination != nil {
		if n.Destination.Net != 0 {
			meta.SetDestination(true)
		}
	}

	// Repeat for source
	if n.Source != nil {
		if n.Source.Net != 0 {
			meta.SetSource(true)
		}
	}
	e.write(meta)

	if meta.HasDestination() {
		e.write(n.Destination.Net)

		// Address
		e.write(n.Destination.Len)
		e.write(n.Destination.Adr)
	}

	if meta.HasSource() {
		e.write(n.Source.Net)

		// Address
		e.write(n.Source.Len)
		e.write(n.Source.Adr)
	}

	// Hop count is after source
	if meta.HasDestination() {
		e.write(n.HopCount)
	}

	if meta.IsNetworkLayerMessage() {
		e.write(n.NetworkLayerMessageType)

		// If the network value is above 0x80, then it should have a vendor id
		if n.NetworkLayerMessageType >= 0x80 {
			e.write(n.VendorId)
		}
	}
}

func (d *Decoder) Address(a *bactype.Address) {
	d.decode(&a.Net)
	d.decode(&a.Len)

	// Make space for address
	a.Adr = make([]uint8, a.Len)
	d.decode(a.Adr)
}

// NPDU encodes the network layer control message
func (d *Decoder) NPDU(n *bactype.NPDU) error {
	d.decode(&n.Version)

	// Prepare metadata into the second byte
	meta := NPDUMetadata(0)
	d.decode(&meta)
	n.ExpectingReply = meta.ExpectingReply()
	n.IsNetworkLayerMessage = meta.IsNetworkLayerMessage()
	n.Priority = meta.Priority()

	if meta.HasDestination() {
		n.Destination = &bactype.Address{}
		d.Address(n.Destination)
	}

	if meta.HasSource() {
		n.Source = &bactype.Address{}
		d.Address(n.Source)
	}

	if meta.HasDestination() {
		d.decode(&n.HopCount)
	} else {
		n.HopCount = 0
	}

	if meta.IsNetworkLayerMessage() {
		d.decode(&n.NetworkLayerMessageType)
		if n.NetworkLayerMessageType > 0x80 {
			d.decode(&n.VendorId)
		}
	}
	return d.Error()
}

// NPDUMetadata includes additional metadata about npdu message
type NPDUMetadata byte

const maskNetworkLayerMessage = 1 << 7
const maskDestination = 1 << 5
const maskSource = 1 << 3
const maskExpectingReply = 1 << 2

// General setter for the info bits using the mask
func (meta *NPDUMetadata) setInfoMask(b bool, mask byte) {
	*meta = NPDUMetadata(setInfoMask(byte(*meta), b, mask))
}

// CheckMask uses mask to check bit position
func (meta *NPDUMetadata) checkMask(mask byte) bool {
	return (*meta & NPDUMetadata(mask)) > 0

}

// IsNetworkLayerMessage returns true if it is a network layer message
func (n *NPDUMetadata) IsNetworkLayerMessage() bool {
	return n.checkMask(maskNetworkLayerMessage)
}

func (n *NPDUMetadata) SetNetworkLayerMessage(b bool) {
	n.setInfoMask(b, maskNetworkLayerMessage)
}

// Priority returns priority
func (n *NPDUMetadata) Priority() bactype.NPDUPriority {
	// Encoded in bit 0 and 1
	return bactype.NPDUPriority(byte(*n) & 3)
}

// SetPriority for NPDU
func (n *NPDUMetadata) SetPriority(p bactype.NPDUPriority) {
	// Clear the first two bits
	//*n &= (0xF - 3)
	*n |= NPDUMetadata(p)
}

func (n *NPDUMetadata) HasDestination() bool {
	return n.checkMask(maskDestination)
}

func (n *NPDUMetadata) SetDestination(b bool) {
	n.setInfoMask(b, maskDestination)
}

func (n *NPDUMetadata) HasSource() bool {
	return n.checkMask(maskSource)
}

func (n *NPDUMetadata) SetSource(b bool) {
	n.setInfoMask(b, maskSource)
}

// IsNetworkLayerMessage returns true if it is a network layer message
func (n *NPDUMetadata) ExpectingReply() bool {
	return n.checkMask(maskExpectingReply)
}

func (n *NPDUMetadata) SetExpectingReply(b bool) {
	n.setInfoMask(b, maskExpectingReply)
}

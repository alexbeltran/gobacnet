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
	"net"

	"github.com/alexbeltran/gobacnet/encoding"
	"github.com/alexbeltran/gobacnet/types"
)

// WhoIs finds all devices with ids between the provided low and high values.
// Use constant ArrayAll for both fields to scan the entire network at once.
// Using ArrayAll is highly discouraged for most networks since it can lead
// to a high congested network.
func (c *Client) WhoIs(low, high int) ([]types.Device, error) {
	dest := types.UDPToAddress(&net.UDPAddr{
		IP:   c.broadcastAddress,
		Port: DefaultPort,
	})
	src, _ := c.localAddress()

	dest.SetBroadcast(true)

	enc := encoding.NewEncoder()
	npdu := types.NPDU{
		Version:               types.ProtocolVersion,
		Destination:           &dest,
		Source:                &src,
		IsNetworkLayerMessage: false,

		// We are not expecting a direct reply from a single destination
		ExpectingReply: false,
		Priority:       types.Normal,
		HopCount:       types.DefaultHopCount,
	}
	enc.NPDU(npdu)

	err := enc.WhoIs(int32(low), int32(high))
	if err != nil {
		return nil, err
	}

	// Subscribe to any changes in the the range. If it is a broadcast,
	var start, end int
	if low == -1 || high == -1 {
		start = 0
		end = maxInt
	} else {
		start = low
		end = high
	}

	// Run in parallel
	errChan := make(chan error)
	go func() {
		_, err = c.send(dest, enc.Bytes())
		errChan <- err
	}()
	values, err := c.utsm.Subscribe(start, end)
	if err != nil {
		return nil, err
	}
	err = <-errChan
	if err != nil {
		return nil, err
	}

	// Weed out values that are not important such as non object type
	// and that are not
	uniqueMap := make(map[types.ObjectInstance]types.Device)
	uniqueList := make([]types.Device, len(uniqueMap))
	for _, v := range values {
		r, ok := v.(types.IAm)

		// Skip non I AM responses
		if !ok {
			continue
		}

		// Check to see if we are in the map before inserting
		if _, ok := uniqueMap[r.ID.Instance]; !ok {
			dev := types.Device{
				Addr:         r.Addr,
				ID:           r.ID,
				MaxApdu:      r.MaxApdu,
				Segmentation: r.Segmentation,
				Vendor:       r.Vendor,
			}
			uniqueMap[r.ID.Instance] = types.Device(dev)
			uniqueList = append(uniqueList, dev)
		}
	}
	return uniqueList, err
}

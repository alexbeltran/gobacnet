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
	"context"
	"fmt"
	"time"

	"github.com/alexbeltran/gobacnet/encoding"
	bactype "github.com/alexbeltran/gobacnet/types"
)

const maxReattempt = 2

// ReadMultipleProperty uses the given device and read property request to read
// from a device. Along with being able to read multiple properties from a
// device, it can also read these properties from multiple objects. This is a
// good feature to read all present values of every object in the device. This
// is a batch operation compared to a ReadProperty and should be used in place
// when reading more than two objects/properties.
func (c *Client) ReadMultiProperty(dev bactype.Device, rp bactype.ReadMultipleProperty) (bactype.ReadMultipleProperty, error) {
	var out bactype.ReadMultipleProperty

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, err := c.tsm.ID(ctx)
	if err != nil {
		return out, fmt.Errorf("unable to get transaction id: %v", err)
	}
	defer c.tsm.Put(id)

	udp, err := c.localUDPAddress()
	if err != nil {
		return out, err
	}
	src := bactype.UDPToAddress(udp)

	enc := encoding.NewEncoder()
	enc.NPDU(bactype.NPDU{
		Version:               bactype.ProtocolVersion,
		Destination:           &dev.Addr,
		Source:                &src,
		IsNetworkLayerMessage: false,
		ExpectingReply:        true,
		Priority:              bactype.Normal,
		HopCount:              bactype.DefaultHopCount,
	})
	enc.ReadMultipleProperty(uint8(id), rp)
	if enc.Error() != nil {
		return out, fmt.Errorf("encoding read multiple property failed: %v", err)
	}

	pack := enc.Bytes()
	if dev.MaxApdu < uint32(len(pack)) {
		return out, fmt.Errorf("read multiple property is too large (max: %d given: %d)", dev.MaxApdu, len(pack))
	}
	// the value filled doesn't matter. it just needs to be non nil
	err = fmt.Errorf("go")

	for count := 0; err != nil && count < maxReattempt; count++ {
		out, err = c.sendReadMultipleProperty(id, dev, pack)
		if err == nil {
			return out, nil
		}
	}
	return out, fmt.Errorf("failed %d tries: %v", maxReattempt, err)
}

func (c *Client) sendReadMultipleProperty(id int, dev bactype.Device, request []byte) (bactype.ReadMultipleProperty, error) {
	var out bactype.ReadMultipleProperty
	_, err := c.send(dev.Addr, request)
	if err != nil {
		return out, err
	}

	raw, err := c.tsm.Receive(id, time.Duration(5)*time.Second)
	if err != nil {
		return out, fmt.Errorf("unable to receive id %d: %v", id, err)
	}

	var b []byte
	switch v := raw.(type) {
	case error:
		return out, err
	case []byte:
		b = v
	default:
		return out, fmt.Errorf("received unknown datatype %T", raw)
	}

	dec := encoding.NewDecoder(b)

	var apdu bactype.APDU
	dec.APDU(&apdu)
	err = dec.ReadMultiplePropertyAck(&out)
	if err != nil {
		c.log.Debugf("WEIRD PACKET: %v: %v", err, b)
		return out, err
	}
	return out, err
}

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
	"log"
	"time"

	"github.com/alexbeltran/gobacnet/encoding"
	bactype "github.com/alexbeltran/gobacnet/types"
)

func (c *Client) sendRequest() error {
	id, err := c.tsm.GetFree()
	if err != nil {
		return err
	}
	log.Printf("id:%d", id)
	return nil
}

func (c *Client) ReadProperty(dest bactype.Device, rp bactype.ReadPropertyData) (bactype.ReadPropertyData, error) {
	id, err := c.tsm.GetFree()
	if err != nil {
		return bactype.ReadPropertyData{}, err
	}
	udp, err := c.LocalUDPAddress()
	if err != nil {
		return bactype.ReadPropertyData{}, err
	}
	src := bactype.UDPToAddress(udp)

	enc := encoding.NewEncoder()
	enc.NPDU(bactype.NPDU{
		Version:               bactype.ProtocolVersion,
		Destination:           &dest.Addr,
		Source:                &src,
		IsNetworkLayerMessage: false,
		ExpectingReply:        true,
		Priority:              bactype.Normal,
		HopCount:              bactype.DefaultHopCount,
	})

	enc.ReadProperty(uint8(id), rp)
	if enc.Error() != nil {
		return bactype.ReadPropertyData{}, err
	}

	// the value filled doesn't matter. it just needs to be non nil
	err = fmt.Errorf("go")
	for count := 0; err != nil && count < 2; count++ {
		var b []byte
		_, err = c.Send(dest.Addr, enc.Bytes())
		if err != nil {
			log.Print(err)
			continue
		}

		b, err = c.tsm.Receive(id, time.Duration(5)*time.Second)
		if err != nil {
			continue
		}
		var out bactype.ReadPropertyData
		dec := encoding.NewDecoder(b)

		var apdu bactype.APDU
		dec.APDU(&apdu)
		dec.ReadProperty(&out)
		if err = dec.Error(); err != nil {
			continue
		}
		return out, err
	}
	return bactype.ReadPropertyData{}, err
}

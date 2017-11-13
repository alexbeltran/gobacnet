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
	"log"

	"github.com/alexbeltran/gobacnet/encoding"
	bactype "github.com/alexbeltran/gobacnet/types"
)

//Close closes all inbound connections
func (c *Client) Close() {
	if c.listener == nil {
		return
	}

	c.listener.Close()
	c.listener = nil
}

func (c *Client) handleMsg(b []byte) {
	var header bactype.BVLC
	var npdu bactype.NPDU
	var apdu bactype.APDU

	dec := encoding.NewDecoder(b)
	err := dec.BVLC(&header)
	if err != nil {
		log.Print(err)
		return
	}

	if header.Function == bactype.BacFuncBroadcast || header.Function == bactype.BacFuncUnicast || header.Function == bactype.BacFuncForwardedNPDU {
		// Remove the header information
		b = b[mtuHeaderLength:]
		err = dec.NPDU(&npdu)
		if err != nil {
			return
		}

		if npdu.IsNetworkLayerMessage {
			//log.Print("Network Layer Message Discarded")
			return
		}

		// We want to keep the APDU intact so we will get a snapshot before decoding
		// further
		send := dec.Bytes()
		err = dec.APDU(&apdu)
		if err != nil {
			log.Print(err)
			return
		}
		switch apdu.DataType {
		case bactype.UnconfirmedServiceRequest:
			if apdu.UnconfirmedService == bactype.ServiceUnconfirmedIAm {
				log.Printf("I AM:%v", apdu.RawData)
				dec = encoding.NewDecoder(apdu.RawData)
				ids := make([]bactype.ObjectID, 1, 64)
				err = dec.IAm(ids)
				if err != nil {
					log.Print(err)
					return
				}
				for _, id := range ids {
					log.Printf("Instance: %d, Type: %d", id.Instance, id.Type)
				}
			} else if apdu.UnconfirmedService == bactype.ServiceUnconfirmedWhoIs {
				dec := encoding.NewDecoder(apdu.RawData)
				var low, high int32
				dec.WhoIs(&low, &high)
				log.Printf("WHO IS Request Low: %d, High:%d", low, high)
			} else {
				log.Printf("Unconfirmed: %d %v", apdu.UnconfirmedService, apdu.RawData)
			}
		case bactype.ComplexAck:
			err := c.tsm.Send(int(apdu.InvokeId), send)
			if err != nil {
				return
			}
		case bactype.ConfirmedServiceRequest:
			err := c.tsm.Send(int(apdu.InvokeId), send)
			if err != nil {
				return
			}
		default:
			// Ignore it
		}
	}

	if header.Function == bactype.BacFuncForwardedNPDU {
		// Right now we are ignoring the NPDU data that is stored in the packet. Eventually
		// we will need to check it for any additional information we can gleam.
		// NDPU has source
		b = b[forwardHeaderLength:]
	}

}

// Receive
func (c *Client) listen() error {
	for c.listener != nil {
		b := make([]byte, 1024)
		i, _, err := c.listener.ReadFrom(b)
		if err != nil {
			log.Println(err)
			continue
		}
		go c.handleMsg(b[:i])
	}
	return nil
}

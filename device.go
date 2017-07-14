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
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/alexbeltran/gobacnet/encoding"
	"github.com/alexbeltran/gobacnet/tsm"
	bactype "github.com/alexbeltran/gobacnet/types"
)

const DefaultStateSize = 20

type Client struct {
	Interface        *net.Interface
	MyAddress        string
	BroadcastAddress net.IP
	Port             int
	tsm              *tsm.TSM
	listener         *net.UDPConn
}

// getBroadcast uses the given address with subnet to return the broadcast address
func getBroadcast(addr string) (net.IP, error) {
	_, ipnet, err := net.ParseCIDR(addr)
	if err != nil {
		return net.IP{}, err
	}
	broadcast := net.IP(make([]byte, 4))
	for i := range broadcast {
		broadcast[i] = ipnet.IP[i] | ^ipnet.Mask[i]
	}
	return broadcast, nil
}

func NewClient(inter string) (*Client, error) {
	c := &Client{}
	i, err := net.InterfaceByName(inter)
	if err != nil {
		return c, err
	}
	c.Interface = i
	c.Port = defaultIPPort
	uni, err := i.Addrs()
	if err != nil {
		return c, err
	}

	if len(uni) == 0 {
		return c, fmt.Errorf("interface %s has no addresses", inter)
	}
	c.MyAddress = uni[0].String()

	broadcast, err := getBroadcast(uni[0].String())
	if err != nil {
		return c, err
	}
	c.BroadcastAddress = broadcast

	c.tsm = tsm.New(DefaultStateSize)

	src, err := c.LocalUDPAddress()
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", src)
	if err != nil {
		return nil, err
	}

	c.listener = conn

	go c.listen()
	return c, nil
}

func (c *Client) LocalAddress() (la bactype.Address, err error) {
	uni, err := c.Interface.Addrs()
	if err != nil {
		return
	}

	if len(uni) == 0 {
		err = fmt.Errorf("interface %s has no addresses", c.Interface.Name)
		return
	}
	ip, _, _ := net.ParseCIDR(c.MyAddress)

	buff := bytes.NewBuffer([]byte(ip))
	binary.Write(buff, encoding.EncodingEndian, c.Port)

	la.Adr = buff.Bytes()
	return la, nil
}

func (c *Client) LocalUDPAddress() (*net.UDPAddr, error) {
	/*
		addrs, _ := c.Interface.Addrs()
		var ip net.IP
		for _, ad := range addrs {
			switch v := ad.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			break
		}
	*/
	netstr := fmt.Sprintf("%s:%d", "0.0.0.0", c.Port)
	log.Printf(netstr)
	return net.ResolveUDPAddr("udp4", netstr)
}

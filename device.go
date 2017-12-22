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
	"os"
	"time"

	"github.com/alexbeltran/gobacnet/tsm"
	bactype "github.com/alexbeltran/gobacnet/types"
	"github.com/alexbeltran/gobacnet/utsm"
	"github.com/sirupsen/logrus"
)

const defaultStateSize = 20

type Client struct {
	netInterface     *net.Interface
	myAddress        string
	broadcastAddress net.IP
	port             int
	tsm              *tsm.TSM
	utsm             *utsm.Manager
	listener         *net.UDPConn
	log              *logrus.Logger
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

// NewClient creates a new client with the given interface and
// port.
func NewClient(inter string, port int) (*Client, error) {
	c := &Client{}
	i, err := net.InterfaceByName(inter)
	if err != nil {
		return c, err
	}
	c.netInterface = i
	if port == 0 {
		c.port = DefaultPort
	} else {
		c.port = port
	}
	uni, err := i.Addrs()
	if err != nil {
		return c, err
	}

	if len(uni) == 0 {
		return c, fmt.Errorf("interface %s has no addresses", inter)
	}

	// Clear out the value
	c.myAddress = ""
	// Find the first IP4 ip
	for _, adr := range uni {
		IP, _, _ := net.ParseCIDR(adr.String())

		// To4 is non nil when the type is ip4
		if IP.To4() != nil {
			c.myAddress = adr.String()
			break
		}
	}
	if len(c.myAddress) == 0 {
		// We couldn't find a interface or all of them are ip6
		return nil, fmt.Errorf("No valid broadcasting address was found on interface %s", inter)
	}

	broadcast, err := getBroadcast(c.myAddress)
	if err != nil {
		return c, err
	}
	c.broadcastAddress = broadcast

	c.tsm = tsm.New(defaultStateSize)
	options := []utsm.ManagerOption{
		utsm.DefaultSubscriberTimeout(time.Second * time.Duration(10)),
		utsm.DefaultSubscriberLastReceivedTimeout(time.Second * time.Duration(2)),
	}
	c.utsm = utsm.NewManager(options...)
	udp, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", c.port))
	conn, err := net.ListenUDP("udp", udp)
	if err != nil {
		return nil, err
	}

	c.listener = conn
	c.log = logrus.New()
	c.log.Formatter = &logrus.TextFormatter{}
	c.log.SetLevel(logrus.DebugLevel)

	// open a debug file
	f, err := os.Create("gobacnet.log")
	if err != nil {
		return c, fmt.Errorf("Could not create a log file")
	}
	c.log.Out = f

	// Print out relevant information
	c.log.Debug(fmt.Sprintf("Broadcast Address: %v", c.broadcastAddress))
	c.log.Debug(fmt.Sprintf("Local Address: %s", c.myAddress))
	c.log.Debug(fmt.Sprintf("Port: %x", c.port))
	go c.listen()
	return c, nil
}

func (c *Client) localAddress() (la bactype.Address, err error) {
	ip, _, _ := net.ParseCIDR(c.myAddress)
	ad := ip.To4()
	udp := net.UDPAddr{
		IP:   ad,
		Port: c.port,
	}
	la = bactype.UDPToAddress(&udp)
	return la, nil
}

func (c *Client) localUDPAddress() (*net.UDPAddr, error) {
	ip, _, _ := net.ParseCIDR(c.myAddress)
	netstr := fmt.Sprintf("%s:%d", ip.String(), c.port)
	return net.ResolveUDPAddr("udp4", netstr)
}

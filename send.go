package gobacnet

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/alexbeltran/gobacnet/encoding"
	bactype "github.com/alexbeltran/gobacnet/types"
)

// address returns the address given
func (c *Client) address(addr bactype.Address) (net.UDPAddr, error) {
	if addr.IsBroadcast() {
		return net.UDPAddr{
			IP:   c.BroadcastAddress,
			Port: c.Port,
		}, nil
	} else if addr.IsSubBroadcast() {
		// Network specific
		if addr.IsUnicast() {
			return addr.UDPAddr(), nil
		}

		// Broadcast
		return net.UDPAddr{
			IP:   c.BroadcastAddress,
			Port: c.Port,
		}, nil
	} else if addr.IsUnicast() {
		return addr.UDPAddr(), nil
	}
	return net.UDPAddr{}, fmt.Errorf("Unable to parse bacnet address")
}

// Sets the udp version used to transfer data
// See https://golang.org/pkg/net/#DialUDP
const udpVersion = "udp"
const mtuHeaderLength = 4
const forwardHeaderLength = 10

// Send packet to destination
func Send(dest bactype.Address, data []byte) (int, error) {
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
	d := dest.UDPAddr()

	// use default udp type, src = local address (nil)
	conn, err := net.DialUDP("udp", nil, &d)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(10) * time.Second))

	return conn.Write(e.Bytes())
}

//Close closes all inbound connections
func (c *Client) Close() {
	if c == nil {
		return
	}

	c.Close()
	c.listener = nil
}

// Receive
func (c *Client) listen() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: defaultIPPort,
	})
	if err != nil {
		return
	}

	c.listener = conn
	defer c.Close()

	var b []byte
	length, _, err := c.listener.ReadFromUDP(b)
	if err != nil {
		log.Print(err)
	}

	var header bactype.BVLC
	dec := encoding.NewDecoder(b)
	err = dec.BVLC(&header)
	if err != nil {
		return
	}

	/*
		if src.IP.Equal(net.ParseIP(conn.LocalAddr())) {
			// We accidentally got the packet back
			// It is not considered an error
			length = 0
			return
		}
	*/

	if header.Function == bactype.BacFuncBroadcast || header.Function == bactype.BacFuncUnicast {
		// Remove the header information
		b = b[mtuHeaderLength:]
		length = length - mtuHeaderLength
		return
	}

	if header.Function == bactype.BacFuncForwardedNPDU {
		// Right now we are ignoring the NPDU data that is stored in the packet. Eventually
		// we will need to check it for any additional information we can gleam.
		// NDPU has source
		b = b[forwardHeaderLength:]
		length = length - forwardHeaderLength
	}
	return
}

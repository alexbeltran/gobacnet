package gobacnet

import (
	"bytes"
	"encoding/binary"
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
	buff := new(bytes.Buffer)

	// Set packet type
	buff.WriteByte(typeBacnetIp)

	if dest.IsBroadcast() || dest.IsSubBroadcast() {
		// SET BROADCAST FLAG
		buff.WriteByte(byte(bacFuncBroadcast))
	} else {
		// SET UNICAST FLAG
		buff.WriteByte(byte(bacFuncUnicast))
	}

	// Write the length of the packet.
	// We add 2 to include this encoded 16 bit length
	l := uint16(buff.Len() + len(data) + 2)
	binary.Write(buff, encoding.EncodingEndian, l)

	// Write main data
	buff.Write(data)

	// Get IP Address
	d := dest.UDPAddr()

	// use default udp type, src = local address (nil)
	conn, err := net.DialUDP("udp", nil, &d)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(10) * time.Second))

	return conn.Write(buff.Bytes())
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

	var function bacFunc
	buff := bytes.NewBuffer(b)
	err = binary.Read(buff, encoding.EncodingEndian, function)
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

	if function == bacFuncBroadcast || function == bacFuncUnicast {
		// Remove the header information
		b = b[mtuHeaderLength:]
		length = length - mtuHeaderLength
		return
	}

	if function == bacFuncForwardedNPDU {
		b = b[forwardHeaderLength:]
		length = length - forwardHeaderLength
	}
	return
}

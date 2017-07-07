package gobacnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
const udpVersion = 'udp'

// Send packet to destination
func Send(dest bactype.Address, data []byte) error {
	buff := bytes.NewBuffer()

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
	// We add 2 to include this encoded piece and.... (idk?)
	l := uint16(buff.Writebuff.Len() + len(data) + 2)
	binary.Write(buff, encoding.EncodingEndian, l)

	// Write main data
	buff.Write(data)

	// Get IP Address
	d, err := dest.Address()
	if err != nil {
		return err
	}

	// use default udp type, src = local address (nil)
	conn, err := net.DialUDP("udp", nil, d)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(10) * time.Second))

	return conn.Write(buff.Bytes())
}

// Receive
func Receive(b []byte, deadline time.Time)(length int, src *net.IP, err error){
	conn, err := net.ListenUdp("udp", net.UDPAddr{
		Port: defaultIPPort,
	})
	if err != nil{
		return 
	}
	defer conn.Close()

	conn.SetReadDeadline(deadline)
	length, src, err := conn.ReadFromUDP(b)
	if err != nil{
		return
	}

	var function bacFunc 
	buff := bytes.NewBuffer(b)
	err = binary.Read(buff, encoding.EncodingEndian, function)
	if err != nil{
		return
	}

	

}

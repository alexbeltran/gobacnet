package gobacnet

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/alexbeltran/gobacnet/tsm"
)

const DefaultStateSize = 20

type Client struct {
	Interface        *net.Interface
	MyAddress        string
	BroadcastAddress string
	Port             uint16
	tsm              *tsm.TSM
}

// getBroadcast uses the given address with subnet to return the broadcast address
func getBroadcast(addr string) (string, error) {
	split := strings.Split(addr, "/")
	if len(split) != 2 {
		return "", fmt.Errorf("%s is not a valid address. Are you missing a subnet?", addr)
	}
	addr = split[0]
	subnet, err := strconv.Atoi(split[1])
	if err != nil {
		return "", err
	}

	// First we are going to convert the string address to 32 bit address
	parts := strings.Split(addr, ".")
	var b uint32
	b = 0
	for i, p := range parts {
		d, err := strconv.Atoi(p)
		if err != nil {
			return "", err
		}

		b = b + (uint32(d)&0xFF)<<(8*uint(3-i))
	}

	// Now we can apply the mask.
	mask := uint32(0xFFFFFFFF >> uint(subnet))
	b |= mask

	ip := make([]uint8, 4)
	for i := 0; i < 4; i++ {
		ip[i] = uint8(b >> uint8(24-8*i))
	}
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3]), nil
}

func NewClient(inter string) (*Client, error) {
	c := &Client{}
	i, err := net.InterfaceByName(inter)
	if err != nil {
		return c, err
	}
	c.Interface = i
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
	return c, nil
}

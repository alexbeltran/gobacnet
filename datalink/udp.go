package datalink

import (
	"fmt"
	"net"

	"github.com/alexbeltran/gobacnet/types"
)

// DefaultPort that BacnetIP will use if a port is not given. Valid ports for
// the bacnet protocol is between 0xBAC0 and 0xBAC9
const DefaultPort = 0xBAC0 //47808

type udpDataLink struct {
	netInterface                *net.Interface
	myAddress, broadcastAddress *types.Address
	port                        int
	listener                    *net.UDPConn
}

func NewUDPDataLink(inter string, port int) (DataLink, error) {
	i, err := net.InterfaceByName(inter)
	if err != nil {
		return nil, err
	}
	if port == 0 {
		port = DefaultPort
	}
	uni, err := i.Addrs()
	if err != nil {
		return nil, err
	}

	if len(uni) == 0 {
		return nil, fmt.Errorf("interface %s has no addresses", inter)
	}

	// Clear out the value
	var myAddress string
	// Find the first IP4 ip
	for _, adr := range uni {
		IP, _, _ := net.ParseCIDR(adr.String())

		// To4 is non nil when the type is ip4
		if IP.To4() != nil {
			myAddress = adr.String()
			break
		}
	}
	if len(myAddress) == 0 {
		// We couldn't find a interface or all of them are ip6
		return nil, fmt.Errorf("no valid broadcasting address was found on interface %s", inter)
	}

	ip, ipnet, err := net.ParseCIDR(myAddress)
	if err != nil {
		return nil, err
	}

	broadcast := net.IP(make([]byte, 4))
	for i := range broadcast {
		broadcast[i] = ipnet.IP[i] | ^ipnet.Mask[i]
	}

	udp, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	conn, err := net.ListenUDP("udp", udp)
	if err != nil {
		return nil, err
	}

	return &udpDataLink{
		listener:         conn,
		myAddress:        IPPortToAddress(ip, port),
		broadcastAddress: IPPortToAddress(broadcast, DefaultPort),
	}, nil
}

func (c *udpDataLink) Close() error {
	if c.listener != nil {
		return c.listener.Close()
	}
	return nil
}

func (c *udpDataLink) Receive(data []byte) (*types.Address, int, error) {
	n, adr, err := c.listener.ReadFromUDP(data)
	if err != nil {
		return nil, n, err
	}
	adr.IP = adr.IP.To4()
	return UDPToAddress(adr), n, nil
}

func (c *udpDataLink) GetMyAddress() *types.Address {
	return c.myAddress
}

// GetBroadcastAddress uses the given address with subnet to return the broadcast address
func (c *udpDataLink) GetBroadcastAddress() *types.Address {
	return c.broadcastAddress
}

func (c *udpDataLink) Send(data []byte, npdu *types.NPDU, dest *types.Address) (int, error) {
	// Get IP Address
	d, err := dest.UDPAddr()
	if err != nil {
		return 0, err
	}
	return c.listener.WriteTo(data, &d)
}

// IPPortToAddress converts a given udp address into a bacnet address
func IPPortToAddress(ip net.IP, port int) *types.Address {
	return UDPToAddress(&net.UDPAddr{
		IP:   ip.To4(),
		Port: port,
	})
}

// UDPToAddress converts a given udp address into a bacnet address
func UDPToAddress(n *net.UDPAddr) *types.Address {
	a := &types.Address{}
	p := uint16(n.Port)

	// Length of IP plus the port
	length := net.IPv4len + 2
	a.Mac = make([]uint8, length)
	//Encode ip
	for i := 0; i < net.IPv4len; i++ {
		a.Mac[i] = n.IP[i]
	}

	// Encode port
	a.Mac[net.IPv4len+0] = uint8(p >> 8)
	a.Mac[net.IPv4len+1] = uint8(p & 0x00FF)

	a.MacLen = uint8(length)
	return a
}

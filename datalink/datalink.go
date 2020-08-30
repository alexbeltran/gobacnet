package datalink

import (
	"github.com/alexbeltran/gobacnet/types"
)

type DataLink interface {
	GetMyAddress() *types.Address
	GetBroadcastAddress() *types.Address
	Send(data []byte, npdu *types.NPDU, dest *types.Address) (int, error)
	Receive(data []byte) (*types.Address, int, error)
	Close() error
}

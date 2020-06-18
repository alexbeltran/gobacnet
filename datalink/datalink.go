package datalink

import "github.com/alexbeltran/gobacnet/types"

type MessageHandler func(src *types.Address, data []byte)

type DataLink interface {
	GetMyAddress() *types.Address
	GetBroadcastAddress() *types.Address
	Run(handler MessageHandler)
	Send(data []byte, dest *types.Address) (int, error)
	Close() error
}

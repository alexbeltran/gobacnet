package types

type NPDUPriority byte

const (
	LifeSafety        NPDUPriority = 3
	CriticalEquipment NPDUPriority = 2
	Urgent            NPDUPriority = 1
	Normal            NPDUPriority = 0
)

type NPDU struct {
	Version uint8

	// Destination (optional)
	Destination *Address

	// Source (optional)
	Source *Address

	VendorId uint16

	IsNetworkLayerMessage   bool
	NetworkLayerMessageType uint8
	ExpectingReply          bool
	Priority                NPDUPriority
	HopCount                uint8
}

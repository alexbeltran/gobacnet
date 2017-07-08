package types

type NPDU struct {
	Version uint8
	// Info stores raw NPDU metadata
	Info byte

	// Destination (optional)
	Destination *Address

	// Source (optional)
	Source *Address

	VendorId uint16
}

type NPDUPriority byte

const (
	LifeSafety        NPDUPriority = 3
	CriticalEquipment NPDUPriority = 2
	Urgent            NPDUPriority = 1
	Normal            NPDUPriority = 0
)

const maskNetworkLayerMessage = 1 << 7
const maskDestination = 1 << 5
const maskSource = 1 << 3

// General setter for the info bits using the mask
func (n *NPDU) setInfoMask(b bool, mask byte) {
	if b {
		n.Info = n.Info | mask
	} else {
		var m byte = 0xFF
		m = m - mask
		n.Info = n.Info & m
	}
}

// CheckMask uses mask to check bit position
func (n *NPDU) checkMask(mask byte) bool {
	return (n.Info & mask) > 0

}

// IsNetworkLayerMessage returns true if it is a network layer message
func (n *NPDU) IsNetworkLayerMessage() bool {
	return n.checkMask(maskNetworkLayerMessage)
}

func (n *NPDU) SetNetworkLayerMessage(b bool) {
	n.setInfoMask(b, maskNetworkLayerMessage)
}

// Priority returns priority
func (n *NPDU) Priority() NPDUPriority {
	// Encoded in bit 0 and 1
	return NPDUPriority(n.Info & 3)
}

// SetPriority for NPDU
func (n *NPDU) SetPriority(p NPDUPriority) {
	// Clear the first two bits
	n.Info &= (0xF - 3)
	n.Info |= byte(p)
}

func (n *NPDU) HasDestination() bool {
	return n.checkMask(maskDestination)
}

func (n *NPDU) SetDestination(b bool) {
	n.setInfoMask(b, maskDestination)
}

func (n *NPDU) HasSource() bool {
	return n.checkMask(maskSource)
}

func (n *NPDU) SetSource(b bool) {
	n.setInfoMask(b, maskSource)
}

package types

// NPDUMetadata includes additional metadata about npdu message
type NPDUMetadata byte
type NPDUPriority byte

const (
	LifeSafety        NPDUPriority = 3
	CriticalEquipment NPDUPriority = 2
	Urgent            NPDUPriority = 1
	Normal            NPDUPriority = 0
)

type NPDU struct {
	Version uint8
	// Info stores raw NPDU metadata
	Info NPDUMetadata

	// Destination (optional)
	Destination *Address

	// Source (optional)
	Source *Address

	VendorId uint16
}

const maskNetworkLayerMessage = 1 << 7
const maskDestination = 1 << 5
const maskSource = 1 << 3

// General setter for the info bits using the mask
func (meta *NPDUMetadata) setInfoMask(b bool, mask byte) {
	if b {
		*meta = *meta | NPDUMetadata(mask)
	} else {
		var m byte = 0xFF
		m = m - mask
		*meta = *meta & NPDUMetadata(m)
	}
}

// CheckMask uses mask to check bit position
func (meta *NPDUMetadata) checkMask(mask byte) bool {
	return (*meta & NPDUMetadata(mask)) > 0

}

// IsNetworkLayerMessage returns true if it is a network layer message
func (n *NPDUMetadata) IsNetworkLayerMessage() bool {
	return n.checkMask(maskNetworkLayerMessage)
}

func (n *NPDUMetadata) SetNetworkLayerMessage(b bool) {
	n.setInfoMask(b, maskNetworkLayerMessage)
}

// Priority returns priority
func (n *NPDUMetadata) Priority() NPDUPriority {
	// Encoded in bit 0 and 1
	return NPDUPriority(byte(*n) & 3)
}

// SetPriority for NPDU
func (n *NPDUMetadata) SetPriority(p NPDUPriority) {
	// Clear the first two bits
	*n &= (0xF - 3)
	*n |= NPDUMetadata(p)
}

func (n *NPDUMetadata) HasDestination() bool {
	return n.checkMask(maskDestination)
}

func (n *NPDUMetadata) SetDestination(b bool) {
	n.setInfoMask(b, maskDestination)
}

func (n *NPDUMetadata) HasSource() bool {
	return n.checkMask(maskSource)
}

func (n *NPDUMetadata) SetSource(b bool) {
	n.setInfoMask(b, maskSource)
}

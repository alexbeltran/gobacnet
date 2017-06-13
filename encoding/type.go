package encoding

type ReadPropertyData struct {
	ObjectType         uint16
	ObjectInstance     uint32
	ObjectProperty     uint32
	ArrayIndex         uint32
	ApplicationData    []uint8
	ApplicationDataLen int
	ErrorClass         uint8
	ErrorCode          uint8
}

package encoding

import bactype "github.com/alexbeltran/gobacnet/types"

// WriteProperty encodes a write property request
func (e *Encoder) WriteProperty(invokeID uint8, data bactype.ReadPropertyData, priority uint) error {
	a := bactype.APDU{
		DataType: bactype.ConfirmedServiceRequest,
		Service:  bactype.ServiceConfirmedWriteProperty,
		MaxSegs:  0,
		MaxApdu:  MaxAPDU,
		InvokeId: invokeID,
	}
	e.APDU(a)

	tagID, err := e.readPropertyHeader(0, data)
	if err != nil {
		return err
	}

	// Tag 3 - the value (unlike other values, this is just a raw byte array)
	e.openingTag(tagID)
	e.AppData(data.Object.Properties[0].Data)
	e.closingTag(tagID)

	tagID++

	// Tag 4 - Optional priorty tag
	// Priority set
	if priority != 0 {
		e.contextUnsigned(tagID, uint32(priority))
	}
	return e.Error()
}

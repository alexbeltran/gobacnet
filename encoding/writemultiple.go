package encoding

import bactype "github.com/alexbeltran/gobacnet/types"

// WriteProperty encodes a write property request
func (e *Encoder) WriteMultiProperty(invokeID uint8, data bactype.MultiplePropertyData) error {
	a := bactype.APDU{
		DataType: bactype.ConfirmedServiceRequest,
		Service:  bactype.ServiceConfirmedWritePropMultiple,
		MaxSegs:  0,
		MaxApdu:  MaxAPDU,
		InvokeId: invokeID,
	}
	e.APDU(a)

	err := e.objects(data.Objects, true)
	if err != nil {
		return err
	}

	return e.Error()
}

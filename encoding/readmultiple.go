package encoding

import (
	bactype "github.com/alexbeltran/gobacnet/types"
)

func (e *Encoder) ReadMultipleProperty(invokeID uint8, data bactype.MultiplePropertyData) error {
	a := bactype.APDU{
		DataType:         bactype.ConfirmedServiceRequest,
		Service:          bactype.ServiceConfirmedReadPropMultiple,
		MaxSegs:          0,
		MaxApdu:          MaxAPDU,
		InvokeId:         invokeID,
		SegmentedMessage: false,
	}
	e.APDU(a)
	err := e.objects(data.Objects, false)
	if err != nil {
		return err
	}

	return e.Error()
}

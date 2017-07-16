package encoding

import (
	bactype "github.com/alexbeltran/gobacnet/types"
)

func (e *Encoder) WhoIs(low, high int32) error {
	apdu := bactype.APDU{
		DataType:           bactype.UnconfirmedServiceRequest,
		UnconfirmedService: bactype.ServiceUnconfirmedWhoIs,
	}
	e.write(apdu.DataType)
	e.write(apdu.UnconfirmedService)

	// The range is optional. A scan for all objects is done when either low/high
	// are negative or when we are scanning above the max instance
	if low >= 0 && high >= 0 && low < bactype.MaxInstance && high <
		bactype.MaxInstance {
		// Tag 0
		e.contextUnsigned(0, uint32(low))

		// Tag 1
		e.contextUnsigned(1, uint32(high))
	}
	return e.Error()
}

func (d *Decoder) WhoIs(low, high *int32) error {
	// APDU read in a higher level
	if d.len() == 0 {
		*low = bactype.WhoIsAll
		*high = bactype.WhoIsAll
		return nil
	}
	// Tag 0 - Low Value
	var expectedTag uint8
	tag, _, value := d.tagNumberAndValue()
	if tag != expectedTag {
		return &ErrorIncorrectTag{Expected: expectedTag, Given: tag}
	}
	l := d.unsigned(int(value))
	*low = int32(l)

	// Tag 1 - High Value
	expectedTag = 1
	tag, _, value = d.tagNumberAndValue()
	if tag != expectedTag {
		return &ErrorIncorrectTag{Expected: expectedTag, Given: tag}
	}
	h := d.unsigned(int(value))
	*high = int32(h)

	return d.Error()
}

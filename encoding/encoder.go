package encoding

import (
	"bytes"
	"encoding/binary"
)

var encodingEndian binary.ByteOrder = binary.BigEndian

type Encoder struct {
	buff *bytes.Buffer
	err  error
}

func NewEncoder() *Encoder {
	e := Encoder{
		buff: new(bytes.Buffer),
		err:  nil,
	}
	return &e
}

func (e *Encoder) Error() error {
	return e.err
}

func (e *Encoder) Bytes() []byte {
	return e.buff.Bytes()
}

func (e *Encoder) write(p interface{}) {
	if e.err != nil {
		return
	}
	e.err = binary.Write(e.buff, encodingEndian, p)
}

func (e *Encoder) readProperty(invokeID uint8, data ReadPropertyData) {
	e.write(pduTypeConfirmedServiceRequest)
	e.write(encodeMaxSegsMaxApdu(0, maxApdu))
	e.write(invokeID)
	e.write(ReadPropertyService)

	var tagCounter uint8 = 0
	if data.ObjectType < MaxObject {
		e.contextObjectID(tagCounter, data.ObjectType, data.ObjectInstance)
		tagCounter++
	}

	if data.ObjectProperty < MaxPropertyID {
		e.contextEnumerated(tagCounter, data.ObjectProperty)
		tagCounter++
	}

	if data.ArrayIndex != ArrayAll {
		e.contextUnsigned(tagCounter, data.ArrayIndex)
		tagCounter++
	}
	return
}

func (e *Encoder) contextObjectID(tagNum uint8, objectType uint16, instance uint32) {
	/* length of object id is 4 octets, as per 20.2.14 */
	e.tag(tagNum, true, 4)
	e.objectId(objectType, instance)
}

func (e *Encoder) tag(tagNum uint8, contextSpecific bool, lenValueType uint32) {
	var t uint8 = 0
	if contextSpecific {
		t = setContextSpecific(t)
	}

	// I have no idea why this is here.
	if lenValueType <= 4 {
		// TODO: I typecasted this here, but I am not too sure if this is correct
		// since the original code used a 32 bit ORed to a 8 bit array
		t |= uint8(lenValueType)
	} else {
		t |= 5
	}

	// We have enough room to put it with the last value
	if tagNum <= 14 {
		t |= (tagNum << 4)
		e.write(t)

		// We don't have enough space so make it in a new byte
	} else {
		t |= 0xF0
		e.write(t)
		e.write(tagNum)
	}
	if lenValueType > 4 {
		// Depending on the length, we will either write it as an 8 bit, 32 bit, or 64 bit integer
		if lenValueType <= 253 {
			e.write(uint8(lenValueType))
		} else if lenValueType <= 65535 {
			e.write(254)
			e.write(uint16(lenValueType))
		} else {
			e.write(255)
			e.write(lenValueType)
		}
	}
}

/* from clause 20.2.14 Encoding of an Object Identifier Value
returns the number of apdu bytes consumed */
func (e *Encoder) objectId(objectType uint16, instance uint32) {
	var value uint32
	value = ((uint32(objectType) & MaxObject) << InstanceBits) | (instance & MaxInstance)
	e.write(value)
}

func (e *Encoder) contextEnumerated(tagNumber uint8, value uint32) {
	e.contextUnsigned(tagNumber, value)
}

func (e *Encoder) contextUnsigned(tagNumber uint8, value uint32) {
	len := valueLength(value)
	e.tag(tagNumber, true, uint32(len))
	e.unsigned(value)
}

func (e *Encoder) enumerated(value uint32) {
	e.unsigned(value)
}

// weird, huh?
func (e *Encoder) unsigned24(value uint32) {
	e.write(uint8((value & 0xFF0000) >> 16))
	e.write(uint8((value & 0x00FF00) >> 8))
	e.write(uint8(value & 0x0000FF))

}

func (e *Encoder) unsigned(value uint32) {
	if value < 0x100 {
		e.write(uint8(value))
	} else if value < 0x10000 {
		e.write(uint16(value))
	} else if value < 0x1000000 {
		// Really!? 24 bits?
		e.unsigned24(value)
	} else {
		e.write(value)
	}
}

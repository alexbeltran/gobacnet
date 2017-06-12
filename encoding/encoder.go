package encoding

import (
	"bytes"
	"encoding/binary"
)

var encodingEndian binary.ByteOrder = binary.BigEndian

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

type Encoder struct {
	buff *bytes.Buffer
}

func NewEncoder() *Encoder {
	e := Encoder{}
	e.buff = new(bytes.Buffer)
	return &e
}

func (e *Encoder) Bytes() []byte {
	return e.buff.Bytes()
}

func (e *Encoder) readProperty(in []byte, invokeID uint8, data ReadPropertyData) (b []byte, err error) {
	buff := bytes.NewBuffer(in)
	write := func(p interface{}) {
		if err != nil {
			return
		}
		err = binary.Write(buff, binary.LittleEndian, p)
	}
	write(pduTypeConfirmedServiceRequest)
	write(encodeMaxSegsMaxApdu(0, maxApdu))
	write(invokeID)
	write(ReadPropertyService)
	if data.ObjectType < MaxObject {
		e.contextObjectID(0, data.ObjectType, data.ObjectInstance)
	}

	if data.ObjectProperty < MaxPropertyID {
		e.contextEnumerated(1, data.ObjectProperty)
	}

	return buff.Bytes(), err
}

func (e *Encoder) contextObjectID(tagNum uint8, objectType uint16, instance uint32) error {
	/* length of object id is 4 octets, as per 20.2.14 */
	err := e.tag(tagNum, true, 4)
	if err != nil {
		return err
	}
	return e.objectId(objectType, instance)
}

func (e *Encoder) write(p interface{}) error {
	return binary.Write(e.buff, encodingEndian, p)
}

func (e *Encoder) tag(tagNum uint8, contextSpecific bool, lenValueType uint32) (err error) {
	var t uint8 = 0
	write := func(p interface{}) {
		if err != nil {
			return
		}
		err = e.write(p)
	}

	if contextSpecific {
		t = 1 << 3
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
		t |= tagNum
		write(t)

		// We don't have enough space so make it in a new byte
	} else {
		t |= 0xF0
		write(t)
		write(tagNum)
	}
	if lenValueType > 4 {
		// Depending on the length, we will either write it as an 8 bit, 32 bit, or 64 bit integer
		if lenValueType <= 253 {
			write(uint8(lenValueType))
		} else if lenValueType <= 65535 {
			write(254)
			write(uint16(lenValueType))
		} else {
			write(255)
			write(lenValueType)
		}
	}
	return nil
}

/* from clause 20.2.14 Encoding of an Object Identifier Value
returns the number of apdu bytes consumed */
func (e *Encoder) objectId(objectType uint16, instance uint32) error {
	var value uint32 = 0
	value = ((uint32(objectType) & MaxObject) << InstanceBits) | (instance & MaxInstance)
	return e.write(value)
}

func (e *Encoder) contextEnumerated(tagNumber uint8, value uint32) error {
	var len int
	/* length of enumerated is variable, as per 20.2.11 */
	if value < 0x100 {
		len = 1
	} else if value < 0x10000 {
		len = 2
	} else if value < 0x1000000 {
		len = 3
	} else {
		len = 4
	}

	err := e.tag(tagNumber, true, uint32(len))
	if err != nil {
		return err
	}

	e.enumerated(value)
	return nil
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

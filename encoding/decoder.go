package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Decoder used
type Decoder struct {
	buff *bytes.Buffer
}

func NewDecoder(b []byte) *Decoder {
	return &Decoder{
		bytes.NewBuffer(b),
	}
}

func (d *Decoder) decode(data interface{}) error {
	return binary.Read(d.buff, encodingEndian, data)
}

func (d *Decoder) readProperty(in []byte) error {
	// Must have at least 7 bytes
	if len(in) < 7 {
		return fmt.Errorf("Missing parameters")
	}

	d.buff = bytes.NewBuffer(in)

	return nil
}

// contexTag decoder

// Returns both a tag and additional meta data stored in this byte. If it is of
// extended type, then that means that the entire first byte is metadata, else
// the firrst 4 bytes store the tag
func (d *Decoder) tagNumber() (tag uint8, meta uint8) {
	// Read the first value
	d.decode(meta)
	if isExtendedTagNumber(meta) {
		d.decode(tag)
		return tag, meta
	}
	return meta >> 4, meta
}

func (d *Decoder) tagNumberAndValue() (tag uint8, value uint32) {
	tag, meta := d.tagNumber()
	if isExtendedTagNumber(meta) {
		var val uint8
		d.decode(&val)
		// Tagged as an uint32
		if val == 255 {
			var parse uint32
			d.decode(parse)
			return tag, parse

			// Tagged as a uint16
		} else if val == 254 {
			var parse uint16
			d.decode(parse)
			return tag, uint32(parse)

			// No tag, it must be a uint8
		} else {
			return tag, uint32(val)
		}
	} else if isOpeningTag(meta) || isClosingTag(meta) {
		return tag, 0
	}
	// It must be a non extended/small value
	return tag, uint32(meta & 0x07)
}

func (d *Decoder) objectId() (objectType uint16, instance uint32, err error) {
	var value uint32
	err = d.decode(&value)
	objectType = uint16((value >> InstanceBits) & MaxObject)
	instance = value & MaxInstance
	return
}

func isExtendedTagNumber(x uint8) bool {
	return ((x & 0xF0) == 0xF0)
}

/* from clause 20.2.1.3.2 Constructed Data */
/* true if the tag is an opening tag */
func isOpeningTag(x uint8) bool {
	return ((x & 0x07) == 6)
}

/* from clause 20.2.1.3.2 Constructed Data */
/* true if the tag is a closing tag */
func isClosingTag(x uint8) bool {
	return ((x & 0x07) == 7)
}

func (d *Decoder) enumerated(len int) uint32 {
	return d.unsigned(len)
}

func (d *Decoder) unsigned24() uint32 {
	var a, b, c uint8
	d.decode(&a)
	d.decode(&b)
	d.decode(&c)

	var x uint32
	x = uint32((uint32(a) << 16) & 0x00ff0000)
	x |= uint32((uint32(b) << 8) & 0x0000ff00)
	x |= uint32(uint32(c) & 0x000000ff)
	return x
}

func (d *Decoder) unsigned(length int) uint32 {
	switch length {
	case 1:
		var val uint8
		d.decode(&val)
		return uint32(val)
	case 2:
		var val uint16
		d.decode(&val)
		return uint32(val)
	case 3:
		return d.unsigned24()
	case 4:
		var val uint32
		d.decode(&val)
		return val
	default:
		return 0
	}
}

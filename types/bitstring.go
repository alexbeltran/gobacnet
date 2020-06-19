package types

import "encoding/json"

const MaxBitStringBytes = 15

type BitString struct {
	bitUsed uint8
	value   []byte
}

func NewBitString(bufferSize int) *BitString {
	if bufferSize > MaxBitStringBytes {
		bufferSize = MaxBitStringBytes
	}
	return &BitString{
		bitUsed: 0,
		value:   make([]byte, bufferSize),
	}
}

func (bs *BitString) Value() []bool {
	value := make([]bool, bs.bitUsed)
	for i := uint8(0); i <= bs.bitUsed; i++ {
		if bs.Bit(i) {
			value[i] = true
		} else {
			value[i] = false
		}
	}
	return value
}

func (bs *BitString) String() string {
	bin, _ := json.Marshal(bs.Value())
	return string(bin)
}

func (bs *BitString) SetBit(bitNumber uint8, value bool) {
	byteNumber := bitNumber / 8
	var bitMask uint8 = 1
	if byteNumber < MaxBitStringBytes {
		/* set max bits used */
		if bs.bitUsed < (bitNumber + 1) {
			bs.bitUsed = bitNumber + 1
		}
		bitMask = bitMask << (bitNumber - (byteNumber * 8))
		if value {
			bs.value[byteNumber] |= bitMask
		} else {
			bs.value[byteNumber] &= ^bitMask
		}
	}
}

func (bs *BitString) Bit(bitNumber uint8) bool {
	byteNumber := bitNumber / 8
	bitMask := uint8(1)
	if bitNumber < (MaxBitStringBytes * 8) {
		bitMask = bitMask << (bitNumber - (byteNumber * 8))
		return (bs.value[byteNumber] & bitMask) != 0
	}
	return false
}

func (bs *BitString) BitUsed() uint8 {
	return bs.bitUsed
}

func (bs *BitString) BytesUsed() uint8 {
	if bs != nil && bs.bitUsed > 0 {
		return (bs.bitUsed-1)/8 + 1
	}
	return 0
}

func (bs *BitString) Byte(index uint8) byte {
	if bs != nil && index < MaxBitStringBytes {
		return bs.value[index]
	}
	return 0
}

func (bs *BitString) SetByte(index uint8, value byte) bool {
	if bs != nil && index < MaxBitStringBytes {
		bs.value[index] = value
		return true
	}
	return false
}

func (bs *BitString) SetBitsUsed(byteUsed uint8, bitsUnused uint8) bool {
	if bs != nil {
		bs.bitUsed = byteUsed*8 - bitsUnused
		return true
	}
	return false
}

func (bs *BitString) BitsCapacity() uint8 {
	if bs != nil {
		return uint8(len(bs.value) * 8)
	}
	return 0
}

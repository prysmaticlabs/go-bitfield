package bitfield

import (
	"encoding/binary"
	"math/bits"
)

var _ = Bitfield(Bitvector4{})
var _ = Bitfield(Bitvector64{})

// Bitvector4 is a bitfield with a known size of 4. There is no length bit
// present in the underlying byte array.
type Bitvector4 []byte

// Bitvector64 is a bitfield with a fixed defined size of 64. There is no length bit
// present in the underlying byte array.
type Bitvector64 []byte

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitvector, then this method returns false.
func (b Bitvector4) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	if idx >= b.Len() {
		return false
	}

	i := uint8(1 << idx)
	return b[0]&i == i

}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitvector, then this method returns
// false.
func (b Bitvector4) SetBitAt(idx uint64, val bool) {
	// Out of bounds, do nothing.
	if idx >= b.Len() {
		return
	}

	bit := uint8(1 << idx)
	if val {
		b[0] |= bit
	} else {
		b[0] &^= bit
	}
}

// Len returns a constant length 4.
func (b Bitvector4) Len() uint64 {
	return 4
}

// Count returns the number of 1s in the bitvector.
func (b Bitvector4) Count() uint64 {
	if len(b) == 0 {
		return 0
	}
	return uint64(bits.OnesCount8(b.Bytes()[0]))
}

// Bytes returns the bytes data representing the bitvector4. This method
// bitmasks the underlying data to ensure that it is an accurate representation.
func (b Bitvector4) Bytes() []byte {
	if len(b) == 0 {
		return []byte{}
	}
	return []byte{b[0] & 0x0F}
}

// Shift bitvector by i. If i >= 0, perform left shift, otherwise right shift.
func (b Bitvector4) Shift(i int) {
	if len(b) == 0 {
		return
	}

	// Shifting greater than 4 bits is pointless and can have unexpected behavior.
	if i > 4 {
		i = 4
	} else if i < -4 {
		i = -4
	}

	if i >= 0 {
		b[0] <<= uint8(i)
	} else {
		b[0] >>= uint8(i * -1)
	}
	b[0] &= 0x0F
}

// NewBitvector64 creates a new bitvector of size 64.
func NewBitvector64() Bitvector64 {
	byteArray := [8]byte{}
	return byteArray[:]
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitvector, then this method returns false.
func (b Bitvector64) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	if idx >= b.Len() {
		return false
	}

	i := uint8(1 << (idx % 8))
	return b[idx/8]&i == i
}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitvector, then this method returns
// false.
func (b Bitvector64) SetBitAt(idx uint64, val bool) {
	// Out of bounds, do nothing.
	if idx >= b.Len() {
		return
	}

	bit := uint8(1 << (idx % 8))
	if val {
		b[idx/8] |= bit
	} else {
		b[idx/8] &^= bit
	}
}

// Len returns a constant length 64.
func (b Bitvector64) Len() uint64 {
	return uint64(len(b) * 8)
}

// Count returns the number of 1s in the bitvector.
func (b Bitvector64) Count() uint64 {
	if len(b) == 0 {
		return 0
	}
	c := 0
	for _, bt := range b {
		c += bits.OnesCount8(bt)
	}
	return uint64(c)
}

// Bytes returns the bytes data representing the bitvector64. This method
// bitmasks the underlying data to ensure that it is an accurate representation.
func (b Bitvector64) Bytes() []byte {
	if len(b) == 0 {
		return []byte{}
	}
	ret := make([]byte, len(b))
	copy(ret, b[:])
	return ret[:]
}

// Shift bitvector by i. If i >= 0, perform left shift, otherwise right shift.
func (b Bitvector64) Shift(i int) {
	if len(b) == 0 {
		return
	}

	// Shifting greater than 64 bits is pointless and can have unexpected behavior.
	if i > 64 {
		i = 64
	} else if i < -64 {
		i = -64
	}
	if i >= 0 {
		num := binary.BigEndian.Uint64(b)
		num <<= uint8(i)
		binary.BigEndian.PutUint64(b, num)
	} else {
		num := binary.BigEndian.Uint64(b)
		num >>= uint8(i * -1)
		binary.BigEndian.PutUint64(b, num)
	}
}

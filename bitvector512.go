package bitfield

import (
	"encoding/binary"
	"math/bits"
)

var _ = Bitfield(Bitvector512{})

// Bitvector512 is a bitfield with a fixed defined size of 512. There is no length bit
// present in the underlying byte array.
type Bitvector512 []byte

// NewBitvector512 creates a new bitvector of size 512.
func NewBitvector512() Bitvector512 {
	byteArray := [64]byte{}
	return byteArray[:]
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitvector, then this method returns false.
func (b Bitvector512) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	if idx >= b.Len() || idx >= 512 {
		return false
	}

	i := uint8(1 << (idx % 8))
	return b[idx/8]&i == i
}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitvector, then this method returns
// false.
func (b Bitvector512) SetBitAt(idx uint64, val bool) {
	// Out of bounds, do nothing.
	if idx >= b.Len() || idx >= 512 {
		return
	}

	bit := uint8(1 << (idx % 8))
	if val {
		b[idx/8] |= bit
	} else {
		b[idx/8] &^= bit
	}
}

// Len returns the number of bits in the bitvector.
func (b Bitvector512) Len() uint64 {
	return uint64(len(b) * 8)
}

// Count returns the number of 1s in the bitvector.
func (b Bitvector512) Count() uint64 {
	if len(b) == 0 {
		return 0
	}
	c := 0
	for _, bt := range b {
		c += bits.OnesCount8(bt)
	}
	return uint64(c)
}

// Bytes returns the bytes data representing the Bitvector512. This method
// bitmasks the underlying data to ensure that it is an accurate representation.
func (b Bitvector512) Bytes() []byte {
	if len(b) == 0 {
		return []byte{}
	}
	ret := make([]byte, len(b))
	copy(ret, b[:])
	return ret[:]
}

// Shift bitvector by i. If i >= 0, perform left shift, otherwise right shift.
func (b Bitvector512) Shift(i int) {
	if len(b) == 0 {
		return
	}

	// Shifting greater than 512 bits is pointless and can have unexpected behavior.
	if i > 512 {
		i = 512
	} else if i < -512 {
		i = -512
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

// BitIndices returns the list of indices that are set to 1.
func (b Bitvector512) BitIndices() []int {
	indices := make([]int, 0, 512)
	for i, bt := range b {
		for j := 0; j < 8; j++ {
			bit := byte(1 << uint(j))
			if bt&bit == bit {
				indices = append(indices, i*8+j)
			}
		}
	}

	return indices
}

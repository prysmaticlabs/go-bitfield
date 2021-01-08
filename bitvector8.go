package bitfield

import (
	"math/bits"
)

var _ = Bitfield(Bitvector8{})

// Bitvector8 is a bitfield with a fixed defined size of 8. There is no length bit
// present in the underlying byte array.
type Bitvector8 []byte

// NewBitvector8 creates a new bitvector of size 8.
func NewBitvector8() Bitvector8 {
	byteArray := [1]byte{}
	return byteArray[:]
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitvector, then this method returns false.
func (b Bitvector8) BitAt(idx uint64) bool {
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
func (b Bitvector8) SetBitAt(idx uint64, val bool) {
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

// Len returns the number of bits in the bitvector.
func (b Bitvector8) Len() uint64 {
	return 8
}

// Count returns the number of 1s in the bitvector.
func (b Bitvector8) Count() uint64 {
	if len(b) == 0 {
		return 0
	}
	return uint64(bits.OnesCount8(b[0])
}

// Bytes returns the bytes data representing the Bitvector8. This method
// bitmasks the underlying data to ensure that it is an accurate representation.
func (b Bitvector8) Bytes() []byte {
	if len(b) == 0 {
		return []byte{}
	}
	return []byte{b[0]}
}

// BitIndices returns the list of indices that are set to 1.
func (b Bitvector8) BitIndices() []int {
	indices := make([]int, 0, 8)
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

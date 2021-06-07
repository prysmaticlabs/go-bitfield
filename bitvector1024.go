package bitfield

import (
	"encoding/binary"
	"math/bits"
)

var _ = Bitfield(Bitvector1024{})

// Bitvector1024 is a bitfield with a fixed defined size of 1024. There is no length bit
// present in the underlying byte array.
type Bitvector1024 []byte

const bitvector1024ByteSize = 128
const bitvector1024BitSize = bitvector1024ByteSize * 8

// NewBitvector1024 creates a new bitvector of size 1024.
func NewBitvector1024() Bitvector1024 {
	byteArray := [bitvector1024ByteSize]byte{}
	return byteArray[:]
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitvector, then this method returns false.
func (b Bitvector1024) BitAt(idx uint64) bool {
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
func (b Bitvector1024) SetBitAt(idx uint64, val bool) {
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

// Len returns the number of bits in the bitvector.
func (b Bitvector1024) Len() uint64 {
	return bitvector1024BitSize
}

// Count returns the number of 1s in the bitvector.
func (b Bitvector1024) Count() uint64 {
	if len(b) == 0 {
		return 0
	}
	c := 0
	for i, bt := range b {
		if i >= bitvector1024ByteSize {
			break
		}
		c += bits.OnesCount8(bt)
	}
	return uint64(c)
}

// Bytes returns the bytes data representing the Bitvector1024.
func (b Bitvector1024) Bytes() []byte {
	if len(b) == 0 {
		return []byte{}
	}
	ln := min(len(b), bitvector1024ByteSize)
	ret := make([]byte, ln)
	copy(ret, b[:ln])
	return ret[:]
}

// Shift bitvector by i. If i >= 0, perform left shift, otherwise right shift.
func (b Bitvector1024) Shift(i int) {
	if len(b) == 0 {
		return
	}

	// Shifting greater than 1024 bits is pointless and can have unexpected behavior.
	if i > bitvector1024BitSize {
		i = bitvector1024BitSize
	} else if i < -bitvector1024BitSize {
		i = -bitvector1024BitSize
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
func (b Bitvector1024) BitIndices() []int {
	indices := make([]int, 0, bitvector1024BitSize)
	for i, bt := range b {
		if i >= bitvector1024ByteSize {
			break
		}
		for j := 0; j < 8; j++ {
			bit := byte(1 << uint(j))
			if bt&bit == bit {
				indices = append(indices, i*8+j)
			}
		}
	}

	return indices
}

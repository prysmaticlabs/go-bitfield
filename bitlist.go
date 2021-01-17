package bitfield

import (
	"math/bits"
)

const (
	// wordSize configures how many bits are there in a single element of bitlist array.
	wordSize = uint64(64)
	// wordSizeLog2 allows optimized division by wordSize using right shift (numBits >> wordSizeLog2).
	// Note: log_2(64) = 6.
	wordSizeLog2 = uint64(6)
	// bytesInWord defines how many bytes are there in a single word i.e. wordSize/8.
	bytesInWord = 8
)

// Bitlist is a bitfield implementation backed by an array of uint64.
type Bitlist struct {
	size uint64
	data []uint64
}

// NewBitlist creates a new bitlist of size N.
func NewBitlist(n uint64) *Bitlist {
	return &Bitlist{
		size: n,
		data: make([]uint64, numWordsRequired(n)),
	}
}

// NewBitlistFrom creates a new bitlist for a given uint64 array.
func NewBitlistFrom(data []uint64) *Bitlist {
	return &Bitlist{
		size: uint64(len(data)) * wordSize,
		data: data,
	}
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitlist, then this method returns false.
func (b *Bitlist) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	if idx >= b.size {
		return false
	}

	i := uint64(1 << (idx % wordSize))
	return b.data[idx>>wordSizeLog2]&i == i
}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitlist, then this method returns
// false.
func (b *Bitlist) SetBitAt(idx uint64, val bool) {
	// Out of bounds, do nothing.
	if idx >= b.size {
		return
	}

	bit := uint64(1 << (idx % wordSize))
	if val {
		b.data[idx>>wordSizeLog2] |= bit
	} else {
		b.data[idx>>wordSizeLog2] &^= bit
	}
}

// Len returns the number of bits in a bitlist (note that underlying array can be bigger).
func (b *Bitlist) Len() uint64 {
	return b.size
}

// Bytes returns underlying array of uint64s as an array of bytes.
// The leading zeros in the bitlist will be trimmed to the smallest byte length
// representation of the bitlist. This may produce an empty byte slice if all
// bits were zero.
func (b *Bitlist) Bytes() []byte {
	if len(b.data) == 0 {
		return []byte{}
	}

	ret := make([]byte, len(b.data)*bytesInWord)
	for idx, word := range b.data {
		for i := uint64(0); i < 8; i++ {
			ret[uint64(idx)*bytesInWord+i] = byte((word >> (i * 8)) & 0xff)
		}
	}

	// Clear any leading zero bytes.
	newLen := len(ret)
	for i := len(ret) - 1; i >= 0; i-- {
		if ret[i] != 0x00 {
			break
		}
		newLen = i
	}

	return ret[:newLen]
}

// Count returns the number of 1s in the bitlist.
func (b *Bitlist) Count() uint64 {
	c := 0
	for _, bt := range b.data {
		c += bits.OnesCount64(bt)
	}
	return uint64(c)
}

// numWordsRequired calculates how many words are required to hold bitlist of n bits.
func numWordsRequired(n uint64) int {
	return int((n + (wordSize - 1)) >> wordSizeLog2)
}

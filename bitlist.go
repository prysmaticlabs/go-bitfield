package bitfield

import (
	"bytes"
	"encoding/binary"
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
	// bytesInWordLog2 = log_2(8)
	bytesInWordLog2 = 3
	// allBitsSet is a word with all bits set.
	allBitsSet = uint64(0xffffffffffffffff)
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

	buf := bytes.NewBuffer(make([]byte, 0, len(b.data)*bytesInWord))
	if err := binary.Write(buf, binary.LittleEndian, b.data); err != nil {
		return []byte{}
	}
	ret := buf.Bytes()

	// Clear any leading zero bytes.
	allLeadingZeroes := 0
	for i := len(b.data) - 1; i >= 0; i-- {
		leadingZeroes := bits.LeadingZeros64(b.data[i])
		allLeadingZeroes += leadingZeroes
		// If the whole word is 0x0, allow to test the next word, break otherwise.
		if uint64(leadingZeroes) != wordSize {
			break
		}
	}

	return ret[:len(ret)-allLeadingZeroes>>bytesInWordLog2]
}

// Count returns the number of 1s in the bitlist.
func (b *Bitlist) Count() uint64 {
	c := 0
	for _, bt := range b.data {
		c += bits.OnesCount64(bt)
	}
	return uint64(c)
}

// Contains returns true if the bitlist contains all of the bits from the provided argument
// bitlist i.e. if `b` is a superset of `c`.
// This method will panic if bitlists are not the same length.
func (b *Bitlist) Contains(c *Bitlist) bool {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	// To ensure all of the bits in c are present in b, we iterate over every word, combine
	// the words from b and c, then XOR them against b. If the result of this is non-zero, then we
	// are assured that a word in c had bits not present in word in b.
	for idx, word := range b.data {
		if word^(word|c.data[idx]) != 0 {
			return false
		}
	}
	return true
}

// Overlaps returns true if the bitlist contains one of the bits from the provided argument
// bitlist. This method will panic if bitlists are not the same length.
func (b *Bitlist) Overlaps(c *Bitlist) bool {
	lenB, lenC := b.Len(), c.Len()
	if lenB != lenC {
		panic("bitlists are different lengths")
	}

	if lenB == 0 || lenC == 0 {
		return false
	}

	// To ensure all of the bits in c are not overlapped in b, we iterate over every word, invert b
	// and xor the word from b and c, then and it against c. If the result is non-zero, then
	// we can be assured that word in c had bits not overlapped in b.
	for idx, word := range b.data {
		if (^word^c.data[idx])&c.data[idx]&allBitsSet != 0 {
			return true
		}
	}
	return false
}

// Or returns the OR result of the two bitfields (union).
// This method will panic if the bitlists are not the same length.
func (b *Bitlist) Or(c *Bitlist) *Bitlist {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	ret := b.Clone()
	b.NoAllocOr(c, ret)

	return ret
}

// NoAllocOr computes the OR result of the two bitfields (union).
// Result is written into provided variable, so no allocation takes place inside the function.
// This method will panic if the bitlists are not the same length.
func (b *Bitlist) NoAllocOr(c, ret *Bitlist) {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	for idx, word := range b.data {
		ret.data[idx] = word | c.data[idx]
	}
}

// And returns the AND result of the two bitfields (intersection).
// This method will panic if the bitlists are not the same length.
func (b *Bitlist) And(c *Bitlist) *Bitlist {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	ret := b.Clone()
	b.NoAllocAnd(c, ret)

	return ret
}

// NoAllocAnd computes the AND result of the two bitfields (intersection).
// Result is written into provided variable, so no allocation takes place inside the function.
// This method will panic if the bitlists are not the same length.
func (b *Bitlist) NoAllocAnd(c, ret *Bitlist) {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	for idx, word := range b.data {
		ret.data[idx] = word & c.data[idx]
	}
}

// Clone safely copies a given bitlist.
func (b *Bitlist) Clone() *Bitlist {
	c := NewBitlist(b.size)
	if b.data != nil {
		copy(c.data, b.data)
	}
	return c
}

// numWordsRequired calculates how many words are required to hold bitlist of n bits.
func numWordsRequired(n uint64) int {
	return int((n + (wordSize - 1)) >> wordSizeLog2)
}

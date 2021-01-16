package bitfield

import (
	"math/bits"
)

var _ = Bitfield(ByteBitlist{})

// ByteBitlist is a bitfield implementation backed by an array of bytes. The most
// significant bit in the array of bytes indicates the start position of the
// bitfield.
//
// Examples of the underlying byte array as bitlist:
//  byte{0b00001000} is a bitlist with 3 bits which are all zero. bits=[0,0,0]
//  byte{0b00011111} is a bitlist with 4 bits which are all one.  bits=[1,1,1,1]
//  byte{0b00011000, 0b00000001} is a bitlist with 8 bits.        bits=[0,0,0,1,1,0,0,0]
//  byte{0b00011000, 0b00000010} is a bitlist with 9 bits.        bits=[0,0,0,0,1,1,0,0,0]
//
// Note: This is the original implementation of bitfield, which is superseded by a more
// optimized version which is based on array of uint64 (see Bitlist for details).
type ByteBitlist []byte

// NewByteBitlist creates a new bitlist of size N.
func NewByteBitlist(n uint64) ByteBitlist {
	ret := make(ByteBitlist, n/8+1)

	// Set most significant bit for length bit.
	i := uint8(1 << (n % 8))
	ret[n/8] |= i

	return ret
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitlist, then this method returns false.
func (b ByteBitlist) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	upperBounds := b.Len()
	if idx >= upperBounds {
		return false
	}

	i := uint8(1 << (idx % 8))
	return b[idx/8]&i == i
}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitlist, then this method returns
// false.
func (b ByteBitlist) SetBitAt(idx uint64, val bool) {
	// Out of bounds, do nothing.
	upperBounds := b.Len()
	if idx >= upperBounds {
		return
	}

	bit := uint8(1 << (idx % 8))
	if val {
		b[idx/8] |= bit
	} else {
		b[idx/8] &^= bit
	}

}

// Len of the bitlist returns the number of bits available in the underlying
// byte array.
func (b ByteBitlist) Len() uint64 {
	if len(b) == 0 {
		return 0
	}
	// The most significant bit is present in the last byte in the array.
	last := b[len(b)-1]

	// Determine the position of the most significant bit.
	msb := bits.Len8(last)
	if msb == 0 {
		return 0
	}

	// The absolute position of the most significant bit will be the number of
	// bits in the preceding bytes plus the position of the most significant
	// bit. Subtract this value by 1 to determine the length of the bitlist.
	return uint64(8*(len(b)-1) + msb - 1)
}

// Bytes returns the trimmed underlying byte array without the length bit. The
// leading zeros in the bitlist will be trimmed to the smallest byte length
// representation of the bitlist. This may produce an empty byte slice if all
// bits were zero.
func (b ByteBitlist) Bytes() []byte {
	if len(b) == 0 {
		return []byte{}
	}

	ret := make([]byte, len(b))
	copy(ret, b)

	// Clear the most significant bit (the length bit).
	msb := uint8(bits.Len8(ret[len(ret)-1])) - 1
	clearBit := uint8(1 << msb)
	ret[len(ret)-1] &^= clearBit

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
func (b ByteBitlist) Count() uint64 {
	c := 0

	for _, bt := range b {
		c += bits.OnesCount8(bt)
	}

	if c > 0 {
		c-- // Remove length bit from count.
	}

	return uint64(c)
}

// Contains returns true if the bitlist contains all of the bits from the provided argument
// bitlist. This method will panic if bitlists are not the same length.
func (b ByteBitlist) Contains(c ByteBitlist) bool {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	// To ensure all of the bits in c are present in b, we iterate over every byte, combine
	// the byte from b and c, then XOR them against b. If the result of this is non-zero, then we
	// are assured that a byte in c had bits not present in b.
	for i := 0; i < len(b); i++ {
		if b[i]^(b[i]|c[i]) != 0 {
			return false
		}
	}

	return true
}

// Overlaps returns true if the bitlist contains one of the bits from the provided argument
// bitlist. This method will panic if bitlists are not the same length.
func (b ByteBitlist) Overlaps(c ByteBitlist) bool {
	lenB, lenC := b.Len(), c.Len()
	if lenB != lenC {
		panic("bitlists are different lengths")
	}

	if lenB == 0 || lenC == 0 {
		return false
	}

	msb := uint8(bits.Len8(b[len(b)-1])) - 1
	lengthBitMask := uint8(1 << msb)

	// To ensure all of the bits in c are not overlapped in b, we iterate over every byte, invert b
	// and xor the byte from b and c, then and it against c. If the result is non-zero, then
	// we can be assured that byte in c had bits not overlapped in b.
	for i := 0; i < len(b); i++ {
		// If this byte is the last byte in the array, mask the length bit.
		mask := uint8(0xFF)
		if i == len(b)-1 {
			mask &^= lengthBitMask
		}

		if (^b[i]^c[i])&c[i]&mask != 0 {
			return true
		}
	}
	return false
}

// Or returns the OR result of the two bitfields. This method will panic if the bitlists are not the same length.
func (b ByteBitlist) Or(c ByteBitlist) ByteBitlist {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	ret := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		ret[i] = b[i] | c[i]
	}

	return ret
}

// And returns the AND result of the two bitfields. This method will panic if the bitlists are not the same length.
func (b ByteBitlist) And(c ByteBitlist) ByteBitlist {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	ret := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		ret[i] = b[i] & c[i]
	}

	return ret
}

// Xor returns the XOR result of the two bitfields. This method will panic if the bitlists are not the same length.
func (b ByteBitlist) Xor(c ByteBitlist) ByteBitlist {
	if b.Len() != c.Len() {
		panic("bitlists are different lengths")
	}

	// Process all bytes but the last.
	ret := make([]byte, len(b))
	for i := 0; i < len(b)-1; i++ {
		ret[i] = b[i] ^ c[i]
	}

	// For the last byte, process only bits smaller than the length bit.
	ret[len(b)-1] = b[len(b)-1]
	msb := uint8(bits.Len8(b[len(b)-1])) - 1
	if msb > 0 {
		mask := uint8(0xff >> (8 - msb))
		ret[len(b)-1] = (b[len(b)-1] ^ c[len(b)-1]) & mask
		ret[len(b)-1] |= uint8(1 << msb)
	}

	return ret
}

// Not returns the NOT result of the bitfield.
func (b ByteBitlist) Not() ByteBitlist {
	if b.Len() == 0 {
		return b
	}

	// Process all bytes but the last.
	ret := make([]byte, len(b))
	for i := 0; i < len(b)-1; i++ {
		ret[i] = ^b[i]
	}

	// For the last byte, process only bits smaller than the length bit.
	ret[len(b)-1] = b[len(b)-1]
	msb := uint8(bits.Len8(b[len(b)-1])) - 1
	if msb > 0 {
		mask := uint8(0xff >> (8 - msb))
		ret[len(b)-1] = (^b[len(b)-1]) & mask
		ret[len(b)-1] |= uint8(1 << msb)
	}

	return ret
}

func (b ByteBitlist) BitIndices() []int {
	indices := make([]int, 0, b.Count())
	for i, bt := range b {
		if i == len(b)-1 {
			// Clear the most significant bit (the length bit).
			msb := uint8(bits.Len8(bt)) - 1
			bt &^= uint8(1 << msb)
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
package bitfield

var _ = Bitfield(Bitvector4{})

// Bitvector4 is a bitfield with a known size of 4. There is no length bit
// present in the underlying byte array.
type Bitvector4 []byte

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitlist, then this method returns false.
func (b Bitvector4) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	if idx >= b.Len() {
		return false
	}

	i := uint8(1 << idx)
	return b[0]&i == i

}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitlist, then this method returns
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

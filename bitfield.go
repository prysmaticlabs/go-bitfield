package bitfield

// Bitfield is the abstraction implemented by Bitlist or BitvectorN.
type Bitfield interface {
	BitAt(idx uint64) bool
	SetBitAt(idx uint64, val bool)
	Len() uint64
	Count() uint64
}

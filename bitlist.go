package bitfield

const (
	// wordSize configures how many bits are there in a single element of bitlist array.
	wordSize = uint64(64)
	// wordSizeLog2 allows optimized division by wordSize using right shift (numBits >> wordSizeLog2).
	wordSizeLog2 = uint64(6)
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

// numWordsRequired calculates how many words are required to hold bitlist of n bits.
func numWordsRequired(n uint64) int {
	return int((n + (wordSize - 1)) >> wordSizeLog2)
}

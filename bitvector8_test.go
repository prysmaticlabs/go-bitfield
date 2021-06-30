package bitfield

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBitvector8_Len(t *testing.T) {
	bvs := []Bitvector8{
		{},
		{0x01},
		{0x01, 0x02},
		{0x0F, 0x0F},
	}
	for _, bv := range bvs {
		if bv.Len() != 8 {
			t.Errorf("(%x).Len() = %d, wanted %d", bv, bv.Len(), 4)
		}
	}
}

func TestBitvector8_BitAt(t *testing.T) {
	tests := []struct {
		bitlist Bitvector8
		idx     uint64
		want    bool
	}{
		{
			bitlist: Bitvector8{0x01}, // 0b00000001
			idx:     55,               // Out of bounds
			want:    false,
		},
		{
			bitlist: Bitvector8{0xFF}, // 0b11111111
			idx:     8,                // Out of bounds
			want:    false,
		},
		{
			bitlist: Bitvector8{0x01}, // 0b00000001
			idx:     0,                //          ^
			want:    true,
		},
		{
			bitlist: Bitvector8{0x0E}, // 0b00001110
			idx:     0,                //          ^
			want:    false,
		},
		{
			bitlist: Bitvector8{0x0E}, // 0b00001110
			idx:     1,                //         ^
			want:    true,
		},
		{
			bitlist: Bitvector8{0x0E}, // 0b00001110
			idx:     2,                //        ^
			want:    true,
		},
		{
			bitlist: Bitvector8{0x0E}, // 0b00001110
			idx:     3,                //       ^
			want:    true,
		},
		{
			bitlist: Bitvector8{0x1E}, // 0b00011110
			idx:     4,                //      ^
			want:    true,
		},
		{
			bitlist: Bitvector8{0x9E}, // 0b10011110
			idx:     7,                //   ^
			want:    true,
		},
	}

	for _, tt := range tests {
		if tt.bitlist.BitAt(tt.idx) != tt.want {
			t.Errorf(
				"(%x).BitAt(%d) = %t, wanted %t",
				tt.bitlist,
				tt.idx,
				tt.bitlist.BitAt(tt.idx),
				tt.want,
			)
		}
	}
}

func TestBitvector8_SetBitAt(t *testing.T) {
	tests := []struct {
		bitvector Bitvector8
		idx       uint64
		val       bool
		want      Bitvector8
	}{
		{
			bitvector: Bitvector8{0x01}, // 0b00000001
			idx:       0,                //          ^
			val:       true,
			want:      Bitvector8{0x01}, // 0b00000001
		},
		{
			bitvector: Bitvector8{0x02}, // 0b00000010
			idx:       0,                //          ^
			val:       true,
			want:      Bitvector8{0x03}, // 0b00000011
		},
		{
			bitvector: Bitvector8{0x00}, // 0b00000000
			idx:       1,                //         ^
			val:       true,
			want:      Bitvector8{0x02}, // 0b00000010
		},
		{
			bitvector: Bitvector8{0x00}, // 0b00000000
			idx:       3,                //       ^
			val:       true,
			want:      Bitvector8{0x08}, // 0b00001000
		},
		{
			bitvector: Bitvector8{0x00}, // 0b00000000
			idx:       4,                //      ^
			val:       true,
			want:      Bitvector8{0x10}, // 0b00010000
		},
		{
			bitvector: Bitvector8{0x00}, // 0b00000000
			idx:       5,                //     ^
			val:       true,
			want:      Bitvector8{0x20}, // 0b00100000
		},
		{
			bitvector: Bitvector8{0x0F}, // 0b00001111
			idx:       0,                //          ^
			val:       true,
			want:      Bitvector8{0x0F}, // 0b00001111
		},
		{
			bitvector: Bitvector8{0x0F}, // 0b00001111
			idx:       0,                //          ^
			val:       false,
			want:      Bitvector8{0x0E}, // 0b00001110
		},
		{
			bitvector: Bitvector8{0x00}, // Out of bound
			idx:       8,
			val:       true,
			want:      Bitvector8{0x00},
		},
	}

	for _, tt := range tests {
		original := make(Bitvector8, len(tt.bitvector))
		copy(original, tt.bitvector)

		tt.bitvector.SetBitAt(tt.idx, tt.val)
		if !bytes.Equal(tt.bitvector, tt.want) {
			t.Errorf(
				"(%x).SetBitAt(%d, %t) = %x, wanted %x",
				original,
				tt.idx,
				tt.val,
				tt.bitvector,
				tt.want,
			)
		}
	}
}

func TestBitvector8_Count(t *testing.T) {
	tests := []struct {
		bitvector Bitvector8
		want      uint64
	}{
		{
			bitvector: Bitvector8{},
			want:      0,
		},
		{
			bitvector: Bitvector8{0x01}, // 0b00000001
			want:      1,
		},
		{
			bitvector: Bitvector8{0x03}, // 0b00000011
			want:      2,
		},
		{
			bitvector: Bitvector8{0x07}, // 0b00000111
			want:      3,
		},
		{
			bitvector: Bitvector8{0x0F}, // 0b00001111
			want:      4,
		},
		{
			bitvector: Bitvector8{0xFF}, // 0b11111111
			want:      8,
		},
		{
			bitvector: Bitvector8{0xF0}, // 0b11110000
			want:      4,
		},
	}

	for _, tt := range tests {
		if tt.bitvector.Count() != tt.want {
			t.Errorf(
				"(%x).Count() = %d, wanted %d",
				tt.bitvector,
				tt.bitvector.Count(),
				tt.want,
			)
		}
	}
}

func TestBitvector8_Bytes(t *testing.T) {
	tests := []struct {
		bitvector Bitvector8
		want      []byte
	}{
		{
			bitvector: Bitvector8{},
			want:      []byte{},
		},
		{
			bitvector: Bitvector8{0x00}, // 0b00000000
			want:      []byte{0x00},     // 0b00000000
		},
		{
			bitvector: Bitvector8{0x01}, // 0b00000001
			want:      []byte{0x01},     // 0b00000001
		},
		{
			bitvector: Bitvector8{0x03}, // 0b00000011
			want:      []byte{0x03},     // 0b00000011
		},
		{
			bitvector: Bitvector8{0x07}, // 0b00000111
			want:      []byte{0x07},     // 0b00000111
		},
		{
			bitvector: Bitvector8{0x0F}, // 0b00001111
			want:      []byte{0x0F},     // 0b00001111
		},
		{
			bitvector: Bitvector8{0xFF}, // 0b11111111
			want:      []byte{0xFF},     // 0b11111111
		},
		{
			bitvector: Bitvector8{0xF0}, // 0b11110000
			want:      []byte{0xF0},     // 0b11110000
		},
		{
			bitvector: Bitvector8{0xF0, 0xFF}, // 0b11110000
			want:      []byte{0xF0},     // 0b11110000
		},
	}

	for _, tt := range tests {
		if !bytes.Equal(tt.bitvector.Bytes(), tt.want) {
			t.Errorf(
				"(%x).Bytes() = %x, wanted %x",
				tt.bitvector,
				tt.bitvector.Bytes(),
				tt.want,
			)
		}
	}
}

func TestBitvector8_BitIndices(t *testing.T) {
	tests := []struct {
		a    Bitvector8
		want []int
	}{
		{
			a:    Bitvector8{0b1001},
			want: []int{0, 3},
		},
		{
			a:    Bitvector8{0b1000},
			want: []int{3},
		},
		{
			a:    Bitvector8{0b10},
			want: []int{1},
		},
		{
			a:    Bitvector8{0b11111111},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			a:    Bitvector8{0b11111111, 0xFF},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
	}

	for _, tt := range tests {
		if !reflect.DeepEqual(tt.a.BitIndices(), tt.want) {
			t.Errorf(
				"(%0.8b).BitIndices() = %x, wanted %x",
				tt.a,
				tt.a.BitIndices(),
				tt.want,
			)
		}
	}
}

func TestBitvector8_Contains(t *testing.T) {
	tests := []struct {
		a    Bitvector8
		b    Bitvector8
		want bool
	}{
		{
			a:    Bitvector8{0x02}, // 0b00000010
			b:    Bitvector8{0x03}, // 0b00000011
			want: false,
		},
		{
			a:    Bitvector8{0x03}, // 0b00000011
			b:    Bitvector8{0x03}, // 0b00000011
			want: true,
		},
		{
			a:    Bitvector8{0x13}, // 0b00010011
			b:    Bitvector8{0x15}, // 0b00010101
			want: false,
		},
		{
			a:    Bitvector8{0x1F}, // 0b00011111
			b:    Bitvector8{0x13}, // 0b00010011
			want: true,
		},
	}

	for _, tt := range tests {
		if tt.a.Contains(tt.b) != tt.want {
			t.Errorf(
				"(%x).Contains(%x) = %t, wanted %t",
				tt.a,
				tt.b,
				tt.a.Contains(tt.b),
				tt.want,
			)
		}
	}
}

func TestBitvector8_Overlaps(t *testing.T) {
	tests := []struct {
		a    Bitvector8
		b    Bitvector8
		want bool
	}{
		{
			a:    Bitvector8{0x06}, // 0b00000110
			b:    Bitvector8{0x01}, // 0b00000101
			want: false,
		},
		{
			a:    Bitvector8{0x06}, // 0b00000110
			b:    Bitvector8{0x05}, // 0b00000101
			want: true,
		},
		{
			a:    Bitvector8{0x1A}, // 0b00011010
			b:    Bitvector8{0x25}, // 0b00100101
			want: false,
		},
		{
			a:    Bitvector8{0x1F}, // 0b00011111
			b:    Bitvector8{0x11}, // 0b00010001
			want: true,
		},
	}

	for _, tt := range tests {
		result := tt.a.Overlaps(tt.b)
		if result != tt.want {
			t.Errorf(
				"(%x).Overlaps(%x) = %t, wanted %t",
				tt.a,
				tt.b,
				result,
				tt.want,
			)
		}
	}
}

func TestBitVector8_Or(t *testing.T) {
	tests := []struct {
		a    Bitvector8
		b    Bitvector8
		want Bitvector8
	}{
		{
			a:    Bitvector8{0x02}, // 0b00000010
			b:    Bitvector8{0x03}, // 0b00000011
			want: Bitvector8{0x03}, // 0b00000011
		},
		{
			a:    Bitvector8{0x03}, // 0b00000011
			b:    Bitvector8{0x03}, // 0b00000011
			want: Bitvector8{0x03}, // 0b00000011
		},
		{
			a:    Bitvector8{0x13}, // 0b00010011
			b:    Bitvector8{0x15}, // 0b00010101
			want: Bitvector8{0x17}, // 0b00010111
		},
		{
			a:    Bitvector8{0x1F}, // 0b00011111
			b:    Bitvector8{0x13}, // 0b00010011
			want: Bitvector8{0x1F}, // 0b00011111
		},
	}

	for _, tt := range tests {
		if !bytes.Equal(tt.a.Or(tt.b), tt.want) {
			t.Errorf(
				"(%x).Or(%x) = %x, wanted %x",
				tt.a,
				tt.b,
				tt.a.Or(tt.b),
				tt.want,
			)
		}
	}
}

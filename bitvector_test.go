package bitfield

import (
	"bytes"
	"testing"
)

func TestBitvector4_Len(t *testing.T) {
	bvs := []Bitvector4{
		{},
		{0x01},
		{0x02},
		{0x03},
		{0x04},
		{0x05},
		{0x06},
		{0x07},
		{0x0F},
	}

	for _, bv := range bvs {
		if bv.Len() != 4 {
			t.Errorf("(%x).Len() = %d, wanted %d", bv, bv.Len(), 4)
		}
	}
}

func TestBitvector4_BitAt(t *testing.T) {
	tests := []struct {
		bitlist Bitvector4
		idx     uint64
		want    bool
	}{
		{
			bitlist: Bitvector4{0x01}, // 0b00000001
			idx:     55,               // Out of bounds
			want:    false,
		},
		{
			bitlist: Bitvector4{0x01}, // 0b00000001
			idx:     0,                //          ^
			want:    true,
		},
		{
			bitlist: Bitvector4{0x0E}, // 0b00001110
			idx:     0,                //          ^
			want:    false,
		},
		{
			bitlist: Bitvector4{0x0E}, // 0b00001110
			idx:     1,                //         ^
			want:    true,
		},
		{
			bitlist: Bitvector4{0x0E}, // 0b00001110
			idx:     2,                //        ^
			want:    true,
		},
		{
			bitlist: Bitvector4{0x0E}, // 0b00001110
			idx:     3,                //       ^
			want:    true,
		},
		{
			bitlist: Bitvector4{0x1E}, // 0b00011110
			idx:     4,                //      ^
			want:    false,
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

func TestBitvector4_SetBitAt(t *testing.T) {
	tests := []struct {
		bitvector Bitvector4
		idx       uint64
		val       bool
		want      Bitvector4
	}{
		{
			bitvector: Bitvector4{0x01}, // 0b00000001
			idx:       0,                //          ^
			val:       true,
			want:      Bitvector4{0x01}, // 0b00000001
		},
		{
			bitvector: Bitvector4{0x02}, // 0b00000010
			idx:       0,                //          ^
			val:       true,
			want:      Bitvector4{0x03}, // 0b00000011
		},
		{
			bitvector: Bitvector4{0x00}, // 0b00000000
			idx:       1,                //         ^
			val:       true,
			want:      Bitvector4{0x02}, // 0b00000010
		},
		{
			bitvector: Bitvector4{0x00}, // 0b00000000
			idx:       3,                //       ^
			val:       true,
			want:      Bitvector4{0x08}, // 0b00001000
		},
		{
			bitvector: Bitvector4{0x00}, // 0b00000000
			idx:       4,                //      ^
			val:       true,
			want:      Bitvector4{0x00}, // 0b00001000
		},
		{
			bitvector: Bitvector4{0x00}, // 0b00000000
			idx:       5,                // Out of bounds
			val:       true,
			want:      Bitvector4{0x00}, // 0b00000000
		},
		{
			bitvector: Bitvector4{0x0F}, // 0b00001111
			idx:       0,                //          ^
			val:       true,
			want:      Bitvector4{0x0F}, // 0b00001111
		},
		{
			bitvector: Bitvector4{0x0F}, // 0b00001111
			idx:       0,                //          ^
			val:       false,
			want:      Bitvector4{0x0E}, // 0b00001110
		},
	}

	for _, tt := range tests {
		original := make(Bitvector4, len(tt.bitvector))
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

func TestBitvector4_Count(t *testing.T) {
	tests := []struct {
		bitvector Bitvector4
		want      uint64
	}{
		{
			bitvector: Bitvector4{},
			want:      0,
		},
		{
			bitvector: Bitvector4{0x01}, // 0b00000001
			want:      1,
		},
		{
			bitvector: Bitvector4{0x03}, // 0b00000011
			want:      2,
		},
		{
			bitvector: Bitvector4{0x07}, // 0b00000111
			want:      3,
		},
		{
			bitvector: Bitvector4{0x0F}, // 0b00001111
			want:      4,
		},
		{
			bitvector: Bitvector4{0xFF}, // 0b11111111
			want:      4,
		},
		{
			bitvector: Bitvector4{0xF0}, // 0b11110000
			want:      0,
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

func TestBitvector4_Bytes(t *testing.T) {
	tests := []struct {
		bitvector Bitvector4
		want      []byte
	}{
		{
			bitvector: Bitvector4{},
			want:      []byte{},
		},
		{
			bitvector: Bitvector4{0x00}, // 0b00000000
			want:      []byte{0x00},     // 0b00000000
		},
		{
			bitvector: Bitvector4{0x01}, // 0b00000001
			want:      []byte{0x01},     // 0b00000001
		},
		{
			bitvector: Bitvector4{0x03}, // 0b00000011
			want:      []byte{0x03},     // 0b00000011
		},
		{
			bitvector: Bitvector4{0x07}, // 0b00000111
			want:      []byte{0x07},     // 0b00000111
		},
		{
			bitvector: Bitvector4{0x0F}, // 0b00001111
			want:      []byte{0x0F},     // 0b00001111
		},
		{
			bitvector: Bitvector4{0xFF}, // 0b11111111
			want:      []byte{0x0F},     // 0b00000001
		},
		{
			bitvector: Bitvector4{0xF0}, // 0b11110000
			want:      []byte{0x00},     // 0b00000000
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

func TestBitvector4_Shift(t *testing.T) {
	tests := []struct {
		bitvector Bitvector4
		shift     int
		want      Bitvector4
	}{
		{
			bitvector: Bitvector4{},
			shift:     1,
			want:      Bitvector4{},
		},
		{
			bitvector: Bitvector4{0x01},
			shift:     1,
			want:      Bitvector4{0x02},
		},
		{
			bitvector: Bitvector4{0x02},
			shift:     1,
			want:      Bitvector4{0x04},
		},
		{
			bitvector: Bitvector4{0x04},
			shift:     1,
			want:      Bitvector4{0x08},
		},
		{
			bitvector: Bitvector4{0x08},
			shift:     1,
			want:      Bitvector4{0x00},
		},
		{
			bitvector: Bitvector4{0x03},
			shift:     1,
			want:      Bitvector4{0x06},
		},
		{
			bitvector: Bitvector4{0x02},
			shift:     -1,
			want:      Bitvector4{0x01},
		},
		{
			bitvector: Bitvector4{0x03},
			shift:     -1,
			want:      Bitvector4{0x01},
		},
		{
			bitvector: Bitvector4{0x03},
			shift:     3,
			want:      Bitvector4{0x08},
		},
		{
			bitvector: Bitvector4{0x0F},
			shift:     -3,
			want:      Bitvector4{0x01},
		},
		{
			bitvector: Bitvector4{0x0F},
			shift:     8,
			want:      Bitvector4{0x00},
		},
		{
			bitvector: Bitvector4{0x0F},
			shift:     256,
			want:      Bitvector4{0x00},
		},
		{
			bitvector: Bitvector4{0x0F},
			shift:     -256,
			want:      Bitvector4{0x00},
		},
	}

	for _, tt := range tests {
		original := make(Bitvector4, len(tt.bitvector))
		copy(original, tt.bitvector)

		tt.bitvector.Shift(tt.shift)
		if !bytes.Equal(tt.bitvector, tt.want) {
			t.Errorf(
				"(%x).Shift(%d) = %x, wanted %x",
				original,
				tt.shift,
				tt.bitvector,
				tt.want,
			)
		}
	}
}

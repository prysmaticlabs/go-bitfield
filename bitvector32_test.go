package bitfield

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBitvector32_Len(t *testing.T) {
	bvs := []Bitvector32{
		{0x00, 0x00, 0x00, 0x00},
		{0x01, 0x00, 0x00, 0x00},
		{0x02, 0x00, 0x00, 0x00},
		{0x03, 0x00, 0x00, 0x00},
	}

	for _, bv := range bvs {
		if bv.Len() != 32 {
			t.Errorf("(%x).Len() = %d, wanted %d", bv, bv.Len(), 32)
		}
	}
}

func TestBitvector32_BitAt(t *testing.T) {
	tests := []struct {
		bitlist Bitvector32
		idx     uint64
		want    bool
	}{
		{
			bitlist: Bitvector32{0x01, 0x23, 0xE2, 0xFE, 0xDD, 0xAC, 0xAD},
			idx:     70, // Out of bounds
			want:    false,
		},
		{
			bitlist: Bitvector32{0x01},
			idx:     0,
			want:    true,
		},
		{
			bitlist: Bitvector32{0x0E, 0xAA, 0x2F},
			idx:     0,
			want:    false,
		},
		{
			bitlist: Bitvector32{0x01, 0x23, 0xE2, 0xFE}, // 00000001 00100011 11100010 11111110
			idx:     35,
			want:    false,
		},
		{
			bitlist: Bitvector32{0x01, 0x23, 0xE2, 0xFE}, // 00000001 00100011 11100010 11111110
			idx:     24,
			want:    false,
		},
		{
			bitlist: Bitvector32{0x0E}, // 0b00001110
			idx:     3,                 //       ^
			want:    true,
		},
		{
			bitlist: Bitvector32{0x1E}, // 0b00011110
			idx:     4,                 //      ^
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

func TestBitvector32_SetBitAt(t *testing.T) {
	tests := []struct {
		bitvector Bitvector32
		idx       uint64
		val       bool
		want      Bitvector32
	}{
		{
			bitvector: Bitvector32{0x01, 0x00, 0x00, 0x00}, // 0b00000001
			idx:       0,                                                           //          ^
			val:       true,
			want:      Bitvector32{0x01, 0x00, 0x00, 0x00}, // 0b00000001
		},
		{
			bitvector: Bitvector32{0x02, 0x00, 0x00, 0x00}, // 0b00000010
			idx:       0,                                                           //          ^
			val:       true,
			want:      Bitvector32{0x03, 0x00, 0x00, 0x00}, // 0b00000011
		},
		{
			bitvector: Bitvector32{0x00, 0x00, 0x00, 0x00}, // 0b00000000
			idx:       1,
			val:       true,
			want:      Bitvector32{0x02, 0x00, 0x00, 0x00}, // 0b00000010
		},
		{
			bitvector: Bitvector32{0x00, 0x00, 0x00, 0x00}, // 0b00000000
			idx:       28,                                                          //       ^
			val:       true,
			want:      Bitvector32{0x00, 0x00, 0x00, 0x10}, // 0b00001000
		},
		{
			bitvector: Bitvector32{0x00, 0x00, 0x00, 0x00}, // 0b00000000
			idx:       30,                                                          //      ^
			val:       true,
			want:      Bitvector32{0x00, 0x00, 0x00, 0x40}, // 0b00001000
		},
		{
			bitvector: Bitvector32{0x00, 0x20, 0x00, 0x00}, // 0b00000000
			idx:       25,
			val:       false,
			want:      Bitvector32{0x00, 0x20, 0x00, 0x00}, // 0b00000000
		},
		{
			bitvector: Bitvector32{0x0F, 0x00, 0x00, 0x00}, // 0b00001111
			idx:       0,                                                           //          ^
			val:       true,
			want:      Bitvector32{0x0F, 0x00, 0x00, 0x00}, // 0b00001111
		},
		{
			bitvector: Bitvector32{0x0F, 0x00, 0x00, 0x00}, // 0b00001111
			idx:       0,                                                           //          ^
			val:       false,
			want:      Bitvector32{0x0E, 0x00, 0x00, 0x00}, // 0b00001110
		},
	}

	for _, tt := range tests {
		original := [8]byte{}
		copy(original[:], tt.bitvector[:])

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

func TestBitvector32_Count(t *testing.T) {
	tests := []struct {
		bitvector Bitvector32
		want      uint64
	}{
		{
			bitvector: Bitvector32{},
			want:      0,
		},
		{
			bitvector: Bitvector32{0x01}, // 0b00000001
			want:      1,
		},
		{
			bitvector: Bitvector32{0x03, 0x00, 0x30, 0x00}, // 0b00000011
			want:      4,
		},
		{
			bitvector: Bitvector32{0x07, 0x40, 0x40, 0x00}, // 0b00000111
			want:      5,
		},
		{
			bitvector: Bitvector32{0x0F, 0x20, 0x00, 0x00}, // 0b00001111
			want:      5,
		},
		{
			bitvector: Bitvector32{0xFF, 0xEE, 0x00, 0x00}, // 0b11111111
			want:      14,
		},
		{
			bitvector: Bitvector32{0x00}, // 0b11110000
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

func TestBitvector32_Bytes(t *testing.T) {
	tests := []struct {
		bitvector Bitvector32
		want      []byte
	}{
		{
			bitvector: Bitvector32{},
			want:      []byte{},
		},
		{
			bitvector: Bitvector32{0x12, 0x34, 0xAB, 0x00},
			want:      []byte{0x12, 0x34, 0xAB, 0x00},
		},
		{
			bitvector: Bitvector32{0x01},
			want:      []byte{0x01},
		},
		{
			bitvector: Bitvector32{0x03},
			want:      []byte{0x03},
		},
		{
			bitvector: Bitvector32{0x07},
			want:      []byte{0x07},
		},
		{
			bitvector: Bitvector32{0x0F},
			want:      []byte{0x0F},
		},
		{
			bitvector: Bitvector32{0xFF},
			want:      []byte{0xFF},
		},
		{
			bitvector: Bitvector32{0xF0},
			want:      []byte{0xF0},
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

func TestBitVector32_BitIndices(t *testing.T) {
	tests := []struct {
		a    Bitvector32
		want []int
	}{
		{
			a:    Bitvector32{0b10010},
			want: []int{1, 4},
		},
		{
			a:    Bitvector32{0b10000},
			want: []int{4},
		},
		{
			a:    Bitvector32{0b10, 0b1},
			want: []int{1, 8},
		},
		{
			a:    Bitvector32{0b11111111, 0b11},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
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


package bitfield

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBitvector1024_Len(t *testing.T) {
	bv := NewBitvector1024()
	if bv.Len() != 1024 {
		t.Errorf("(%x).Len() = %d, wanted %d", bv, bv.Len(), 1024)
	}
}

func TestBitvector1024_BitAt(t *testing.T) {
	tests := []struct {
		bitlist Bitvector1024
		idx     uint64
		want    bool
	}{
		{
			bitlist: Bitvector1024{0x01},
			idx:     0,
			want:    true,
		},
		{
			bitlist: Bitvector1024{0x0E, 0xAA, 0x2F},
			idx:     1,
			want:    true,
		},
		{
			bitlist: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x01, 0x23, 0xE2, 0xFE, 0xDD, 0xAC, 0xAD},
			idx:     87,
			want:    true,
		},
		{
			bitlist: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x01, 0x23, 0xE2, 0xFE, 0xDD, 0xAC, 0xAD},
			idx:     76,
			want:    false,
		},
		{
			bitlist: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0E},
			idx:     131,
			want:    true,
		},
		{
			bitlist: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40},
			idx:  1022,
			want: true,
		},
		{
			bitlist: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40},
			idx:  1023,
			want: false,
		},
		{
			bitlist: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80},
			idx:  1023,
			want: true,
		},
		{
			bitlist: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF},
			idx:  1024, // Out of bound.
			want: false,
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

func TestBitvector1024_SetBitAt(t *testing.T) {
	tests := []struct {
		bitvector Bitvector1024
		idx       uint64
		val       bool
		want      Bitvector1024
	}{
		{
			bitvector: Bitvector1024{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00000001
			idx:       0,                                                             //          ^
			val:       true,
			want:      Bitvector1024{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00000001
		},
		{
			bitvector: Bitvector1024{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00000010
			idx:       0,                                                             //          ^
			val:       true,
			want:      Bitvector1024{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00000011
		},
		{
			bitvector: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00000000
			idx:       1,
			val:       true,
			want:      Bitvector1024{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00000010
		},
		{
			bitvector: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00000000
			idx:       48,                                                            //       ^
			val:       true,
			want:      Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00}, // 0b00001000
		},
		{
			bitvector: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00000000
			idx:       30,                                                            //      ^
			val:       true,
			want:      Bitvector1024{0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00}, // 0b00001000
		},
		{
			bitvector: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00}, // 0b00000000
			idx:       45,
			val:       false,
			want:      Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00000000
		},
		{
			bitvector: Bitvector1024{0x0F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00001111
			idx:       0,                                                             //          ^
			val:       true,
			want:      Bitvector1024{0x0F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00001111
		},
		{
			bitvector: Bitvector1024{0x0F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00001111
			idx:       0,                                                             //          ^
			val:       false,
			want:      Bitvector1024{0x0E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b00001110
		},
		{
			bitvector: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			idx: 1019,
			val: true,
			want: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
		},
		{
			bitvector: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			idx: 1023,
			val: true,
			want: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80},
		},
		{
			bitvector: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80},
			idx: 1023,
			val: false,
			want: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
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

func TestBitvector1024_Count(t *testing.T) {
	tests := []struct {
		bitvector Bitvector1024
		want      uint64
	}{
		{
			bitvector: Bitvector1024{},
			want:      0,
		},
		{
			bitvector: Bitvector1024{0x01}, // 0b00000001
			want:      1,
		},
		{
			bitvector: Bitvector1024{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00}, // 0b00000011
			want:      4,
		},
		{
			bitvector: Bitvector1024{0x07, 0x00, 0x00, 0x00, 0x00, 0x40, 0x40, 0x00}, // 0b00000111
			want:      5,
		},
		{
			bitvector: Bitvector1024{0x0F, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00}, // 0b00001111
			want:      5,
		},
		{
			bitvector: Bitvector1024{0xFF, 0xEE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 0b11111111
			want:      14,
		},
		{
			bitvector: Bitvector1024{0x00},
			want:      0,
		},
		{
			bitvector: Bitvector1024{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			want: 1024,
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

func TestBitvector1024_Bytes(t *testing.T) {
	tests := []struct {
		bitvector Bitvector1024
		want      []byte
	}{
		{
			bitvector: Bitvector1024{},
			want:      []byte{},
		},
		{
			bitvector: Bitvector1024{0x12, 0x34, 0xAB, 0x00},
			want:      []byte{0x12, 0x34, 0xAB, 0x00},
		},
		{
			bitvector: Bitvector1024{0x01},
			want:      []byte{0x01},
		},
		{
			bitvector: Bitvector1024{0x03},
			want:      []byte{0x03},
		},
		{
			bitvector: Bitvector1024{0x07},
			want:      []byte{0x07},
		},
		{
			bitvector: Bitvector1024{0x0F},
			want:      []byte{0x0F},
		},
		{
			bitvector: Bitvector1024{0xFF},
			want:      []byte{0xFF},
		},
		{
			bitvector: Bitvector1024{0xF0},
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

func TestBitvector1024_Shift(t *testing.T) {
	tests := []struct {
		bitvector Bitvector1024
		shift     int
		want      Bitvector1024
	}{
		{
			bitvector: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			shift:     1,
			want:      Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			bitvector: Bitvector1024{0x01, 0x23, 0xE2, 0xFE, 0xDD, 0xAC, 0xAD, 0xAD},
			shift:     1,
			want:      Bitvector1024{0x02, 0x47, 0xC5, 0xFD, 0xBB, 0x59, 0x5B, 0x5A},
		},
		{
			bitvector: Bitvector1024{0x23, 0x01, 0xAD, 0xE2, 0xDD, 0xFE, 0xAC, 0xAD},
			shift:     1,
			want:      Bitvector1024{0x46, 0x03, 0x5b, 0xc5, 0xBB, 0xFD, 0x59, 0x5A},
		},
		{
			bitvector: Bitvector1024{0x01, 0x23, 0xE2, 0xFE, 0xDD, 0xAC, 0xAD, 0xAD},
			shift:     -1,
			want:      Bitvector1024{0x00, 0x91, 0xf1, 0x7f, 0x6e, 0xd6, 0x56, 0xd6},
		},
		{
			bitvector: Bitvector1024{0xd6, 0x23, 0x6e, 0x91, 0xDD, 0xAC, 0x7f, 0xE2},
			shift:     -1,
			want:      Bitvector1024{0x6b, 0x11, 0xb7, 0x48, 0xee, 0xd6, 0x3f, 0xf1},
		},
		{
			bitvector: Bitvector1024{0x01, 0x23, 0xE2, 0xFE, 0xDD, 0xAC, 0xAD, 0xAD},
			shift:     3,
			want:      Bitvector1024{0x09, 0x1f, 0x17, 0xf6, 0xed, 0x65, 0x6d, 0x68},
		},
		{
			bitvector: Bitvector1024{0x17, 0xDD, 0x09, 0x17, 0x1f, 0x17, 0xf6, 0xed},
			shift:     -3,
			want:      Bitvector1024{0x02, 0xfb, 0xa1, 0x22, 0xe3, 0xe2, 0xfe, 0xdd},
		},
		{
			bitvector: Bitvector1024{0x01, 0x23, 0xE2, 0xFE, 0xDD, 0xAC, 0xAD, 0xAD},
			shift:     8,
			want:      Bitvector1024{0x23, 0xe2, 0xfe, 0xdd, 0xac, 0xad, 0xad, 0x00},
		},
	}

	for _, tt := range tests {
		original := make(Bitvector1024, len(tt.bitvector))
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

func TestBitvector1024_BitIndices(t *testing.T) {
	tests := []struct {
		a    Bitvector1024
		want []int
	}{
		{
			a:    Bitvector1024{0b10010},
			want: []int{1, 4},
		},
		{
			a:    Bitvector1024{0b10000},
			want: []int{4},
		},
		{
			a:    Bitvector1024{0b10, 0b1},
			want: []int{1, 8},
		},
		{
			a:    Bitvector1024{0b11111111, 0b11},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			a: Bitvector1024{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF0, 0xF0},
			want: []int{1012, 1013, 1014, 1015, 1020, 1021, 1022, 1023},
		},
	}

	for _, tt := range tests {
		if !reflect.DeepEqual(tt.a.BitIndices(), tt.want) {
			t.Errorf(
				"(%0.8b).BitIndices() = %d, wanted %d",
				tt.a,
				tt.a.BitIndices(),
				tt.want,
			)
		}
	}
}

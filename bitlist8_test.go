package bitfield

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestNewByteBitlist(t *testing.T) {
	tests := []struct {
		size uint64
		want ByteBitlist
	}{
		{
			size: 0,
			want: ByteBitlist{0x01},
		},
		{
			size: 1,
			want: ByteBitlist{0x02},
		},
		{
			size: 2,
			want: ByteBitlist{0x04},
		},
		{
			size: 3,
			want: ByteBitlist{0x08},
		},
		{
			size: 8,
			want: ByteBitlist{0x00, 0x01},
		},
		{
			size: 9,
			want: ByteBitlist{0x00, 0x02},
		},
	}

	for _, tt := range tests {
		got := NewByteBitlist(tt.size)
		if !bytes.Equal(got, tt.want) {
			t.Errorf(
				"NewByteBitlist(%d) = %x, wanted %x",
				tt.size,
				got,
				tt.want,
			)
		}
	}
}

func TestByteBitlist_Len(t *testing.T) {
	tests := []struct {
		bitlist ByteBitlist
		want    uint64
	}{
		{
			bitlist: ByteBitlist{},
			want:    0,
		},
		{
			bitlist: ByteBitlist{0x00}, // 0b00000000, invalid list
			want:    0,
		},
		{
			bitlist: ByteBitlist{0x01}, // 0b00000001
			want:    0,
		},
		{
			bitlist: ByteBitlist{0x02}, // 0b00000010
			want:    1,
		},
		{
			bitlist: ByteBitlist{0x08}, // 0b00001000
			want:    3,
		},
		{
			bitlist: ByteBitlist{0x0E}, // 0b00001110
			want:    3,
		},
		{
			bitlist: ByteBitlist{0x0F}, // 0b00001111
			want:    3,
		},
		{
			bitlist: ByteBitlist{0x10}, // 0b00010000
			want:    4,
		},
		{
			bitlist: ByteBitlist{0x00, 0x01}, // 0b00000000, 0b00000001
			want:    8,
		},
		{
			bitlist: ByteBitlist{0x00, 0x02}, // 0b00000000, 0b00000010
			want:    9,
		},
		{
			bitlist: ByteBitlist{0x00, 0x02, 0x08}, // 0b00000000, 0b00000010, 0b00001000
			want:    19,
		},
	}

	for _, tt := range tests {
		if tt.bitlist.Len() != tt.want {
			t.Errorf("(%x).Len() = %d, wanted %d", tt.bitlist, tt.bitlist.Len(), tt.want)
		}
	}
}

func TestByteBitlist_BitAt(t *testing.T) {
	tests := []struct {
		bitlist ByteBitlist
		idx     uint64
		want    bool
	}{
		{
			bitlist: ByteBitlist{},
			idx:     0,
			want:    false,
		},
		{
			bitlist: ByteBitlist{0x01}, // 0b00000001
			idx:     55,            // Out of bounds
			want:    false,
		},
		{
			bitlist: ByteBitlist{0x01}, // 0b00000001
			idx:     0,             //          ^ (length bit)
			want:    false,
		},
		{
			bitlist: ByteBitlist{0x0E}, // 0b00001110
			idx:     0,             //          ^
			want:    false,
		},
		{
			bitlist: ByteBitlist{0x0E}, // 0b00001110
			idx:     1,             //         ^
			want:    true,
		},
		{
			bitlist: ByteBitlist{0x0E}, // 0b00001110
			idx:     3,             //       ^
			want:    false,
		},
		{
			bitlist: ByteBitlist{0x0E}, // 0b00001110
			idx:     4,             //       ^ (length bit)
			want:    false,
		},
		{
			bitlist: ByteBitlist{0xFF, 0x0F}, // 0b11111111, 0b00001111
			idx:     4,                   //      ^
			want:    true,
		},
		{
			bitlist: ByteBitlist{0xFF, 0x0F}, // 0b11111111, 0b00001111
			idx:     12,                  //                  ^
			want:    false,
		},
		{
			bitlist: ByteBitlist{0xFF, 0x0F}, // 0b11111111, 0b00001111
			idx:     11,                  //                   ^ (length bit)
			want:    false,
		},
		{
			bitlist: ByteBitlist{0x00, 0x0F}, // 0b00000000, 0b00001111
			idx:     10,                  //                    ^
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

func TestByteBitlist_SetBitAt(t *testing.T) {
	tests := []struct {
		bitlist ByteBitlist
		idx     uint64
		val     bool
		want    ByteBitlist
	}{
		{
			bitlist: ByteBitlist{},
			idx:     0,
			val:     true,
			want:    ByteBitlist{},
		},
		{
			bitlist: ByteBitlist{0x01}, // 0b00000001
			idx:     0,             //          ^
			val:     true,
			want:    ByteBitlist{0x01}, // 0b00000001
		},
		{
			bitlist: ByteBitlist{0x02}, // 0b00000010
			idx:     0,             //          ^
			val:     true,
			want:    ByteBitlist{0x03}, // 0b00000011
		},
		{
			bitlist: ByteBitlist{0x10}, // 0b00010000
			idx:     0,             //          ^
			val:     true,
			want:    ByteBitlist{0x11}, // 0b00010001
		},
		{
			bitlist: ByteBitlist{0x10}, // 0b00010000
			idx:     0,             //          ^
			val:     true,
			want:    ByteBitlist{0x11}, // 0b00010001
		},
		{
			bitlist: ByteBitlist{0x10}, // 0b00010000
			idx:     64,            // Out of bounds
			val:     true,
			want:    ByteBitlist{0x10}, // 0b00010001
		},
		{
			bitlist: ByteBitlist{0x1F}, // 0b00011111
			idx:     0,             //          ^
			val:     true,
			want:    ByteBitlist{0x1F}, // 0b00011111
		},
		{
			bitlist: ByteBitlist{0x1F}, // 0b00011111
			idx:     1,             //         ^
			val:     false,
			want:    ByteBitlist{0x1D}, // 0b00011101
		},
		{
			bitlist: ByteBitlist{0x1F}, // 0b00011111
			idx:     4,             //      ^ (length bit)
			val:     false,
			want:    ByteBitlist{0x1F}, // 0b00011111
		},
		{
			bitlist: ByteBitlist{0x1F}, // 0b00011111
			idx:     64,            // Out of bounds
			val:     false,
			want:    ByteBitlist{0x1F}, // 0b00011111
		},
		{
			bitlist: ByteBitlist{0x1F, 0x01}, // 0b00011111, 0b00000001
			idx:     0,                   //          ^
			val:     true,
			want:    ByteBitlist{0x1F, 0x01}, // 0b00011111, 0b00000001
		},
		{
			bitlist: ByteBitlist{0x1F, 0x01}, // 0b00011111, 0b00000001
			idx:     0,                   //          ^
			val:     false,
			want:    ByteBitlist{0x1E, 0x01}, // 0b00011110, 0b00000001
		},
		{
			bitlist: ByteBitlist{0x00, 0x10}, // 0b00000000, 0b00010000
			idx:     8,                   //                      ^
			val:     true,
			want:    ByteBitlist{0x00, 0x11}, // 0b00000000, 0b00010001
		},
		{
			bitlist: ByteBitlist{0x00, 0x11}, // 0b00000000, 0b00010001
			idx:     8,                   //                      ^
			val:     false,
			want:    ByteBitlist{0x00, 0x10}, // 0b00000000, 0b00010000
		},
	}

	for _, tt := range tests {
		original := make(ByteBitlist, len(tt.bitlist))
		copy(original, tt.bitlist)

		tt.bitlist.SetBitAt(tt.idx, tt.val)
		if !bytes.Equal(tt.bitlist, tt.want) {
			t.Errorf(
				"(%x).SetBitAt(%d, %t) = %x, wanted %x",
				original,
				tt.idx,
				tt.val,
				tt.bitlist,
				tt.want,
			)
		}
	}
}

func TestByteBitlist_Bytes(t *testing.T) {
	tests := []struct {
		bitlist ByteBitlist
		want    []byte
	}{
		{
			bitlist: ByteBitlist{},
			want:    []byte{},
		},
		{
			bitlist: ByteBitlist{0x00},
			want:    []byte{},
		},
		{
			bitlist: ByteBitlist{0x01},
			want:    []byte{},
		},
		{
			bitlist: ByteBitlist{0x02},
			want:    []byte{},
		},
		{
			bitlist: ByteBitlist{0x03},
			want:    []byte{0x01},
		},
		{
			bitlist: ByteBitlist{0x12},
			want:    []byte{0x02},
		},
		{
			bitlist: ByteBitlist{0x02, 0x01},
			want:    []byte{0x02},
		},
		{
			bitlist: ByteBitlist{0x02, 0x02},
			want:    []byte{0x02},
		},
		{
			bitlist: ByteBitlist{0x02, 0x01},
			want:    []byte{0x02},
		},
		{
			bitlist: ByteBitlist{0x02, 0x03},
			want:    []byte{0x02, 0x01},
		},
		{
			bitlist: ByteBitlist{0x01, 0x00, 0x08},
			want:    []byte{0x01},
		}, {
			bitlist: ByteBitlist{0x00, 0x00, 0x02},
			want:    []byte{},
		},
	}

	for _, tt := range tests {
		got := tt.bitlist.Bytes()
		if !bytes.Equal(got, tt.want) {
			t.Errorf(
				"(%x).Bytes() = %x, wanted %x",
				tt.bitlist,
				got,
				tt.want,
			)
		}
	}
}

func TestByteBitlist_Count(t *testing.T) {
	tests := []struct {
		bitlist ByteBitlist
		want    uint64
	}{
		{
			bitlist: ByteBitlist{},
			want:    0,
		},
		{
			bitlist: ByteBitlist{0x00}, // 0b00000000, invalid list
			want:    0,
		},
		{
			bitlist: ByteBitlist{0x01}, // 0b00000001
			want:    0,
		},
		{
			bitlist: ByteBitlist{0x03}, // 0b00000011
			want:    1,
		},
		{
			bitlist: ByteBitlist{0x0F}, // 0b00001111
			want:    3,
		},
		{
			bitlist: ByteBitlist{0x0F, 0x01}, // 0b00001111, 0b00000001
			want:    4,
		},
		{
			bitlist: ByteBitlist{0x0F, 0x03}, // 0b00001111, 0b00000011
			want:    5,
		},
	}

	for _, tt := range tests {
		if tt.bitlist.Count() != tt.want {
			t.Errorf(
				"(%x).Count() = %d, wanted %d",
				tt.bitlist,
				tt.bitlist.Count(),
				tt.want,
			)
		}
	}
}

func TestByteBitlist_Contains(t *testing.T) {
	tests := []struct {
		a    ByteBitlist
		b    ByteBitlist
		want bool
	}{
		{
			a:    ByteBitlist{0x02}, // 0b00000010
			b:    ByteBitlist{0x03}, // 0b00000011
			want: false,
		},
		{
			a:    ByteBitlist{0x03}, // 0b00000011
			b:    ByteBitlist{0x03}, // 0b00000011
			want: true,
		},
		{
			a:    ByteBitlist{0x13}, // 0b00010011
			b:    ByteBitlist{0x15}, // 0b00010101
			want: false,
		},
		{
			a:    ByteBitlist{0x1F}, // 0b00011111
			b:    ByteBitlist{0x13}, // 0b00010011
			want: true,
		},
		{
			a:    ByteBitlist{0x1F}, // 0b00011111
			b:    ByteBitlist{0x13}, // 0b00010011
			want: true,
		},
		{
			a:    ByteBitlist{0x1F, 0x03}, // 0b00011111, 0b00000011
			b:    ByteBitlist{0x13, 0x02}, // 0b00010011, 0b00000010
			want: true,
		},
		{
			a:    ByteBitlist{0x1F, 0x01}, // 0b00011111, 0b00000001
			b:    ByteBitlist{0x93, 0x01}, // 0b10010011, 0b00000001
			want: false,
		},
		{
			a:    ByteBitlist{0xFF, 0x02}, // 0b11111111, 0x00000010
			b:    ByteBitlist{0x13, 0x03}, // 0b00010011, 0x00000011
			want: false,
		},
		{
			a:    ByteBitlist{0xFF, 0x85}, // 0b11111111, 0x10000111
			b:    ByteBitlist{0x13, 0x8F}, // 0b00010011, 0x10001111
			want: false,
		},
		{
			a:    ByteBitlist{0xFF, 0x8F}, // 0b11111111, 0x10001111
			b:    ByteBitlist{0x13, 0x83}, // 0b00010011, 0x10000011
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

func TestByteBitlist_Overlaps(t *testing.T) {
	tests := []struct {
		a    ByteBitlist
		b    ByteBitlist
		want bool
	}{
		{
			a:    ByteBitlist{0x06}, // 0b00000110
			b:    ByteBitlist{0x05}, // 0b00000101
			want: false,
		},
		{
			a:    ByteBitlist{0x32}, // 0b00110010
			b:    ByteBitlist{0x21}, // 0b00100001
			want: false,
		},
		{
			a:    ByteBitlist{0x41}, // 0b00100001
			b:    ByteBitlist{0x40}, // 0b00100000
			want: false,
		},
		{
			a:    ByteBitlist{0x1F}, // 0b00011111
			b:    ByteBitlist{0x11}, // 0b00010001
			want: true,
		},
		{
			a:    ByteBitlist{0xFF, 0x85}, // 0b11111111, 0b10000111
			b:    ByteBitlist{0x13, 0x8F}, // 0b00010011, 0b10001111
			want: true,
		},
		{
			a:    ByteBitlist{0x01, 0x40}, // 0b00000001, 0b01000000
			b:    ByteBitlist{0x00, 0x40}, // 0b00000010, 0b01000000
			want: false,
		},
		{
			a:    ByteBitlist{0x01, 0x40}, // 0b00000001, 0b01000000
			b:    ByteBitlist{0x00, 0x40}, // 0b00000010, 0b01000000
			want: false,
		},
		{
			a:    ByteBitlist{0x01, 0x40}, // 0b00000001, 0b01000000
			b:    ByteBitlist{0x02, 0x40}, // 0b00000010, 0b01000000
			want: false,
		},
		{
			a:    ByteBitlist{0x01, 0x40}, // 0b00000001, 0b01000000
			b:    ByteBitlist{0x03, 0x40}, // 0b00000011, 0b01000000
			want: true,
		},
		{
			a:    ByteBitlist{0x01, 0x01, 0x01}, // 0b00000001, 0b00000001, 0b00000001
			b:    ByteBitlist{0x02, 0x00, 0x01}, // 0b00000010, 0b00000000, 0b00000001
			want: false,
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

func TestByteBitlist_Or(t *testing.T) {
	tests := []struct {
		a    ByteBitlist
		b    ByteBitlist
		want ByteBitlist
	}{
		{
			a:    ByteBitlist{0x02}, // 0b00000010
			b:    ByteBitlist{0x03}, // 0b00000011
			want: ByteBitlist{0x03}, // 0b00000011
		},
		{
			a:    ByteBitlist{0x03}, // 0b00000011
			b:    ByteBitlist{0x03}, // 0b00000011
			want: ByteBitlist{0x03}, // 0b00000011
		},
		{
			a:    ByteBitlist{0x13}, // 0b00010011
			b:    ByteBitlist{0x15}, // 0b00010101
			want: ByteBitlist{0x17}, // 0b00010111
		},
		{
			a:    ByteBitlist{0x1F}, // 0b00011111
			b:    ByteBitlist{0x13}, // 0b00010011
			want: ByteBitlist{0x1F}, // 0b00011111
		},
		{
			a:    ByteBitlist{0x1F, 0x03}, // 0b00011111, 0b00000011
			b:    ByteBitlist{0x13, 0x02}, // 0b00010011, 0b00000010
			want: ByteBitlist{0x1F, 0x03}, // 0b00011111, 0b00000011
		},
		{
			a:    ByteBitlist{0x1F, 0x01}, // 0b00011111, 0b00000001
			b:    ByteBitlist{0x93, 0x01}, // 0b10010011, 0b00000001
			want: ByteBitlist{0x9F, 0x01}, // 0b00011111, 0b00000001
		},
		{
			a:    ByteBitlist{0xFF, 0x02}, // 0b11111111, 0x00000010
			b:    ByteBitlist{0x13, 0x03}, // 0b00010011, 0x00000011
			want: ByteBitlist{0xFF, 0x03}, // 0b11111111, 0x00000011
		},
		{
			a:    ByteBitlist{0xFF, 0x85}, // 0b11111111, 0x10000111
			b:    ByteBitlist{0x13, 0x8F}, // 0b00010011, 0x10001111
			want: ByteBitlist{0xFF, 0x8F}, // 0b11111111, 0x10001111
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

func TestByteBitlist_And(t *testing.T) {
	tests := []struct {
		a    ByteBitlist
		b    ByteBitlist
		want ByteBitlist
	}{
		{
			a:    ByteBitlist{0x02}, // 0b00000010
			b:    ByteBitlist{0x03}, // 0b00000011
			want: ByteBitlist{0x02}, // 0b00000010
		},
		{
			a:    ByteBitlist{0x03}, // 0b00000011
			b:    ByteBitlist{0x03}, // 0b00000011
			want: ByteBitlist{0x03}, // 0b00000011
		},
		{
			a:    ByteBitlist{0x13}, // 0b00010011
			b:    ByteBitlist{0x15}, // 0b00010101
			want: ByteBitlist{0x11}, // 0b00010001
		},
		{
			a:    ByteBitlist{0x1F}, // 0b00011111
			b:    ByteBitlist{0x13}, // 0b00010011
			want: ByteBitlist{0x13}, // 0b00010011
		},
		{
			a:    ByteBitlist{0x1F, 0x03}, // 0b00011111, 0b00000011
			b:    ByteBitlist{0x13, 0x02}, // 0b00010011, 0b00000010
			want: ByteBitlist{0x13, 0x02}, // 0b00010011, 0b00000010
		},
		{
			a:    ByteBitlist{0x9F, 0x01}, // 0b10011111, 0b00000001
			b:    ByteBitlist{0x93, 0x01}, // 0b10010011, 0b00000001
			want: ByteBitlist{0x93, 0x01}, // 0b10010011, 0b00000001
		},
		{
			a:    ByteBitlist{0xFF, 0x02}, // 0b11111111, 0x00000010
			b:    ByteBitlist{0x13, 0x03}, // 0b00010011, 0x00000011
			want: ByteBitlist{0x13, 0x02}, // 0b00010011, 0x00000010
		},
		{
			a:    ByteBitlist{0xFF, 0x87}, // 0b11111111, 0x10000111
			b:    ByteBitlist{0x13, 0x8F}, // 0b00010011, 0x10001111
			want: ByteBitlist{0x13, 0x87}, // 0b00010011, 0x10000111
		},
	}

	for _, tt := range tests {
		if !bytes.Equal(tt.a.And(tt.b), tt.want) {
			t.Errorf(
				"(%x).And(%x) = %x, wanted %x",
				tt.a,
				tt.b,
				tt.a.And(tt.b),
				tt.want,
			)
		}
	}
}

func TestByteBitlist_Xor(t *testing.T) {
	tests := []struct {
		a    ByteBitlist
		b    ByteBitlist
		want ByteBitlist
	}{
		{
			a:    ByteBitlist{0x02}, // 0b00000010
			b:    ByteBitlist{0x03}, // 0b00000011
			want: ByteBitlist{0x03}, // 0b00000011
		},
		{
			a:    ByteBitlist{0x03}, // 0b00000011
			b:    ByteBitlist{0x03}, // 0b00000011
			want: ByteBitlist{0x02}, // 0b00000010
		},
		{
			a:    ByteBitlist{0x13}, // 0b00010011
			b:    ByteBitlist{0x15}, // 0b00010101
			want: ByteBitlist{0x16}, // 0b00010110
		},
		{
			a:    ByteBitlist{0x1F}, // 0b00011111
			b:    ByteBitlist{0x13}, // 0b00010011
			want: ByteBitlist{0x1c}, // 0b00011100
		},
		{
			a:    ByteBitlist{0x1F, 0x03}, // 0b00011111, 0b00000011
			b:    ByteBitlist{0x13, 0x02}, // 0b00010011, 0b00000010
			want: ByteBitlist{0x0c, 0x03}, // 0b00001100, 0b00000011
		},
		{
			a:    ByteBitlist{0x9F, 0x01}, // 0b10011111, 0b00000001
			b:    ByteBitlist{0x93, 0x01}, // 0b10010011, 0b00000001
			want: ByteBitlist{0x0c, 0x01}, // 0b00001100, 0b00000001
		},
		{
			a:    ByteBitlist{0xFF, 0x02}, // 0b11111111, 0x00000010
			b:    ByteBitlist{0x13, 0x03}, // 0b00010011, 0x00000011
			want: ByteBitlist{0xec, 0x03}, // 0b11101100, 0x00000011
		},
		{
			a:    ByteBitlist{0xFF, 0x87}, // 0b11111111, 0x10000111
			b:    ByteBitlist{0x13, 0x8F}, // 0b00010011, 0x10001111
			want: ByteBitlist{0xec, 0x88}, // 0b11101100, 0x10001000
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("(%x).Xor(%x)", tt.a, tt.b), func(t *testing.T) {
			if !bytes.Equal(tt.a.Xor(tt.b), tt.want) {
				t.Errorf(
					"(%x).Xor(%x) = %x, wanted %x",
					tt.a,
					tt.b,
					tt.a.Xor(tt.b),
					tt.want,
				)
			}
		})
	}
}

func TestByteBitlist_Not(t *testing.T) {
	tests := []struct {
		a    ByteBitlist
		want ByteBitlist
	}{
		{
			a:    ByteBitlist{0x01}, // 0b00000001
			want: ByteBitlist{0x01}, // 0b00000001
		},
		{
			a:    ByteBitlist{0x02}, // 0b00000010
			want: ByteBitlist{0x03}, // 0b00000011
		},
		{
			a:    ByteBitlist{0x03}, // 0b00000011
			want: ByteBitlist{0x02}, // 0b00000010
		},
		{
			a:    ByteBitlist{0x05}, // 0b00000101
			want: ByteBitlist{0x06}, // 0b00000110
		},
		{
			a:    ByteBitlist{0x06}, // 0b00000110
			want: ByteBitlist{0x05}, // 0b00000101
		},
		{
			a:    ByteBitlist{0x83}, // 0b10000011
			want: ByteBitlist{0xfc}, // 0b11111100
		},
		{
			a:    ByteBitlist{0x13}, // 0b00010011
			want: ByteBitlist{0x1c}, // 0b00011100
		},
		{
			a:    ByteBitlist{0x1F}, // 0b00011111
			want: ByteBitlist{0x10}, // 0b00010000
		},
		{
			a:    ByteBitlist{0x1F, 0x03}, // 0b00011111, 0b00000011
			want: ByteBitlist{0xe0, 0x02}, // 0b11100000, 0b00000010
		},
		{
			a:    ByteBitlist{0x9F, 0x01}, // 0b10011111, 0b00000001
			want: ByteBitlist{0x60, 0x01}, // 0b01100000, 0b00000001
		},
		{
			a:    ByteBitlist{0xFF, 0x02}, // 0b11111111, 0x00000010
			want: ByteBitlist{0x00, 0x03}, // 0b00000000, 0x00000011
		},
		{
			a:    ByteBitlist{0xFF, 0x87}, // 0b11111111, 0x10000111
			want: ByteBitlist{0x00, 0xf8}, // 0b00000000, 0x11111000
		},
		{
			a:    ByteBitlist{0xFF, 0x07}, // 0b11111111, 0x00000111
			want: ByteBitlist{0x00, 0x04}, // 0b00000000, 0x00000100
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("(%#x).Not()", tt.a), func(t *testing.T) {
			if !bytes.Equal(tt.a.Not(), tt.want) {
				t.Errorf(
					"(%x).Not() = %x, wanted %x",
					tt.a,
					tt.a.Not(),
					tt.want,
				)
			}
		})
	}
}

func TestByteBitlist_BitIndices(t *testing.T) {
	tests := []struct {
		a    ByteBitlist
		want []int
	}{
		{
			a:    ByteBitlist{0b10010},
			want: []int{1},
		},
		{
			a:    ByteBitlist{0b10000},
			want: []int{},
		},
		{
			a:    ByteBitlist{0b10, 0b1},
			want: []int{1},
		},
		{
			a:    ByteBitlist{0b11111111, 0b11},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
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
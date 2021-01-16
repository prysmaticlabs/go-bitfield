package bitfield

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBitlist_NewBitlist(t *testing.T) {
	makeData := func(n uint64) []uint64 {
		return make([]uint64, n, n)
	}
	tests := []struct {
		size uint64
		want *Bitlist
	}{
		{
			size: 0,
			want: &Bitlist{size: 0, data: []uint64{}},
		},
		{
			size: 1,
			want: &Bitlist{size: 1, data: []uint64{0x00}},
		},
		{
			size: 2,
			want: &Bitlist{size: 2, data: []uint64{0x00}},
		},
		{
			size: 3,
			want: &Bitlist{size: 3, data: []uint64{0x00}},
		},
		{
			size: 8,
			want: &Bitlist{size: 8, data: []uint64{0x00}},
		},
		{
			size: 9,
			want: &Bitlist{size: 9, data: []uint64{0x00}},
		},
		{
			size: 31,
			want: &Bitlist{size: 31, data: []uint64{0x00}},
		},
		{
			size: 32,
			want: &Bitlist{size: 32, data: []uint64{0x00}},
		},
		{
			size: 63,
			want: &Bitlist{size: 63, data: []uint64{0x00}},
		},
		{
			size: 64,
			want: &Bitlist{size: 64, data: []uint64{0x00}},
		},
		{
			size: 65,
			want: &Bitlist{size: 65, data: []uint64{0x00, 0x00}},
		},
		{
			size: 128,
			want: &Bitlist{size: 128, data: []uint64{0x00, 0x00}},
		},
		{
			size: 256,
			want: &Bitlist{size: 256, data: []uint64{0x00, 0x00, 0x00, 0x00}},
		},
		{
			size: 512,
			want: &Bitlist{size: 512, data: []uint64{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		},
		{
			size: 1024,
			want: &Bitlist{size: 1024, data: makeData(1024 / wordSize)},
		},
		{
			size: 2048,
			want: &Bitlist{size: 2048, data: makeData(2048 / wordSize)},
		},
		{
			size: 4096,
			want: &Bitlist{size: 4096, data: makeData(4096 / wordSize)},
		},
		{
			// 10000/wordSizeLog2 = 156,7 ~ 157 (where wordSizeLog2 = log_2(wordSize = 64) = 6.
			size: 10000,
			want: &Bitlist{size: 10000, data: makeData(157)},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("size:%d", tt.size), func(t *testing.T) {
			got := NewBitlist(tt.size)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBitlist(%d) = %+v, wanted %+v", tt.size, got, tt.want)
			}
		})
	}
}

func TestBitlist_NewBitlistFrom(t *testing.T) {
	tests := []struct {
		from []uint64
		want *Bitlist
	}{
		{
			from: []uint64{},
			want: &Bitlist{size: 0, data: []uint64{}},
		},
		{
			from: []uint64{0x0000000000000000},
			want: &Bitlist{size: 64, data: []uint64{0x0000000000000000}},
		},
		{
			from: []uint64{0x001002000c002000},
			want: &Bitlist{size: 64, data: []uint64{0x001002000c002000}},
		},
		{
			from: []uint64{0xFFFFFFFFFFFFFFFF},
			want: &Bitlist{size: 64, data: []uint64{0xFFFFFFFFFFFFFFFF}},
		},
		{
			from: []uint64{0x0000000000000000, 0x0000000000000000},
			want: &Bitlist{size: 128, data: []uint64{0x00, 0x00}},
		},
		{
			from: []uint64{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF},
			want: &Bitlist{size: 128, data: []uint64{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF}},
		},
		{
			from: []uint64{0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000},
			want: &Bitlist{size: 256, data: []uint64{0x00, 0x00, 0x00, 0x00}},
		},
		{
			from: []uint64{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF},
			want: &Bitlist{
				size: 256,
				data: []uint64{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF},
			},
		},
		{
			from: []uint64{
				0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000,
				0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000,
			},
			want: &Bitlist{
				size: 512,
				data: []uint64{
					0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000,
					0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000,
				},
			},
		},
		{
			from: []uint64{
				0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
				0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
			},
			want: &Bitlist{
				size: 512,
				data: []uint64{
					0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
					0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
				},
			},
		},
		{
			from: []uint64{
				0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
				0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
				0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
				0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
				0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
				0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
				0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
				0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
			},
			want: &Bitlist{
				size: 2048,
				data: []uint64{
					0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
					0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
					0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
					0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
					0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
					0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
					0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
					0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0x1111FFFFFFFFCCCC, 0X1111FFFFFFFFCCCC,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("data:%#x", tt.from), func(t *testing.T) {
			got := NewBitlistFrom(tt.from)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBitlistFrom(%#x) = %+v, wanted %+v", tt.from, got, tt.want)
			}
		})
	}
}

func TestBitlist_Len(t *testing.T) {
	tests := []struct {
		bitlist *Bitlist
		want    uint64
	}{
		{
			bitlist: NewBitlist(0),
			want:    0,
		},
		{
			bitlist: NewBitlistFrom([]uint64{}),
			want:    0,
		},
		{
			bitlist: NewBitlistFrom([]uint64{0x00}),
			want:    wordSize,
		},
		{
			bitlist: NewBitlistFrom([]uint64{0x01}),
			want:    wordSize,
		},
		{
			bitlist: NewBitlistFrom([]uint64{0x02}),
			want:    wordSize,
		},
		{
			bitlist: NewBitlistFrom([]uint64{0x08}),
			want:    wordSize,
		},
		{
			bitlist: NewBitlistFrom([]uint64{0x0E}),
			want:    wordSize,
		},
		{
			bitlist: NewBitlistFrom([]uint64{0x0F}),
			want:    wordSize,
		},
		{
			bitlist: NewBitlistFrom([]uint64{0x00, 0x01}),
			want:    wordSize * 2,
		},
		{
			bitlist: NewBitlistFrom([]uint64{0x00, 0x02}),
			want:    wordSize * 2,
		},
		{
			bitlist: NewBitlistFrom([]uint64{0x00, 0x02, 0x08}),
			want:    wordSize * 3,
		},
		{
			bitlist: NewBitlistFrom(make([]uint64, 2048)),
			want:    wordSize * 2048,
		},
	}

	for _, tt := range tests {
		if tt.bitlist.Len() != tt.want {
			t.Errorf("(%+v).Len() = %d, wanted %d", tt.bitlist, tt.bitlist.Len(), tt.want)
		}
	}
}

func TestBitlist_BitAt(t *testing.T) {
	tests := []struct {
		bitlist []uint64
		idx     uint64
		want    bool
	}{
		{
			bitlist: []uint64{},
			idx:     0,
			want:    false,
		},
		{
			bitlist: []uint64{0x01},
			idx:     64, // Out of bounds
			want:    false,
		},
		{
			bitlist: []uint64{0x01},
			idx:     163465, // Out of bounds
			want:    false,
		},
		{
			bitlist: []uint64{0x01},
			idx:     0,
			want:    true,
		},
		{
			bitlist: []uint64{0x0E}, // 0b00001110
			idx:     0,              //          ^
			want:    false,
		},
		{
			bitlist: []uint64{0x0E}, // 0b00001110
			idx:     1,              //         ^
			want:    true,
		},
		{
			bitlist: []uint64{0x0E}, // 0b00001110
			idx:     2,              //        ^
			want:    true,
		},
		{
			bitlist: []uint64{0x0E}, // 0b00001110
			idx:     3,              //       ^
			want:    true,
		},
		{
			bitlist: []uint64{0x0E}, // 0b00001110
			idx:     4,              //      ^
			want:    false,
		},
		{
			bitlist: []uint64{0xFF, 0x0F}, // 0b11111111, 0b00001111
			idx:     4,                    //      ^
			want:    true,
		},
		{
			bitlist: []uint64{0x00, 0x0F}, // 0b00000000, 0b00001111
			idx:     67,                   //                   ^
			want:    true,
		},
		{
			bitlist: []uint64{0xFF, 0x0F}, // 0b11111111, 0b00001111
			idx:     68,                   //                  ^
			want:    false,
		},
		{
			bitlist: []uint64{0x00, 0x00, 0b00000100}, // 0b0, 0b0, 0b00000100
			idx:     130,                              //                  ^
			want:    true,
		},
		{
			bitlist: []uint64{0x00, 0x00, 0b00000100}, // 0b0, 0b0, 0b00000100
			idx:     129,                              //                   ^
			want:    false,
		},
		{
			bitlist: []uint64{0x00, 0x00, 0b00000100}, // 0b0, 0b0, 0b00000100
			idx:     131,                              //                 ^
			want:    false,
		},
	}

	for _, tt := range tests {
		if NewBitlistFrom(tt.bitlist).BitAt(tt.idx) != tt.want {
			t.Errorf(
				"(%#b).BitAt(%d) = %t, wanted %t",
				tt.bitlist,
				tt.idx,
				NewBitlistFrom(tt.bitlist).BitAt(tt.idx),
				tt.want,
			)
		}
	}
}

func TestBitlist_SetBitAt(t *testing.T) {
	tests := []struct {
		bitlist []uint64
		idx     uint64
		val     bool
		want    []uint64
	}{
		{
			bitlist: []uint64{},
			idx:     0,
			val:     true,
			want:    []uint64{},
		},
		{
			bitlist: []uint64{0x01}, // 0b00000001
			idx:     0,              //          ^
			val:     true,
			want:    []uint64{0x01}, // 0b00000001
		},
		{
			bitlist: []uint64{0x01}, // 0b00000001
			idx:     1,              //         ^
			val:     true,
			want:    []uint64{0x03}, // 0b00000011
		},
		{
			bitlist: []uint64{0x01}, // 0b00000001
			idx:     2,              //        ^
			val:     true,
			want:    []uint64{0x05}, // 0b00000101
		},
		{
			bitlist: []uint64{0x02}, // 0b00000010
			idx:     0,              //          ^
			val:     true,
			want:    []uint64{0x03}, // 0b00000011
		},
		{
			bitlist: []uint64{0x10}, // 0b00010000
			idx:     0,              //          ^
			val:     true,
			want:    []uint64{0x11}, // 0b00010001
		},
		{
			bitlist: []uint64{0x10}, // 0b00010000
			idx:     64,             // Out of bounds
			val:     true,
			want:    []uint64{0x10}, // 0b00010000
		},
		{
			bitlist: []uint64{0x10}, // 0b00010000
			idx:     63,
			val:     true,
			want:    []uint64{0x8000000000000010}, // 0b1000..010000
		},
		{
			bitlist: []uint64{0x1F}, // 0b00011111
			idx:     0,              //          ^
			val:     true,
			want:    []uint64{0x1F}, // 0b00011111
		},
		{
			bitlist: []uint64{0x1F}, // 0b00011111
			idx:     1,              //         ^
			val:     false,
			want:    []uint64{0x1D}, // 0b00011101
		},
		{
			bitlist: []uint64{0x1F}, // 0b00011111
			idx:     4,              //      ^
			val:     false,
			want:    []uint64{0x0F}, // 0b00001111
		},
		{
			bitlist: []uint64{0x1F}, // 0b00011111
			idx:     64,             // Out of bounds
			val:     false,
			want:    []uint64{0x1F}, // 0b00011111
		},
		{
			bitlist: []uint64{0x1F, 0x01}, // 0b00011111, 0b00000001
			idx:     0,                    //          ^
			val:     true,
			want:    []uint64{0x1F, 0x01}, // 0b00011111, 0b00000001
		},
		{
			bitlist: []uint64{0x1F, 0x01}, // 0b00011111, 0b00000001
			idx:     0,                    //          ^
			val:     false,
			want:    []uint64{0x1E, 0x01}, // 0b00011110, 0b00000001
		},
		{
			bitlist: []uint64{0x00, 0x10}, // 0b00000000, 0b00010000
			idx:     64,                   //                      ^
			val:     true,
			want:    []uint64{0x00, 0x11}, // 0b00000000, 0b00010001
		},
		{
			bitlist: []uint64{0x00, 0x11}, // 0b00000000, 0b00010001
			idx:     64,                   //                      ^
			val:     false,
			want:    []uint64{0x00, 0x10}, // 0b00000000, 0b00010000
		},
	}

	for _, tt := range tests {
		s := NewBitlistFrom(tt.bitlist)
		s.SetBitAt(tt.idx, tt.val)
		if !reflect.DeepEqual(tt.want, s.data) {
			t.Errorf("(%+v).SetBitAt(%d, %t) = %x, wanted %x", s, tt.idx, tt.val, tt.bitlist, tt.want)
		}
	}
}

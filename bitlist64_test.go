package bitfield

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestBitlist64_NewBitlist(t *testing.T) {
	makeData := func(n uint64) []uint64 {
		return make([]uint64, n, n)
	}
	tests := []struct {
		size uint64
		want *Bitlist64
	}{
		{
			size: 0,
			want: &Bitlist64{size: 0, data: []uint64{}},
		},
		{
			size: 1,
			want: &Bitlist64{size: 1, data: []uint64{0x00}},
		},
		{
			size: 2,
			want: &Bitlist64{size: 2, data: []uint64{0x00}},
		},
		{
			size: 3,
			want: &Bitlist64{size: 3, data: []uint64{0x00}},
		},
		{
			size: 8,
			want: &Bitlist64{size: 8, data: []uint64{0x00}},
		},
		{
			size: 9,
			want: &Bitlist64{size: 9, data: []uint64{0x00}},
		},
		{
			size: 31,
			want: &Bitlist64{size: 31, data: []uint64{0x00}},
		},
		{
			size: 32,
			want: &Bitlist64{size: 32, data: []uint64{0x00}},
		},
		{
			size: 63,
			want: &Bitlist64{size: 63, data: []uint64{0x00}},
		},
		{
			size: 64,
			want: &Bitlist64{size: 64, data: []uint64{0x00}},
		},
		{
			size: 65,
			want: &Bitlist64{size: 65, data: []uint64{0x00, 0x00}},
		},
		{
			size: 128,
			want: &Bitlist64{size: 128, data: []uint64{0x00, 0x00}},
		},
		{
			size: 256,
			want: &Bitlist64{size: 256, data: []uint64{0x00, 0x00, 0x00, 0x00}},
		},
		{
			size: 512,
			want: &Bitlist64{size: 512, data: []uint64{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		},
		{
			size: 1024,
			want: &Bitlist64{size: 1024, data: makeData(1024 / wordSize)},
		},
		{
			size: 2048,
			want: &Bitlist64{size: 2048, data: makeData(2048 / wordSize)},
		},
		{
			size: 4096,
			want: &Bitlist64{size: 4096, data: makeData(4096 / wordSize)},
		},
		{
			// 10000/wordSizeLog2 = 156,7 ~ 157 (where wordSizeLog2 = log_2(wordSize = 64) = 6.
			size: 10000,
			want: &Bitlist64{size: 10000, data: makeData(157)},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("size:%d", tt.size), func(t *testing.T) {
			got := NewBitlist64(tt.size)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBitlist64(%d) = %+v, wanted %+v", tt.size, got, tt.want)
			}
		})
	}
}

func TestBitlist64_NewBitlistFrom(t *testing.T) {
	tests := []struct {
		from []uint64
		want *Bitlist64
	}{
		{
			from: []uint64{},
			want: &Bitlist64{size: 0, data: []uint64{}},
		},
		{
			from: []uint64{0x0000000000000000},
			want: &Bitlist64{size: 64, data: []uint64{0x0000000000000000}},
		},
		{
			from: []uint64{0x001002000c002000},
			want: &Bitlist64{size: 64, data: []uint64{0x001002000c002000}},
		},
		{
			from: []uint64{0xFFFFFFFFFFFFFFFF},
			want: &Bitlist64{size: 64, data: []uint64{0xFFFFFFFFFFFFFFFF}},
		},
		{
			from: []uint64{0x0000000000000000, 0x0000000000000000},
			want: &Bitlist64{size: 128, data: []uint64{0x00, 0x00}},
		},
		{
			from: []uint64{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF},
			want: &Bitlist64{size: 128, data: []uint64{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF}},
		},
		{
			from: []uint64{0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000},
			want: &Bitlist64{size: 256, data: []uint64{0x00, 0x00, 0x00, 0x00}},
		},
		{
			from: []uint64{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF},
			want: &Bitlist64{
				size: 256,
				data: []uint64{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF},
			},
		},
		{
			from: []uint64{
				0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000,
				0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000,
			},
			want: &Bitlist64{
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
			want: &Bitlist64{
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
			want: &Bitlist64{
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
			got := NewBitlist64From(tt.from)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBitlist64From(%#x) = %+v, wanted %+v", tt.from, got, tt.want)
			}
		})
	}
}

func TestBitlist64_Len(t *testing.T) {
	tests := []struct {
		bitlist *Bitlist64
		want    uint64
	}{
		{
			bitlist: NewBitlist64(0),
			want:    0,
		},
		{
			bitlist: NewBitlist64From([]uint64{}),
			want:    0,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x00}),
			want:    wordSize,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x01}),
			want:    wordSize,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x02}),
			want:    wordSize,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x08}),
			want:    wordSize,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x0E}),
			want:    wordSize,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x0F}),
			want:    wordSize,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x00, 0x01}),
			want:    wordSize * 2,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x00, 0x02}),
			want:    wordSize * 2,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x00, 0x02, 0x08}),
			want:    wordSize * 3,
		},
		{
			bitlist: NewBitlist64From(make([]uint64, 2048)),
			want:    wordSize * 2048,
		},
	}

	for _, tt := range tests {
		if tt.bitlist.Len() != tt.want {
			t.Errorf("(%+v).Len() = %d, wanted %d", tt.bitlist, tt.bitlist.Len(), tt.want)
		}
	}
}

func TestBitlist64_BitAt(t *testing.T) {
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
		if NewBitlist64From(tt.bitlist).BitAt(tt.idx) != tt.want {
			t.Errorf(
				"(%#b).BitAt(%d) = %t, wanted %t",
				tt.bitlist,
				tt.idx,
				NewBitlist64From(tt.bitlist).BitAt(tt.idx),
				tt.want,
			)
		}
	}
}

func TestBitlist64_SetBitAt(t *testing.T) {
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
		s := NewBitlist64From(tt.bitlist)
		s.SetBitAt(tt.idx, tt.val)
		if !reflect.DeepEqual(tt.want, s.data) {
			t.Errorf("(%+v).SetBitAt(%d, %t) = %x, wanted %x", s, tt.idx, tt.val, tt.bitlist, tt.want)
		}
	}
}

func TestBitlist64_Bytes(t *testing.T) {
	tests := []struct {
		bitlist *Bitlist64
		want    []byte
	}{
		{
			bitlist: NewBitlist64From([]uint64{}),
			want:    []byte{},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x00}),
			want:    []byte{},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x01}),
			want:    []byte{0x01},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x02}),
			want:    []byte{0x02},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x03}),
			want:    []byte{0x03},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x12}),
			want:    []byte{0x12},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x02, 0x01}),
			want:    []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x02, 0x02}),
			want:    []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x02, 0x03}),
			want:    []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x01, 0x00, 0x00}),
			want:    []byte{0x01},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x01, 0x00, 0x001F00}),
			want: []byte{
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x1F,
			},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x00, 0x00, 0x00}),
			want:    []byte{},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x00, 0x01, 0x00}),
			want: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x01,
			},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x01, 0x00}),
			want:    []byte{0x01},
		},
		{
			bitlist: NewBitlist64From([]uint64{0x0807060504030201, 0x02}),
			want:    []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x02},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("bitlist:%+v", tt.bitlist), func(t *testing.T) {
			got := tt.bitlist.Bytes()
			if !bytes.Equal(got, tt.want) {
				t.Errorf("(%+v).Bytes() = %x, wanted %x", tt.bitlist, got, tt.want)
			}
		})
	}
}

func TestBitlist64_Count(t *testing.T) {
	tests := []struct {
		bitlist *Bitlist64
		want    uint64
	}{
		{
			bitlist: NewBitlist64From([]uint64{}),
			want:    0,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x00}), // 0b00000000
			want:    0,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x01}), // 0b00000001
			want:    1,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x03}), // 0b00000011
			want:    2,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x0F}), // 0b00001111
			want:    4,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x0F, 0x01}), // 0b00001111, 0b00000001
			want:    5,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x0F, 0x03}), // 0b00001111, 0b00000011
			want:    6,
		},
		{
			bitlist: NewBitlist64From([]uint64{0x0F, 0x00, 0x03}), // 0b00001111, 0b00000011
			want:    6,
		},
	}

	for _, tt := range tests {
		if tt.bitlist.Count() != tt.want {
			t.Errorf(
				"(%+v).Count() = %d, wanted %d",
				tt.bitlist,
				tt.bitlist.Count(),
				tt.want,
			)
		}
	}
}

func TestBitlist64_Contains(t *testing.T) {
	tests := []struct {
		a    *Bitlist64
		b    *Bitlist64
		want bool
	}{
		{
			a:    NewBitlist64From([]uint64{0x02}), // 0b00000010
			b:    NewBitlist64From([]uint64{0x03}), // 0b00000011
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0x03}), // 0b00000011
			b:    NewBitlist64From([]uint64{0x03}), // 0b00000011
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x13}), // 0b00010011
			b:    NewBitlist64From([]uint64{0x15}), // 0b00010101
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0x1F}), // 0b00011111
			b:    NewBitlist64From([]uint64{0x13}), // 0b00010011
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x1F}), // 0b00011111
			b:    NewBitlist64From([]uint64{0x13}), // 0b00010011
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x1F, 0x03}), // 0b00011111, 0b00000011
			b:    NewBitlist64From([]uint64{0x13, 0x02}), // 0b00010011, 0b00000010
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x1F, 0x01}), // 0b00011111, 0b00000001
			b:    NewBitlist64From([]uint64{0x93, 0x01}), // 0b10010011, 0b00000001
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0xFF, 0x02}), // 0b11111111, 0x00000010
			b:    NewBitlist64From([]uint64{0x13, 0x03}), // 0b00010011, 0x00000011
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0xFF, 0x85}), // 0b11111111, 0x10000111
			b:    NewBitlist64From([]uint64{0x13, 0x8F}), // 0b00010011, 0x10001111
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0xFF, 0x8F}), // 0b11111111, 0x10001111
			b:    NewBitlist64From([]uint64{0x13, 0x83}), // 0b00010011, 0x10000011
			want: true,
		},
	}

	for _, tt := range tests {
		if tt.a.Contains(tt.b) != tt.want {
			t.Errorf("(%+v).Contains(%+v) = %t, wanted %t", tt.a, tt.b, tt.a.Contains(tt.b), tt.want)
		}
	}

	t.Run("check panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic not thrown")
			}
		}()
		a := NewBitlist64(64)
		b := NewBitlist64(128)
		a.Contains(b)
	})
}

func TestBitlist64_Overlaps(t *testing.T) {
	tests := []struct {
		a    *Bitlist64
		b    *Bitlist64
		want bool
	}{
		{
			a:    NewBitlist64From([]uint64{}), // zero-length bitlist
			b:    NewBitlist64From([]uint64{}), // zero-length bitlist
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0x06}), // 0b00000110
			b:    NewBitlist64From([]uint64{0x05}), // 0b00000101
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x06}), // 0b00000110
			b:    NewBitlist64From([]uint64{0x01}), // 0b00000001
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0x32}), // 0b00110010
			b:    NewBitlist64From([]uint64{0x21}), // 0b00100001
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x32}), // 0b00110010
			b:    NewBitlist64From([]uint64{0x01}), // 0b00000001
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0x41}), // 0b00100001
			b:    NewBitlist64From([]uint64{0x40}), // 0b00100000
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x41}), // 0b00100001
			b:    NewBitlist64From([]uint64{0x00}), // 0b00000000
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0x1F}), // 0b00011111
			b:    NewBitlist64From([]uint64{0x11}), // 0b00010001
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0xFF, 0x85}), // 0b11111111, 0b10000111
			b:    NewBitlist64From([]uint64{0x13, 0x8F}), // 0b00010011, 0b10001111
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x01, 0x40}), // 0b00000001, 0b01000000
			b:    NewBitlist64From([]uint64{0x00, 0x40}), // 0b00000010, 0b01000000
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x01, 0x40}), // 0b00000001, 0b01000000
			b:    NewBitlist64From([]uint64{0x00, 0x00}), // 0b00000010, 0b00000000
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0x01, 0x40}), // 0b00000001, 0b01000000
			b:    NewBitlist64From([]uint64{0x00, 0x80}), // 0b00000010, 0b10000000
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0x01, 0x40}), // 0b00000001, 0b01000000
			b:    NewBitlist64From([]uint64{0x02, 0x80}), // 0b00000010, 0b10000000
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0x01, 0x40}), // 0b00000001, 0b01000000
			b:    NewBitlist64From([]uint64{0x03, 0x80}), // 0b00000011, 0b10000000
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x01, 0x40}), // 0b00000001, 0b01000000
			b:    NewBitlist64From([]uint64{0x02, 0x50}), // 0b00000010, 0b01010000
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x01, 0x40}), // 0b00000001, 0b01000000
			b:    NewBitlist64From([]uint64{0x02, 0x40}), // 0b00000010, 0b01000000
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x01, 0x00}), // 0b00000001, 0b01000000
			b:    NewBitlist64From([]uint64{0x02, 0x00}), // 0b00000010, 0b01000000
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0x01, 0x80}), // 0b00000001, 0b10000000
			b:    NewBitlist64From([]uint64{0x03, 0x40}), // 0b00000011, 0b01000000
			want: true,
		},
		{
			a:    NewBitlist64From([]uint64{0x01, 0x01, 0x02}), // 0b00000001, 0b00000001, 0b00000010
			b:    NewBitlist64From([]uint64{0x02, 0x00, 0x01}), // 0b00000010, 0b00000000, 0b00000001
			want: false,
		},
		{
			a:    NewBitlist64From([]uint64{0x01, 0x01, 0x02}), // 0b00000001, 0b00000001, 0b00000010
			b:    NewBitlist64From([]uint64{0x02, 0x03, 0x01}), // 0b00000010, 0b00000000, 0b00000001
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("bitlist:%+v,%+v", tt.a, tt.b), func(t *testing.T) {
			result := tt.a.Overlaps(tt.b)
			if result != tt.want {
				t.Errorf("(%+v).Overlaps(%+v) = %t, wanted %t", tt.a, tt.b, result, tt.want)
			}
		})
	}

	t.Run("check panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic not thrown")
			}
		}()
		a := NewBitlist64(64)
		b := NewBitlist64(128)
		a.Overlaps(b)
	})
}

func TestBitlist64_Or(t *testing.T) {
	tests := []struct {
		a    *Bitlist64
		b    *Bitlist64
		want *Bitlist64
	}{
		{
			a:    NewBitlist64From([]uint64{0x02}), // 0b00000010
			b:    NewBitlist64From([]uint64{0x03}), // 0b00000011
			want: NewBitlist64From([]uint64{0x03}), // 0b00000011
		},
		{
			a:    NewBitlist64From([]uint64{0x03}), // 0b00000011
			b:    NewBitlist64From([]uint64{0x03}), // 0b00000011
			want: NewBitlist64From([]uint64{0x03}), // 0b00000011
		},
		{
			a:    NewBitlist64From([]uint64{0x13}), // 0b00010011
			b:    NewBitlist64From([]uint64{0x15}), // 0b00010101
			want: NewBitlist64From([]uint64{0x17}), // 0b00010111
		},
		{
			a:    NewBitlist64From([]uint64{0x1F}), // 0b00011111
			b:    NewBitlist64From([]uint64{0x13}), // 0b00010011
			want: NewBitlist64From([]uint64{0x1F}), // 0b00011111
		},
		{
			a:    NewBitlist64From([]uint64{0x1F, 0x03}), // 0b00011111, 0b00000011
			b:    NewBitlist64From([]uint64{0x13, 0x02}), // 0b00010011, 0b00000010
			want: NewBitlist64From([]uint64{0x1F, 0x03}), // 0b00011111, 0b00000011
		},
		{
			a:    NewBitlist64From([]uint64{0x1F, 0x01}), // 0b00011111, 0b00000001
			b:    NewBitlist64From([]uint64{0x93, 0x01}), // 0b10010011, 0b00000001
			want: NewBitlist64From([]uint64{0x9F, 0x01}), // 0b00011111, 0b00000001
		},
		{
			a:    NewBitlist64From([]uint64{0xFF, 0x02}), // 0b11111111, 0x00000010
			b:    NewBitlist64From([]uint64{0x13, 0x03}), // 0b00010011, 0x00000011
			want: NewBitlist64From([]uint64{0xFF, 0x03}), // 0b11111111, 0x00000011
		},
		{
			a:    NewBitlist64From([]uint64{0xFF, 0x85}), // 0b11111111, 0x10000111
			b:    NewBitlist64From([]uint64{0x13, 0x8F}), // 0b00010011, 0x10001111
			want: NewBitlist64From([]uint64{0xFF, 0x8F}), // 0b11111111, 0x10001111
		},
	}

	t.Run("Or()", func(t *testing.T) {
		for _, tt := range tests {
			if !reflect.DeepEqual(tt.a.Or(tt.b).data, tt.want.data) {
				t.Errorf("(%+v).Or(%+v) = %+v, wanted %x", tt.a, tt.b, tt.a.Or(tt.b), tt.want)
			}
		}
	})
	t.Run("NoAllocOr()", func(t *testing.T) {
		for _, tt := range tests {
			res := tt.a.Clone()
			// Make sure that no existing bits set interfere with operation. This is done to simulate
			// the case when res variable is already populated from the previous run.
			for i := uint64(0); i < res.Len(); i += 10 {
				res.SetBitAt(i, true)
			}
			tt.a.NoAllocOr(tt.b, res)
			if !reflect.DeepEqual(res.data, tt.want.data) {
				t.Errorf("(%+v).NoAllocOr(%+v) = %+v, wanted %x", tt.a, tt.b, res.data, tt.want)
			}
		}
	})
	t.Run("OrCount()", func(t *testing.T) {
		for _, tt := range tests {
			if tt.a.OrCount(tt.b) != tt.want.Count() {
				t.Errorf("(%+v).OrCount(%+v) = %d, wanted %d", tt.a, tt.b, tt.a.OrCount(tt.b), tt.want.Count())
			}
		}
	})
	t.Run("check panics", func(t *testing.T) {
		t.Run("Or()", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic not thrown")
				}
			}()
			a := NewBitlist64(64)
			b := NewBitlist64(128)
			a.Or(b)
		})
		t.Run("NoAllocOr()", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic not thrown")
				}
			}()
			a := NewBitlist64(64)
			b := NewBitlist64(128)
			ret := NewBitlist64(64)
			a.NoAllocOr(b, ret)
		})
		t.Run("NoAllocOr() wrong length of result param", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic not thrown")
				}
			}()
			a := NewBitlist64(64)
			b := NewBitlist64(64)
			ret := NewBitlist64(128)
			a.NoAllocOr(b, ret)
		})
		t.Run("OrCount()", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic not thrown")
				}
			}()
			a := NewBitlist64(64)
			b := NewBitlist64(128)
			a.OrCount(b)
		})
	})
}

func TestBitlist64_And(t *testing.T) {
	tests := []struct {
		a    *Bitlist64
		b    *Bitlist64
		want *Bitlist64
	}{
		{
			a:    NewBitlist64From([]uint64{0x02}), // 0b00000010
			b:    NewBitlist64From([]uint64{0x03}), // 0b00000011
			want: NewBitlist64From([]uint64{0x02}), // 0b00000010
		},
		{
			a:    NewBitlist64From([]uint64{0x03}), // 0b00000011
			b:    NewBitlist64From([]uint64{0x03}), // 0b00000011
			want: NewBitlist64From([]uint64{0x03}), // 0b00000011
		},
		{
			a:    NewBitlist64From([]uint64{0x13}), // 0b00010011
			b:    NewBitlist64From([]uint64{0x15}), // 0b00010101
			want: NewBitlist64From([]uint64{0x11}), // 0b00010001
		},
		{
			a:    NewBitlist64From([]uint64{0x1F}), // 0b00011111
			b:    NewBitlist64From([]uint64{0x13}), // 0b00010011
			want: NewBitlist64From([]uint64{0x13}), // 0b00010011
		},
		{
			a:    NewBitlist64From([]uint64{0x1F, 0x03}), // 0b00011111, 0b00000011
			b:    NewBitlist64From([]uint64{0x13, 0x02}), // 0b00010011, 0b00000010
			want: NewBitlist64From([]uint64{0x13, 0x02}), // 0b00010011, 0b00000010
		},
		{
			a:    NewBitlist64From([]uint64{0x9F, 0x01}), // 0b10011111, 0b00000001
			b:    NewBitlist64From([]uint64{0x93, 0x01}), // 0b10010011, 0b00000001
			want: NewBitlist64From([]uint64{0x93, 0x01}), // 0b10010011, 0b00000001
		},
		{
			a:    NewBitlist64From([]uint64{0xFF, 0x02}), // 0b11111111, 0x00000010
			b:    NewBitlist64From([]uint64{0x13, 0x03}), // 0b00010011, 0x00000011
			want: NewBitlist64From([]uint64{0x13, 0x02}), // 0b00010011, 0x00000010
		},
		{
			a:    NewBitlist64From([]uint64{0xFF, 0x87}), // 0b11111111, 0x10000111
			b:    NewBitlist64From([]uint64{0x13, 0x8F}), // 0b00010011, 0x10001111
			want: NewBitlist64From([]uint64{0x13, 0x87}), // 0b00010011, 0x10000111
		},
	}

	t.Run("And()", func(t *testing.T) {
		for _, tt := range tests {
			if !reflect.DeepEqual(tt.a.And(tt.b).data, tt.want.data) {
				t.Errorf("(%+v).And(%+v) = %+v, wanted %x", tt.a, tt.b, tt.a.And(tt.b), tt.want)
			}
		}
	})
	t.Run("NoAllocAnd()", func(t *testing.T) {
		for _, tt := range tests {
			res := tt.a.Clone()
			// Make sure that no existing bits set interfere with operation. This is done to simulate
			// the case when res variable is already populated from the previous run.
			for i := uint64(0); i < res.Len(); i += 10 {
				res.SetBitAt(i, true)
			}
			tt.a.NoAllocAnd(tt.b, res)
			if !reflect.DeepEqual(res.data, tt.want.data) {
				t.Errorf("(%+v).NoAllocAnd(%+v) = %+v, wanted %x", tt.a, tt.b, res.data, tt.want)
			}
		}
	})
	t.Run("AndCount()", func(t *testing.T) {
		for _, tt := range tests {
			if tt.a.AndCount(tt.b) != tt.want.Count() {
				t.Errorf("(%+v).AndCount(%+v) = %d, wanted %d", tt.a, tt.b, tt.a.AndCount(tt.b), tt.want.Count())
			}
		}
	})
	t.Run("check panics", func(t *testing.T) {
		t.Run("And()", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic not thrown")
				}
			}()
			a := NewBitlist64(64)
			b := NewBitlist64(128)
			a.And(b)
		})
		t.Run("NoAllocAnd()", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic not thrown")
				}
			}()
			a := NewBitlist64(64)
			b := NewBitlist64(128)
			ret := NewBitlist64(64)
			a.NoAllocAnd(b, ret)
		})
		t.Run("NoAllocAnd() wrong length of result param", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic not thrown")
				}
			}()
			a := NewBitlist64(64)
			b := NewBitlist64(64)
			ret := NewBitlist64(128)
			a.NoAllocAnd(b, ret)
		})
		t.Run("AndCount()", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic not thrown")
				}
			}()
			a := NewBitlist64(64)
			b := NewBitlist64(128)
			a.AndCount(b)
		})
	})
}

func TestBitlist64_Xor(t *testing.T) {
	tests := []struct {
		a    *Bitlist64
		b    *Bitlist64
		want *Bitlist64
	}{
		{
			a:    NewBitlist64From([]uint64{0x02}), // 0b00000010
			b:    NewBitlist64From([]uint64{0x03}), // 0b00000011
			want: NewBitlist64From([]uint64{0x01}), // 0b00000001
		},
		{
			a:    NewBitlist64From([]uint64{0x03}), // 0b00000011
			b:    NewBitlist64From([]uint64{0x03}), // 0b00000011
			want: NewBitlist64From([]uint64{0x00}), // 0b00000000
		},
		{
			a:    NewBitlist64From([]uint64{0x13}), // 0b00010011
			b:    NewBitlist64From([]uint64{0x15}), // 0b00010101
			want: NewBitlist64From([]uint64{0x06}), // 0b00000110
		},
		{
			a:    NewBitlist64From([]uint64{0x33}), // 0b00110011
			b:    NewBitlist64From([]uint64{0x15}), // 0b00010101
			want: NewBitlist64From([]uint64{0x26}), // 0b00100110
		},
		{
			a:    NewBitlist64From([]uint64{0x1F}), // 0b00011111
			b:    NewBitlist64From([]uint64{0x13}), // 0b00010011
			want: NewBitlist64From([]uint64{0x0c}), // 0b00001100
		},
		{
			a:    NewBitlist64From([]uint64{0x1F, 0x03}), // 0b00011111, 0b00000011
			b:    NewBitlist64From([]uint64{0x13, 0x02}), // 0b00010011, 0b00000010
			want: NewBitlist64From([]uint64{0x0c, 0x01}), // 0b00001100, 0b00000001
		},
		{
			a:    NewBitlist64From([]uint64{0x9F, 0x01}), // 0b10011111, 0b00000001
			b:    NewBitlist64From([]uint64{0x93, 0x01}), // 0b10010011, 0b00000001
			want: NewBitlist64From([]uint64{0x0c, 0x00}), // 0b00001100, 0b00000000
		},
		{
			a:    NewBitlist64From([]uint64{0xFF, 0x02}), // 0b11111111, 0x00000010
			b:    NewBitlist64From([]uint64{0x13, 0x03}), // 0b00010011, 0x00000011
			want: NewBitlist64From([]uint64{0xec, 0x01}), // 0b11101100, 0x00000001
		},
		{
			a:    NewBitlist64From([]uint64{0xFF, 0x87}), // 0b11111111, 0x10000111
			b:    NewBitlist64From([]uint64{0x13, 0x8F}), // 0b00010011, 0x10001111
			want: NewBitlist64From([]uint64{0xec, 0x08}), // 0b11101100, 0x00001000
		},
	}

	t.Run("Xor()", func(t *testing.T) {
		for _, tt := range tests {
			if !reflect.DeepEqual(tt.a.Xor(tt.b).data, tt.want.data) {
				t.Errorf("(%+v).Xor(%+v) = %+v, wanted %x", tt.a, tt.b, tt.a.Xor(tt.b), tt.want)
			}
		}
	})
	t.Run("NoAllocXor()", func(t *testing.T) {
		for _, tt := range tests {
			res := tt.a.Clone()
			// Make sure that no existing bits set interfere with operation. This is done to simulate
			// the case when res variable is already populated from the previous run.
			for i := uint64(0); i < res.Len(); i += 10 {
				res.SetBitAt(i, true)
			}
			tt.a.NoAllocXor(tt.b, res)
			if !reflect.DeepEqual(res.data, tt.want.data) {
				t.Errorf("(%+v).NoAllocXor(%+v) = %+v, wanted %x", tt.a, tt.b, res.data, tt.want)
			}
		}
	})
	t.Run("check panics", func(t *testing.T) {
		t.Run("Xor()", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic not thrown")
				}
			}()
			a := NewBitlist64(64)
			b := NewBitlist64(128)
			a.Xor(b)
		})
		t.Run("NoAllocXor()", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic not thrown")
				}
			}()
			a := NewBitlist64(64)
			b := NewBitlist64(128)
			ret := NewBitlist64(64)
			a.NoAllocXor(b, ret)
		})
		t.Run("NoAllocXor() wrong length of result param", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic not thrown")
				}
			}()
			a := NewBitlist64(64)
			b := NewBitlist64(64)
			ret := NewBitlist64(128)
			a.NoAllocXor(b, ret)
		})
	})
}

func TestBitlist64_Not(t *testing.T) {
	tests := []struct {
		a    *Bitlist64
		want *Bitlist64
	}{
		{
			a:    NewBitlist64From([]uint64{}), // zero-length bitlist
			want: NewBitlist64From([]uint64{}),
		},
		{
			a:    NewBitlist64From([]uint64{0x01}),               // 0b00000001
			want: NewBitlist64From([]uint64{0xFFFFFFFFFFFFFFFE}), // 0b11111110
		},
		{
			a:    NewBitlist64From([]uint64{0x02}),               // 0b00000010
			want: NewBitlist64From([]uint64{0xFFFFFFFFFFFFFFFD}), // 0b11111101
		},
		{
			a:    NewBitlist64From([]uint64{0x03}),               // 0b00000011
			want: NewBitlist64From([]uint64{0xFFFFFFFFFFFFFFFC}), // 0b11111100
		},
		{
			a:    NewBitlist64From([]uint64{0x05}),               // 0b00000101
			want: NewBitlist64From([]uint64{0xFFFFFFFFFFFFFFFA}), // 0b11111010
		},
		{
			a:    NewBitlist64From([]uint64{0x06}),               // 0b00000110
			want: NewBitlist64From([]uint64{0xFFFFFFFFFFFFFFF9}), // 0b11111001
		},
		{
			a:    NewBitlist64From([]uint64{0x83}),               // 0b10000011
			want: NewBitlist64From([]uint64{0xFFFFFFFFFFFFFF7C}), // 0b01111100
		},
		{
			a:    NewBitlist64From([]uint64{0x13}),               // 0b00010011
			want: NewBitlist64From([]uint64{0xFFFFFFFFFFFFFFEC}), // 0b11101100
		},
		{
			a:    NewBitlist64From([]uint64{0x1F}),               // 0b00011111
			want: NewBitlist64From([]uint64{0xFFFFFFFFFFFFFFE0}), // 0b11100000
		},
		{
			a:    NewBitlist64From([]uint64{0x1F, 0x03}),                             // 0b00011111, 0b00000011
			want: NewBitlist64From([]uint64{0xFFFFFFFFFFFFFFE0, 0xFFFFFFFFFFFFFFFC}), // 0b11100000, 0b11111100
		},
		{
			a:    NewBitlist64From([]uint64{0x9F, 0x01}),                             // 0b10011111, 0b00000001
			want: NewBitlist64From([]uint64{0xFFFFFFFFFFFFFF60, 0xFFFFFFFFFFFFFFFE}), // 0b01100000, 0b11111110
		},
		{
			a:    NewBitlist64From([]uint64{allBitsSet, 0x02}),         // 0b11111111, 0x00000010
			want: NewBitlist64From([]uint64{0x00, 0xFFFFFFFFFFFFFFFD}), // 0b00000000, 0x11111101
		},
		{
			a:    NewBitlist64From([]uint64{allBitsSet, 0x87}),         // 0b11111111, 0x10000111
			want: NewBitlist64From([]uint64{0x00, 0xFFFFFFFFFFFFFF78}), // 0b00000000, 0x01111000
		},
		{
			a:    NewBitlist64From([]uint64{allBitsSet, 0x07}),         // 0b11111111, 0x00000111
			want: NewBitlist64From([]uint64{0x00, 0xFFFFFFFFFFFFFFF8}), // 0b00000000, 0x11111000
		},
	}

	t.Run("Not()", func(t *testing.T) {
		for _, tt := range tests {
			if !reflect.DeepEqual(tt.a.Not().data, tt.want.data) {
				t.Errorf("(%+v).Not() = %x, wanted %x", tt.a, tt.a.Not().data, tt.want)
			}
		}
	})
	t.Run("NoAllocNot()", func(t *testing.T) {
		for _, tt := range tests {
			res := tt.a.Clone()
			// Make sure that no existing bits set interfere with operation. This is done to simulate
			// the case when res variable is already populated from the previous run.
			for i := uint64(0); i < res.Len(); i += 10 {
				res.SetBitAt(i, true)
			}
			tt.a.NoAllocNot(res)
			if !reflect.DeepEqual(res.data, tt.want.data) {
				t.Errorf("(%+v).NoAllocNot() = %+v, wanted %x", tt.a, res.data, tt.want)
			}
		}
	})
}

func TestBitlist64_BitIndices(t *testing.T) {
	tests := []struct {
		a    *Bitlist64
		want []int
	}{
		{
			a:    NewBitlist64From([]uint64{}),
			want: []int{},
		},
		{
			a:    NewBitlist64From([]uint64{0b10010}),
			want: []int{1, 4},
		},
		{
			a:    NewBitlist64From([]uint64{0b10000}),
			want: []int{4},
		},
		{
			a:    NewBitlist64From([]uint64{0b10, 0b1}),
			want: []int{1, int(wordSize)},
		},
		{
			a: NewBitlist64From([]uint64{0x10, 0x01, 0xF0, 0xE0}),
			want: []int{
				4,
				int(wordSize),
				int(wordSize)*2 + 4, int(wordSize)*2 + 5, int(wordSize)*2 + 6, int(wordSize)*2 + 7,
				int(wordSize)*3 + 5, int(wordSize)*3 + 6, int(wordSize)*3 + 7,
			},
		},
		{
			a:    NewBitlist64From([]uint64{0b11111111, 0b0}),
			want: []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			a:    NewBitlist64From([]uint64{0b11111111, 0b1}),
			want: []int{0, 1, 2, 3, 4, 5, 6, 7, int(wordSize)},
		},
	}

	for _, tt := range tests {
		got := tt.a.BitIndices()
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("(%0.8b).BitIndices() = %v, wanted %v", tt.a, got, tt.want)
		}
	}
}

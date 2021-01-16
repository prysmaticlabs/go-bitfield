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
			from: []uint64{0xffffffffffffffff},
			want: &Bitlist{size: 64, data: []uint64{0xffffffffffffffff}},
		},
		{
			from: []uint64{0x0000000000000000, 0x0000000000000000},
			want: &Bitlist{size: 128, data: []uint64{0x00, 0x00}},
		},
		{
			from: []uint64{0xffffffffffffffff, 0xffffffffffffffff},
			want: &Bitlist{size: 128, data: []uint64{0xffffffffffffffff, 0xffffffffffffffff}},
		},
		{
			from: []uint64{0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000},
			want: &Bitlist{size: 256, data: []uint64{0x00, 0x00, 0x00, 0x00}},
		},
		{
			from: []uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff},
			want: &Bitlist{
				size: 256,
				data: []uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff},
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
				0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
				0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
			},
			want: &Bitlist{
				size: 512,
				data: []uint64{
					0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
					0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
				},
			},
		},
		{
			from: []uint64{
				0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
				0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
				0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
				0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
				0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
				0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
				0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
				0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
			},
			want: &Bitlist{
				size: 2048,
				data: []uint64{
					0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
					0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
					0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
					0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
					0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
					0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
					0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
					0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc, 0x1111ffffffffcccc,
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

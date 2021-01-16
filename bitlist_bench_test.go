package bitfield

import (
	"fmt"
	"testing"
)

func BenchmarkBitlist_New(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 256 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte new", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					NewByteBitlist(n)
				}
			})
			b.Run("[]uint64 new", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					NewBitlist(n)
				}
			})
			b.Run("[]uint64 new+from", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					NewBitlistFrom(NewBitlist(n).data)
				}
			})
		})
	}
}

func BenchmarkBitlist_Len(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 1024 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
				b.StopTimer()
				set := NewByteBitlist(n)
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					set.Len()
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				set := NewBitlist(n)
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					set.Len()
				}
			})
		})
	}
}

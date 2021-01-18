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
				s := NewByteBitlist(n)
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Len()
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist(n)
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Len()
				}
			})
		})
	}
}

func BenchmarkBitlist_SetBitAt(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 1024 {
		idx := n / 2
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
				b.StopTimer()
				s := NewByteBitlist(n)
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.BitAt(idx)
					s.SetBitAt(idx, true)
					s.SetBitAt(idx, false)
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist(n)
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.BitAt(idx)
					s.SetBitAt(idx, true)
					s.SetBitAt(idx, false)
				}
			})
		})
	}
}

func BenchmarkBitlist_Count(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 512 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
				b.StopTimer()
				s := NewByteBitlist(n)
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Count()
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist(n)
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Count()
				}
			})
		})
	}
}

func BenchmarkBitlist_Bytes(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 512 {
		b.Run("bitlist non empty", func(b *testing.B) {
			b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
				b.Run("[]byte", func(b *testing.B) {
					b.StopTimer()
					s := NewByteBitlist(n)
					for i := uint64(0); i < n; i += 10 {
						s.SetBitAt(i, true)
					}
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.Bytes()
					}
				})
				b.Run("[]uint64", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist(n)
					for i := uint64(0); i < n; i += 10 {
						s.SetBitAt(i, true)
					}
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.Bytes()
					}
				})
			})
		})
		b.Run("up to half bitlist non empty", func(b *testing.B) {
			b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
				b.Run("[]byte", func(b *testing.B) {
					b.StopTimer()
					s := NewByteBitlist(n)
					for i := uint64(0); i < n/2; i += 10 {
						s.SetBitAt(i, true)
					}
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.Bytes()
					}
				})
				b.Run("[]uint64", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist(n)
					for i := uint64(0); i < n/2; i += 10 {
						s.SetBitAt(i, true)
					}
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.Bytes()
					}
				})
			})
		})
		b.Run("only single bit set", func(b *testing.B) {
			b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
				b.Run("[]byte", func(b *testing.B) {
					b.StopTimer()
					s := NewByteBitlist(n)
					s.SetBitAt(n, true)
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.Bytes()
					}
				})
				b.Run("[]uint64", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist(n)
					s.SetBitAt(n, true)
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.Bytes()
					}
				})
			})
		})
	}
}

func BenchmarkBitlist_Contains(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 512 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
				b.StopTimer()
				s := NewByteBitlist(n)
				s1 := NewByteBitlist(n) // subset
				s2 := NewByteBitlist(n) // not a subset
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i, true)
				}
				s2.SetBitAt(1, true)
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Contains(s1)
					s.Contains(s2)
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist(n)
				s1 := NewBitlist(n) // subset
				s2 := NewBitlist(n) // not a subset
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i, true)
				}
				s2.SetBitAt(1, true)
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Contains(s1)
					s.Contains(s2)
				}
			})
		})
	}
}

func BenchmarkBitlist_Overlaps(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 512 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
				b.StopTimer()
				s := NewByteBitlist(n)
				s1 := NewByteBitlist(n) // has overlaps
				s2 := NewByteBitlist(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Overlaps(s1)
					s.Overlaps(s2)
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist(n)
				s1 := NewBitlist(n) // has overlaps
				s2 := NewBitlist(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Overlaps(s1)
					s.Overlaps(s2)
				}
			})
		})
	}
}

func BenchmarkBitlist_Or(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 256 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
				b.StopTimer()
				s := NewByteBitlist(n)
				s1 := NewByteBitlist(n) // has overlaps
				s2 := NewByteBitlist(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Or(s1)
					s.Or(s2)
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist(n)
				s1 := NewBitlist(n) // has overlaps
				s2 := NewBitlist(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Or(s1)
					s.Or(s2)
				}
			})
			b.Run("[]uint64 (noalloc)", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist(n)
				s1 := NewBitlist(n) // has overlaps
				s2 := NewBitlist(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				result := s.Clone()
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.NoAllocOr(s1, result)
					s.NoAllocOr(s2, result)
				}
			})
		})
	}
}

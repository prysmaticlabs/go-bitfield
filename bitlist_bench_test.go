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
					NewBitlist(n)
				}
			})
			b.Run("[]uint64 new", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					NewBitlist64(n)
				}
			})
			b.Run("[]uint64 new+from", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					NewBitlist64From(NewBitlist64(n).data)
				}
			})
		})
	}
}

func BenchmarkBitlist_Len(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 512 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist(n)
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Len()
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Len()
				}
			})
		})
	}
}

func BenchmarkBitlist_SetBitAt(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 512 {
		idx := n / 2
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist(n)
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.BitAt(idx)
					s.SetBitAt(idx, true)
					s.SetBitAt(idx, false)
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
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
				s := NewBitlist(n)
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
				s := NewBitlist64(n)
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
					s := NewBitlist(n)
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
					s := NewBitlist64(n)
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
					s := NewBitlist(n)
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
					s := NewBitlist64(n)
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
					s := NewBitlist(n)
					s.SetBitAt(n, true)
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.Bytes()
					}
				})
				b.Run("[]uint64", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist64(n)
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
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // subset
				s2 := NewBitlist64(n) // not a subset
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
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
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
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
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
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
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

func BenchmarkBitlist_OrCount(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 256 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
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
					a, _ := s.Or(s1)
					a.Count()
					b, _ := s.Or(s2)
					b.Count()
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					a, _ := s.Or(s1)
					a.Count()
					b, _ := s.Or(s2)
					b.Count()
				}
			})
			b.Run("[]uint64 (noalloc)", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				result := s.Clone()
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.NoAllocOr(s1, result)
					result.Count()
					s.NoAllocOr(s2, result)
					result.Count()
				}
			})
			b.Run("[]uint64 (OrCount)", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.OrCount(s1)
					s.OrCount(s2)
				}
			})
		})
	}
}

func BenchmarkBitlist_And(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 256 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
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
					s.And(s1)
					s.And(s2)
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.And(s1)
					s.And(s2)
				}
			})
			b.Run("[]uint64 (noalloc)", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				result := s.Clone()
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.NoAllocAnd(s1, result)
					s.NoAllocAnd(s2, result)
				}
			})
		})
	}
}

func BenchmarkBitlist_AndCount(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 256 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
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
					a, _ := s.And(s1)
					a.Count()
					b, _ := s.And(s2)
					b.Count()
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					a, _ := s.And(s1)
					a.Count()
					b, _ := s.And(s2)
					b.Count()
				}
			})
			b.Run("[]uint64 (noalloc)", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				result := s.Clone()
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.NoAllocAnd(s1, result)
					result.Count()
					s.NoAllocAnd(s2, result)
					result.Count()
				}
			})
			b.Run("[]uint64 (AndCount)", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.AndCount(s1)
					s.AndCount(s2)
				}
			})
		})
	}
}

func BenchmarkBitlist_Xor(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 256 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
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
					s.Xor(s1)
					s.Xor(s2)
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Xor(s1)
					s.Xor(s2)
				}
			})
			b.Run("[]uint64 (noalloc)", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				result := s.Clone()
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.NoAllocXor(s1, result)
					s.NoAllocXor(s2, result)
				}
			})
		})
	}
}

func BenchmarkBitlist_XorCount(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 256 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
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
					a, _ := s.Xor(s1)
					a.Count()
					b, _ := s.Xor(s2)
					b.Count()
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					a, _ := s.Xor(s1)
					a.Count()
					b, _ := s.Xor(s2)
					b.Count()
				}
			})
			b.Run("[]uint64 (noalloc)", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				result := s.Clone()
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.NoAllocXor(s1, result)
					result.Count()
					s.NoAllocXor(s2, result)
					result.Count()
				}
			})
			b.Run("[]uint64 (XorCount)", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				s1 := NewBitlist64(n) // has overlaps
				s2 := NewBitlist64(n) // no overlaps
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
					s1.SetBitAt(i, true)
					s2.SetBitAt(i+1, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.XorCount(s1)
					s.XorCount(s2)
				}
			})
		})
	}
}

func BenchmarkBitlist_Not(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 256 {
		b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
			b.Run("[]byte", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist(n)
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Not()
				}
			})
			b.Run("[]uint64", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
				}
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.Not()
				}
			})
			b.Run("[]uint64 (noalloc)", func(b *testing.B) {
				b.StopTimer()
				s := NewBitlist64(n)
				for i := uint64(0); i < n; i += 100 {
					s.SetBitAt(i, true)
				}
				result := s.Clone()
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					s.NoAllocNot(result)
				}
			})
		})
	}
}

func BenchmarkBitlist_BitIndices(b *testing.B) {
	for n := uint64(0); n <= 2048; n += 512 {
		b.Run("bitlist non empty", func(b *testing.B) {
			b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
				b.Run("[]byte", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist(n)
					for i := uint64(0); i < n; i += 10 {
						s.SetBitAt(i, true)
					}
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.BitIndices()
					}
				})
				b.Run("[]uint64", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist64(n)
					for i := uint64(0); i < n; i += 10 {
						s.SetBitAt(i, true)
					}
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.BitIndices()
					}
				})
				b.Run("[]uint64 (noalloc)", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist64(n)
					for i := uint64(0); i < n; i += 10 {
						s.SetBitAt(i, true)
					}
					indices := make([]int, s.Count())
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.NoAllocBitIndices(indices)
					}
				})
			})
		})
		b.Run("up to half bitlist non empty", func(b *testing.B) {
			b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
				b.Run("[]byte", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist(n)
					for i := uint64(0); i < n/2; i += 10 {
						s.SetBitAt(i, true)
					}
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.BitIndices()
					}
				})
				b.Run("[]uint64", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist64(n)
					for i := uint64(0); i < n/2; i += 10 {
						s.SetBitAt(i, true)
					}
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.BitIndices()
					}
				})
				b.Run("[]uint64 (noalloc)", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist64(n)
					for i := uint64(0); i < n/2; i += 10 {
						s.SetBitAt(i, true)
					}
					indices := make([]int, s.Count())
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.NoAllocBitIndices(indices)
					}
				})
			})
		})
		b.Run("only single bit set", func(b *testing.B) {
			b.Run(fmt.Sprintf("size:%d", n), func(b *testing.B) {
				b.Run("[]byte", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist(n)
					s.SetBitAt(n, true)
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.BitIndices()
					}
				})
				b.Run("[]uint64", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist64(n)
					s.SetBitAt(n, true)
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.BitIndices()
					}
				})
				b.Run("[]uint64 (noalloc)", func(b *testing.B) {
					b.StopTimer()
					s := NewBitlist64(n)
					s.SetBitAt(n, true)
					indices := make([]int, s.Count())
					b.StartTimer()
					for i := 0; i < b.N; i++ {
						s.NoAllocBitIndices(indices)
					}
				})
			})
		})
	}
}

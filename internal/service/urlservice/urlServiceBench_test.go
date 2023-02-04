package urlservice

import "testing"

func BenchmarkReducing(b *testing.B) {
	url := "example"
	b.Run("reducing", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Reducing(url)
		}
	})
}

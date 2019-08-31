// Sample benchmarks to test the performance of search.
package search_test

import (
	"fly/search"
	"testing"
)

// BenchmarkSearchRun provides performance numbers for the search.Run function.
func BenchmarkSearchRun(b *testing.B) {
	b.ResetTimer()
	q := search.Query{}
	for i := 0; i < b.N; i++ {
		search.Run(&q)
	}
}

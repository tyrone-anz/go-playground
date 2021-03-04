package awesome

import "testing"

func benchmarkSumIt(a1, a2 int, b *testing.B) {
    for n := 0; n < b.N; n++ {
        SumIt(a1, a2)
    }
}

func BenchmarkSumIt1(b *testing.B)  {
    // b.N = 100
    benchmarkSumIt(1, 2, b)
}
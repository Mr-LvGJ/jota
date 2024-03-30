package id

import (
	"testing"
)

func Benchmark_Snowflake1(b *testing.B) {
	sf := NewSnowflake1(WithWorkerID(1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sf.NextID()
	}
}
func Benchmark_Snowflake2(b *testing.B) {
	sf := NewSnowflake2(WithWorkerID2(1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sf.NextID()
	}
}

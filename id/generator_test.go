package id

import (
	"testing"
)

func TestSnowflake1_NextID(t *testing.T) {
	sf := NewSnowflake1(WithWorkerID(100))
	for i := 0; i < 100; i++ {
		t.Log(sf.NextID())
	}
}

func TestSnowflake2_NextID(t *testing.T) {
	sf := NewSnowflake2(WithWorkerID2(1))
	for i := 0; i < 100; i++ {
		t.Log(sf.NextID())
	}
}

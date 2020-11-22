package utils

import "testing"

func setupTestCase(t *testing.T) func(t *testing.T) {
	count = 0
	return func(t *testing.T) {
		count = 0
	}
}

func Test_increaseCount(t *testing.T) {
	t.Run("first", func(t *testing.T) {
		teardownTestCase := setupTestCase(t) // setup before test
		defer teardownTestCase(t)
		want := uint32(1)
		if got := increaseCount(); got != want {
			t.Errorf("increaseCount() = %v, want %v", got, want)
		}
	})
	t.Run("increase", func(t *testing.T) {
		teardownTestCase := setupTestCase(t) // setup before test
		defer teardownTestCase(t)
		want := uint32(102)
		got := uint32(0)
		for i := 0; i < 102; i++ {
			got = increaseCount()
		}
		if got != want {
			t.Errorf("increaseCount() = %v, want %v", got, want)
		}
	})

	t.Run("reset", func(t *testing.T) {
		teardownTestCase := setupTestCase(t) // setup before test
		defer teardownTestCase(t)
		want := uint32(0)
		got := uint32(0)
		for i := 0; i < 1000; i++ {
			got = increaseCount()
		}
		if got != want {
			t.Errorf("increaseCount() = %v, want %v", got, want)
		}
	})

	t.Run("roll out", func(t *testing.T) {
		teardownTestCase := setupTestCase(t) // setup before test
		defer teardownTestCase(t)
		want := uint32(200)
		got := uint32(0)
		for i := 0; i < 1200; i++ {
			got = increaseCount()
		}
		if got != want {
			t.Errorf("increaseCount() = %v, want %v", got, want)
		}
	})
}

func Benchmark_increaseCount(b *testing.B) {
	b.Run("random", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			increaseCount()
		}
	})
}

func TestGenerateId(t *testing.T) {
	t.Run("generate id", func(t *testing.T) {
		t.Log(CurrentMillis())
		t.Log(GenerateId())
		t.Log(GenerateId())
		t.Log(GenerateId())
		t.Log(GenerateId())
	})
}

package store

import (
	"strconv"
	"testing"
	"time"
)

func Test_aggr(t *testing.T) {
	t.Run("wait max millis", func(t *testing.T) {
		ch := make(chan string, 10)
		count := 3

		go func() {
			start := time.Now()
			items := aggr(ch, 5, 20)
			duration := time.Since(start)

			if duration > 25*time.Millisecond*25 {
				t.Errorf("duration: %v\n", duration)
			}
			if len(items) != count {
				t.Errorf("items size: %d\n", len(items))
			}

		}()

		for i := 0; i < count; i++ {
			item := strconv.Itoa(i)
			ch <- item
		}
		time.Sleep(time.Millisecond * 50)
	})

	t.Run("encounter max count", func(t *testing.T) {
		ch := make(chan string, 10)
		maxCount := 5

		go func() {
			items := aggr(ch, maxCount, 200)
			if len(items) != maxCount {
				t.Errorf("items size: %d\n", len(items))
			}
		}()

		for i := 0; i < maxCount+4; i++ {
			item := strconv.Itoa(i)
			ch <- item
		}
		time.Sleep(time.Millisecond * 50)
	})
}

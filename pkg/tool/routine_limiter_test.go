package tool_test

import (
	"sync/atomic"
	"testing"
	"time"

	"git.papkovda.ru/library/gokit/pkg/tool"
)

func TestRoutineLimiter(t *testing.T) {
	t.Run("Respects max concurrent routines", func(t *testing.T) {
		maxRoutines := 3
		limiter := tool.NewRoutineLimiter(maxRoutines)

		var activeRoutines int32
		var maxObservedRoutines int32

		for i := 0; i < 10; i++ {
			limiter.Run(func() {
				current := atomic.AddInt32(&activeRoutines, 1)
				if current > maxObservedRoutines {
					atomic.StoreInt32(&maxObservedRoutines, current)
				}
				time.Sleep(10 * time.Millisecond)
				atomic.AddInt32(&activeRoutines, -1)
			})
		}

		limiter.Wait()

		if maxObservedRoutines != int32(maxRoutines) {
			t.Errorf("Expected max %d concurrent routines, but observed %d", maxRoutines, maxObservedRoutines)
		}
	})

	t.Run("Waits for all routines to complete", func(t *testing.T) {
		limiter := tool.NewRoutineLimiter(5)
		var counter int32

		for i := 0; i < 10; i++ {
			limiter.Run(func() {
				time.Sleep(10 * time.Millisecond)
				atomic.AddInt32(&counter, 1)
			})
		}

		limiter.Wait()

		if counter != 10 {
			t.Errorf("Expected all 10 routines to complete, but only %d completed", counter)
		}
	})

	t.Run("Handles zero concurrent routines", func(t *testing.T) {
		limiter := tool.NewRoutineLimiter(0)
		var counter int32

		for i := 0; i < 5; i++ {
			limiter.Run(func() {
				atomic.AddInt32(&counter, 1)
			})
		}

		limiter.Wait()

		if counter != 5 {
			t.Errorf("Expected all 5 routines to complete even with 0 limit, but %d completed", counter)
		}
	})
}

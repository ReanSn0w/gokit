package tool

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimiter_Do(t *testing.T) {
	rl := NewRateLimiter(time.Second)

	// Тест успешного выполнения
	err := rl.Do("test_key", func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Тест ограничения частоты
	err = rl.Do("test_key", func() error {
		return nil
	})
	if err != ErrRateLimited {
		t.Errorf("Expected ErrRateLimited, got %v", err)
	}

	// Тест после истечения времени ограничения
	time.Sleep(time.Second)
	err = rl.Do("test_key", func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error after timeout, got %v", err)
	}
}

func TestRateLimiter_MultipleKeys(t *testing.T) {
	rl := NewRateLimiter(time.Second)

	// Тест с разными ключами
	err1 := rl.Do("key1", func() error {
		return nil
	})
	err2 := rl.Do("key2", func() error {
		return nil
	})

	if err1 != nil || err2 != nil {
		t.Errorf("Expected no errors for different keys, got %v and %v", err1, err2)
	}
}

func TestRateLimiter_Concurrency(t *testing.T) {
	rl := NewRateLimiter(time.Second)
	const concurrency = 10

	var wg sync.WaitGroup
	errChan := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := rl.Do("concurrent_key", func() error {
				return nil
			})
			errChan <- err
		}()
	}

	wg.Wait()
	close(errChan)

	successCount := 0
	for err := range errChan {
		if err == nil {
			successCount++
		} else if err != ErrRateLimited {
			t.Errorf("Unexpected error: %v", err)
		}
	}

	if successCount != 1 {
		t.Errorf("Expected exactly one success in concurrent scenario, got %d", successCount)
	}
}

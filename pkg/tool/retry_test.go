package tool_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ReanSn0w/gokit/pkg/tool"
	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
)

func TestRetry_Do(t *testing.T) {
	log := lgr.New(lgr.Msec, lgr.Debug, lgr.CallerFile, lgr.CallerFunc)

	t.Run("successful on first attempt", func(t *testing.T) {
		retry := tool.NewRetry(log, 3, time.Millisecond)
		counter := 0
		err := retry.Do(func() error {
			counter++
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, counter)
	})

	t.Run("successful after retries", func(t *testing.T) {
		retry := tool.NewRetry(log, 3, time.Millisecond)
		counter := 0
		err := retry.Do(func() error {
			counter++
			if counter < 3 {
				return errors.New("temporary error")
			}
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 3, counter)
	})

	t.Run("failure after max retries", func(t *testing.T) {
		retry := tool.NewRetry(log, 3, time.Millisecond)
		counter := 0
		err := retry.Do(func() error {
			counter++
			return errors.New("persistent error")
		})
		assert.Error(t, err)
		assert.Equal(t, 3, counter)
	})

	t.Run("respects delay", func(t *testing.T) {
		retry := tool.NewRetry(log, 3, 100*time.Millisecond)
		start := time.Now()
		_ = retry.Do(func() error {
			return errors.New("error")
		})
		duration := time.Since(start)
		assert.GreaterOrEqual(t, duration, 300*time.Millisecond)
	})
}

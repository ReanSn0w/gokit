package base

import (
	"encoding/binary"
	"encoding/hex"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── NewID ────────────────────────────────────────────────────────────────────

func TestNewID_Length(t *testing.T) {
	id := NewID()
	assert.Len(t, id, 24, "ID должен состоять ровно из 24 символов")
}

func TestNewID_IsHexString(t *testing.T) {
	id := NewID()
	matched, err := regexp.MatchString(`^[0-9a-f]{24}$`, id)
	require.NoError(t, err)
	assert.True(t, matched, "ID должен содержать только строчные hex-символы: %s", id)
}

func TestNewID_Unique(t *testing.T) {
	const n = 1000
	seen := make(map[string]struct{}, n)

	for _ = range n {
		id := NewID()
		_, exists := seen[id]
		assert.False(t, exists, "обнаружен дубликат ID: %s", id)
		seen[id] = struct{}{}
	}
}

func TestNewID_UniqueUnderConcurrency(t *testing.T) {
	const goroutines = 20
	const perGoroutine = 500

	var (
		mu   sync.Mutex
		seen = make(map[string]struct{}, goroutines*perGoroutine)
		wg   sync.WaitGroup
	)

	wg.Add(goroutines)
	for _ = range goroutines {
		go func() {
			defer wg.Done()
			ids := make([]string, perGoroutine)
			for i := range ids {
				ids[i] = NewID()
			}
			mu.Lock()
			defer mu.Unlock()
			for _, id := range ids {
				_, exists := seen[id]
				assert.False(t, exists, "дубликат ID при конкурентной генерации: %s", id)
				seen[id] = struct{}{}
			}
		}()
	}
	wg.Wait()
}

// ── NewIDFromTimestamp ───────────────────────────────────────────────────────

func TestNewIDFromTimestamp_EncodesTimestamp(t *testing.T) {
	ts := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	id := NewIDFromTimestamp(ts)

	raw, err := hex.DecodeString(id)
	require.NoError(t, err)
	require.Len(t, raw, 12)

	// Первые 4 байта — Unix timestamp в big-endian
	encoded := binary.BigEndian.Uint32(raw[0:4])
	assert.Equal(t, uint32(ts.Unix()), encoded,
		"первые 4 байта ID должны кодировать Unix timestamp")
}

func TestNewIDFromTimestamp_DifferentTimestamps_DifferentPrefix(t *testing.T) {
	t1 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	id1 := NewIDFromTimestamp(t1)
	id2 := NewIDFromTimestamp(t2)

	// Первые 8 символов hex = 4 байта timestamp
	assert.NotEqual(t, id1[:8], id2[:8], "разные timestamps должны давать разные префиксы")
}

func TestNewIDFromTimestamp_SameTimestamp_MonotonicCounter(t *testing.T) {
	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	id1 := NewIDFromTimestamp(ts)
	id2 := NewIDFromTimestamp(ts)

	// Одинаковый timestamp → одинаковый префикс [0:8]
	assert.Equal(t, id1[:8], id2[:8], "одинаковый timestamp → одинаковый префикс")

	// Счётчик (последние 6 символов hex = 3 байта) должен возрастать
	raw1, err1 := hex.DecodeString(id1)
	raw2, err2 := hex.DecodeString(id2)
	require.NoError(t, err1)
	require.NoError(t, err2)

	counter1 := uint32(raw1[9])<<16 | uint32(raw1[10])<<8 | uint32(raw1[11])
	counter2 := uint32(raw2[9])<<16 | uint32(raw2[10])<<8 | uint32(raw2[11])

	assert.Greater(t, counter2, counter1, "счётчик должен монотонно возрастать")
}

func TestNewIDFromTimestamp_SameTimestamp_UniqueIDs(t *testing.T) {
	ts := time.Now()
	id1 := NewIDFromTimestamp(ts)
	id2 := NewIDFromTimestamp(ts)

	assert.NotEqual(t, id1, id2, "два ID с одним timestamp должны быть уникальны за счёт счётчика")
}

func TestNewIDFromTimestamp_Length(t *testing.T) {
	id := NewIDFromTimestamp(time.Now())
	assert.Len(t, id, 24)
}

// ── putUint24 ────────────────────────────────────────────────────────────────

func TestPutUint24_Zero(t *testing.T) {
	b := make([]byte, 3)
	putUint24(b, 0)
	assert.Equal(t, []byte{0, 0, 0}, b)
}

func TestPutUint24_One(t *testing.T) {
	b := make([]byte, 3)
	putUint24(b, 1)
	assert.Equal(t, []byte{0, 0, 1}, b)
}

func TestPutUint24_MaxByte(t *testing.T) {
	b := make([]byte, 3)
	putUint24(b, 0xFF)
	assert.Equal(t, []byte{0, 0, 0xFF}, b)
}

func TestPutUint24_SecondByte(t *testing.T) {
	b := make([]byte, 3)
	putUint24(b, 0x0100)
	assert.Equal(t, []byte{0, 1, 0}, b)
}

func TestPutUint24_Max(t *testing.T) {
	b := make([]byte, 3)
	putUint24(b, 0xFFFFFF)
	assert.Equal(t, []byte{0xFF, 0xFF, 0xFF}, b)
}

func TestPutUint24_BigEndianOrder(t *testing.T) {
	b := make([]byte, 3)
	// 0x010203: старший байт → b[0], младший → b[2]
	putUint24(b, 0x010203)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, b)
}

// ── byteToHex ────────────────────────────────────────────────────────────────

func TestByteToHex_ZeroBytes(t *testing.T) {
	result := byteToHex([12]byte{})
	assert.Equal(t, "000000000000000000000000", result)
}

func TestByteToHex_Length(t *testing.T) {
	result := byteToHex([12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
	assert.Len(t, result, 24)
}

func TestByteToHex_KnownValue(t *testing.T) {
	var b [12]byte
	b[0] = 0xDE
	b[1] = 0xAD
	b[11] = 0xFF

	result := byteToHex(b)

	assert.Equal(t, "dead", result[:4], "первые два байта должны кодироваться как 'dead'")
	assert.Equal(t, "ff", result[22:], "последний байт должен кодироваться как 'ff'")
}

func TestByteToHex_IsLowercase(t *testing.T) {
	b := [12]byte{0xAB, 0xCD, 0xEF}
	result := byteToHex(b)

	matched, err := regexp.MatchString(`^[0-9a-f]+$`, result)
	require.NoError(t, err)
	assert.True(t, matched, "hex должен быть в нижнем регистре: %s", result)
}

func TestByteToHex_RoundTrip(t *testing.T) {
	// hex.DecodeString(byteToHex(b)) должен вернуть исходный массив
	original := [12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	encoded := byteToHex(original)

	decoded, err := hex.DecodeString(encoded)
	require.NoError(t, err)

	assert.Equal(t, original[:], decoded)
}

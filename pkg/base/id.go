package base

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

var (
	// idCounter — глобальный атомарный счётчик, инициализируется случайным
	// значением при старте процесса, чтобы исключить коллизии между перезапусками.
	idCounter = readRandomUint32()

	// processUnique — 5 случайных байт, уникальных для текущего процесса.
	// Вместе с timestamp и счётчиком гарантируют глобальную уникальность ID
	// даже при параллельной работе нескольких экземпляров приложения.
	processUnique = processUniqueBytes()
)

// NewID генерирует новый уникальный 24-символьный hex-идентификатор
// на основе текущего времени. Структура ID (12 байт → 24 hex-символа):
//
//	[0:4]  — Unix-timestamp (big-endian uint32)
//	[4:9]  — уникальные байты процесса (processUnique)
//	[9:12] — монотонно возрастающий счётчик (big-endian uint24)
//
// Формат аналогичен MongoDB ObjectID.
func NewID() string {
	return NewIDFromTimestamp(time.Now())
}

// NewIDFromTimestamp генерирует уникальный 24-символьный hex-идентификатор
// с использованием переданного timestamp вместо текущего времени.
// Удобно для воспроизводимой генерации ID в тестах или при импорте данных.
func NewIDFromTimestamp(timestamp time.Time) string {
	var b [12]byte

	binary.BigEndian.PutUint32(b[0:4], uint32(timestamp.Unix()))
	copy(b[4:9], processUnique[:])
	putUint24(b[9:12], atomic.AddUint32(&idCounter, 1))

	return byteToHex(b)
}

// processUniqueBytes генерирует 5 криптографически случайных байт,
// уникальных для текущего процесса. Вызывается один раз при инициализации пакета.
// Паникует, если crypto/rand недоступен — это критическая ошибка окружения.
func processUniqueBytes() [5]byte {
	var b [5]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %w", err))
	}

	return b
}

// readRandomUint32 читает 4 криптографически случайных байта и собирает из них
// uint32 в little-endian порядке. Используется для инициализации idCounter.
// Паникует, если crypto/rand недоступен — это критическая ошибка окружения.
func readRandomUint32() uint32 {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %w", err))
	}

	return (uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
}

// putUint24 записывает младшие 24 бита значения v в срез b (минимум 3 байта)
// в big-endian порядке: b[0] — старший байт, b[2] — младший.
func putUint24(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

// byteToHex кодирует массив из 12 байт в строку из 24 hex-символов (lowercase).
func byteToHex(b [12]byte) string {
	var buf [24]byte
	hex.Encode(buf[:], b[:])
	return string(buf[:])
}

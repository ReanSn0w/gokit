package base

import "log"

// Logger — интерфейс для логирования с форматированием в стиле fmt.Sprintf.
// Позволяет подменять реализацию в тестах или использовать любой совместимый backend.
type Logger interface {
	Logf(format string, args ...any)
}

// NewStdLogger создаёт StdLogger, оборачивающий стандартный *log.Logger.
func NewStdLogger(l *log.Logger) StdLogger {
	return StdLogger{
		log: l,
	}
}

// StdLogger — реализация Logger поверх *log.Logger из стандартной библиотеки.
type StdLogger struct {
	log *log.Logger
}

// Logf пишет форматированное сообщение через хранимый logger.
func (l *StdLogger) Logf(format string, args ...any) {
	l.log.Printf(format, args...)
}

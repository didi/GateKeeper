package zerolog

import (
	"io"
	"strings"
	"sync"
	"time"
)

// FormatWriter define format method
type FormatWriter interface {
	FormatCaller(string) string
	FormatTimestamp() string
	FormatMessage(string) string
	FormatLevel(string) string
}

// LevelWriter defines as interface a writer may implement in order
// to receive level information with payload.
type LevelWriter interface {
	FormatWriter
	io.Writer
	WriteLevel(level Level, p []byte) (n int, err error)
}

type levelWriterAdapter struct {
	io.Writer
}

func (lw levelWriterAdapter) WriteLevel(l Level, p []byte) (n int, err error) {
	//finalwrite
	return lw.Write(p)
}

func (lw levelWriterAdapter) FormatCaller(i string) string {
	return defaultFormatCaller(i)
}

func (lw levelWriterAdapter) FormatTimestamp() string {
	return defaultFormatTimestamp()
}

func (lw levelWriterAdapter) FormatMessage(i string) string {
	return defaultFormatMessage(i)
}

func (lw levelWriterAdapter) FormatLevel(i string) string {
	return defaultFormatLevel(i)
}

type syncWriter struct {
	mu sync.Mutex
	lw LevelWriter
}

// SyncWriter wraps w so that each call to Write is synchronized with a mutex.
// This syncer can be the call to writer's Write method is not thread safe.
// Note that os.File Write operation is using write() syscall which is supposed
// to be thread-safe on POSIX systems. So there is no need to use this with
// os.File on such systems as zerolog guaranties to issue a single Write call
// per log event.
func SyncWriter(w io.Writer) io.Writer {
	if lw, ok := w.(LevelWriter); ok {
		return &syncWriter{lw: lw}
	}
	return &syncWriter{lw: levelWriterAdapter{w}}
}

// Write implements the io.Writer interface.
func (s *syncWriter) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lw.Write(p)
}

// WriteLevel implements the LevelWriter interface.
func (s *syncWriter) WriteLevel(l Level, p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lw.WriteLevel(l, p)
}

type multiLevelWriter struct {
	writers []LevelWriter
}

func (t multiLevelWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

func (t multiLevelWriter) WriteLevel(l Level, p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.WriteLevel(l, p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

func (t multiLevelWriter) FormatCaller(i string) string {
	return defaultFormatCaller(i)
}

func (t multiLevelWriter) FormatTimestamp() string {
	return defaultFormatTimestamp()
}

func (t multiLevelWriter) FormatMessage(i string) string {
	return defaultFormatMessage(i)
}

func (t multiLevelWriter) FormatLevel(i string) string {
	return defaultFormatLevel(i)
}

// MultiLevelWriter creates a writer that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command. If some writers
// implement LevelWriter, their WriteLevel method will be used instead of Write.
func MultiLevelWriter(writers ...io.Writer) LevelWriter {
	lwriters := make([]LevelWriter, 0, len(writers))
	for _, w := range writers {
		if lw, ok := w.(LevelWriter); ok {
			lwriters = append(lwriters, lw)
		} else {
			lwriters = append(lwriters, levelWriterAdapter{w})
		}
	}
	return multiLevelWriter{lwriters}
}

// ----- DEFAULT FORMATTERS ---------------------------------------------------

var (
	defaultFormatTimestamp = func() string {
		t := time.Now().Format("2006-01-02T15:04:05.000-0700")
		t = "[" + t + "]"
		return t
	}

	defaultFormatLevel = func(i string) string {
		var l string

		switch i {
		case "debug":
			l = "[DEBUG]"
		case "info":
			l = "[INFO]"
		case "warn":
			l = "[WARNING]"
		case "error":
			l = "[ERROR]"
		case "fatal":
			l = "[FATAL]"
		case "panic":
			l = "[PANIC]"
		default:
			l = "[???]"
		}
		return l
	}

	defaultFormatCaller = func(i string) string {
		var c string

		if len(i) > 0 {
			if cwd != "" {
				c = strings.TrimPrefix(i, cwd)
				c = strings.TrimPrefix(c, "/")
			}
			c = "[" + c + "]"
		}
		return c
	}

	defaultFormatMessage = func(i string) string {
		return i
	}
)

package zerolog

import (
	"os"
)

type StdoutWriter struct {
	AccessWriter
}

// NewStdoutWriter will write to stdout.
func NewStdoutWriter(options ...Option) *StdoutWriter {
	w := AccessWriter{}

	for _, opt := range options {
		opt(&w)
	}

	w.FileOut = os.Stdout
	return &StdoutWriter{w}
}

func (w *StdoutWriter) WriteLevel(l Level, p []byte) (n int, err error) {

	return w.Write(p)
}

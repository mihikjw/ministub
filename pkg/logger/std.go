package logger

import (
	"io"
	"log"
	"os"
)

// StdLogger logs to StdOut, StdErr
type StdLogger struct {
	stdOut *log.Logger
	stdErr *log.Logger
}

// NewStdLogger creates a new instance of a logger to StdOut, StdErr
func NewStdLogger() (result *StdLogger) {
	result = new(StdLogger)
	result.stdOut = createOutput(os.Stdout, "INFO: ")
	result.stdErr = createOutput(os.Stderr, "ERROR: ")
	return result
}

// Info logs a message to StdOut
func (l *StdLogger) Info(msg string) {
	l.stdOut.Print(msg)
}

// Error logs a message to StdErr
func (l *StdLogger) Error(msg string) {
	l.stdErr.Print(msg)
}

// createOutput creates a logger output suitable for writing too
func createOutput(output io.Writer, prefix string) *log.Logger {
	return log.New(output, prefix, log.LstdFlags)
}

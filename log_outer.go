// This file definse the logOuter interface and several types of logOuter.
//
// emptyOuter = logOuter where both Out and Outf are noops
// lineOuter = logOuter where a newline is inserted after every call to
//			   Out and Outf
// fatalLineOuter = logOuter that logs message with inserted newline then
//					exits with call to os.EXIT(1)

package golog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type LogMessage struct {
	Level       int
	Nanoseconds int64
	Message     string
	// A map from the type of metadata to the metadata, if present.
	Metadata map[string]string
}

type LogOuter interface {
	// Output a LogMessage (to a file, to stderr, to a tester, etc). Output
	// must be safe to call from multiple threads.
	Output(*LogMessage)
}

// Render a formatted LogLocation to the buffer. If all present, format is 
// "{pack}.{func}/{file}:{line}". If some fields omitted, intelligently
// delimits the remaining fields.
func renderLogLocation(buf *bytes.Buffer, m *LogMessage) {
	if m == nil {
		return
	}

	packName, packPresent := m.Metadata["package"]
	file, filePresent := m.Metadata["file"]
	funcName, funcPresent := m.Metadata["function"]
	line, linePresent := m.Metadata["line"]

	if packPresent || filePresent || funcPresent || linePresent {
		buf.WriteString(" ")
	}

	// TODO(awreece) This logic is terrifying.
	if packPresent {
		buf.WriteString(packName)
	}
	if funcPresent {
		if packPresent {
			buf.WriteString(".")
		}
		buf.WriteString(funcName)
	}
	if (packPresent || funcPresent) && (filePresent || linePresent) {
		buf.WriteString("/")
	}
	if filePresent {
		buf.WriteString(file)
	}
	if linePresent {
		if filePresent {
			buf.WriteString(":")
		}
		buf.WriteString(line)
	}
}

// Format the message as a string, optionally inserting a newline.
// Format is: "L{level} {time} {pack}.{func}/{file}:{line}] {message}"
func formatLogMessage(m *LogMessage, insertNewline bool) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("L%d", m.Level))
	t := time.NanosecondsToLocalTime(m.Nanoseconds)
	buf.WriteString(t.Format(" 15:04:05.000000"))
	renderLogLocation(&buf, m)
	buf.WriteString("] ")
	buf.WriteString(m.Message)
	if insertNewline {
		buf.WriteString("\n")
	}
	return buf.String()
}

type writerLogOuter struct {
	lock sync.Mutex
	io.Writer
}

func (f *writerLogOuter) Output(m *LogMessage) {
	f.lock.Lock()
	defer f.lock.Unlock()

	// TODO(awreece) Handle short write?
	// Make sure to insert a newline.
	f.Write([]byte(formatLogMessage(m, true)))
}

// Returns a LogOuter wrapping the io.Writer.
func NewWriterLogOuter(f io.Writer) LogOuter {
	return &writerLogOuter{io.Writer: f}
}

// Returns a LogOuter wrapping the file, or an error if the file cannot be
// opened.
func NewFileLogOuter(filename string) (LogOuter, os.Error) {
	if file, err := os.Create(filename); err != nil {
		return nil, err
	} else {
		return NewWriterLogOuter(file), nil
	}

	panic("Code never reaches here, this mollifies the compiler.")
}

// We want to allow an abitrary testing framework.
type TestController interface {
	// We will assume that testers insert newlines in manner similar to 
	// the FEATURE of testing.T where it inserts extra newlines. >.<
	Log(...interface{})
	FailNow()
}

type testLogOuter struct {
	TestController
}

func (t *testLogOuter) Output(m *LogMessage) {
	// Don't insert an additional log message since the tester inserts them
	// for us.
	t.Log(formatLogMessage(m, false))
}

// Return a LogOuter wrapping the TestControlller.
func NewTestLogOuter(t TestController) LogOuter {
	return &testLogOuter{t}
}

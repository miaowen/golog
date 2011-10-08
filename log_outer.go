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
	"goprotobuf.googlecode.com/hg/proto"
	"os"
	"strconv"
	"time"
)

type LogOuter interface {
	Output(*LogMessage)
	FailNow()
}

func formatLogMessage(m *LogMessage) string {
	var buf bytes.Buffer
	buf.WriteString(levelStrings[int(proto.GetInt32(m.Level))])
	t := time.NanosecondsToLocalTime(proto.GetInt64(m.Nanoseconds))
	buf.WriteString(t.Format(" 15:04:05.000000 "))
	if m.Location != nil {
		l := *m.Location
		if l.Package != nil {
			buf.WriteString(*l.Package)
		}
		if l.File != nil {
			buf.WriteString(*l.File)
		}
		if l.Function != nil {
			buf.WriteString(*l.Function)
		}
		if l.Line != nil {
			buf.WriteString(strconv.Itoa(
				int(proto.GetInt32(l.Line))))
		}
	}
	buf.WriteString("] ")
	buf.WriteString(proto.GetString(m.Message))
	return buf.String()
}

type fileLogOuter struct {
	// TODO Insert mutex?
	*os.File
}

func (f *fileLogOuter) Output(m *LogMessage) {
	// TODO Grab mutex?
	s := proto.GetString(m.Message)
	l := len(s)
	if l > 0 {
		if s[l-1] == '\n' {
			f.WriteString(s)
		} else {
			f.WriteString(s + "\n")
		}
	}

	f.Sync()
}

func (f *fileLogOuter) FailNow() {
	// TODO Grab mutex?
	f.Close()
	os.Exit(1)
}

func NewFileLogOuter(f *os.File) LogOuter {
	return &fileLogOuter{f}
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
	s := proto.GetString(m.Message)
	l := len(s)
	if l > 0 {
		// Since testers insert newlines, we strip the newline
		// in our string.
		if s[l-1] == '\n' {
			t.Log(s[:l-1])
		} else {
			t.Log(s)
		}
	}
}

func NewTestLogOuter(t TestController) LogOuter {
	return &testLogOuter{t}
}

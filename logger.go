package mails

import (
	"fmt"
	"io"
)

type MockLogger struct {
	io io.Writer
}

func (m MockLogger) Print(log string) {
	_, _ = m.io.Write([]byte(log))
}

func (m MockLogger) Printf(format string, args ...interface{}) {
	_, _ = m.io.Write([]byte(fmt.Sprintf(format, args...)))
}

func NewLogger(writer io.Writer) MockLogger {
	return MockLogger{io: writer}
}

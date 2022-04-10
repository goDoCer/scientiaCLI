package logging

import (
	"fmt"
	"io"
)

var _ io.Writer = &logger{}

//L is the logger singleton
var L *logger = new(logger)

type logger struct {
	logs [][]byte
}

func (l *logger) Write(buf []byte) (n int, err error) {
	l.logs = append(l.logs, buf)
	return len(buf), nil
}

func (l *logger) Print() {
	fmt.Println("LOGS:")
	for _, buf := range l.logs {
		fmt.Println(string(buf))
	}
}

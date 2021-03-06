package logger

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh/terminal"
)

type ColorLoggerFactory struct {
	maxLength int
	tty       bool
}

type ColorLogger struct {
	name        string
	colorPrefix string
	factory     *ColorLoggerFactory
}

func NewColorLoggerFactory() *ColorLoggerFactory {
	return &ColorLoggerFactory{
		tty: terminal.IsTerminal(int(os.Stdout.Fd())),
	}
}

func (c *ColorLoggerFactory) Create(name string) Logger {
	if c.maxLength < len(name) {
		c.maxLength = len(name)
	}

	return &ColorLogger{
		name:        name,
		factory:     c,
		colorPrefix: <-colorPrefix,
	}
}

func (c *ColorLogger) Out(bytes []byte) {
	logFmt, name := c.getLogFmt()
	message := fmt.Sprintf(logFmt, name, string(bytes))
	fmt.Print(message)
}

func (c *ColorLogger) Err(bytes []byte) {
	logFmt, name := c.getLogFmt()
	message := fmt.Sprintf(logFmt, name, string(bytes))
	fmt.Fprint(os.Stderr, message)
}

func (c *ColorLogger) getLogFmt() (string, string) {
	pad := c.factory.maxLength

	logFmt := "%s | %s"
	if c.factory.tty {
		logFmt = c.colorPrefix + " %s"
	}

	name := fmt.Sprintf("%-"+strconv.Itoa(pad)+"s", c.name)

	return logFmt, name
}

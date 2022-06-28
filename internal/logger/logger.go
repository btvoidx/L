package logger

import (
	"fmt"
	"io"
	"strings"

	"github.com/muesli/termenv"
)

type Logger struct {
	Stdout  io.Writer
	Stderr  io.Writer
	Verbose bool
	Silent  bool
	prevEph bool
}

var p = termenv.ColorProfile()
var colorOk = p.Color("105")
var colorErr = p.Color("9")

func colorizeL(s string, c termenv.Color) string {
	if !strings.HasPrefix(s, "L") {
		return s
	}

	return strings.Replace(s, "L", termenv.String().Foreground(c).Styled("L"), 1)
}

func colorizeArguments(c termenv.Color, a []any) []any {
	for i, v := range a {
		if s, ok := v.(termenv.Style); ok {
			a[i] = s.Foreground(c)
		}
	}

	return a
}

func (l *Logger) Write(format string, a ...any) {
	if l.Silent {
		return
	}

	if l.prevEph {
		termenv.CursorPrevLine(1)
		termenv.ClearLine()
		l.prevEph = false
	}

	a = colorizeArguments(colorOk, a)
	format = colorizeL(format, colorOk)
	fmt.Fprintf(l.Stdout, format+"\n", a...)
}

func (l *Logger) WriteEphemeral(format string, a ...any) {
	if l.Silent {
		return
	}

	if l.prevEph {
		termenv.CursorPrevLine(1)
		termenv.ClearLine()
	}

	a = colorizeArguments(colorOk, a)
	format = colorizeL(format, colorOk)
	fmt.Fprintf(l.Stdout, format+"\n", a...)
	l.prevEph = true
}

func (l *Logger) Err(format string, a ...any) {
	if l.Silent {
		return
	}

	if l.prevEph {
		termenv.CursorPrevLine(1)
		termenv.ClearLine()
		l.prevEph = false
	}

	a = colorizeArguments(colorErr, a)
	format = colorizeL(format, colorErr)
	fmt.Fprintf(l.Stdout, format+"\n", a...)
}

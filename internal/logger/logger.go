package logger

import (
	"fmt"
	"io"
	"strings"

	"github.com/muesli/termenv"
)

var p = termenv.ColorProfile()

type Color func(format string, a ...any) string

var (
	None Color = func(s string, a ...any) string {
		return fmt.Sprintf(termenv.String().Styled(s), a...)
	}

	Magenta Color = func(s string, a ...any) string {
		return fmt.Sprintf(termenv.String().Foreground(p.Color("105")).Styled(s), a...)
	}

	Red Color = func(s string, a ...any) string {
		return fmt.Sprintf(termenv.String().Foreground(p.Color("9")).Styled(s), a...)
	}
)

type Logger struct {
	Stdout  io.Writer
	Stderr  io.Writer
	Verbose bool
	Silent  bool
	prevEph bool
}

func replaceL(s string, c Color) string {
	if !strings.HasPrefix(s, "L") {
		return s
	}

	return strings.Replace(s, "L", c("L"), 1)
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

	format = replaceL(format, Magenta)
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

	format = replaceL(format, Magenta)
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

	format = replaceL(format, Red)
	fmt.Fprintf(l.Stdout, format+"\n", a...)
}

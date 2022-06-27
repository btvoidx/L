package logger

import (
	"fmt"
	"io"
	"strings"

	"github.com/btvoidx/L/internal/color"
	"github.com/muesli/termenv"
)

type Logger struct {
	Stdout  io.Writer
	Stderr  io.Writer
	Verbose bool
	Silent  bool
	prevEph bool
}

func colorizeL(s string, c color.Color) string {
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

	format = colorizeL(format, color.Magenta)
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

	format = colorizeL(format, color.Magenta)
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

	format = colorizeL(format, color.Red)
	fmt.Fprintf(l.Stdout, format+"\n", a...)
}

package color

import (
	"fmt"

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

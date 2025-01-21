package shandler

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

func printer(src []io.Writer, data ...any) {
	for _, s := range src {
		fmt.Fprintln(s, data...)
	}
}

func printerf(src []io.Writer, pid string, format string, data ...any) {
	for _, s := range src {
		if pid == "" {
			fmt.Fprintf(s, format, data...)
		} else {
			fmt.Fprintf(s, "["+pid+"] "+format, data...)
		}
	}
}

func printerrj(src []io.Writer, g, pid, format string, data ...any) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
	}
	var left string
	if pid == "" {
		left = fmt.Sprintf(strings.TrimSpace(format), data...)
	} else {
		left = fmt.Sprintf("["+pid+"] "+strings.TrimSpace(format), data...)
	}

	rightWidth := width - len(left)
	if rightWidth < 0 {
		rightWidth = 0
	}

	for _, s := range src {
		fmt.Fprintf(s, "%s%*s\n", strings.TrimSpace(left), rightWidth, g)
	}
}

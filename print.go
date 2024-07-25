package shandler

import (
	"fmt"
	"io"
)

func printer(src []io.Writer, data ...any) {
	for _, s := range src {
		fmt.Fprintln(s, data...)
	}
}

func printerf(src []io.Writer, format string, data ...any) {
	for _, s := range src {
		fmt.Fprintf(s, format, data...)
	}
}

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

func printerf(src []io.Writer, pid string, format string, data ...any) {
	for _, s := range src {
		if pid == "" {
			fmt.Fprintf(s, format, data...)
		} else {
			fmt.Fprintf(s, "["+pid+"] "+format, data...)
		}
	}
}

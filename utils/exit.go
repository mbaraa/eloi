package utils

import (
	"os"

	"codeberg.org/lordbaraa/eloi/cli/cfmt"
)

func Exit(msg string) {
	_, _ = cfmt.Red().Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func ExitWhite(msg string) {
	_, _ = cfmt.Bold().Fprintln(os.Stderr, msg)
	os.Exit(1)
}

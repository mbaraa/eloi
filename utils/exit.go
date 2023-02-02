package utils

import (
	"fmt"
	"os"

	"github.com/mbaraa/eloi/cli"
)

func Exit(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, cli.ColorRed.StringColored(msg))
	os.Exit(1)
}

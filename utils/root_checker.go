package utils

import (
	"fmt"
	"os"
	"os/user"

	"github.com/mbaraa/eloi/cli"
)

func IsRoot() bool {
	u, err := user.Current()
	return err == nil && u.Uid == "0"
}

func AssertRoot() {
	if !IsRoot() {
		_, _ = fmt.Fprintln(os.Stderr, cli.ColorRed.StringColored("this action requires superuser access..."))
		os.Exit(1)
	}
}

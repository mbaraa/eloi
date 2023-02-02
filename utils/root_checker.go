package utils

import (
	"os/user"
)

func IsRoot() bool {
	u, err := user.Current()
	return err == nil && u.Uid == "0"
}

func AssertRoot() {
	if !IsRoot() {
		Exit("this action requires superuser access...")
	}
}

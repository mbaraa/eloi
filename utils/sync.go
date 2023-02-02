package utils

import (
	"os"
	"os/exec"
)

func Sync(repoName ...string) error {
	c := exec.Command("emerge", "--sync")
	c.Stdout = os.Stdout

	if len(repoName) > 0 && len(repoName[0]) > 0 {
		c.Args = append(c.Args, repoName[0])
	}

	return c.Run()
}

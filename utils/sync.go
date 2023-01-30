package utils

import (
	"os"
	"os/exec"
)

func Sync(repoName ...string) error {
	var _repoName string
	if len(repoName) > 0 {
		_repoName = repoName[0]
	}
	c := exec.Command("emerge", "--sync", _repoName)
	c.Stdout = os.Stdout
	return c.Run()
}

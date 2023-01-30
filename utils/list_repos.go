package utils

import (
	"os"
	"regexp"
	"strings"

	"github.com/mbaraa/eloi/globals"
)

func GetEnabledRepos() (map[string]bool, error) {
	reposDir, err := os.ReadDir(globals.ReposDirectory)
	if err != nil {
		return nil, err
	}

	repos := make([]string, 0)
	patt := regexp.MustCompile(`.*\.conf$`)
	for _, file := range reposDir {
		if patt.MatchString(file.Name()) {
			repos = append(repos, file.Name())
		}
	}

	repoNamePatt := regexp.MustCompile(`^\[.*\]$`)
	realReposNames := make(map[string]bool)
	for _, repoFileName := range repos {
		repoFile, err := os.ReadFile(globals.ReposDirectory + repoFileName)
		if err != nil {
			return nil, err
		}

		for _, line := range strings.Split(string(repoFile), "\n") {
			if repoNamePatt.MatchString(line) {
				realReposNames[line[1:len(line)-1]] = true
			}
		}
	}

	return realReposNames, nil
}

func IsRepoEnabled(repoName string) bool {
	repos, err := GetEnabledRepos()
	return err == nil && repos[repoName]
}

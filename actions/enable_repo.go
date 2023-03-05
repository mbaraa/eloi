package actions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/mbaraa/eloi/config"
	"github.com/mbaraa/eloi/globals"
	"github.com/mbaraa/eloi/models"
)

var _ Action = new(EnableRepoAction)

type EnableRepoAction struct{}

func (e *EnableRepoAction) Exec(output io.Writer, args ...any) error {
	name := args[0].(string)
	if e.isRepoEnabled(name) {
		return nil
	}

	_, err := output.Write([]byte("adding new repo: " + name + "...\n"))
	if err != nil {
		return err
	}

	resp, err := http.Get(fmt.Sprintf("%s/overlays/single?name=%s&simple=true", config.BackendAddress(), name))
	if err != nil {
		return err
	}

	overlay := models.ServerOverlay{}
	err = json.NewDecoder(resp.Body).Decode(&overlay)
	if err != nil {
		return err
	}

	_ = os.Mkdir(globals.ReposDirectory, 0755)

	err = e.createReposFile(e.getRepoString(overlay))
	if err != nil {
		return err
	}

	_, _ = output.Write([]byte("done ✓\n"))
	return Sync(name)
}

func (e *EnableRepoAction) NeedsRoot() bool {
	return true
}

func (e *EnableRepoAction) HasArgs() bool {
	return true
}

func (e *EnableRepoAction) getRepoString(overlay models.ServerOverlay) string {
	sb := new(strings.Builder)

	_, _ = fmt.Fprintf(sb, "[%s]\n", overlay.Name)
	_, _ = fmt.Fprintf(sb, "location = /var/db/repos/%s\n", overlay.Name)
	source := overlay.Source[0]
	_, _ = fmt.Fprintf(sb, "sync-type = %s\n", source.Type)
	_, _ = fmt.Fprintf(sb, "sync-uri = %s\n\n", source.Link)

	return sb.String()
}

func (e *EnableRepoAction) createReposFile(content string) error {
	f, err := os.OpenFile(globals.ReposDirectory+"/eloi-repos.conf", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if os.IsNotExist(err) {
		f, err = os.Create(globals.ReposDirectory + "/eloi-repos.conf")
		f.WriteString("# created by eloi\n")
		if err != nil {
			return err
		}
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

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

func (e *EnableRepoAction) isRepoEnabled(repoName string) bool {
	repos, err := GetEnabledRepos()
	return err == nil && repos[repoName]
}

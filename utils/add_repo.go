package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mbaraa/eloi/config"
	"github.com/mbaraa/eloi/globals"
	"github.com/mbaraa/eloi/models"
)

func AddOverlayRepo(name string) error {
	if IsRepoEnabled(name) {
		return nil
	}

	resp, err := http.Get(fmt.Sprintf("%s/overlays/single?name=%s", config.BackendAddress(), name))
	if err != nil {
		return err
	}

	overlay := models.Overlay{}
	err = json.NewDecoder(resp.Body).Decode(&overlay)
	if err != nil {
		return err
	}

	_ = os.Mkdir(globals.ReposDirectory, 0755)

	err = createReposFile(getRepoString(overlay))
	if err != nil {
		return err
	}

	return Sync(name)
}

func getRepoString(overlay models.Overlay) string {
	sb := new(strings.Builder)

	fmt.Fprintf(sb, "[%s]\n", overlay.Name)
	fmt.Fprintf(sb, "location = /var/db/repos/%s\n", overlay.Name)
	// TODO
	// get sync type from overlay
	fmt.Fprintln(sb, "sync-type = git")
	fmt.Fprintf(sb, "sync-uri = %s\n\n", overlay.Homepage)

	return sb.String()
}

func createReposFile(content string) error {
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

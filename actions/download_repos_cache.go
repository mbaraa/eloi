package actions

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mbaraa/eloi/config"
	"github.com/mbaraa/eloi/globals"
)

var _ Action = new(DownloadReposCacheAction)

type DownloadReposCacheAction struct {
}

func (d *DownloadReposCacheAction) Exec(output io.Writer, _ ...any) error {
	_, err := output.Write([]byte("Synchronizing overlays locally...\n"))
	if err != nil {
		return err
	}
	resp, err := http.Get(config.BackendAddress() + "/overlays/ebuilds")
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&globals.Ebuilds)
	if err != nil {
		return err
	}

	err = d.saveEbuildToLocalFile()
	if err != nil {
		return err
	}

	_, err = output.Write([]byte("All overlays has been synchronized ✓\n"))
	if err != nil {
		return err
	}

	return nil
}

func (d *DownloadReposCacheAction) NeedsRoot() bool {
	return true
}

func (d *DownloadReposCacheAction) HasArgs() bool {
	return false
}

func (d *DownloadReposCacheAction) createEloiDirectory() error {
	if _, err := os.Stat(cacheDirectory); !os.IsNotExist(err) {
		return nil
	}
	return os.Mkdir(cacheDirectory, 0755)
}

func (d *DownloadReposCacheAction) saveEbuildToLocalFile() error {
	err := d.createEloiDirectory()
	if err != nil {
		return err
	}

	ebuildsFile, err := os.Create(ebuildsFilePath)
	if err != nil {
		return err
	}
	defer ebuildsFile.Close()

	err = json.NewEncoder(ebuildsFile).Encode(globals.Ebuilds)
	if err != nil {
		return err
	}

	extractEbuildsNames()
	ebuildsNamesFile, err := os.Create(ebuildsNamesFilePath)
	if err != nil {
		return err
	}
	defer ebuildsNamesFile.Close()

	err = json.NewEncoder(ebuildsNamesFile).Encode(globals.EbuildsWithNamesOnly)
	if err != nil {
		return err
	}

	return nil
}

func extractEbuildsNames() {
	globals.EbuildsWithNamesOnly = make(map[string][]string)

	for _name := range globals.Ebuilds {
		name := _name[strings.Index(_name, "/")+1:]
		group := _name[:strings.Index(_name, "/")]

		globals.EbuildsWithNamesOnly[name] = append(globals.EbuildsWithNamesOnly[name], group+"/"+name)
	}
}

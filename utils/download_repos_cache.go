package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mbaraa/eloi/config"
	"github.com/mbaraa/eloi/globals"
)

const (
	eloiPath             = "/var/cache/eloi"
	ebuildsFilePath      = eloiPath + "/ebuilds.json"
	ebuildsNamesFilePath = eloiPath + "/ebuilds-names.json"
)

func DownloadReposCache() error {
	fmt.Println("Synchronizing overlays locally...")
	resp, err := http.Get(config.BackendAddress() + "/overlays/ebuilds")
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&globals.Ebuilds)
	if err != nil {
		return err
	}

	err = saveEbuildToLocalFile()
	if err != nil {
		return err
	}

	fmt.Println("All overlays has been synchronized ✓")
	return nil
}

func LoadLocalOverlays() error {
	err := loadLocalOverlays()
	if err != nil {
		return err
	}

	return nil
}

func createEloiDirectory() error {
	if _, err := os.Stat(eloiPath); !os.IsNotExist(err) {
		return nil
	}
	return os.Mkdir(eloiPath, 0755)
}

func saveEbuildToLocalFile() error {
	err := createEloiDirectory()
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

func loadLocalOverlays() error {
	ebuildsFile, err := os.Open(ebuildsFilePath)
	if os.IsNotExist(err) {
		err = Sync()
		if err != nil {
			return err
		}
	}
	defer ebuildsFile.Close()

	err = json.NewDecoder(ebuildsFile).Decode(&globals.Ebuilds)
	if err != nil {
		return err
	}

	ebuildsNamesFile, err := os.Open(ebuildsNamesFilePath)
	if err != nil {
		return err
	}
	defer ebuildsNamesFile.Close()

	err = json.NewDecoder(ebuildsNamesFile).Decode(&globals.EbuildsWithNamesOnly)
	if err != nil {
		return err
	}

	for name, ebuild := range globals.EbuildsWithNamesOnly {
		if name != strings.ToLower(name) {
			globals.EbuildsWithNamesOnly[strings.ToLower(name)] = ebuild
			delete(globals.EbuildsWithNamesOnly, name)
		}
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

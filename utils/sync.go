package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mbaraa/eloi/config"
	"github.com/mbaraa/eloi/globals"
)

const (
	eloiPath        = "/var/cache/eloi"
	ebuildsFilePath = eloiPath + "/ebuilds.json"
)

func Sync() error {
	fmt.Println("Synchronizing overlays locally...")
	resp, err := http.Get(config.BackendAddress() + "/overlays/ebuilds")
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&globals.Ebuilds)
	if err != nil {
		return err
	}

	err = createEloiDirectory()
	if err != nil {
		return err
	}

	ebuildsFile, err := os.Create(ebuildsFilePath)
	if err != nil {
		return err
	}

	err = json.NewEncoder(ebuildsFile).Encode(globals.Ebuilds)
	if err != nil {
		return err
	}

	fmt.Println("All overlays has been synchronized ✓")
	return nil
}

func createEloiDirectory() error {
	if _, err := os.Stat(eloiPath); !os.IsNotExist(err) {
		return nil
	}
	return os.Mkdir(eloiPath, 0755)
}

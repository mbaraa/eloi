package actions

import (
	"encoding/json"
	"github.com/mbaraa/eloi/db"
	"github.com/mbaraa/eloi/models"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/mbaraa/eloi/config"
	"github.com/mbaraa/eloi/globals"
)

var _ Action = new(DownloadReposCacheAction)

type DownloadReposCacheAction struct {
	db *gorm.DB
	// ebuildGroup-ebuildName => ebuildVersion => ebuild
	ebuilds map[string]map[string]models.ServerEbuild
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

	err = json.NewDecoder(resp.Body).Decode(&d.ebuilds)
	if err != nil {
		return err
	}

	err = d.createEloiDB()
	if err != nil {
		return err
	}

	err = d.saveEbuildToLocalDB()
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

func (d *DownloadReposCacheAction) createEloiDB() error {
	err := os.Mkdir(globals.CacheDirectory, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	_ = os.Remove(db.EloiDBPath)

	d.db, err = db.GetInstance()
	if err != nil {
		return err
	}

	return db.InitTables()
}

func (d *DownloadReposCacheAction) saveEbuildToLocalDB() error {
	return d.convertModelsAndPersist()
}

func (d *DownloadReposCacheAction) convertModelsAndPersist() error {
	var ebuilds []models.Ebuild

	for _, versions := range d.ebuilds {
		extraData := make([]models.ExtraData, 0)

		var ebuild models.Ebuild
		for version, ebuildsWithVersions := range versions {
			if len(ebuild.Name) == 0 {
				ebuild = models.Ebuild{
					Name:      versions[version].Name,
					GroupName: versions[version].GroupName,
					Homepage:  versions[version].Homepage,
					ExtraData: extraData,
				}
			}
			extraData = append(extraData, models.ExtraData{
				Version:      version,
				OverlayName:  ebuildsWithVersions.OverlayName,
				Flags:        ebuildsWithVersions.Flags,
				License:      ebuildsWithVersions.License,
				Architecture: ebuildsWithVersions.Architecture,
				EbuildID:     0,
			})
		}
		ebuild.ExtraData = extraData
		ebuilds = append(ebuilds, ebuild)
	}

	sort.Slice(ebuilds, func(i, j int) bool {
		return strings.Compare(ebuilds[i].Name, ebuilds[j].Name) < 0
	})

	for _, ebuild := range ebuilds {
		err := d.db.Create(&ebuild).Error
		if err != nil {
			return err
		}
	}

	return nil
}

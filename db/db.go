package db

import (
	"os"

	"github.com/mbaraa/eloi/globals"
	"github.com/mbaraa/eloi/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const EloiDBPath = globals.CacheDirectory + "eloi.db"

var instance *gorm.DB = nil

func GetInstance() (*gorm.DB, error) {
	if _, err := os.Stat(globals.CacheDirectory); !os.IsExist(err) {
		_ = os.Mkdir(globals.CacheDirectory, 0755)
	}
	if _, err := os.Stat(EloiDBPath); instance == nil || err != nil {
		var err error
		instance, err = gorm.Open(sqlite.Open(EloiDBPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func InitTables() error {
	instance, err := GetInstance()
	if err != nil {
		return err
	}
	if instance != nil {
		return instance.AutoMigrate(
			new(models.Ebuild),
			new(models.Overlay),
			new(models.ExtraData),
		)
	} else {
		panic("No DB connection was initialized")
	}
}

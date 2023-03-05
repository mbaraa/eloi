package db

import (
	"errors"
	"github.com/mbaraa/eloi/models"
	"gorm.io/gorm"
)

var ebuildDBInstance *EbuildDB = nil

type EbuildDB struct {
	db *gorm.DB
}

func GetEbuildDB() (*EbuildDB, error) {
	if ebuildDBInstance == nil {
		db, err := GetInstance()
		if err != nil {
			return nil, err
		}
		ebuildDBInstance = &EbuildDB{db: db}
	}
	return ebuildDBInstance, nil
}

// Add creates a new record of the given object, and returns an occurring error
// the new object is a pointer, so it updates the object's id after creation
func (e *EbuildDB) Add(ebuild *models.Ebuild) error {
	if ebuild == nil {
		return errors.New("ebuild's object is nil")
	}

	return e.db.
		Model(new(models.Ebuild)).
		Create(ebuild).
		Error
}

// AddMany is same as Add but for numerous objects
func (e *EbuildDB) AddMany(ebuilds []*models.Ebuild) error {
	if len(ebuilds) == 0 {
		return errors.New("ebuilds' slice is empty")
	}

	return e.db.
		Model(new(models.Ebuild)).
		Create(ebuilds).
		Error
}

// Get retrieves the object which has the given id
func (e *EbuildDB) Get(id uint) (models.Ebuild, error) {
	var ebuild models.Ebuild

	err := e.db.
		Model(new(models.Ebuild)).
		Preload("ExtraData").
		First(&ebuild, "id = ?", id).
		Error

	if err != nil {
		return models.Ebuild{}, errors.New("record was not found")
	}

	return ebuild, nil
}

// GetByConds is the extended version of Get,
// which uses a given search condition and retrieves every record with the given condition
func (e *EbuildDB) GetByConds(conds ...any) ([]models.Ebuild, error) {
	if !checkConds(conds...) {
		return nil, errors.New("invalid query conditions")
	}

	var foundRecords []models.Ebuild

	err := e.db.
		Model(new(models.Ebuild)).
		Preload("ExtraData").
		Find(&foundRecords, conds...).
		Error

	if err != nil || len(foundRecords) == 0 {
		return nil, errors.New("no records were found")
	}

	return foundRecords, nil
}

// GetAll retrieves all the records of the given model
func (e *EbuildDB) GetAll() ([]models.Ebuild, error) {
	return e.GetByConds("id != ?", 0)
}

// Count returns the number of records of the given model
func (e *EbuildDB) Count() int64 {
	var count int64

	err := e.db.
		Model(new(models.Ebuild)).
		Count(&count).
		Error

	if err != nil {
		return 0
	}

	return count
}

// checkConds reports whether the provided conditions are valid or not
func checkConds(conds ...any) bool {
	return len(conds) > 1 && checkCondsMeaning(conds...)
}

func checkCondsMeaning(conds ...any) bool {
	ok := false

	switch conds[0].(type) {
	case string:
		ok = true
	default:
		return false
	}

	for _, cond := range conds[1:] {
		switch cond.(type) {
		case bool,
			int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64,
			complex64, complex128,
			string:
			ok = true
		default:
			return false
		}
	}

	return ok
}

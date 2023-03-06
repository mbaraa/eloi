package models

import "gorm.io/gorm"

type Overlay struct {
	gorm.Model `json:"-"`
	Id         uint `gorm:"primaryKey;autoIncrement"`
	OwnerEmail string
	OwnerName  string
	OwnerType  string

	SourceType string
	SourceLink string

	Name        string `gorm:"unique"`
	Description string
	Homepage    string
	Feed        string
}

type ExtraData struct {
	gorm.Model   `json:"-"`
	Id           uint `gorm:"primaryKey;autoIncrement"`
	Version      string
	OverlayName  string
	Flags        string
	License      string
	Architecture string
	EbuildID     uint
}

type Ebuild struct {
	gorm.Model  `json:"-"`
	Id          uint `gorm:"primaryKey;autoIncrement"`
	Name        string
	GroupName   string
	Homepage    string
	Description string `json:"description"`

	ExtraData []ExtraData // version => providers overlays
}

func (e Ebuild) FullName() string {
	return e.GroupName + "/" + e.Name
}

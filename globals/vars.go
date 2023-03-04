package globals

import "github.com/mbaraa/eloi/models"

var (
	// ebuildGroup-ebuildName => ebuildVersion => ebuild
	Ebuilds map[string]map[string]*models.Ebuild
	// ebuildName => ebuildFullName
	// used for indexing
	EbuildsWithNamesOnly map[string][]string
)

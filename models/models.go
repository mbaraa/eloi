package models

type Overlay struct {
	EbuildGroups map[string]*EbuildGroup `json:"ebuildGroups"`
	Name         string                  `json:"name"`
	URL          string                  `json:"url"`
	Homepage     string                  `json:"homepage"`
	Description  string                  `json:"description"`
	Source       string                  `json:"source"`
	Mail         string                  `json:"mail"`
	Feed         string                  `json:"feed"`
	NumEbuilds   int                     `json:"numEbuilds"`
}

type EbuildGroup struct {
	Ebuilds map[string]*Ebuild `json:"ebuilds"`
	Name    string             `json:"name"`
	Link    string             `json:"link"`
}

type Ebuild struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Homepage     string `json:"homepage"`
	Flags        string `json:"flags"`
	License      string `json:"license"`
	OverlayName  string `json:"overlayName"`
	GroupName    string `json:"groupName"`
	Architecture string `json:"architecture"`
}

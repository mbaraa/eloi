package models

type Overlay struct {
	Owner struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Type  string `json:"type"`
	} `json:"owner"`

	Source []struct {
		Type string `json:"type"`
		Link string `json:"link"`
	} `json:"source"`

	EbuildGroups map[string]*EbuildGroup `json:"ebuildGroups"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Feed        string `json:"feed"`
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

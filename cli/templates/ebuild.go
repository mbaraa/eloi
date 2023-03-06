package templates

import (
	"fmt"
	"strings"

	"github.com/mbaraa/eloi/cli/cfmt"

	"github.com/mbaraa/eloi/models"
)

func EbuildTemplate(ebuild models.Ebuild) string {
	/*
	   * group-ebuild
	       Homepage:    	   https://link.com
	       Description:   	   something really descriptive
	       Available overlays: overlay1, overlay2
	       Available versions: 1.2, 1.3, 1.4, **9999
	*/

	sb := new(strings.Builder)

	_, _ = cfmt.Green().Fprint(sb, ebuild.GroupName)
	_, _ = cfmt.White().Fprint(sb, "/")
	_, _ = cfmt.White().Fprintln(sb, ebuild.Name)

	if len(ebuild.Homepage) != 0 {
		_, _ = cfmt.Green().Bold().Fprintf(sb, "\tHomepage: ")
		_, _ = cfmt.White().Fprintln(sb, ebuild.Homepage)
	}

	if len(ebuild.Description) != 0 {
		_, _ = cfmt.Green().Bold().Fprintf(sb, "\tDescription: ")
		_, _ = cfmt.White().Fprintln(sb, ebuild.Description)
	}

	versions := make(map[string]bool)
	overlays := make(map[string]bool)
	versionsStr := new(strings.Builder)
	overlaysStr := new(strings.Builder)

	for _, extraData := range ebuild.ExtraData {
		if !versions[extraData.Version] {
			_, _ = fmt.Fprint(versionsStr, extraData.Version)
			_, _ = fmt.Fprint(versionsStr, ", ")
		}

		if !overlays[extraData.OverlayName] {
			_, _ = fmt.Fprint(overlaysStr, extraData.OverlayName)
			_, _ = fmt.Fprint(overlaysStr, ", ")
		}

		versions[extraData.Version] = true
		overlays[extraData.OverlayName] = true
	}

	_, _ = cfmt.Green().Bold().Fprintf(sb, "\tAvailable Overlays: ")
	_, _ = cfmt.Yellow().Fprintln(sb, removeTrailingCommaSpace(overlaysStr.String()))

	_, _ = cfmt.Green().Bold().Fprintf(sb, "\tAvailable Versions: ")
	_, _ = cfmt.Yellow().Fprintln(sb, removeTrailingCommaSpace(versionsStr.String()))

	return sb.String()
}

func removeTrailingCommaSpace(s string) string {
	if len(s) < 2 {
		return ""
	}
	if s[len(s)-2:] == ", " {
		return s[:len(s)-2]
	}
	return s
}

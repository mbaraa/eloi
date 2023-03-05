package templates

import (
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
	_, _ = cfmt.White().Bold().Fprintln(sb, ebuild.Name)

	if len(ebuild.Homepage) != 0 {
		_, _ = cfmt.Green().Bold().Fprintf(sb, "\tHomepage: ")
		_, _ = cfmt.White().Fprintln(sb, ebuild.Homepage)
	}

	// if len(ebuild.Description) != 0 {
	// 	_, _ = cfmt.Green().Bold().Fprintf(sb, "\tDescription: ")
	// 	_, _ = cfmt.White().Fprintln(sb, ebuild.Description)
	// }

	versions := new(strings.Builder)
	overlays := new(strings.Builder)
	for i, extraData := range ebuild.ExtraData {
		_, _ = cfmt.Yellow().Fprint(versions, extraData.Version)
		_, _ = cfmt.Yellow().Fprint(overlays, extraData.OverlayName)
		if i < len(ebuild.ExtraData)-1 {
			_, _ = cfmt.Yellow().Fprint(versions, ", ")
			_, _ = cfmt.Yellow().Fprint(overlays, ", ")
		}
	}

	_, _ = cfmt.Green().Bold().Fprintf(sb, "\tAvailable Overlays: ")
	_, _ = cfmt.Yellow().Fprintln(sb, overlays.String())
	_, _ = cfmt.Green().Bold().Fprintf(sb, "\tAvailable Versions: ")
	_, _ = cfmt.Yellow().Fprintln(sb, versions.String())

	return sb.String()
}

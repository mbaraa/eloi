package templates

import (
	"fmt"
	"strings"

	"github.com/mbaraa/eloi/cli/cfmt"
	"github.com/mbaraa/eloi/models"
)

func EbuildTemplate(ebuildWithVersions map[string]models.Ebuild) string {
	/*

	   * group-ebuild
	       Available versions: 1.2 1.3 1.4
	       Homepage:  https://link.com
	       Description:   something really descriptive
	       Available Overlays:    overlay1, overlay2, ...

	       or

	       Available versions:    overlay1(1.2 1.3) overlay2(1.4 **9999)
	*/

	sb := new(strings.Builder)

	providers := make(map[string][]string) // overlayName => version
	var name, groupName, license string
	for _, ebuild := range ebuildWithVersions {
		name, groupName, license = ebuild.Name, ebuild.GroupName, ebuild.License
		providers[ebuild.OverlayName] = append(providers[ebuild.OverlayName], ebuild.Version)
	}

	versionsStr := new(strings.Builder) // overlayN(versionN...)...
	for name, versions := range providers {
		cfmt.Yellow().Fprint(versionsStr, name+"(")
		for i, version := range versions {
			cfmt.Yellow().Fprint(versionsStr, version)
			if i < len(versions)-1 {
				versionsStr.WriteString(" ")
			}
		}
		cfmt.Yellow().Fprint(versionsStr, ") ")
	}

	fmt.Fprintf(sb, "%s/%s\n", cfmt.Green().Sprint(groupName), cfmt.Bold().Sprint(name))
	fmt.Fprintf(sb, "\t%s %s\n", cfmt.Green().Sprint("Available versions:"), versionsStr.String())
	if len(license) != 0 {
		fmt.Fprintf(sb, "\t%s %s\n", cfmt.Green().Sprint("License:"), license)
	}
	return sb.String()
}

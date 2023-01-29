package templates

import (
	"fmt"
	"strings"

	"github.com/mbaraa/eloi/cli"
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
		versionsStr.WriteString(cli.ColorLightyellow.StringColored(name[len("Overlay: "):] + "("))
		for i, version := range versions {
			versionsStr.WriteString(cli.ColorLightyellow.StringColored(version))
			if i < len(versions)-1 {
				versionsStr.WriteString(" ")
			}
		}
		versionsStr.WriteString(cli.ColorLightyellow.StringColored(") "))
	}

	sb.WriteString(fmt.Sprintf("%s %s/%s\n", cli.ColorLightgreen.StringColored("*"), groupName, cli.ColorBold.StringColored(name)))
	sb.WriteString(fmt.Sprintf("\t%s %s\n", cli.ColorLightgreen.StringColored("Available versions:"), versionsStr.String()))
	if len(license) != 0 {
		sb.WriteString(fmt.Sprintf("\t%s %s\n", cli.ColorLightgreen.StringColored("License:"), license))
	}
	return sb.String()
}

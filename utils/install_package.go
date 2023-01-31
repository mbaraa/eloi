package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mbaraa/eloi/cli"
	"github.com/mbaraa/eloi/globals"
	"github.com/mbaraa/eloi/models"
)

func InstallPackage(ebuild map[string]*models.Ebuild) error {
	versions := make([]PackageEntity, 0)

	var name, group string
	for _, details := range ebuild {
		name, group = details.Name, details.GroupName
		versions = append(versions, PackageEntity{
			Display:  fmt.Sprintf("%s::%s", cli.ColorYellow.StringColored(details.Version), cli.ColorGreen.StringColored(details.OverlayName)),
			FullName: fmt.Sprintf("%s#%s", details.Version, details.OverlayName),
		})
	}

	reallySelected := ebuild[versions[0].FullName[:strings.Index(versions[0].FullName, "#")]]

	if len(ebuild) > 1 {
	selectVersion:
		fmt.Printf("\nthere are %s versions available of the package %s\n", cli.ColorYellow.StringColored(fmt.Sprint(len(ebuild))), cli.ColorGreen.StringColored(group+"/"+name))
		for i, version := range versions {
			fmt.Printf("(%s) %s\n", cli.ColorPurple.StringColored(fmt.Sprint(i+1)), version.Display)
		}

		prompt := cli.ColorGreen.StringColored("==>")
		fmt.Printf("%s Select a version to install\n%s ", prompt, prompt)
		selection := 0
		fmt.Scan(&selection)

		if selection-1 > len(versions) || selection-1 < 0 {
			fmt.Println("invalid selection")
			goto selectVersion
		}

		selected := versions[selection-1]
		version := selected.FullName[:strings.Index(selected.FullName, "#")]
		overlay := selected.FullName[strings.Index(selected.FullName, "#")+1:]

		ebuildVersions := globals.Ebuilds[group+"/"+name]

		// TODO
		// handle multiple overlays providers for the same version
		for _version, ebuild := range ebuildVersions {
			if version == _version && ebuild.OverlayName == overlay {
				reallySelected = ebuild
			}
		}
	}

	return installPackage(reallySelected)
}

func installPackage(ebuild *models.Ebuild) error {
	err := AddOverlayRepo(ebuild.OverlayName)
	if err != nil {
		return err
	}

	err = AcceptPackageLicense(*ebuild)
	if err != nil {
		return err
	}

	// TODO
	// add license and flags enablers

	c := exec.Command("emerge", "-qav", fmt.Sprintf("=%s/%s-%s::%s", ebuild.GroupName, ebuild.Name, ebuild.Version, ebuild.OverlayName))
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	return c.Run()
}

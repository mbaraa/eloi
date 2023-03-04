package actions

import (
	"errors"
	"fmt"
	"github.com/mbaraa/eloi/utils"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/mbaraa/eloi/cli/cfmt"
	"github.com/mbaraa/eloi/globals"
	"github.com/mbaraa/eloi/models"
)

var _ Action = new(EbuildInstallAction)

type EbuildInstallAction struct {
	output io.Writer
}

func (i *EbuildInstallAction) Exec(output io.Writer, args ...any) error {
	i.output = output
	return i.promptSelectPackage(args[0].([]PackageEntity))
}

func (i *EbuildInstallAction) NeedsRoot() bool {
	return true
}

func (i *EbuildInstallAction) HasArgs() bool {
	return false
}

func (i *EbuildInstallAction) promptSelectPackage(pkgs []PackageEntity) error {
	if len(pkgs) == 0 {
		return errors.New("no packages were found")
	}
	prompt := cfmt.Green().Sprint("==>")
	_, err := i.output.Write([]byte(fmt.Sprintf("%s Select a package to install\n%s ", prompt, prompt)))
	if err != nil {
		return err
	}
	selection := 0
	_, err = fmt.Scan(&selection)
	if err != nil {
		return err
	}

	if selection-1 > len(pkgs) || selection-1 < 0 {
		_, err = i.output.Write([]byte("there is nothing to do\n"))
	} else {
		err = i.listEbuildsForInstallation(globals.Ebuilds[pkgs[selection-1].FullName])
	}
	return err
}

func (i *EbuildInstallAction) listEbuildsForInstallation(ebuild map[string]*models.Ebuild) error {
	versions := make([]PackageEntity, 0)

	var name, group string
	for _, details := range ebuild {
		name, group = details.Name, details.GroupName
		versions = append(versions, PackageEntity{
			Display:  fmt.Sprintf("%s::%s", cfmt.Yellow().Sprint(details.Version), cfmt.Green().Sprint(details.OverlayName)),
			FullName: fmt.Sprintf("%s#%s", details.Version, details.OverlayName),
		})
	}

	reallySelected := ebuild[versions[0].FullName[:strings.Index(versions[0].FullName, "#")]]

	if len(ebuild) > 1 {
	selectVersion:
		_, err := i.output.Write([]byte(fmt.Sprintf("\nthere are %s versions available of the package %s\n", cfmt.Yellow().Sprint(len(ebuild)), cfmt.Green().Sprint(group+"/"+name))))
		if err != nil {
			return err
		}
		for index, version := range versions {
			_, err := i.output.Write([]byte(fmt.Sprintf("(%s) %s\n", cfmt.Magenta().Sprint(index+1), version.Display)))
			if err != nil {
				return err
			}
		}

		prompt := cfmt.Green().Sprint("==>")
		_, err = i.output.Write([]byte(fmt.Sprintf("%s Select a version to install\n%s ", prompt, prompt)))
		if err != nil {
			return err
		}
		selection := 0
		_, err = fmt.Scan(&selection)
		if err != nil {
			return err
		}

		if selection-1 > len(versions) || selection-1 < 0 {
			_, err = i.output.Write([]byte("invalid selection\n"))
			if err != nil {
				return err
			}
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

	return i.installEbuild(reallySelected)
}

func (i *EbuildInstallAction) installEbuild(ebuild *models.Ebuild) error {
	err := new(EnableRepoAction).Exec(i.output, ebuild.OverlayName)
	if err != nil {
		return err
	}

	err = utils.AcceptPackageLicense(*ebuild)
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

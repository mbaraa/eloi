package actions

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/mbaraa/eloi/cli/cfmt"
	"github.com/mbaraa/eloi/cli/templates"
	"github.com/mbaraa/eloi/models"
)

var _ Action = new(EbuildInstallAction)

type EbuildInstallAction struct {
	output io.Writer
}

func (i *EbuildInstallAction) Exec(output io.Writer, args ...any) error {
	i.output = output
	return i.promptSelectPackage(args[0].([]models.Ebuild))
}

func (i *EbuildInstallAction) NeedsRoot() bool {
	return true
}

func (i *EbuildInstallAction) HasArgs() bool {
	return false
}

func (i *EbuildInstallAction) promptSelectPackage(pkgs []models.Ebuild) error {
	if len(pkgs) == 0 {
		return errors.New("no packages were found")
	}

	for index, ebuild := range pkgs {
		_, err := i.output.Write([]byte(fmt.Sprintf("(%s) %s\n", cfmt.Magenta().Sprint(index+1), templates.EbuildTemplate(ebuild))))
		if err != nil {
			return err
		}
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
		err = i.listEbuildsForInstallation(pkgs[selection-1])
	}
	return err
}

func (i *EbuildInstallAction) listEbuildsForInstallation(ebuild models.Ebuild) error {
	selectedProviderIndex := 1

	if len(ebuild.ExtraData) > 1 {
	selectVersion:
		_, _ = i.output.Write([]byte(fmt.Sprintf("there are %s versions available of the package %s\n",
			cfmt.Yellow().Sprint(len(ebuild.ExtraData)),
			cfmt.Green().Sprint(ebuild.FullName()))))

		for index, provider := range ebuild.ExtraData {
			_, _ = fmt.Fprintf(i.output, "(%s) %s\n",
				cfmt.Magenta().Sprint(index+1),
				fmt.Sprintf("%s::%s", cfmt.Yellow().Sprint(provider.Version), cfmt.Green().Sprint(provider.OverlayName)))
		}

		prompt := cfmt.Green().Sprint("==>")
		_, _ = i.output.Write([]byte(fmt.Sprintf("%s Select a version to install\n%s ", prompt, prompt)))
		_, _ = fmt.Scan(&selectedProviderIndex)

		if selectedProviderIndex-1 > len(ebuild.ExtraData) || selectedProviderIndex-1 < 0 {
			_, _ = i.output.Write([]byte("invalid selection\n"))
			goto selectVersion
		}
	}

	return i.installEbuild(ebuild, selectedProviderIndex-1)
}

func (i *EbuildInstallAction) installEbuild(ebuild models.Ebuild, providerIndex int) error {
	err := NewEnableRepoForPackageAction().Exec(i.output, ebuild, providerIndex)
	if err != nil {
		return err
	}

	err = NewLicenceUnmaskAction().Exec(i.output, ebuild, providerIndex)
	if err != nil {
		return err
	}

	// TODO
	// add option for enabling use flags always for an ebuild, or just for the current install
	err = NewFlagsUnmaskAction().Exec(i.output, ebuild, providerIndex)
	if err != nil {
		return err
	}

	c := exec.Command("emerge", "-qav", fmt.Sprintf("=%s-%s::%s",
		ebuild.FullName(), ebuild.ExtraData[providerIndex].Version, ebuild.ExtraData[providerIndex].OverlayName))
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	return c.Run()
}

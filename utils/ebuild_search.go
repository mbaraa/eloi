package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/mbaraa/eloi/cli"
	"github.com/mbaraa/eloi/cli/templates"
	"github.com/mbaraa/eloi/globals"
	"github.com/mbaraa/eloi/models"
)

type PackageEntity struct {
	Display, FullName string
}

func FindEbuild(name string) {
	allMatchingEbuilds := make([]string, 0)

	for _name := range globals.Ebuilds {
		if strings.Contains(_name, name) {
			allMatchingEbuilds = append(allMatchingEbuilds, globals.EbuildsWithNamesOnly[_name[strings.Index(_name, "/")+1:]]...)
		}
		//        globals.EbuildsWithNamesOnly[name]
	}

	ebuildsToDisplay := make([]PackageEntity, 0)

	for _, ebuild := range allMatchingEbuilds {
		name := ebuild[strings.Index(ebuild, "/")+1:]
		group := ebuild[:strings.Index(ebuild, "/")]

		versions := make(map[string]models.Ebuild)
		for _, ebuildVersion := range globals.Ebuilds[group+"/"+name] {
			versions[ebuildVersion.Version] = *ebuildVersion
		}

		ebuildsToDisplay = append(ebuildsToDisplay, PackageEntity{templates.EbuildTemplate(versions), group + "/" + name})
	}

	for i, ebuild := range ebuildsToDisplay {
		fmt.Printf("(%s) %s\n", cli.ColorPurple.StringColored(fmt.Sprint(i+1)), ebuild.Display)
	}

	PromptSelectPackage(ebuildsToDisplay)
}

func PromptSelectPackage(pkgs []PackageEntity) {
	prompt := cli.ColorGreen.StringColored("==>")
	fmt.Printf("%s Select a package to install\n%s ", prompt, prompt)
	selection := 0
	fmt.Scan(&selection)

	if selection-1 > len(pkgs) {
		fmt.Println("there is nothing to do")
	} else {
		err := InstallPackage(pkgs[selection-1].FullName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		}
	}
}

func InstallPackage(ebuildName string) error {
	fmt.Printf("%+v\n", globals.Ebuilds[ebuildName])
	return nil
}

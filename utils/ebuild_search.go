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
		if strings.Contains(strings.ToLower(_name), strings.ToLower(name)) {
			allMatchingEbuilds = append(allMatchingEbuilds, globals.EbuildsWithNamesOnly[strings.ToLower(_name[strings.Index(_name, "/")+1:])]...)
		}
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
	if len(pkgs) == 0 {
		fmt.Fprintln(os.Stderr, "no packages were found!")
		return
	}
	prompt := cli.ColorGreen.StringColored("==>")
	fmt.Printf("%s Select a package to install\n%s ", prompt, prompt)
	selection := 0
	fmt.Scan(&selection)

	if selection-1 > len(pkgs) || selection-1 < 0 {
		fmt.Println("there is nothing to do")
	} else {
		err := InstallPackage(globals.Ebuilds[pkgs[selection-1].FullName])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		}
	}
}

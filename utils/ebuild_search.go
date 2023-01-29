package utils

import (
	"fmt"
	"strings"

	"github.com/mbaraa/eloi/cli/templates"
	"github.com/mbaraa/eloi/globals"
	"github.com/mbaraa/eloi/models"
)

func FindEbuild(name string) {
	allMatchingEbuilds := make([]string, 0)

	for _name := range globals.Ebuilds {
		if strings.Contains(_name, name) {
			allMatchingEbuilds = append(allMatchingEbuilds, globals.EbuildsWithNamesOnly[_name[strings.Index(_name, "/")+1:]]...)
		}
		//        globals.EbuildsWithNamesOnly[name]
	}

	for _, ebuild := range allMatchingEbuilds {
		name := ebuild[strings.Index(ebuild, "/")+1:]
		group := ebuild[:strings.Index(ebuild, "/")]

		versions := make(map[string]models.Ebuild)
		for _, ebuildVersion := range globals.Ebuilds[group+"/"+name] {
			versions[ebuildVersion.Version] = *ebuildVersion
		}
		fmt.Println(templates.EbuildTemplate(versions))
	}
}

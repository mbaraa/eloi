package utils

import (
	"fmt"

	"github.com/mbaraa/eloi/globals"
)

func FindEbuild(name string) {
	fmt.Printf("%+v\n", globals.EbuildsWithNamesOnly[name])
}

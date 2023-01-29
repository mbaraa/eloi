package cmd

import (
	"flag"
	"fmt"

	"github.com/mbaraa/eloi/utils"
)

func Start() {
	var meow bool
	flag.BoolVar(&meow, "meow", false, "meow meow meow")

	var sync bool
	flag.BoolVar(&sync, "sync", false, "sync local repos")

	var ebuildName string
	flag.StringVar(&ebuildName, "S", "", "search for an ebuild")
	flag.StringVar(&ebuildName, "search", "", "search for an ebuild")

	flag.Parse()

	////

	if meow {
		fmt.Println("meow meow meow")
	}

	if sync {
		err := utils.Sync()
		if err != nil {
			panic(err)
		}
	}

	if ebuildName != "" {
		err := utils.LoadLocalOverlays()
		if err != nil {
			panic("Local repos are not synced, run with --sync to sync them")
		}
		utils.FindEbuild(ebuildName)
	}
}

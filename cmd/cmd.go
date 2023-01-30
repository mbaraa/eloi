package cmd

import (
	"flag"
	"fmt"

	"github.com/mbaraa/eloi/utils"
)

func Start() {
	var meow bool
	flag.BoolVar(&meow, "meow", false, "meow meow meow")

	var download bool
	flag.BoolVar(&download, "download", false, "download overlays repos cache")

	var ebuildName string
	flag.StringVar(&ebuildName, "S", "", "search for an ebuild")
	flag.StringVar(&ebuildName, "search", "", "search for an ebuild")

	var overlayName string
	flag.StringVar(&overlayName, "enable", "", "add a specific repository")

	var sync bool
	flag.BoolVar(&sync, "sync", false, "sync portage repos")

	flag.Parse()

	////

	if meow {
		fmt.Println("meow meow meow")
	}

	if len(overlayName) != 0 {
		utils.AssertRoot()
		err := utils.AddOverlayRepo(overlayName)
		if err != nil {
			panic(err)
		}
	}

	if sync {
		utils.AssertRoot()
		err := utils.Sync()
		if err != nil {
			panic(err)
		}
	}

	if download {
		utils.AssertRoot()
		err := utils.DownloadReposCache()
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

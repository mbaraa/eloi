package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/mbaraa/eloi/utils"
)

var (
	flags        = flag.NewFlagSet("Eloi", flag.ExitOnError)
	meow         = new(bool)
	download     = new(bool)
	syncRepoName = new(string)
	ebuildName   = new(string)
	repoName     = new(string)
)

func Start() {
	registerFlags()
	flags.Parse(os.Args[1:])
	runWithGivenArgs()
}

func runWithGivenArgs() {
	if *meow {
		fmt.Println("meow meow meow")
	} else if *download {
		utils.AssertRoot()
		err := utils.DownloadReposCache()
		if err != nil {
			panic(err)
		}
	} else if len(*ebuildName) > 0 {
		err := utils.LoadLocalOverlays()
		if err != nil {
			panic("Local repos are not synced, run with --sync to sync them")
		}
		utils.FindEbuild(*ebuildName)
	} else if len(*repoName) > 0 {
		utils.AssertRoot()
		err := utils.AddOverlayRepo(*repoName)
		if err != nil {
			panic(err)
		}
	} else if len(*syncRepoName) > 0 {
		utils.AssertRoot()
		err := utils.Sync(*syncRepoName)
		if err != nil {
			panic(err)
		}
	} else {
		flags.Usage()
	}
}

func registerFlags() {
	flags.BoolVar(meow, "meow", false, "meow meow meow")
	flags.BoolVar(download, "download", false, "download overlays repos cache")
	flags.StringVar(ebuildName, "search", "", "find ebuild to install")
	flags.StringVar(ebuildName, "S", "", "find ebuild to install")
	flags.StringVar(repoName, "enable", "", "add a specific repository")
	flags.StringVar(syncRepoName, "sync", "", "sync portage repos")
}

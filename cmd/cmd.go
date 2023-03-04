package cmd

import (
	"fmt"
	"github.com/mbaraa/eloi/actions"
	"os"

	"github.com/mbaraa/eloi/utils"
)

const (
	usageStr = `Usage of Eloi:
  -S string
  --search string
        find an ebuild to install
  --download
        download overlays repos cache
  --enable string
        add a specific repository`
)

var argsActions = map[string]actions.ActionType{
	"--download": actions.DownloadReposCacheActionType,
	"-S":         actions.EbuildSearchActionType,
	"--search":   actions.EbuildSearchActionType,
	"--enable":   actions.EnableRepoActionType,
}

func Start() {
	if len(os.Args) < 2 {
		utils.ExitWhite(usageStr)
	}

	actionType, ok := argsActions[os.Args[1]]
	if !ok {
		utils.ExitWhite(usageStr)
	}

	var arg string
	if len(os.Args) > 2 {
		arg = os.Args[2]
	}

	runWithGivenArgs(actionType, arg)
}

func runWithGivenArgs(actionType actions.ActionType, arg string) {
	action := actions.GetActionFactory(actionType)
	if action.NeedsRoot() {
		utils.AssertRoot()
	}

	if action.HasArgs() && len(arg) == 0 {
		fmt.Println("this flag needs an argument")
		utils.ExitWhite(usageStr)
	}

	err := action.Exec(os.Stdout, arg)
	if err != nil {
		utils.Exit(err.Error())
	}
}

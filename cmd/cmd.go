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

	flag.Parse()

	if meow {
		fmt.Println("meow meow meow")
	}

	if sync {
		err := utils.Sync()
		if err != nil {
			panic(err)
		}
	}
}

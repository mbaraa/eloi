package actions

import (
	"github.com/mbaraa/eloi/db"
	"io"
)

// Action is an operation that can be used when running the program
// any implementation will be called directly from the command line
type Action interface {
	// Exec executes the selected action, and prints action's output on the given io.Writer
	// or exits with an error
	Exec(output io.Writer, args ...any) error

	// NeedsRoot return true if the action needs root to be executed
	NeedsRoot() bool

	// HasArgs returns true if the action requires arguments to run
	HasArgs() bool
}

type ActionType int

const (
	DownloadReposCacheActionType ActionType = iota + 1
	EbuildSearchActionType
	EnableRepoActionType
)

func GetActionFactory(at ActionType) Action {
	ebuildDB, err := db.GetEbuildDB()
	if err != nil {
		return nil
	}

	switch at {
	case DownloadReposCacheActionType:
		return new(DownloadReposCacheAction)
	case EbuildSearchActionType:
		return NewEbuildSearchAction(ebuildDB)
	case EnableRepoActionType:
		return new(EnableRepoAction)
	default:
		return nil
	}
}

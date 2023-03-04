package actions

import "io"

// Action is an operation that can be used when running the program
// any implementation will be called directly from the command line
type Action interface {
	// Exec executes the selected action, and prints action's output on the given io.Writer
	// or exits with an error
	Exec(output io.Writer, args ...string) error

	// NeedsRoot return true if the action needs root to be executed
	NeedsRoot() bool

	// HasArgs returns true if the action requires arguments to run
	HasArgs() bool
}

type ActionType int

const (
	DownloadReposCacheActionType ActionType = iota + 1
)

func GetActionFactory(at ActionType) Action {
	switch at {
	case DownloadReposCacheActionType:
		return new(DownloadReposCacheAction)
	default:
		return nil
	}
}

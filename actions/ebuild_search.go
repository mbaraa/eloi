package actions

import (
	"io"

	"github.com/mbaraa/eloi/db"
)

var _ Action = new(EbuildSearchAction)

type EbuildSearchAction struct {
	ebuildDB *db.EbuildDB
	output   io.Writer
}

func NewEbuildSearchAction(db *db.EbuildDB) *EbuildSearchAction {
	return &EbuildSearchAction{ebuildDB: db}
}

func (s *EbuildSearchAction) Exec(output io.Writer, args ...any) error {
	s.output = output
	return s.findEbuild(args[0].(string))
}

func (s *EbuildSearchAction) NeedsRoot() bool {
	// TODO
	// separate search from install
	return true
}

func (s *EbuildSearchAction) HasArgs() bool {
	return true
}

func (s *EbuildSearchAction) findEbuild(name string) error {
	name = "%" + name + "%"
	result, err := s.ebuildDB.GetByConds("name like ? OR group_name like ?", name, name)
	if err != nil {
		return err
	}
	return new(EbuildInstallAction).Exec(s.output, result)
}

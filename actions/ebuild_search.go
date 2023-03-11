package actions

import (
	"fmt"
	"io"
	"strings"

	"github.com/mbaraa/eloi/db"
	"github.com/mbaraa/eloi/models"
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
	var err error
	var result []models.Ebuild

	if slashIndex := strings.Index(name, "/"); slashIndex != -1 {
		groupName := name[:slashIndex]
		name = name[slashIndex+1:]
		result, err = s.ebuildDB.GetByConds("name like ? AND group_name like ?", likeQuery(name), likeQuery(groupName))
	} else {
		result, err = s.ebuildDB.GetByConds("name like ? OR group_name like ?", likeQuery(name), likeQuery(name))
	}
	if err != nil {
		return err
	}

	return new(EbuildInstallAction).Exec(s.output, result)
}

func likeQuery(s string) string {
	return fmt.Sprintf("%%%s%%", s)
}

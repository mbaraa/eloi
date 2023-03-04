package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mbaraa/eloi/cli/cfmt"
	"github.com/mbaraa/eloi/cli/templates"
	"github.com/mbaraa/eloi/globals"
	"github.com/mbaraa/eloi/models"
)

var _ Action = new(EbuildSearchAction)

type EbuildSearchAction struct {
	output io.Writer
}

func (s *EbuildSearchAction) Exec(output io.Writer, args ...any) error {
	s.output = output
	err := s.loadLocalOverlays()
	if err != nil {
		return err
	}
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

type PackageEntity struct {
	Display, FullName string
}

func (s *EbuildSearchAction) findEbuild(name string) error {
	allMatchingEbuilds := make([]string, 0)

	for _name := range globals.Ebuilds {
		if strings.Contains(strings.ToLower(_name), strings.ToLower(name)) {
			allMatchingEbuilds = append(allMatchingEbuilds, globals.EbuildsWithNamesOnly[strings.ToLower(_name[strings.Index(_name, "/")+1:])]...)
		}
	}

	dubs := make(map[string]bool)
	for _, matchingEbuild := range allMatchingEbuilds {
		dubs[matchingEbuild] = true
	}

	allMatchingEbuilds = make([]string, 0)
	for name := range dubs {
		allMatchingEbuilds = append(allMatchingEbuilds, name)
	}

	ebuildsToDisplay := make([]PackageEntity, 0)

	for _, ebuild := range allMatchingEbuilds {
		name := ebuild[strings.Index(ebuild, "/")+1:]
		group := ebuild[:strings.Index(ebuild, "/")]

		versions := make(map[string]models.Ebuild)
		for _, ebuildVersion := range globals.Ebuilds[group+"/"+name] {
			versions[ebuildVersion.Version] = *ebuildVersion
		}

		ebuildsToDisplay = append(ebuildsToDisplay, PackageEntity{templates.EbuildTemplate(versions), group + "/" + name})
	}

	for i, ebuild := range ebuildsToDisplay {
		_, err := s.output.Write([]byte(fmt.Sprintf("(%s) %s\n", cfmt.Magenta().Sprint(i+1), ebuild.Display)))
		if err != nil {
			return err
		}
	}

	return new(EbuildInstallAction).Exec(s.output, ebuildsToDisplay)
}

func (s *EbuildSearchAction) loadLocalOverlays() error {
	ebuildsFile, err := os.Open(ebuildsFilePath)
	if os.IsNotExist(err) {
		return errors.New("local repos are not synced, run with --sync to sync them")
	}
	defer ebuildsFile.Close()

	err = json.NewDecoder(ebuildsFile).Decode(&globals.Ebuilds)
	if err != nil {
		return err
	}

	ebuildsNamesFile, err := os.Open(ebuildsNamesFilePath)
	if err != nil {
		return err
	}
	defer ebuildsNamesFile.Close()

	err = json.NewDecoder(ebuildsNamesFile).Decode(&globals.EbuildsWithNamesOnly)
	if err != nil {
		return err
	}

	for name, ebuild := range globals.EbuildsWithNamesOnly {
		if name != strings.ToLower(name) {
			globals.EbuildsWithNamesOnly[strings.ToLower(name)] = ebuild
			delete(globals.EbuildsWithNamesOnly, name)
		}
	}

	return nil
}

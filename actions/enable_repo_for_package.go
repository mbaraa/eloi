package actions

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mbaraa/eloi/models"
	"github.com/mbaraa/eloi/utils/gentoolkit/depgraph"
)

var _ Action = new(EnableRepoForPackageAction)

type EnableRepoForPackageAction struct {
	output                     io.Writer
	ebuild                     models.Ebuild
	providerIndex              int
	packageMaskDirectoryPath   string
	packageMaskFilePath        string
	packageUnmaskDirectoryPath string
	packageUnmaskFilePath      string
}

func NewEnableRepoForPackageAction() *EnableRepoForPackageAction {
	return new(EnableRepoForPackageAction).init()
}

func (e *EnableRepoForPackageAction) init() *EnableRepoForPackageAction {
	e.packageMaskDirectoryPath = "/etc/portage/package.mask/"
	e.packageMaskFilePath = e.packageMaskDirectoryPath + "maskedByEloi"
	e.packageUnmaskDirectoryPath = "/etc/portage/package.unmask/"
	e.packageUnmaskFilePath = e.packageUnmaskDirectoryPath + "unmaskedByEloi"
	return e
}

func (e *EnableRepoForPackageAction) Exec(output io.Writer, args ...any) error {
	e.ebuild = args[0].(models.Ebuild)
	e.providerIndex = args[1].(int)

	e.output = output
	return e.enableRepo()
}

func (e *EnableRepoForPackageAction) NeedsRoot() bool {
	return true
}

func (e *EnableRepoForPackageAction) HasArgs() bool {
	return true
}

func (e *EnableRepoForPackageAction) enableRepo() error {
	enableRepoAction := new(EnableRepoAction)
	if !enableRepoAction.isRepoEnabled(e.ebuild.ExtraData[e.providerIndex].OverlayName) {
		err := enableRepoAction.Exec(e.output, e.ebuild.ExtraData[e.providerIndex].OverlayName)
		if err != nil {
			return err
		}
		err = e.maskAllPackages()
		if err != nil {
			return err
		}
	}
	return e.unmaskDeps()
}

func (e *EnableRepoForPackageAction) maskAllPackages() error {
	_ = os.Mkdir(e.packageMaskDirectoryPath, 0755)
	maskFile, err := os.OpenFile(e.packageMaskFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		if os.IsNotExist(err) {
			maskFile, err = os.Create(e.packageMaskFilePath)
			_, _ = maskFile.WriteString("# created by eloi\n\n")
		}
	}
	if err != nil {
		return err
	}
	defer maskFile.Close()

	_, _ = fmt.Fprintf(maskFile, "*/*::%s\n", e.ebuild.ExtraData[e.providerIndex].OverlayName)
	return nil
}

func (e *EnableRepoForPackageAction) unmaskDeps() error {
	_ = os.Mkdir(e.packageUnmaskDirectoryPath, 0755)
	unmaskFile, err := os.OpenFile(e.packageUnmaskFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		if os.IsNotExist(err) {
			unmaskFile, err = os.Create(e.packageUnmaskFilePath)
			_, _ = unmaskFile.WriteString("# created by eloi\n\n")
		}
	}
	if err != nil {
		return err
	}
	defer unmaskFile.Close()

	_, _ = fmt.Fprintf(unmaskFile, "=%s-%s::%s\n", e.ebuild.FullName(), e.ebuild.ExtraData[e.providerIndex].Version, e.ebuild.ExtraData[e.providerIndex].OverlayName)

	deps := depgraph.GetDeps(e.ebuild.FullName() + "-" + e.ebuild.ExtraData[e.providerIndex].Version)
	_, _ = fmt.Fprintf(unmaskFile, "# dependency for %s\n", e.ebuild.FullName()+"-"+e.ebuild.ExtraData[e.providerIndex].Version)
	overlayName := e.ebuild.ExtraData[e.providerIndex].OverlayName
	for _, dep := range deps {
		if strings.Contains(dep, "[") {
			dep = dep[:strings.Index(dep, "[")]
		}
		// Skip deps already installed from another repo to avoid conflicts
		if isInstalledFromOtherRepo(e.ebuild.GroupName, e.ebuild.Name, overlayName) {
			continue
		}
		_, _ = fmt.Fprintf(unmaskFile, "%s::%s\n", dep, overlayName)
	}

	return nil
}

// isInstalledFromOtherRepo checks if a package is installed from a repo other than the given overlay
func isInstalledFromOtherRepo(groupName, pkgName, overlayName string) bool {
	data, err := os.ReadFile(fmt.Sprintf("/var/db/pkg/%s/%s/repository", groupName, pkgName))
	if err != nil {
		return false
	}
	repo := strings.TrimSpace(string(data))
	return repo != "" && repo != overlayName
}

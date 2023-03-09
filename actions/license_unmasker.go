package actions

import (
	"errors"
	"fmt"
	"github.com/mbaraa/eloi/cli/cfmt"
	"github.com/mbaraa/eloi/models"
	"io"
	"os"
	"strings"
)

type LicenseUnmaskAction struct {
	output                      io.Writer
	packageLicenseDirectoryPath string
	packageLicenseFilePath      string
	opensourceLicenses          map[string]bool
	ebuild                      models.Ebuild
	providerIndex               int
}

func NewLicenceUnmaskAction() *LicenseUnmaskAction {
	return new(LicenseUnmaskAction).init()
}

func (l *LicenseUnmaskAction) init() *LicenseUnmaskAction {
	l.packageLicenseDirectoryPath = "/etc/portage/package.license"
	l.packageLicenseFilePath = l.packageLicenseDirectoryPath + "/acceptedByEloi"
	l.opensourceLicenses =
		// all in lowercase for easy fetching
		map[string]bool{
			"academic free license 2.0":        true,
			"academic free license":            true,
			"adaptive public license":          true,
			"afl 2.1":                          true,
			"afl 3.0":                          true,
			"afl":                              true,
			"agpl 3.0":                         true,
			"agpl":                             true,
			"apache license 2.0":               true,
			"apache license":                   true,
			"apsl 1.0":                         true,
			"apsl 1.1":                         true,
			"apsl 1.2":                         true,
			"apsl 2.0":                         true,
			"apsl":                             true,
			"artistic license 1.0 clarified":   true,
			"artistic license 1.0":             true,
			"artistic license 2.0":             true,
			"artistic license":                 true,
			"bsd 2-clause license":             true,
			"bsd 3-clause license":             true,
			"bsd":                              true,
			"bsl 1.0":                          true,
			"bsl":                              true,
			"cc-by 4.0":                        true,
			"cc-by-sa 4.0":                     true,
			"cc0 1.0":                          true,
			"cddl 1.0":                         true,
			"cddl 1.1":                         true,
			"cnri python license":              true,
			"cpal 1.0":                         true,
			"cpal":                             true,
			"cua opl 1.0":                      true,
			"ecl 1.0":                          true,
			"ecl 2.0":                          true,
			"ecl":                              true,
			"efl 1.0":                          true,
			"efl 2.0":                          true,
			"efl":                              true,
			"eu datagrid software license":     true,
			"eupl 1.0":                         true,
			"eupl 1.1":                         true,
			"eupl":                             true,
			"expat":                            true,
			"fair license":                     true,
			"frameworx 1.0":                    true,
			"frameworx":                        true,
			"gpl":                              true,
			"gpl-2":                            true,
			"gpl-3":                            true,
			"hpnd":                             true,
			"ibm public license 1.0":           true,
			"ibm public license":               true,
			"icu":                              true,
			"ipl 1.0":                          true,
			"isc":                              true,
			"json":                             true,
			"lgpl":                             true,
			"lgpl-2.0":                         true,
			"lgpl-2.1":                         true,
			"lgpl-3":                           true,
			"libpng":                           true,
			"lppl-1.2":                         true,
			"lppl-1.3":                         true,
			"lppl-1.3c":                        true,
			"mit":                              true,
			"mpl":                              true,
			"mpl-1.0":                          true,
			"mpl-1.1":                          true,
			"mpl-2.0":                          true,
			"ms-pl":                            true,
			"ms-rl":                            true,
			"ncsa":                             true,
			"ngpl":                             true,
			"nosl":                             true,
			"npl 1.0":                          true,
			"npl 1.1":                          true,
			"npl":                              true,
			"oclc research public license 2.0": true,
			"oclc research public license":     true,
			"ofl 1.1":                          true,
			"ofl":                              true,
			"openssl":                          true,
			"php license 3.0":                  true,
			"php license":                      true,
			"postgresql license":               true,
			"psf license":                      true,
			"python license 2.0":               true,
			"python license":                   true,
			"qpl 1.0":                          true,
			"qpl":                              true,
			"rpl 1.1":                          true,
			"rpl":                              true,
			"ruby license":                     true,
			"sgi free software license b":      true,
			"sil open font license 1.1":        true,
			"sil open font license":            true,
			"sleepycat":                        true,
			"spl 1.0":                          true,
			"spl":                              true,
			"unlicense":                        true,
			"uoi-ncsa":                         true,
			"vsl 1.0":                          true,
			"vsl":                              true,
			"w3c software notice and license":  true,
			"wtfpl":                            true,
			"x11":                              true,
			"xfree86 1.1":                      true,
			"xfree86":                          true,
			"zlib":                             true,
		}
	return l
}

func (l *LicenseUnmaskAction) Exec(output io.Writer, args ...any) error {
	l.ebuild = args[0].(models.Ebuild)
	l.providerIndex = args[1].(int)
	if l.providerIndex < 0 || l.providerIndex > len(l.ebuild.ExtraData)-1 {
		return errors.New("invalid choice")
	}
	l.output = output
	return l.promptAcceptLicense()
}

func (l *LicenseUnmaskAction) NeedsRoot() bool {
	return true
}

func (l *LicenseUnmaskAction) HasArgs() bool {
	return true
}

func (l *LicenseUnmaskAction) promptAcceptLicense() error {
	_, _ = l.output.Write([]byte("Updating enabled licenses...\n"))
	err := l.acceptEbuildLicense()
	if err != nil {
		return err
	}
	_, _ = l.output.Write([]byte("Done ✓\n"))
	return nil
}

func (l *LicenseUnmaskAction) acceptEbuildLicense() error {
	licenses := l.getNonOpensourceLicenses()
	if len(licenses) == 0 {
		return nil
	}

	_, _ = cfmt.Bold().Fprintf(l.output, "%s requires the following license `%s` to be accepted, do you want to proceed? [%s/%s]",
		cfmt.Magenta().Bold().Sprint(l.ebuild.FullName()),
		cfmt.Yellow().Bold().Sprint(licenses),
		cfmt.Green().Bold().Sprint("Yes"),
		cfmt.Red().Bold().Sprint("No"))

	prompt := ""
	_, _ = fmt.Scan(&prompt)

	switch strings.ToLower(prompt) {
	case "no", "n":
		return errors.New("user has canceled the operation")
	case "yes", "y":
		_ = os.Mkdir(l.packageLicenseDirectoryPath, 0755)
		licenseFile, err := os.OpenFile(l.packageLicenseFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			if os.IsNotExist(err) {
				licenseFile, err = os.Create(l.packageLicenseFilePath)
				_, _ = licenseFile.WriteString("# created by eloi\n\n")
			}
		}
		if err != nil {
			return err
		}
		defer licenseFile.Close()
		_, err = fmt.Fprintf(licenseFile, "%s::%s %s\n", l.ebuild.FullName(), l.ebuild.ExtraData[l.providerIndex].OverlayName, licenses)
	}

	return nil
}

func (l *LicenseUnmaskAction) getNonOpensourceLicenses() string {
	licenses := new(strings.Builder)
	for _, license := range strings.Split(l.ebuild.ExtraData[l.providerIndex].License, " ") {
		if !l.opensourceLicenses[strings.ToLower(license)] {
			licenses.WriteString(license + " ")
			delete(l.opensourceLicenses, strings.ToLower(license))
		}
	}
	return strings.TrimSpace(licenses.String())
}

package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/mbaraa/eloi/models"
)

const (
	packageUseDirectoryPath     = "/etc/portage/package.use"
	packageUseFilePath          = packageUseDirectoryPath + "/unmaskedByEloi"
	packageLicenseDirectoryPath = "/etc/portage/package.license"
	packageLicenseFilePath      = packageLicenseDirectoryPath + "/acceptedByEloi"
)

var acceptedLicenses = []string{"GPL", "LGPL", "BSD", "MIT", "Apache", "Artistic"}

func AddPackageFlags(ebuildFullName, flags string) error {
	os.Mkdir(packageUseDirectoryPath, 0755)
	useFlagsFile, err := openFileOrCreateIfNoExist(packageUseFilePath)
	if err != nil {
		return err
	}
	defer useFlagsFile.Close()

	_, err = fmt.Fprintf(useFlagsFile, "%s %s\n", ebuildFullName, flags)
	return err
}

func AcceptPackageLicense(ebuild models.Ebuild) error {
	for _, license := range acceptedLicenses {
		if strings.Contains(ebuild.License, license) {
			return nil
		}
	}

	os.Mkdir(packageLicenseDirectoryPath, 0755)
	licenseFile, err := openFileOrCreateIfNoExist(packageLicenseFilePath)
	if err != nil {
		return err
	}
	defer licenseFile.Close()

	_, err = fmt.Fprintf(licenseFile, "%s/%s %s\n", ebuild.GroupName, ebuild.Name, ebuild.License)
	return err
}

func openFileOrCreateIfNoExist(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(path)
			f.WriteString("# created by eloi\n\n")
		}
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

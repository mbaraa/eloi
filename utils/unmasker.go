package utils

import (
	"fmt"
	"os"

	"github.com/mbaraa/eloi/models"
)

const (
	packageUseDirectoryPath     = "/etc/portage/package.use"
	packageUseFilePath          = packageUseDirectoryPath + "/unmaskedByEloi"
	packageLicenseDirectoryPath = "/etc/portage/package.license"
	packageLicenseFilePath      = packageLicenseDirectoryPath + "/ecceptedByEloi"
)

func AddPackageFlags(ebuildFullName, flags string) error {
	os.Mkdir(packageUseDirectoryPath, 0755)
	useFlagsFile, err := openFileOrCreateIfNoExist(packageUseFilePath)
	if err != nil {
		return err
	}
	defer useFlagsFile.Close()

	_, err = fmt.Fprintf(useFlagsFile, "%s %s", ebuildFullName, flags)
	return err
}

func AcceptPackageLicense(ebuild models.Ebuild) error {
	os.Mkdir(packageLicenseDirectoryPath, 0755)
	licenseFile, err := openFileOrCreateIfNoExist(packageLicenseFilePath)
	if err != nil {
		return err
	}
	defer licenseFile.Close()

	_, err = fmt.Fprintf(licenseFile, "%s/%s %s", ebuild.GroupName, ebuild.Name, ebuild.License)
	return err
}

func openFileOrCreateIfNoExist(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(path)
		}
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

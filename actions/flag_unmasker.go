package actions

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/mbaraa/eloi/cli/cfmt"
	"github.com/mbaraa/eloi/models"
)

type FlagsUnmaskAction struct {
	output                  io.Writer
	packageUseDirectoryPath string
	packageUseFilePath      string
	ebuild                  models.Ebuild
	providerIndex           int
}

func NewFlagsUnmaskAction() *FlagsUnmaskAction {
	return new(FlagsUnmaskAction).init()
}

func (f *FlagsUnmaskAction) init() *FlagsUnmaskAction {
	f.packageUseDirectoryPath = "/etc/portage/package.use"
	f.packageUseFilePath = f.packageUseDirectoryPath + "/unmaskedByEloi"
	return f
}

func (f *FlagsUnmaskAction) Exec(output io.Writer, args ...any) error {
	f.ebuild = args[0].(models.Ebuild)
	f.providerIndex = args[1].(int)
	if f.providerIndex < 0 || f.providerIndex > len(f.ebuild.ExtraData)-1 {
		return errors.New("invalid choice")
	}
	f.output = output
	return f.addPackageFlags()
}

func (f *FlagsUnmaskAction) NeedsRoot() bool {
	return true
}

func (f *FlagsUnmaskAction) HasArgs() bool {
	return true
}

func (f *FlagsUnmaskAction) addPackageFlags() error {
	_ = os.Mkdir(f.packageUseDirectoryPath, 0755)
	useFlagsFile, err := os.OpenFile(f.packageUseFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		if os.IsNotExist(err) {
			useFlagsFile, err = os.Create(f.packageUseFilePath)
			_, _ = useFlagsFile.WriteString("# created by eloi\n\n")
		}
	}
	if err != nil {
		return err
	}
	defer useFlagsFile.Close()

	flags := removeFlagsDuplicates(f.ebuild.ExtraData[f.providerIndex].Flags)

	_, _ = fmt.Fprintf(f.output, "%s has the following use flags:\n", cfmt.Magenta().Sprint(f.ebuild.FullName()))
	for index, flag := range flags {
		_, _ = fmt.Fprintf(f.output, "(%s) %s\n", cfmt.Magenta().Sprint(index+1), flag)
	}

	prompt := cfmt.Green().Sprint("==>")
	_, _ = f.output.Write([]byte(fmt.Sprintf("%s Flag/s to enable %s \n%s ", prompt, cfmt.Bold().Sprint("e.g. (1, 2, 0 for all, or -1 for none)"), prompt)))

	selectionsInput, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	selectionsIndices := make([]int, 0)
	selectAll := false
	for _, selection := range bytes.Split(selectionsInput, []byte(" ")) {
		intSelection, err := strconv.ParseInt(string(selection), 10, 32)
		if err != nil {
			intSelection = -2
		}
		if intSelection > 0 && int(intSelection) <= len(flags) {
			selectionsIndices = append(selectionsIndices, int(intSelection)-1)
		}
		if intSelection == 0 {
			selectAll = true
		}
		if intSelection == -1 {
			selectionsIndices = []int{}
			break
		}
	}
	if selectAll {
		selectionsIndices = make([]int, 0)
		for i := 0; i < len(flags); i++ {
			selectionsIndices = append(selectionsIndices, i)
		}
	}
	selectedFlags := new(strings.Builder)
	for _, flagIndex := range selectionsIndices {
		selectedFlags.WriteString(flags[flagIndex] + " ")
	}

	if len(selectionsIndices) == 0 {
		return nil
	}

	_, err = fmt.Fprintf(useFlagsFile, "%s::%s %s\n", f.ebuild.FullName(),
		f.ebuild.ExtraData[f.providerIndex].OverlayName, selectedFlags.String())
	return err
}

func removeFlagsDuplicates(flags string) []string {
	flagsSplit := strings.Split(flags, " ")
	dubs := map[string]bool{}
	for _, flag := range flagsSplit {
		dubs[flag] = true
	}
	flagsSplit = make([]string, 0)
	for key := range dubs {
		flagsSplit = append(flagsSplit, key)
	}
	sort.Slice(flagsSplit, func(i, j int) bool {
		return flagsSplit[i] < flagsSplit[j]
	})
	return flagsSplit
}

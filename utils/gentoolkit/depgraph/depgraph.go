package depgraph

// ldflags was generated using python3.10-config --ldflags --embed

// #cgo LDFLAGS: -L/usr/lib64 -lpython3.10 -ldl  -lm
// #include "./depgraph.h"
import "C"
import (
	"strings"
	"unsafe"
)

func GetDeps(pkg string) []string {
	return fixDepsVersions(getDepgraph(pkg))
}

func fixDepsVersions(deps []string) []string {
	fixed := make([]string, 0)
	for _, dep := range deps {
		if strings.Contains(dep, ":") {
			dep = dep[:strings.Index(dep, ":")]
		}
		fixed = append(fixed, dep)
	}
	return removeDuplicates(fixed)
}

func getDepgraph(pkg string) []string {
	pkgName := C.CString(pkg)
	defer C.free(unsafe.Pointer(pkgName))
	deps := C.get_pkg_depgraph(pkgName)
	depsStr := C.GoString(deps)
	C.cleanup()

	return strings.Split(depsStr, " ")
}

func removeDuplicates(s []string) []string {
	dubs := map[string]bool{}
	for _, str := range s {
		dubs[str] = true
	}
	nonDub := make([]string, 0)
	for str := range dubs {
		nonDub = append(nonDub, str)
	}
	return nonDub
}

package files

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetVirtualRoot(paths ...string) (string, error) {
	vroot, err := GetVirtualPath(VirtualRoot, paths...)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(vroot, 0766)
	if err != nil {
		return "", err
	}
	return vroot, nil
}

func GetVirtualPath(virtualRoot string, paths ...string) (string, error) {
	absRoot, err := filepath.Abs(virtualRoot)
	if err != nil {
		return "", err
	}

	p := path.Join(paths...)
	finalPath := filepath.Join(absRoot, p)

	dir := path.Dir(finalPath)
	if !strings.HasPrefix(dir, absRoot) && finalPath != absRoot {
		return "", fmt.Errorf("Invalid path: %s", p)
	}

	return finalPath, nil
}

func PathJoin(ps ...string) string {
	return path.Join(ps...)
}

func PathSplit(p string) []string {
	return strings.Split(p, string(filepath.Separator))
}

func PathParent(p string) string {
	return path.Dir(p)
}

func PathName(p string) string {
	return path.Base(p)
}

package loader

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

func isURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func joinURL(base, p string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	if path.IsAbs(p) {
		u.Path = p
		return u.String(), nil
	}

	ref, err := url.Parse(p)
	if err != nil {
		return "", err
	}
	return u.ResolveReference(ref).String(), nil
}

func joinFinalPath(root, currentDir, file string) (string, error) {
	if isURL(file) {
		return file, nil
	}

	if isURL(currentDir) {
		return joinURL(currentDir, file)
	}

	if path.IsAbs(file) {
		return path.Join(root, file), nil
	}

	return path.Join(currentDir, file), nil
}

func getFinalPath(root, currentDir, file string, fetch bool) (string, error) {
	finalPath, err := joinFinalPath(root, currentDir, file)
	if err != nil {
		return "", err
	}

	if isURL(finalPath) {
		if !fetch {
			return "", fmt.Errorf("fetching is not allowed for %s", finalPath)
		}
		return finalPath, nil
	}

	dir := path.Dir(finalPath)
	if !strings.HasPrefix(dir, root) && dir != root {
		return "", fmt.Errorf("access denied: %s is outside root %s", dir, root)
	}

	return finalPath, nil
}

func getKeyPath(root, p string) string {
	if isURL(p) {
		return p
	}
	p = strings.TrimPrefix(p, root)
	p = strings.ReplaceAll(p, "\\", "/")
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}

func getPathDir(p string) string {
	if !isURL(p) {
		return path.Dir(p)
	}
	u, _ := url.Parse(p)
	u.Path = path.Dir(u.Path)
	return u.String()
}

func setExtension(p string, ext string) string {
	if isURL(p) {
		return p
	}

	base := path.Base(p)
	if path.Ext(base) != "" {
		return p
	}
	return p + ext
}

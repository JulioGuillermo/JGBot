package loader

import (
	"path/filepath"
	"regexp"
	"strings"
)

func LoadCode(root, file string, fetch bool) (map[string]string, error) {
	return LoadModule(root, root, file, fetch)
}

func LoadModule(root, dir, file string, fetch bool) (map[string]string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	if !isURL(dir) {
		dir, err = filepath.Abs(dir)
		if err != nil {
			return nil, err
		}
	}

	codes := make(map[string]string)
	_, err = loadRecursive(root, dir, file, codes, fetch)
	if err != nil {
		return nil, err
	}

	return codes, nil
}

func loadRecursive(root, currentDir, file string, codes map[string]string, allowFetch bool) (string, error) {
	finalPath, err := getFinalPath(root, currentDir, file, allowFetch)
	if err != nil {
		return "", err
	}
	finalPath = setExtension(finalPath, ".js")

	key := getKeyPath(root, finalPath)
	if _, exists := codes[key]; exists {
		return key, nil
	}

	var code string
	if isURL(finalPath) {
		code, err = fetchCode(finalPath)
	} else {
		code, err = readFileCode(finalPath)
	}
	if err != nil {
		return "", err
	}

	newDir := getPathDir(finalPath)

	re := regexp.MustCompile(`import\s*(?:[^\w][\w\s{},*]*[^\w]\s*from\s*)?['"]([^'"]+)['"]`)
	matches := re.FindAllStringSubmatch(code, -1)

	replacements := make(map[string]string)
	for _, match := range matches {
		importPath := match[1]

		key, err := loadRecursive(
			root,
			newDir,
			importPath,
			codes,
			allowFetch,
		)
		if err != nil {
			return "", err
		}

		replacements[importPath] = key
	}

	code = re.ReplaceAllStringFunc(code, func(match string) string {
		submatch := re.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}
		pathArg := submatch[1]
		if newKey, ok := replacements[pathArg]; ok {
			return strings.Replace(match, pathArg, newKey, 1)
		}
		return match
	})

	codes[key] = code
	return key, nil
}

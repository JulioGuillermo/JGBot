package loader

import (
	"path/filepath"
	"regexp"
	"strings"
)

type Code struct {
	Key  string
	Code string
}

func LoadCode(root, file string, fetch bool) ([]Code, error) {
	return LoadModule(root, root, file, fetch)
}

func LoadModule(root, dir, file string, fetch bool) ([]Code, error) {
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

	codes := make([]Code, 0)
	_, err = loadRecursive(root, dir, file, &codes, fetch)
	if err != nil {
		return nil, err
	}

	return codes, nil
}

func loadRecursive(root, currentDir, file string, codes *[]Code, allowFetch bool) (string, error) {
	finalPath, err := getFinalPath(root, currentDir, file, allowFetch)
	if err != nil {
		return "", err
	}
	finalPath = setExtension(finalPath, ".js")

	key := getKeyPath(root, finalPath)
	if code := GetCode(*codes, key); code != nil {
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
	// log.Info("Loading module...", "key", key, "path", finalPath)

	newDir := getPathDir(finalPath)

	re := regexp.MustCompile(`(?:import\s*(?:[^\w][\w\s{},*]*[^\w]\s*from\s*)?|export.*from\s*)['"]([^'"]+)['"]`)
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

	*codes = append(*codes, Code{
		Key:  key,
		Code: code,
	})
	return key, nil
}

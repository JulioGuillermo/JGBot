package loader

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func LoadCode(root string) (map[string]string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	files, err := loadFiles(root)
	if err != nil {
		return nil, err
	}

	codes, err := readCode(files)
	if err != nil {
		return nil, err
	}

	codes = removePath(root, codes)
	for file, _ := range codes {
		fmt.Println(file)
	}
	return codes, nil
}

func removePath(root string, codes map[string]string) map[string]string {
	cleanCode := map[string]string{}
	for file, code := range codes {
		cleanName := strings.TrimPrefix(file, root)
		cleanCode[cleanName] = code
	}
	return cleanCode
}

func readCode(files []string) (map[string]string, error) {
	codes := map[string]string{}
	for _, file := range files {
		code, err := readFileCode(file)
		if err != nil {
			return nil, err
		}
		codes[file] = code
	}
	return codes, nil
}

func readFileCode(path string) (string, error) {
	code, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(code), nil
}

func loadFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return loadDirFiles(path)
	}

	return []string{path}, nil
}

func loadDirFiles(dir string) ([]string, error) {
	elements, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string

	for _, element := range elements {
		if element.IsDir() {
			subFiles, err := loadDirFiles(path.Join(dir, element.Name()))
			if err != nil {
				return nil, err
			}

			files = append(files, subFiles...)
			continue
		}

		if strings.HasSuffix(element.Name(), ".js") {
			files = append(files, path.Join(dir, element.Name()))
		}
	}

	return files, nil
}

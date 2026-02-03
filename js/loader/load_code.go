package loader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

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

func readFileCode(path string) (string, error) {
	code, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(code), nil
}

func fetchCode(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Fail to fetch the code: status code %d", resp.StatusCode)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

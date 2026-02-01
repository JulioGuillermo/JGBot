package tools

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

func CreateParentDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "" || dir == "." || dir == "/" {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

func ReadFile(filePath string) ([]byte, error) {
	err := CreateParentDir(filePath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func WriteFile(filePath string, bytes []byte) error {
	err := CreateParentDir(filePath)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func ReadJSONFile(filePath string, v any) error {
	bytes, err := ReadFile(filePath)
	if err != nil && os.IsNotExist(err) {
		return WriteJSONFile(filePath, v)
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func WriteJSONFile(filePath string, v any) error {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return WriteFile(filePath, jsonBytes)
}

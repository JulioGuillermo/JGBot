package files

import (
	"os"
	"path"
)

func ReadDir(virtualRoot string, p string) ([]*FileInfo, error) {
	finalPath, err := GetVirtualPath(virtualRoot, p)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(finalPath)
	if err != nil {
		return nil, CheckError(virtualRoot, err)
	}

	result := []*FileInfo{}
	for _, file := range files {
		fileInfo := &FileInfo{
			Name:  file.Name(),
			Size:  0,
			IsDir: file.IsDir(),
		}
		if !fileInfo.IsDir {
			stat, err := os.Stat(path.Join(finalPath, file.Name()))
			if err != nil {
				return nil, CheckError(virtualRoot, err)
			}
			fileInfo.Size = stat.Size()
		}
		result = append(result, fileInfo)
	}

	return result, nil
}

func CreateDir(virtualRoot string, p string) error {
	finalPath, err := GetVirtualPath(virtualRoot, p)
	if err != nil {
		return err
	}

	err = os.MkdirAll(finalPath, 0755)
	if err != nil {
		return CheckError(virtualRoot, err)
	}

	return nil
}

func DeleteDir(virtualRoot string, p string) error {
	finalPath, err := GetVirtualPath(virtualRoot, p)
	if err != nil {
		return err
	}

	err = os.RemoveAll(finalPath)
	if err != nil {
		return CheckError(virtualRoot, err)
	}

	return nil
}

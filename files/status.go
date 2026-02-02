package files

import "os"

type FileInfo struct {
	Name  string
	Size  int64
	IsDir bool
}

func GetFileInfo(virtualRoot string, path string) (*FileInfo, error) {
	finalPath, err := GetVirtualPath(virtualRoot, path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(finalPath)
	if err != nil {
		return nil, CheckError(virtualRoot, err)
	}

	fInfo := &FileInfo{
		Name:  info.Name(),
		Size:  0,
		IsDir: info.IsDir(),
	}

	if !fInfo.IsDir {
		fInfo.Size = info.Size()
	}

	return fInfo, nil
}

func Exists(virtualRoot string, path string) (bool, error) {
	finalPath, err := GetVirtualPath(virtualRoot, path)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(finalPath)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, CheckError(virtualRoot, err)
}

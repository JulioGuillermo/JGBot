package files

import "os"

func ReadFile(virtualRoot string, p string) ([]byte, error) {
	finalPath, err := GetVirtualPath(virtualRoot, p)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(finalPath)
	if err != nil {
		return nil, CheckError(virtualRoot, err)
	}

	return content, nil
}

func WriteFile(virtualRoot string, p string, content []byte) error {
	finalPath, err := GetVirtualPath(virtualRoot, p)
	if err != nil {
		return err
	}

	err = os.WriteFile(finalPath, content, 0755)
	if err != nil {
		return CheckError(virtualRoot, err)
	}

	return nil
}

func DeleteFile(virtualRoot string, p string) error {
	finalPath, err := GetVirtualPath(virtualRoot, p)
	if err != nil {
		return err
	}

	err = os.Remove(finalPath)
	if err != nil {
		return CheckError(virtualRoot, err)
	}

	return nil
}

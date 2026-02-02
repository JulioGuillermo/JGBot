package files

import (
	"errors"
	"strings"
)

func CheckError(virtualRoot string, err error) error {
	errMsg := strings.ReplaceAll(err.Error(), virtualRoot, "")
	return errors.New(errMsg)
}

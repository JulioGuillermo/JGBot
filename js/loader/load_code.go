package loader

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

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
		return "", fmt.Errorf("Fail to fetch the code from %s: status code %d", url, resp.StatusCode)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func GetCode(codes []Code, key string) *Code {
	for _, code := range codes {
		if code.Key == key {
			return &code
		}
	}
	return nil
}

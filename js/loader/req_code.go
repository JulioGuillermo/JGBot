package loader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

func GetReqCode(root, p string) (string, error) {
	p = getPath(root, p)

	if isURL(p) {
		return fetchCode(p)
	}

	return readFileCode(p)
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

func getPath(root, p string) string {
	if isURL(p) {
		return p
	}
	return path.Join(root, p)
}

func isURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

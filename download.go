package iplocation

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func downloadOrCache(cachePath string, u *url.URL, force bool) (string, time.Time, error) {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", time.Time{}, err
	}

	var isGzip bool
	if strings.HasSuffix(u.Path, ".gz") {
		isGzip = true
		u.Path = strings.TrimSuffix(u.Path, ".gz")
	}

	file := path.Join(cachePath, u.Host, u.Path)

	if !force {
		fi, err := os.Stat(file)
		if err == nil && fi.Size() != 0 {
			return file, fi.ModTime(), nil
		}
	}

	req.Header.Set("User-Agent", "wzshiming/iplocation")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = os.MkdirAll(path.Dir(file), 0700)
	if err != nil {
		return "", time.Time{}, err
	}

	f, err := os.OpenFile(file+".tmp", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", time.Time{}, err
	}
	defer f.Close()

	var body io.Reader = resp.Body

	if isGzip {
		r, err := gzip.NewReader(body)
		if err != nil {
			return "", time.Time{}, err
		}
		defer r.Close()

		body = r
	}

	_, err = io.Copy(f, body)
	if err != nil {
		return "", time.Time{}, err
	}

	err = os.Rename(file+".tmp", file)
	if err != nil {
		return "", time.Time{}, err
	}
	return file, time.Now(), nil
}

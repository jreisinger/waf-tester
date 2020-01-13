package target

import (
	"net/http"
	"time"
)

type Target struct {
	Scheme     string
	Host       string
	Path       string
	URL        string
	Err        error
	StatusCode int
}

func (t *Target) getStatusCode() {
	t.URL = t.Scheme + "://" + t.Host + "/" + t.Path

	req, err := http.NewRequest("GET", t.URL, nil)
	if err != nil {
		t.Err = err
		return
	}

	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		t.Err = err
		return
	}
	defer resp.Body.Close()

	t.StatusCode = resp.StatusCode

	return
}

func TargetChecker(ch chan Target, scheme string, host string, path string) {
	t := Target{Scheme: scheme, Host: host, Path: path}
	t.getStatusCode()
	ch <- t
}

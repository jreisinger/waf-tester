package main

import (
	"fmt"
	"net/http"
	_ "net/url"
	"os"
	"time"
)

type Target struct {
	Host       string
	Path       string
	StatusCode int
	Err        error
	Url        string
}

func (t *Target) getStatusCode() {
	t.Url = "http://" + t.Host + "/" + t.Path

	req, err := http.NewRequest("GET", t.Url, nil)
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

func TargetChecker(ch chan Target, host string, path string) {
	t := Target{Host: host, Path: path}
	t.getStatusCode()
	ch <- t
}

var paths = []string{
	"etc/passwd",
	"?page=/etc/passwd",
	"?exec=/bin/bash",
	"?id=1' or '1' = '1'",
	"?<script>",
}

func main() {
	hosts := os.Args[1:]
	ch := make(chan Target)

	for _, host := range hosts {
		for _, path := range paths {
			//path = url.PathEscape(path)
			go TargetChecker(ch, host, path)
		}
	}

	for range hosts {
		for range paths {
			t := <-ch
			if t.Err != nil {
				fmt.Printf("ERR  (%03.0f): %s\n", float64(t.StatusCode), t.Err)
			} else if t.StatusCode != 403 {
				fmt.Printf("FAIL (%03.0f): %s\n", float64(t.StatusCode), t.Url)
			} else {
				fmt.Printf("OK   (%03.0f): %s\n", float64(t.StatusCode), t.Url)
			}
		}
	}
}

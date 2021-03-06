package service

import (
	"io"
	"net/http"
)

func httpGet(url string, fn func(r *http.Request)) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	if fn != nil {
		fn(req)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	body := resp.Body
	defer body.Close()
	bs, err := io.ReadAll(body)
	if err != nil {
		return nil, nil, err
	}

	return bs, resp, nil
}

func githubToken(token string) func(*http.Request) {
	return func(r *http.Request) {
		if token != "" {
			r.Header.Add("Authorization", "Token "+token)
		}
	}
}

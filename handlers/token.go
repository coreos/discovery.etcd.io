package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/coreos/discovery/pkg/lockstring"
)

var (
	currentLeader lockstring.LockString
)

func init() {
	currentLeader.Set("127.0.0.1:4001")
}

func proxyRequest(r *http.Request) (*http.Response, error) {
	for i := 0; i <= 10; i++ {
		u := url.URL{
			Host: currentLeader.String(),
			Scheme: "http",
			Path: path.Join("v2", "keys", "_etcd", "registry", r.URL.Path),
			RawQuery: r.URL.RawQuery,
		}

		req, err := http.NewRequest(r.Method, u.String(), r.Body)
		if err != nil {
			return nil, err
		}

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		// Try again on the next host
		if resp.StatusCode == 307 && r.Method == "PUT" {
			u, err := resp.Location()
			if err != nil {
				return nil, err
			}
			currentLeader.Set(u.Host)
			continue
		}

		return resp, nil
	}

	return nil, errors.New("All attempts at proxying to etcd failed")
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := proxyRequest(r)
	if err != nil {
		log.Printf("Error making request: %v", err)
		http.Error(w, "", 500)
	}

	// Copy all of the headers, set the status code and copy the body
	for k, v := range resp.Header {
		for _, q := range v {
			w.Header().Add(k, q)
		}
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

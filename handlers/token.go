package handlers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/coreos/discovery.etcd.io/handlers/httperror"
	"github.com/coreos/etcd/clientv3"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	tokenCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "endpoint_token_requests_total",
			Help: "How many /token requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(tokenCounter)
}

var tokenCounter *prometheus.CounterVec

func (st *State) proxyRequest(r *http.Request, useV3 bool) (resp *http.Response, err error) {
	body, _ := ioutil.ReadAll(r.Body)
	key := path.Join("_etcd", "registry", r.URL.Path)

	eps := []string{st.getCurrentLeader()}
	for i := 0; i <= 10; i++ {
		if useV3 {
			c, err := clientv3.New(clientv3.Config{Endpoints: eps})
			if err != nil {
				return nil, err
			}
			defer c.Close()

			var result string
			switch r.Method {
			case http.MethodGet:
				var gresp *clientv3.GetResponse
				gresp, err = c.Get(context.Background(), key)
				if gresp != nil && len(gresp.Kvs) > 0 {
					result = string(gresp.Kvs[0].Value)
				}
			case http.MethodPut:
				_, err = c.Put(context.Background(), key, string(body))
			case http.MethodDelete:
				_, err = c.Delete(context.Background(), key)
			}
			if err == nil {
				resp = &http.Response{
					Status: "200 OK",
					Body:   ioutil.NopCloser(strings.NewReader(result)),
				}
				return resp, err
			}

			log.Printf("failed with %v", err)
			// sync endpoints, and try other endpoints
			if err = c.Sync(context.Background()); err != nil {
				return nil, err
			}
			eps = c.Endpoints()
		} else {
			u := url.URL{
				Scheme:   "http",
				Host:     eps[0],
				Path:     path.Join("v2", "keys", key),
				RawQuery: r.URL.RawQuery,
			}

			buf := bytes.NewBuffer(body)
			outreq, rerr := http.NewRequest(r.Method, u.String(), buf)
			if err != nil {
				return nil, rerr
			}
			copyHeader(outreq.Header, r.Header)

			client := http.Client{}
			resp, err = client.Do(outreq)
			if err != nil {
				return nil, err
			}

			// Try again on the next host
			if resp.StatusCode == 307 &&
				(r.Method == "PUT" || r.Method == "DELETE") {
				u, err := resp.Location()
				if err != nil {
					return nil, err
				}
				st.setCurrentLeader(u.Host)
				eps[0] = st.getCurrentLeader()
				continue
			}
		}

		return resp, nil
	}

	return nil, errors.New("All attempts at proxying to etcd failed")
}

// copyHeader copies all of the headers from dst to src.
func copyHeader(dst, src http.Header) {
	for k, v := range src {
		for _, q := range v {
			dst.Add(k, q)
		}
	}
}

func TokenHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	st := ctx.Value(stateKey).(*State)

	resp, err := st.proxyRequest(r, st.isV3())
	if err != nil {
		log.Printf("Error making request: %v", err)
		httperror.Error(w, r, "", 500, tokenCounter)
	}

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	tokenCounter.WithLabelValues(strconv.Itoa(resp.StatusCode), r.Method).Add(1)
}

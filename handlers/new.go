package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/coreos/discovery.etcd.io/handlers/httperror"
	"github.com/coreos/etcd/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

var newCounter *prometheus.CounterVec
var cfg *client.Config

func init() {
	newCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "endpoint_new_requests_total",
			Help: "How many /new requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(newCounter)
}

func generateCluster() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

func Setup() {
	cfg = &client.Config{
		Endpoints: []string{viper.GetString("etcd")},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	u, _ := url.Parse(viper.GetString("etcd"))
	currentLeader.Set(u.Host)
}

func setupToken(size int) (string, error) {
	token := generateCluster()
	if token == "" {
		return "", errors.New("Couldn't generate a token")
	}

	c, _ := client.New(*cfg)
	kapi := client.NewKeysAPI(c)

	key := path.Join("_etcd", "registry", token)

	resp, err := kapi.Create(context.Background(), path.Join(key, "_config", "size"), strconv.Itoa(size))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Couldn't setup state %v %v", resp, err))
	}

	return token, nil
}

func deleteToken(token string) error {
	c, _ := client.New(*cfg)
	kapi := client.NewKeysAPI(c)

	if token == "" {
		return errors.New("No token given")
	}

	_, err := kapi.Delete(
		context.Background(),
		path.Join("_etcd", "registry", token),
		&client.DeleteOptions{Recursive: true},
	)

	return err
}

func NewTokenHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	size := 3
	s := r.FormValue("size")
	if s != "" {
		size, err = strconv.Atoi(s)
		if err != nil {
			httperror.Error(w, r, err.Error(), http.StatusBadRequest, newCounter)
			return
		}
	}
	token, err := setupToken(size)

	if err != nil {
		log.Printf("setupToken returned: %v", err)
		httperror.Error(w, r, "Unable to generate token", 400, newCounter)
		return
	}

	log.Println("New cluster created", token)

	fmt.Fprintf(w, "%s/%s", bytes.TrimRight([]byte(viper.GetString("host")), "/"), token)
	newCounter.WithLabelValues("200", r.Method).Add(1)
}

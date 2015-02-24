package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/coreos/go-etcd/etcd"
)

func generateCluster() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

func setupToken(size int) (string, error) {
	token := generateCluster()
	if token == "" {
		return "", errors.New("Couldn't generate a token")
	}

	client := etcd.NewClient(nil)
	key := path.Join("_etcd", "registry", token)
	resp, err := client.CreateDir(key, 0)

	if err != nil || resp.Node == nil || resp.Node.Key != "/"+key || resp.Node.Dir != true {
		return "", errors.New(fmt.Sprintf("Couldn't setup state %v %v", resp, err))
	}

	resp, err = client.Create(path.Join(key, "_config", "size"), strconv.Itoa(size), 0)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Couldn't setup state %v %v", resp, err))
	}

	return token, nil
}

func deleteToken(token string) error {
	client := etcd.NewClient(nil)

	if token == "" {
		return errors.New("No token given")
	}

	_, err := client.Delete(path.Join("_etcd", "registry", token), true)

	return err
}

func getRootUrl() (*url.URL, error) {
	var root_url = os.Getenv("ROOT_URL")
	if root_url == "" {
		root_url = "https://discovery.etcd.io"
	}

	u, err := url.Parse(root_url)
	if len(u.Path) > 1 || ( len(u.Path) == 1 && u.Path != "/" ) {
		return u, errors.New(
			fmt.Sprintf("Expected URL without path (%v)", u.Path))
	}

	if u.RawQuery != "" {
		return u, errors.New(
			fmt.Sprintf("Expected URL without query (%v)", u.RawQuery))
	}

	if u.Fragment != "" {
		return u, errors.New(
			fmt.Sprintf("Expected URL without fragment (%v)", u.Fragment))
	}

	return u, err
}

func NewTokenHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	size := 3
	s := r.FormValue("size")
	if s != "" {
		size, err = strconv.Atoi(s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	token, err := setupToken(size)

	if err != nil {
		log.Printf("setupToken returned: %v", err)
		http.Error(w, "Unable to generate token", 400)
		return
	}

	log.Println("New cluster created", token)

	u, err := getRootUrl()
	if err != nil {
		log.Printf("getRootUrl returned: %v", err)
		http.Error(w, "Unable to generate URL with given token", 400)
		return
	}

	u.Path = token

	fmt.Fprintf(w, u.String())
}

package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"os"
	"github.com/coreos/go-etcd/etcd"
)

func getDomain() string {
	domain := os.Getenv("ETCD_DISCOVERY_DOMAIN")
	if domain != "" { 
		return domain
	}

	return "https://discovery.etcd.io/"
}

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

	domain := getDomain()

	fmt.Fprintf(w, domain+"/"+token)
}

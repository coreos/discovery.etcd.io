package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

func generateCluster() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(b)
}

func NewHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, generateCluster())
}

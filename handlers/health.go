package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	token, err := setupToken(0)

	if err != nil || token == "" {
		log.Printf("health failed to setupToken %v", err)
		http.Error(w, "health failed to setupToken", 400)
		return
	}

	err = deleteToken(token)
	if err != nil {
		log.Printf("health failed to deleteToken %v", err)
		http.Error(w, "health failed to deleteToken", 400)
		return
	}

	fmt.Fprintf(w, "OK")
}

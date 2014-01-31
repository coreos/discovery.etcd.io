package main

import (
	"fmt"
	"net/http"

	"github.com/coreos/go-systemd/activation"
	"github.com/gorilla/mux"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "OK")
}

func main() {
	listeners, err := activation.Listeners(true)
	if err != nil {
		panic(err)
	}

	if len(listeners) != 1 {
		panic("Unexpected number of socket activation fds")
	}

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/new", NewHandler)
	r.HandleFunc("/health", HealthHandler)
	http.Handle("/", r)

	http.Serve(listeners[0], nil)
}

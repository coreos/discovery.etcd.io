package main

import (
	"net/http"

	"github.com/coreos/go-systemd/activation"
	"github.com/gorilla/mux"
)

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
	http.Handle("/", r)
}

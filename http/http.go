package http

import (
	"net/http"
	"os"

	gorillaHandlers "github.com/gorilla/handlers"

	"github.com/coreos/discovery.etcd.io/handlers"
	"github.com/gorilla/mux"

)

func init() {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("/new", handlers.NewTokenHandler)
	r.HandleFunc("/health", handlers.HealthHandler)

	// Only allow exact tokens with GETs and PUTs
	r.HandleFunc("/{token:[a-f0-9]{32}}", handlers.TokenHandler).
		Methods("GET", "PUT")
	r.HandleFunc("/{token:[a-f0-9]{32}}/", handlers.TokenHandler).
		Methods("GET", "PUT")
	r.HandleFunc("/{token:[a-f0-9]{32}}/{machine}", handlers.TokenHandler).
		Methods("GET", "PUT", "DELETE")

	logH := gorillaHandlers.LoggingHandler(os.Stdout, r)

	http.Handle("/", logH)
}

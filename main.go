package main

import (
	"net/http"

	"github.com/ellied33/is-that-murphy/handlers"
)

func main() {
	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.VerifyHandler(w, r)
		case http.MethodPost:
			handlers.AddHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":8080", nil)
}
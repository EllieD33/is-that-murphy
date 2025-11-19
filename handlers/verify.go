package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ellied33/is-that-murphy/models"
	"github.com/ellied33/is-that-murphy/store"
)

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	value := r.URL.Query().Get("value")
	if value == "" {
		http.Error(w, "missing value", http.StatusBadRequest)
		return
	}

	if v, ok := store.IsVerified(value); ok {
		json.NewEncoder(w).Encode(v)
	} else {
		json.NewEncoder(w).Encode(map[string]string{
			"value": "value",
			"type": "not verfied",
		})
	}
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	var v models.VerifiedValue
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	store.Add(v)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}
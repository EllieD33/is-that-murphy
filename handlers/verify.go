package handlers

import (
	"encoding/json"
	"html"
	"net/http"

	"github.com/ellied33/is-that-murphy/models"
	"github.com/ellied33/is-that-murphy/store"
	"github.com/ellied33/is-that-murphy/utils"
)

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	const MaxInputLength = 1024

	value := r.URL.Query().Get("value")

	if value == "" {
		http.Error(w, "missing value", http.StatusBadRequest)
		return
	}
	if len(value) > MaxInputLength {
		http.Error(w, "input too long", http.StatusBadRequest)
		return
	}

	c := utils.Canonical(value)

	if v, ok := store.IsVerified(c); ok {
		json.NewEncoder(w).Encode(v)
	} else {
		json.NewEncoder(w).Encode(map[string]string{
			"value": html.EscapeString(value),
			"type":  "not verified",
		})
	}
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var v models.VerifiedValue
	if err := decoder.Decode(&v); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := v.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	store.Add(v)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

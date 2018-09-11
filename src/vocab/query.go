package main

import (
	"encoding/json"
	"net/http"
)

func query(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			http.Error(w, e.(error).Error(), http.StatusInternalServerError)
		}
	}()
	word := r.URL.Path[7:]
	qr, err := lookup(word)
	assert(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	assert(json.NewEncoder(w).Encode(qr))
}

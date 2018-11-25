package main

import (
	"encoding/json"
	"net/http"
)

func related(w http.ResponseWriter, r *http.Request) {
	rels, err := oc.QueryRelated(r.URL.Path[9:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rels)
}

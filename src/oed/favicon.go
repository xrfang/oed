package main

import (
	"net/http"
)

func favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	ico, _ := Asset("img/favicon.png")
	w.Write(ico)
}

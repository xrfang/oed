package main

import (
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		sendAsset(w, r.URL.Path)
		return
	}
	renderTemplate(w, "home.html", nil)
}

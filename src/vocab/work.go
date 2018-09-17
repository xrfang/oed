package main

import (
	"net/http"
)

func work(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "work.html", nil)
}

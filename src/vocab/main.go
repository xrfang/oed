package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	oed "oedcli"
	"os"
	"time"
)

var oc *oed.Client

func main() {
	ver := flag.Bool("version", false, "show version info")
	id := flag.String("id", "", "OED app_id")
	key := flag.String("key", "", "OED app_key")
	port := flag.String("port", "63300", "service port")
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}
	oc = oed.NewClient(*id, *key, 30)
	qr, err := oc.Query("ace")
	if err != nil {
		fmt.Println("query failed:", err)
	} else {
		je := json.NewEncoder(os.Stdout)
		je.SetIndent("", "    ")
		je.Encode(qr)
	}
	return
	http.HandleFunc("/", home)
	http.HandleFunc("/favicon.ico", favicon)
	svr := http.Server{
		Addr:         ":" + *port,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	svr.ListenAndServe()
}

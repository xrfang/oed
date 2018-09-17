package main

import (
	"flag"
	"fmt"
	"net/http"
	oed "oedcli"
	"os"
	"time"
)

var oc *oed.Client
var cache string

func main() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("%v\n", e)
		}
	}()
	ver := flag.Bool("version", false, "show version info")
	id := flag.String("id", "", "OED app_id")
	key := flag.String("key", "", "OED app_key")
	port := flag.String("port", "3528", "service port")
	flag.StringVar(&cache, "cache", "cache", "cache directory")
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}
	if *id == "" || *key == "" {
		fmt.Println("ERROR: missing OED credentials (-id or -key)")
		return
	}
	assert(os.MkdirAll(cache, 0755))
	oc = oed.NewClient(*id, *key, cache, 30)
	http.HandleFunc("/", home)
	http.HandleFunc("/work", work)
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/query/", query)
	svr := http.Server{
		Addr:         ":" + *port,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	assert(svr.ListenAndServe())
}

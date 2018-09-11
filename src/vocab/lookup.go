package main

import (
	"encoding/json"
	"oedcli"
	"os"
	"path"
)

var cache string

func lookup(word string) ([]oed.LexicalEntry, error) {
	cacheEntry := path.Join(cache, word+".json")
	var les []oed.LexicalEntry
	f, err := os.Open(cacheEntry)
	if err == nil {
		defer f.Close()
		jd := json.NewDecoder(f)
		err := jd.Decode(&les)
		return les, err
	}
	rep, err := oc.Query(word)
	if err != nil {
		return nil, err
	}
	for _, r := range rep.Results {
		les = append(les, r.LexicalEntries...)
	}
	g, err := os.Create(cacheEntry)
	if err == nil {
		defer g.Close()
		je := json.NewEncoder(g)
		je.Encode(les)
	}
	return les, nil
}

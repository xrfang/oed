package main

import (
	"encoding/json"
	"fmt"
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
	tls := make(map[string]map[string]bool)
	for _, r := range rep.Results {
		for _, le := range r.LexicalEntries {
			for _, e := range le.Entries {
				for _, s := range e.Senses {
					for _, tl := range s.ThesaurusLinks {
						l := tls[tl.EntryID]
						if l == nil {
							l = make(map[string]bool)
						}
						l[tl.SenseID] = true
						tls[tl.EntryID] = l
					}
				}
			}
			les = append(les, le)
		}
	}
	fmt.Printf("TLS: %+v\n", tls)
	g, err := os.Create(cacheEntry)
	if err == nil {
		defer g.Close()
		je := json.NewEncoder(g)
		je.Encode(les)
	}
	return les, nil
}

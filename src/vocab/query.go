package main

import (
	"encoding/json"
	"net/http"
	"oedcli"
	"regexp"
	"strings"
)

var nx *regexp.Regexp

func init() {
	nx = regexp.MustCompile(`\s+`)
}

func lookup(word, feature string) (les []oed.LexicalEntry, err error) {
	qr, err := oc.Query(word, feature)
	if err != nil {
		return nil, err
	}
	for _, r := range qr.Results {
		les = append(les, r.LexicalEntries...)
	}
	return
}

type RelatedWords struct {
	Examples []string `json:",omitempty"`
	Tags     string   `json:",omitempty"`
	Type     string
	Words    []string
	tags     []string
}

type Sense struct {
	Definition string
	Domains    string            `json:",omitempty"`
	Examples   []string          `json:",omitempty"`
	Notes      map[string]string `json:",omitempty"`
	Registers  string            `json:",omitempty"`
	Thesaurus  []RelatedWords    `json:",omitempty"`
	SubSenses  []Sense           `json:",omitempty"`
}

type LexicalEntry struct {
	Category       string
	Senses         []Sense
	Pronunciations []oed.Pronunciation
	Text           string
}

func query(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			http.Error(w, e.(error).Error(), http.StatusInternalServerError)
		}
	}()
	word := nx.ReplaceAllString(strings.TrimSpace(r.URL.Path[7:]), "_")
	rels := make(map[string]map[string]RelatedWords)
	getRelatedWords := func(feature string) {
		extractWords := func(s oed.Sense, relation string) {
			rwg := rels[s.ID]
			if rwg == nil {
				rwg = make(map[string]RelatedWords)
			}
			rw := rwg[relation]
			for _, r := range s.Registers {
				var exists bool
				for _, t := range rw.tags {
					exists = t == r
					if exists {
						break
					}
				}
				if !exists {
					rw.tags = append(rw.tags, r)
				}
			}
			for _, ex := range s.Examples {
				rw.Examples = append(rw.Examples, ex.Text)
			}
			rw.Type = relation
			switch relation {
			case "synonyms":
				for _, w := range s.Synonyms {
					rw.Words = append(rw.Words, w.Text)
				}
			case "antonyms":
				for _, w := range s.Antonyms {
					rw.Words = append(rw.Words, w.Text)
				}
			}
			rw.Tags = strings.Join(rw.tags, ",")
			rwg[relation] = rw
			rels[s.ID] = rwg
		}
		rels, err := lookup(word, feature)
		if err != nil {
			return
		}
		for _, rel := range rels {
			for _, e := range rel.Entries {
				for _, s := range e.Senses {
					extractWords(s, feature)
					for _, ss := range s.SubSenses {
						extractWords(ss, feature)
					}
				}
			}

		}
	}
	getRelatedWords("synonyms")
	getRelatedWords("antonyms")
	entries, err := lookup(word, "")
	assert(err)
	var les []LexicalEntry
	getSense := func(s oed.Sense) Sense {
		var ss Sense
		ss.Definition = strings.Join(s.Definitions, "\n")
		ss.Domains = strings.Join(s.Domains, ",")
		for _, x := range s.Examples {
			ss.Examples = append(ss.Examples, x.Text)
		}
		ss.Notes = make(map[string]string)
		for _, n := range s.Notes {
			ss.Notes[n.Type] = n.Text
		}
		ss.Registers = strings.Join(s.Registers, ",")
		for _, tl := range s.ThesaurusLinks {
			rw := rels[tl.SenseID]
			if rw == nil {
				continue
			}
			if len(rw["antonyms"].Words) > 0 {
				ss.Thesaurus = append(ss.Thesaurus, rw["antonyms"])
			}
			if len(rw["synonyms"].Words) > 0 {
				ss.Thesaurus = append(ss.Thesaurus, rw["synonyms"])
			}
		}
		return ss
	}
	for _, entry := range entries {
		var le LexicalEntry
		le.Category = entry.LexicalCategory
		for _, e := range entry.Entries {
			for _, s := range e.Senses {
				ss := getSense(s)
				for _, sub := range s.SubSenses {
					ss.SubSenses = append(ss.SubSenses, getSense(sub))
				}
				le.Senses = append(le.Senses, ss)
			}
		}
		le.Pronunciations = entry.Pronunciations
		le.Text = entry.Text
		les = append(les, le)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	je := json.NewEncoder(w)
	assert(je.Encode(les))
}

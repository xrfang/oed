package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"oedcli"
	"strings"
)

func extractEntries(qr oed.QueryReply, err error) ([]oed.LexicalEntry, error) {
	if err != nil {
		return nil, err
	}
	if qr.Error != "" {
		return nil, errors.New(qr.Error)
	}
	var les []oed.LexicalEntry
	for _, r := range qr.Results {
		les = append(les, r.LexicalEntries...)
	}
	return les, nil
}

type Thesaurus struct {
	Examples  []string    `json:",omitempty"`
	Registers []string    `json:",omitempty"`
	Synonyms  []string    `json:",omitempty"`
	Antonyms  []string    `json:",omitempty"`
	SubSenses []Thesaurus `json:",omitempty"`
}

type Sense struct {
	Definition string
	Domains    []string          `json:",omitempty"`
	Examples   []string          `json:",omitempty"`
	Notes      map[string]string `json:",omitempty"`
	Registers  []string          `json:",omitempty"`
	Thesaurus  []*Thesaurus      `json:",omitempty"`
	SubSenses  []Sense           `json:",omitempty"`
}

type LexicalEntry struct {
	Category       string
	Senses         []Sense
	Pronunciations []oed.Pronunciation
	Text           string
}

func getThesaurus(tl oed.ThesaurusLink) (*Thesaurus, error) {
	entries, err := extractEntries(oc.QueryThesaurus(tl.EntryID))
	if err != nil {
		return nil, err
	}
	collectThes := func(s oed.Sense) Thesaurus {
		var th Thesaurus
		for _, x := range s.Examples {
			th.Examples = append(th.Examples, x.Text)
		}
		th.Registers = s.Registers
		for _, w := range s.Synonyms {
			th.Synonyms = append(th.Synonyms, w.Text)
		}
		for _, w := range s.Antonyms {
			th.Antonyms = append(th.Antonyms, w.Text)
		}
		return th
	}
	var th Thesaurus
	for _, ent := range entries {
		for _, e := range ent.Entries {
			for _, s := range e.Senses {
				if s.ID != tl.SenseID {
					continue
				}
				th = collectThes(s)
				for _, ss := range s.SubSenses {
					th.SubSenses = append(th.SubSenses, collectThes(ss))
				}
				return &th, nil
			}
		}
	}
	return nil, nil
}

func getSense(s oed.Sense) Sense {
	var ss Sense
	ss.Definition = strings.Join(s.Definitions, "\n")
	ss.Domains = s.Domains
	for _, x := range s.Examples {
		ss.Examples = append(ss.Examples, x.Text)
	}
	ss.Notes = make(map[string]string)
	for _, n := range s.Notes {
		ss.Notes[n.Type] = n.Text
	}
	ss.Registers = s.Registers
	for _, tl := range s.ThesaurusLinks {
		t, err := getThesaurus(tl)
		if err != nil || t == nil {
			continue
		}
		ss.Thesaurus = append(ss.Thesaurus, t)
	}
	return ss
}

func query(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			http.Error(w, e.(error).Error(), http.StatusInternalServerError)
		}
	}()
	word := strings.TrimSpace(r.URL.Path[7:])
	var les []LexicalEntry
	entries, err := extractEntries(oc.QueryDictionary(word))
	assert(err)
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

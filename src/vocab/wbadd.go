package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func wbadd(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			http.Error(w, e.(error).Error(), http.StatusInternalServerError)
		}
	}()
	e := strings.Split(r.URL.Path[8:], "/")
	entry := e[0]
	category := e[1]
	sense, _ := strconv.Atoi(e[2])
	subSense, _ := strconv.Atoi(e[3])
	url := fmt.Sprintf("http://127.0.0.1:%s/query/%s", port, entry)
	resp, err := http.Get(url)
	assert(err)
	if resp.StatusCode != http.StatusOK {
		panic(errors.New(resp.Status))
	}
	defer resp.Body.Close()
	var les []LexicalEntry
	assert(json.NewDecoder(resp.Body).Decode(&les))
	var le *LexicalEntry
	for _, e := range les {
		if e.Category == category {
			le = &e
			break
		}
	}
	if le == nil {
		http.Error(w, fmt.Sprintf("not found: entry=%s;category=%s", entry, category), http.StatusNotFound)
		return
	}
	if sense < 0 || sense >= len(le.Senses) {
		http.Error(w, "invalid sense index", http.StatusNotFound)
		return
	}
	s := le.Senses[sense]
	if subSense >= len(s.SubSenses) {
		http.Error(w, "invalid subsense index", http.StatusNotFound)
		return
	}
	if subSense >= 0 {
		s = s.SubSenses[subSense]
	}
	//TODO: save sense to workbook
	json.NewEncoder(w).Encode(s)
}

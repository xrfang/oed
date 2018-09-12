package oed

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

type Client struct {
	timeout time.Duration
	AppID   string
	AppKey  string
	url     string
	cache   string
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func NewClient(appID, appKey, cache string, queryTimeout int) *Client {
	return &Client{
		timeout: time.Duration(queryTimeout) * time.Second,
		AppID:   appID,
		AppKey:  appKey,
		url:     "https://od-api.oxforddictionaries.com/api/v1/entries/en/",
		cache:   cache,
	}
}

func (c Client) Query(word, feature string) (qr QueryReply, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	var cacheEntry, oedEntry string
	switch feature {
	case "":
		cacheEntry = path.Join(c.cache, word+".json")
		oedEntry = word
	case "synonyms":
		cacheEntry = path.Join(c.cache, word+".synonyms.json")
		oedEntry = word + "/synonyms"
	case "antonyms":
		cacheEntry = path.Join(c.cache, word+".antonyms.json")
		oedEntry = word + "/antonyms"
	default:
		panic(errors.New("invalid feature request: " + feature))
	}
	f, e := os.Open(cacheEntry)
	if e == nil {
		jd := json.NewDecoder(f)
		e := jd.Decode(&qr)
		f.Close()
		if e == nil {
			return
		}
	}
	hc := http.Client{Timeout: c.timeout}
	req, err := http.NewRequest("GET", c.url+oedEntry, nil)
	assert(err)
	req.Header.Set("app_id", c.AppID)
	req.Header.Set("app_key", c.AppKey)
	resp, err := hc.Do(req)
	assert(err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(errors.New(resp.Status))
	}
	f, err = os.OpenFile(cacheEntry, os.O_RDWR|os.O_CREATE, 0644)
	assert(err)
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	assert(err)
	f.Seek(0, 0)
	assert(json.NewDecoder(f).Decode(&qr))
	return
}

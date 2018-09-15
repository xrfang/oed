package oed

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
	"unicode"
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

func (c Client) doQuery(url, cache string) (qr QueryReply, err error) {
	f, e := os.Open(cache)
	if e == nil {
		jd := json.NewDecoder(f)
		e := jd.Decode(&qr)
		f.Close()
		if e == nil {
			return
		}
	}
	hc := http.Client{Timeout: c.timeout}
	req, err := http.NewRequest("GET", url, nil)
	assert(err)
	req.Header.Set("app_id", c.AppID)
	req.Header.Set("app_key", c.AppKey)
	resp, err := hc.Do(req)
	assert(err)
	defer resp.Body.Close()
	f, err = os.OpenFile(cache, os.O_RDWR|os.O_CREATE, 0644)
	assert(err)
	defer f.Close()
	if resp.StatusCode == http.StatusOK {
		_, err = io.Copy(f, resp.Body)
		assert(err)
		f.Seek(0, 0)
		assert(json.NewDecoder(f).Decode(&qr))
	} else {
		qr.Error = resp.Status
		var buf bytes.Buffer
		_, err = io.Copy(&buf, resp.Body)
		if err == nil && buf.Len() > 0 {
			qr.Error += "\n\n" + buf.String()
		}
		je := json.NewEncoder(f)
		je.SetIndent("", "    ")
		je.Encode(qr)
	}
	return
}

func (c Client) entry(word string) string {
	r := regexp.MustCompile(`\W+`)
	word = strings.ToLower(r.ReplaceAllString(word, "_"))
	return strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
}

func (c Client) QueryDictionary(word string) (qr QueryReply, err error) {
	cache := path.Join(c.cache, c.entry(word)+".json")
	url := c.url + word
	return c.doQuery(url, cache)
}

func (c Client) QueryThesaurus(word string) (qr QueryReply, err error) {
	cache := path.Join(c.cache, c.entry(word)+".thesaurus.json")
	url := c.url + word + "/synonyms;antonyms"
	return c.doQuery(url, cache)
}

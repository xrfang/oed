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
	"sync"
	"time"
	"unicode"
)

type Client struct {
	timeout time.Duration
	AppID   string
	AppKey  string
	url     string
	cache   string
	delay   int //in milliseconds
	next    time.Time
	sync.Mutex
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
		delay:   1000, //free account
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
	c.Lock()
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
		c.next = time.Now().Add(time.Duration(c.delay) * time.Millisecond)
		c.Unlock()
	}()
	for {
		if time.Now().After(c.next) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	hc := http.Client{Timeout: c.timeout}
	var req *http.Request
	var resp *http.Response
	req, err = http.NewRequest("GET", url, nil)
	assert(err)
	req.Header.Set("app_id", c.AppID)
	req.Header.Set("app_key", c.AppKey)
	resp, err = hc.Do(req)
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
		if resp.StatusCode == 404 {
			qr.Error = "ERR_NO_ENTRY"
		} else {
			qr.Error = resp.Status
			var buf bytes.Buffer
			_, err = io.Copy(&buf, resp.Body)
			if err == nil && buf.Len() > 0 {
				qr.Error += "\n\n" + buf.String()
			}
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

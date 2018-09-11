package oed

import (
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	timeout time.Duration
	AppID   string
	AppKey  string
	url     string
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func NewClient(appID, appKey string, queryTimeout int) *Client {
	return &Client{
		timeout: time.Duration(queryTimeout) * time.Second,
		AppID:   appID,
		AppKey:  appKey,
		url:     "https://od-api.oxforddictionaries.com/api/v1/entries/en/",
	}
}

func (c Client) Query(word string) (qr QueryReply, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	hc := http.Client{Timeout: c.timeout}
	req, err := http.NewRequest("GET", c.url+word, nil)
	assert(err)
	req.Header.Set("app_id", c.AppID)
	req.Header.Set("app_key", c.AppKey)
	resp, err := hc.Do(req)
	assert(err)
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return
	/*
		if resp.StatusCode != http.StatusOK {
			panic(errors.New(resp.Status))
		}
		assert(json.NewDecoder(resp.Body).Decode(&qr))
		return
	*/
}

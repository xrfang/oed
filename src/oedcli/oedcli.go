package oed

type Client struct {
	app_id  string
	app_key string
}

func NewClient(appID, appKey string) *Client {
	return &Client{app_id: appID, app_key: appKey}
}

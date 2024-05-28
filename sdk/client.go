package sdk

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	url string
}

func Post(url string, req any, res any) (*http.Response, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpRes, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return httpRes, err
	}

	if httpRes.Body != http.NoBody {
		defer httpRes.Body.Close()

		err = json.NewDecoder(httpRes.Body).Decode(res)
		if err != nil {
			return httpRes, err
		}
	}

	return httpRes, nil
}

func New(url string) *Client {
	return &Client{url: url}
}

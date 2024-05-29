package sdk

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

type Client struct {
	url          string
	sessionToken string
}

func New(url string) *Client {
	return &Client{url: url}
}

func (client *Client) Authenticate(sessionToken string) {
	client.sessionToken = sessionToken
}

func (client *Client) Post(path string, req any, res any) (*http.Response, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var httpClient http.Client
	httpReq, err := http.NewRequest("POST", client.url+path, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Application-Type", "application/json")
	if client.sessionToken != "" {
		httpReq.Header.Set(api.SessionTokenHeader, client.sessionToken)
	}
	httpRes, err := httpClient.Do(httpReq)
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

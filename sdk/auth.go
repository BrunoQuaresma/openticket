package sdk

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) Login(req api.LoginRequest, res *api.LoginResponse) (*http.Response, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpRes, err := http.Post(c.url+"/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	if httpRes.Body != http.NoBody {
		defer httpRes.Body.Close()

		err = json.NewDecoder(httpRes.Body).Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	return httpRes, nil
}

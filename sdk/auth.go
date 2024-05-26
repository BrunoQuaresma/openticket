package sdk

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

type LoginRequestResult = RequestResult[api.LoginResponse]

func (c *Client) Login(req api.LoginRequest) (LoginRequestResult, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return LoginRequestResult{}, err
	}

	r, err := http.Post(c.url+"/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return LoginRequestResult{}, err
	}

	result := LoginRequestResult{StatusCode: r.StatusCode}
	if r.Body != http.NoBody {
		defer r.Body.Close()

		err = json.NewDecoder(r.Body).Decode(&result.Response.Data)
		if err != nil {
			return LoginRequestResult{}, err
		}
	}

	return result, nil
}

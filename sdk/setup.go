package sdk

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/BrunoQuaresma/openticket/api"
)

func (c *Client) Setup(req api.SetupRequest) (RequestResult, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return RequestResult{}, err
	}

	r, err := http.Post(c.url+"/setup", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return RequestResult{}, err
	}

	result := RequestResult{StatusCode: r.StatusCode}
	if r.Body != http.NoBody {
		defer r.Body.Close()
		err = json.NewDecoder(r.Body).Decode(&result.Response)
		if err != nil {
			return RequestResult{}, err
		}
	}

	return result, nil
}

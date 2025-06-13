package validation

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client struct {
	url        string
	httpClient *http.Client
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) CatBreedExists(breed string) (bool, error) {
	resp, err := c.httpClient.Get(c.url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var breeds []map[string]interface{}
	if err = json.Unmarshal(bodyBytes, &breeds); err != nil {
		return false, err
	}

	for _, b := range breeds {
		if name, ok := b["name"].(string); ok {
			if name == breed {
				return true, nil
			}
		}
	}
	return false, nil
}

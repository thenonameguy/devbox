package searcher

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

const searchAPIEndpoint = "https://search.devbox.sh/search"

type client struct {
	endpoint string
}

func NewClient() *client {
	endpoint := searchAPIEndpoint
	if os.Getenv("DEVBOX_SEARCH_ENDPOINT") != "" {
		endpoint = os.Getenv("DEVBOX_SEARCH_ENDPOINT")
	}
	return &client{
		endpoint: endpoint,
	}
}

func (c *client) Search(query string) (*SearchResult, error) {
	response, err := http.Get(c.endpoint + "?q=" + query)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var result SearchResult
	return &result, json.Unmarshal(data, &result)
}

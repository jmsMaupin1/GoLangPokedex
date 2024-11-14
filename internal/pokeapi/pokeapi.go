package pokeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jmsMaupin1/pokedex/internal/pokecache"
)

type client struct {
	LocationResults BatchLocationAreaResults
	CaughtPokemon map[string]Pokemon
	baseURL *url.URL
	httpClient *http.Client
	cache *pokecache.Cache
}

func NewClient(baseURL string) (*client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	cache := pokecache.NewCache(5 * time.Second)

	client := &client {
		LocationResults: BatchLocationAreaResults{},
		httpClient: &http.Client{},
		baseURL: u,
		cache: &cache,
		CaughtPokemon: map[string]Pokemon{},
	}

	return client, nil
}

func (c *client) newRequest(method string, fullURL string, body interface{}) (*http.Request, error) {
	var data []byte
	var err error

	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("json Marshal error: %v", err)
		}
	}

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(data))
	if err != nil {
		return  nil, fmt.Errorf("New Request error: %v", err) 
	}

	return req, nil
}

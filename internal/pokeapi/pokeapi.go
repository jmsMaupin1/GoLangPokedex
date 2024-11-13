package pokeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jmsMaupin1/pokedex/internal/pokecache"
)

type LocationAreaResults struct {
	Count    int             `json:"count"`
	Next     *string          `json:"next"`
	Previous *string          `json:"previous"`
	Results  []LocationArea  `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type client struct {
	LocationResults LocationAreaResults
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
		LocationResults: LocationAreaResults{},
		httpClient: &http.Client{},
		baseURL: u,
		cache: &cache,
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

func (c *client) GetNextLocationBatch() (error) {
	var locationURL string
	var data io.ReadCloser

	if c.LocationResults.Next != nil {
		locationURL = *c.LocationResults.Next
	} else {
		locationURL = c.baseURL.String() + "/location-area/"
	}

	val, ok := c.cache.Get(locationURL)

	if ok {
		data = io.NopCloser(bytes.NewReader(val))
	} else {
		req, err := c.newRequest("GET", locationURL, nil)
		if err != nil {
			return err
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			return nil
		}

		defer res.Body.Close()
		data = res.Body
	}	

	decoder := json.NewDecoder(data)
	if err := decoder.Decode(&c.LocationResults); err != nil {
		return err
	}

	return nil
}

func (c *client) GetPreviousLocationBatch() (error) {	
	var data io.ReadCloser

	if c.LocationResults.Previous == nil {
		return fmt.Errorf("Already on first page")
	}
	
	val, ok := c.cache.Get(*c.LocationResults.Previous)

	if ok {
		data = io.NopCloser(bytes.NewReader(val))
	} else {
		req, err := c.newRequest("GET", *c.LocationResults.Previous, nil)
		if err != nil {
			return err
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		data = res.Body
	}	

	decoder := json.NewDecoder(data)
	if err := decoder.Decode(&c.LocationResults); err != nil {
		return err
	}

	return nil
}

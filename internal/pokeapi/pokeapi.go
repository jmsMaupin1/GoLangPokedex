package pokeapi

import (
	"net/http"
	"net/url"
	"encoding/json"
	"bytes"
	"fmt"
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
}

func NewClient(baseURL string) (*client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	client := &client {
		LocationResults: LocationAreaResults{},
		httpClient: &http.Client{},
		baseURL: u,
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
	var req *http.Request
	var err error

	if c.LocationResults.Next != nil {
		req, err = c.newRequest("GET", *c.LocationResults.Next, nil)
	} else {
		req, err = c.newRequest("GET", c.baseURL.String() + "/location-area/", nil)
	}
	
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&c.LocationResults); err != nil {
		return err
	}

	return nil
}

func (c *client) GetPreviousLocationBatch() (error) {	
	if c.LocationResults.Previous == nil {
		return fmt.Errorf("Already on first page")
	}

	req, err := c.newRequest("GET", *c.LocationResults.Previous, nil)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&c.LocationResults); err != nil {
		return err
	}

	return nil
}


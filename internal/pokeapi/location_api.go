package pokeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type LocationAreaInfo struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type BatchLocationAreaResults struct {
	Count    int             `json:"count"`
	Next     *string          `json:"next"`
	Previous *string          `json:"previous"`
	Results  []LocationArea  `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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

func (c *client) GetLocationInformation(area string) (LocationAreaInfo, error) {
	var data io.ReadCloser

	fullURL := c.baseURL.String() + fmt.Sprintf("/location-area/%s", area)

	val, ok := c.cache.Get(fullURL)

	if ok {
		data = io.NopCloser(bytes.NewReader(val))
	} else {
		req, err := c.newRequest("GET", fullURL, nil)
		if err != nil {
			return LocationAreaInfo{}, err
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			return LocationAreaInfo{}, err
		}

		defer res.Body.Close()

		data = res.Body
	}
	
	var locationAreaInfo LocationAreaInfo
	decoder := json.NewDecoder(data)
	if err := decoder.Decode(&locationAreaInfo); err != nil {
		return LocationAreaInfo{}, err
	}

	return locationAreaInfo, nil
}

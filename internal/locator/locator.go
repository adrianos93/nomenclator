package locator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Locator is an interface for interacting with the geolocation API
type Locator struct {
	apikey string
	limit  int
}

// Location is a custom type used to describe the data returned by the geolocation API to other packages
type Location struct {
	City    string
	Country string
	Date    time.Time
	Weather string
}

// locationData is used to store the response from the geolocation API
type locationData struct {
	Data []data `json:"data"`
}

// data is a struct used to unmarshal the JSON array returned by the geolocation API
type data struct {
	Latitude           float64     `json:"latitude"`
	Longitude          float64     `json:"longitude"`
	Type               string      `json:"type"`
	Distance           float64     `json:"distance"`
	Name               string      `json:"name"`
	Number             string      `json:"number"`
	PostalCode         string      `json:"postal_code"`
	Street             string      `json:"street"`
	Confidence         float64     `json:"confidence"`
	Region             string      `json:"region"`
	RegionCode         string      `json:"region_code"`
	County             string      `json:"county"`
	Locality           string      `json:"locality"`
	AdministrativeArea interface{} `json:"administrative_area"`
	Neighbourhood      string      `json:"neighbourhood"`
	Country            string      `json:"country"`
	CountryCode        string      `json:"country_code"`
	Continent          string      `json:"continent"`
	Label              string      `json:"label"`
}

type LocatorOptions func(*Locator)

func WithDataLimit(max int) LocatorOptions {
	return func(l *Locator) {
		l.limit = max
	}
}

// The URL of the geolocation API
var url = "http://api.positionstack.com/v1/reverse"

// New returns a new Locator
func New(apikey string, options ...LocatorOptions) *Locator {
	locator := &Locator{apikey: apikey, limit: 0}
	for _, option := range options {
		option(locator)
	}
	return locator
}

// Locate is used to return geographical data based on geographical coordinates
func (l *Locator) Locate(latitude, longitude float64) (Location, error) {
	query := fmt.Sprintf("%f,%f", latitude, longitude)
	req := l.reverseGeoRequestBuilder(query)
	resp, err := http.Get(req)
	if err != nil {
		return Location{}, fmt.Errorf("failed to retrieve geospatial data: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return Location{}, errors.New("non 2xx response from location API")
	}
	defer resp.Body.Close()
	locationData := locationData{}
	if err := json.NewDecoder(resp.Body).Decode(&locationData); err != nil {
		return Location{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return Location{
		City:    locationData.Data[0].Region,
		Country: locationData.Data[0].Country,
	}, nil
}

func (l *Locator) reverseGeoRequestBuilder(query string) string {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	q := req.URL.Query()
	q.Add("access_key", l.apikey)
	q.Add("query", query)
	if l.limit != 0 {
		q.Add("limit", fmt.Sprintf("%d", l.limit))
	}
	req.URL.RawQuery = q.Encode()
	return req.URL.String()
}

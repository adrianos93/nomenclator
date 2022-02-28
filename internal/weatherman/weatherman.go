package weatherman

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Weatherman is an interface for interacting with the Weather API
type Weatherman struct {
	apikey  string
	filters []string
}

type WeatherOptions func(*Weatherman)

func WithFilter(filters []string) WeatherOptions {
	return func(w *Weatherman) {
		w.filters = filters
	}
}

// Forecast is a custom type used to communicate weather data to other packages.
type Forecast struct {
	Conditions string
}

// apiData struct is used for unmarshalling JSON returned by the API into a type this program can parse.
type apiData struct {
	Days []struct {
		Date       string `json:"datetime"`
		UnixEpoch  int    `json:"datetimeEpoch"`
		Conditions string `json:"conditions"`
	} `json:"days"`
}

var (
	url    = "weather.visualcrossing.com"
	scheme = "HTTPS"
)

// New returns a new Weatherman
func New(apikey string, options ...WeatherOptions) *Weatherman {
	weatherman := &Weatherman{apikey: apikey}
	for _, option := range options {
		option(weatherman)
	}
	return weatherman
}

// CheckWeather is a function that will return weather data for a certain date based on geographical coordinates and a date.
func (w *Weatherman) CheckWeather(latitude, longitude float64, date time.Time) (Forecast, error) {
	locationQuery := fmt.Sprintf("%f,%f", latitude, longitude)
	// golang uses some constant dates for formatting datetime. see: https://pkg.go.dev/time#pkg-constants
	revisedDate := date.Format("2006-01-02")
	req, err := w.weatherRequestBuilder(locationQuery, revisedDate)
	if err != nil {
		return Forecast{}, err
	}
	resp, err := http.Get(req)
	if err != nil {
		return Forecast{}, fmt.Errorf("request to %s failed: %w", url, err)
	}
	defer resp.Body.Close()
	weatherData := apiData{}
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return Forecast{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return Forecast{
		Conditions: weatherData.Days[0].Conditions,
	}, nil
}

func (w *Weatherman) weatherRequestBuilder(query, date string) (string, error) {
	r := mux.NewRouter()
	s := r.Host(url).
		Schemes(scheme).
		Path(`/VisualCrossingWebServices/rest/services/timeline/{location}/{from:\d{4}-\d{2}-\d{2}}/{to:\d{4}-\d{2}-\d{2}}`).
		Queries("unitgroup", "{unitgroup:metric|UK|US}", "elements", "{elements}", "include", "{include:obs,days}",
			"key", "{key}", "options", "{options:nonulls}", "contentType", "{contentType:json}")

	url, _ := s.URL("location", query,
		"from", date,
		"to", date,
		"unitgroup", "metric",
		"elements", strings.Join(w.filters, ","),
		"include", "obs,days",
		"key", w.apikey,
		"options", "nonulls",
		"contentType", "json",
	)

	return url.String(), nil
}

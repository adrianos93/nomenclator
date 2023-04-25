package processor

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/adrianos93/nomenclator/internal/locator"
	"github.com/adrianos93/nomenclator/internal/weatherman"
)

// Locator is an interface for the locator package
type Locator interface {
	Locate(latitude, longitude float64) (locator.Location, error)
}

// Weatherman is an interface for the weatherman package
type Weatherman interface {
	CheckWeather(latitude, longitude float64, date time.Time) (weatherman.Forecast, error)
}

// Processor is an interface for interacting with the processing piece of analysing photo metadata
type Processor struct {
	locator    Locator
	weatherman Weatherman
}

// Metadata is a custom type for storing photo metadata
type Metadata struct {
	latitude, longitude float64
	date                time.Time
}

// New returns a new Processor
func New(l Locator, w Weatherman) *Processor {
	return &Processor{
		locator:    l,
		weatherman: w,
	}
}

// Process is used to process the data obtained from a source file and returning
// a title based on common features of the pictures composing an album.
func (p *Processor) Process(data [][]string) (string, []error) {
	errs := make([]error, 0, len(data))
	albumMetadata := make([]locator.Location, 0, len(data))
	for _, row := range data {
		metadata, err := mapDataRowToStruct(row)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid photo: %w", err))
			continue
		}
		photoMetadata, err := p.locator.Locate(metadata.latitude, metadata.longitude)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid photo: %w", err))
			continue
		}
		photoMetadata.Date = metadata.date
		weatherCondition, err := p.weatherman.CheckWeather(metadata.latitude, metadata.longitude, metadata.date)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid photo: %w", err))
			continue
		}
		photoMetadata.Weather = weatherCondition.Conditions

		albumMetadata = append(albumMetadata, photoMetadata)
	}
	if len(albumMetadata) == 0 {
		return "", errs
	}
	title := "A " + weatherConditions(albumMetadata) + " " + albumPeriod(albumMetadata) + " in " + albumCity(albumMetadata)
	return title, errs
}

// mapDataRowToStruct is a helper function used to map raw photo metadata to the Metadata custom type
func mapDataRowToStruct(metadata []string) (Metadata, error) {
	if len(metadata) < 1 {
		return Metadata{}, errors.New("empty row")
	}
	date, err := time.Parse("2006-01-02T15:04:05Z", metadata[0])
	if err != nil {
		return Metadata{}, fmt.Errorf("invalid date: %w", err)
	}

	latitude, err := strconv.ParseFloat(metadata[1], 64)
	if err != nil {
		return Metadata{}, err
	}
	longitude, err := strconv.ParseFloat(metadata[2], 64)
	if err != nil {
		return Metadata{}, err
	}
	return Metadata{
		latitude:  latitude,
		longitude: longitude,
		date:      date,
	}, nil
}

// albumCity is a helper function used to return the most commonly occuring city from an album
func albumCity(album []locator.Location) string {
	m := make(map[string]int, len(album))
	var count int
	var city string
	for _, photo := range album {
		m[photo.City]++
		if m[photo.City] > count {
			count = m[photo.City]
			city = photo.City
		}
	}
	return city
}

// albumPeriod is a helper function used to return the period of time the album was composed in.
func albumPeriod(album []locator.Location) string {
	if len(album) < 1 {
		return ""
	}
	var period string
	weekendRegExp := regexp.MustCompile(`Saturday|Sunday|Friday`)
	minDate, maxDate := album[0].Date, album[0].Date
	for _, photo := range album {
		if photo.Date.Before(minDate) {
			minDate = photo.Date
		}
		if photo.Date.After(maxDate) {
			maxDate = photo.Date
		}
	}
	days := maxDate.Sub(minDate).Hours() / 24
	switch {
	case days > 3:
		period = "week"
	case days > 1:
		period = "few days"
		if weekendRegExp.MatchString(maxDate.Weekday().String()) || weekendRegExp.MatchString(minDate.Weekday().String()) {
			period = "weekend"
		}
	default:
		period = "day"
	}
	return period
}

// weatherConditions is a helper function used to analyse the data returned by Weatherman and translate it to something more title friendly.
func weatherConditions(album []locator.Location) string {
	m := make(map[string]int, len(album))
	var count int
	var weather string
	rainyRegExp := regexp.MustCompile(`rain|drizzle|shower`)
	snowyRegExp := regexp.MustCompile(`snow`)
	stormyRegExp := regexp.MustCompile(`storm|thunder|tornado`)
	icyRegExp := regexp.MustCompile(`ice|icy`)
	foggyRegExp := regexp.MustCompile(`mist|overcast|fog`)
	for _, photo := range album {
		switch {
		case rainyRegExp.MatchString(strings.ToLower(photo.Weather)):
			m["rainy"]++
			if m["rainy"] > count {
				count = m["rainy"]
				weather = "rainy"
			}
		case snowyRegExp.MatchString(strings.ToLower(photo.Weather)):
			m["snowy"]++
			if m["snowy"] > count {
				count = m["snowy"]
				weather = "snowy"
			}
		case stormyRegExp.MatchString(strings.ToLower(photo.Weather)):
			m["stormy"]++
			if m["stormy"] > count {
				count = m["stormy"]
				weather = "stormy"
			}
		case icyRegExp.MatchString(strings.ToLower(photo.Weather)):
			m["chilly"]++
			if m["chilly"] > count {
				count = m["icy"]
				weather = "chilly"
			}
		case foggyRegExp.MatchString(strings.ToLower(photo.Weather)):
			m["foggy"]++
			if m["foggy"] > count {
				count = m["foggy"]
				weather = "foggy"
			}
		default:
			m["sunny"]++
			if m["sunny"] > count {
				count = m["sunny"]
				weather = "sunny"
			}
		}

	}
	return weather
}

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/adrianos93/nomenclator/internal/locator"
	"github.com/adrianos93/nomenclator/internal/processor"
	"github.com/adrianos93/nomenclator/internal/weatherman"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "nomenclator FILE_PATH")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "You must specify a csv file to process\n\nUsage:")
		flag.Usage()
		os.Exit(1)
	}
	file := flag.Arg(0)

	mapsAPIKey, found := os.LookupEnv("LOCATOR_API_KEY")
	if !found {
		fmt.Fprintln(os.Stderr, "LOCATOR_API_KEY env var not set. Please set a valid API Key")
		os.Exit(1)
	}

	weatherAPIKey, found := os.LookupEnv("WEATHER_API_KEY")
	if !found {
		fmt.Fprintln(os.Stderr, "WEATHER_API_KEY env var not set. Please set a valid API Key")
		os.Exit(1)
	}
	data, err := readCSV(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	locator := locator.New(mapsAPIKey, locator.WithDataLimit(1))
	weatherman := weatherman.New(weatherAPIKey, weatherman.WithFilter([]string{"datetime", "datetimeEpoch", "conditions"}))
	processor := processor.New(locator, weatherman)

	title, errs := processor.Process(data)
	if len(errs) > 0 {
		fmt.Fprintln(os.Stderr, errs)
	}
	fmt.Printf("Album title: %s", title)
}

func readCSV(file string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

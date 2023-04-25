# nomenclator

nomenclator is a commmand line tool that will read in a provided CSV file which contains rows of photo metadata, and will attempt to return an appropriate title for the collection.

nomenclator is dependant on 2 APIs in order to function:

- [Positionstack geolocation API](https://positionstack.com/documentation)
- [VisualCrossing Weather Data Services](https://www.visualcrossing.com/weather/weather-data-services)

To use this tool, API keys for each API is required, meaning it is necessary to create accounts with the 2 providers above. They do have a free tier, so no credit card data is required.

## Installation

Get the package with:

`go get github.com/adrianos93/nomenclator`

Install the binary under go/bin folder:

`go install github.com/adrianos93/nomenclator/cmd/nomenclator@latest`

Ensure that your go/bin folder is in your PATH, otherwise this will not work.

## Requirements

As mentioned above, in order to run, `nomenclator` requires 2 API keys from the API providers mentioned above.

Once obtained, you need to set the API keys in 2 environment variables named:

- LOCATOR_API_KEY (for the geolocation API)
- WEATHER_API_KEY (for the weather API)

If these keys are not set, the program will fail.

## Data requirements

This program ingests CSV files to produce an output.

The accepted scheme for a CSV files is:

```CSV
2020-03-30T14:12:19Z,40.728808,-73.996106
2020-03-30T14:20:10Z,40.728656,-73.998790
2020-03-30T14:32:02Z,40.727160,-73.996044
```

where the first column contains a date in the [RFC3339 format](https://www.ietf.org/rfc/rfc3339.txt).
The second and third columns have geographical coordinates data, latitude and longitude respectively.

Other schemas or formats are not supported.

Sample files can be found in the `data` folder provided.

## Usage

`cd path/to/app`

Run nomenclator:

`nomenclator path/to/csv_file`

Alternative to running the binary:

`cd path/to/nomenclator_main.go`

`go run main.go path/to/csv_file`

Example output:

Print to stdout:

`Album title: A rainy weekend in New York`

# weatherman
--
    import "github.com/adrianos93/nomenclator/internal/weatherman"


## Usage

#### type Forecast

```go
type Forecast struct {
	Conditions string
}
```

Forecast is a custom type used to communicate weather data to other packages.

#### type WeatherOptions

```go
type WeatherOptions func(*Weatherman)
```


#### func  WithFilter

```go
func WithFilter(filters []string) WeatherOptions
```

#### type Weatherman

```go
type Weatherman struct {
}
```

Weatherman is an interface for interacting with the Weather API

#### func  New

```go
func New(apikey string, options ...WeatherOptions) *Weatherman
```
New returns a new Weatherman

#### func (*Weatherman) CheckWeather

```go
func (w *Weatherman) CheckWeather(latitude, longitude float64, date time.Time) (Forecast, error)
```
CheckWeather is a function that will return weather data for a certain date
based on geographical coordinates and a date.

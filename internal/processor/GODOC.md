# processor
--
    import "github.com/adrianos93/nomenclator/internal/processor"


## Usage

#### type Locator

```go
type Locator interface {
	Locate(latitude, longitude float64) (locator.Location, error)
}
```

Locator is an interface for the locator package

#### type Metadata

```go
type Metadata struct {
}
```

Metadata is a custom type for storing photo metadata

#### type Processor

```go
type Processor struct {
}
```

Processor is an interface for interacting with the processing piece of analysing
photo metadata

#### func  New

```go
func New(l Locator, w Weatherman) *Processor
```
New returns a new Processor

#### func (*Processor) Process

```go
func (p *Processor) Process(data [][]string) (string, []error)
```
Process is used to process the data obtained from a source file and returning a
title based on common features of the pictures composing an album.

#### type Weatherman

```go
type Weatherman interface {
	CheckWeather(latitude, longitude float64, date time.Time) (weatherman.Forecast, error)
}
```

Weatherman is an interface for the weatherman package

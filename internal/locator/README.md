# locator
--
    import "github.com/adrianos93/nomenclator/internal/locator"


## Usage

#### type Location

```go
type Location struct {
	City    string
	Country string
	Date    time.Time
	Weather string
}
```

Location is a custom type used to describe the data returned by the geolocation
API to other packages

#### type Locator

```go
type Locator struct {
}
```

Locator is an interface for interacting with the geolocation API

#### func  New

```go
func New(apikey string, options ...LocatorOptions) *Locator
```
New returns a new Locator

#### func (*Locator) Locate

```go
func (l *Locator) Locate(latitude, longitude float64) (Location, error)
```
Locate is used to return geographical data based on geographical coordinates

#### type LocatorOptions

```go
type LocatorOptions func(*Locator)
```


#### func  WithDataLimit

```go
func WithDataLimit(max int) LocatorOptions
```

package locator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeGeoAPI struct {
	T *testing.T
}

func (f *fakeGeoAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/v1/error":
		_ = json.NewEncoder(w).Encode(`{"mystuff:hello"}`)
	case r.URL.Path == "/v1/reverse":
		location := locationData{
			Data: []data{
				{Region: "London", Country: "United Kingdom"},
			},
		}
		_ = json.NewEncoder(w).Encode(location)
	}
}

func TestLocator_New(t *testing.T) {
	for name, test := range map[string]struct {
		apikey  string
		expect  *Locator
		limit   int
		fields  []string
		options []LocatorOptions
	}{
		"returns a new Locator": {
			apikey: "iamapikey",
			limit:  0,
			expect: &Locator{},
		},
		"returns a new locator with limit": {
			apikey: "iamapikey",
			limit:  1,
			options: []LocatorOptions{
				WithDataLimit(1),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := New(test.apikey, test.options...)
			require.IsType(t, test.expect, got)
			require.Equal(t, test.apikey, got.apikey)
			require.Equal(t, test.limit, got.limit)
		})
	}
}

func TestLocator_Locate(t *testing.T) {
	for name, test := range map[string]struct {
		latitude, longitude float64

		want Location

		wantReqErr     bool
		wantDecoderErr bool
		wantErr        bool
	}{
		"successfully retrieve geospatial data": {
			latitude:  40.728808,
			longitude: -73.996106,
			want: Location{
				City:    "London",
				Country: "United Kingdom",
			},
		},
		"fail to generate request": {
			latitude:   40.728808,
			longitude:  -73.996106,
			wantReqErr: true,
			wantErr:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			testGeoAPI := &fakeGeoAPI{T: t}
			ts := httptest.NewServer(http.HandlerFunc(testGeoAPI.ServeHTTP))
			defer func(current string) { url = current }(url)
			url = ts.URL + "/v1/reverse"
			if test.wantReqErr {
				url = "invalid"
			}
			if test.wantDecoderErr {
				url = ts.URL + "/v1/error"
			}

			l := Locator{
				apikey: "iamapikey",
				limit:  1,
			}

			got, err := l.Locate(test.latitude, test.longitude)
			if (err != nil) != test.wantErr {
				t.Errorf("Locator.Locate() error = %v, wantErr %v", err, test.wantErr)
			}
			require.Equal(t, got.City, test.want.City)
			require.Equal(t, got.Country, test.want.Country)

		})
	}
}

func TestLocator_Locater(t *testing.T) {
	type fields struct {
		apikey string
		limit  int
	}
	type args struct {
		latitude  float64
		longitude float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Location
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Locator{
				apikey: tt.fields.apikey,
				limit:  tt.fields.limit,
			}
			got, err := l.Locate(tt.args.latitude, tt.args.longitude)
			if (err != nil) != tt.wantErr {
				t.Errorf("Locator.Locate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Locator.Locate() = %v, want %v", got, tt.want)
			}
		})
	}
}

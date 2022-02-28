package weatherman

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeWeatherAPI struct {
	T *testing.T
}

func (f *fakeWeatherAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.Contains(r.URL.Path, "services/timeline"):
		_ = json.NewEncoder(w).Encode(apiData{
			Days: []struct {
				Date       string "json:\"datetime\""
				UnixEpoch  int    "json:\"datetimeEpoch\""
				Conditions string "json:\"conditions\""
			}{
				{
					Conditions: "Rain,Overcast",
				},
			},
		})
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func TestWeatherman_New(t *testing.T) {
	for name, test := range map[string]struct {
		apikey  string
		expect  *Weatherman
		filters []string
		options []WeatherOptions
	}{
		"returns a new Weatherman": {
			apikey: "iamapikey",
			expect: &Weatherman{},
		},
		"returns a new Weatherman with filters": {
			apikey:  "iamapikey",
			filters: []string{"conditions"},
			options: []WeatherOptions{
				WithFilter([]string{"conditions"}),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := New(test.apikey, test.options...)
			require.IsType(t, test.expect, got)
			require.Equal(t, test.apikey, got.apikey)
			require.Equal(t, test.filters, got.filters)
		})
	}
}

func TestWeatherman_CheckWeather(t *testing.T) {
	for name, test := range map[string]struct {
		latitude, longitude float64
		date                time.Time

		want Forecast

		wantBuildError, wantReqErr, wantErr bool
	}{
		"successfully retrieve weather data": {
			latitude:  40.728808,
			longitude: -73.996106,
			date:      time.Now(),
			want: Forecast{
				Conditions: "Rain,Overcast",
			},
		},
		"request failed": {
			latitude:   40.728808,
			longitude:  -73.996106,
			date:       time.Now(),
			wantReqErr: true,
			wantErr:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			testWeatherAPI := &fakeWeatherAPI{T: t}
			ts := httptest.NewServer(http.HandlerFunc(testWeatherAPI.ServeHTTP))
			defer func(current string) { url = current }(url)
			defer func(current string) { scheme = current }(scheme)
			scheme = "HTTP"
			url = strings.TrimPrefix(ts.URL, "http://")
			if test.wantBuildError {
				url = ts.URL
			}
			if test.wantReqErr {
				url = "invalid"
			}
			w := Weatherman{
				apikey:  "iamapikey",
				filters: []string{"datetime", "datetimeEpoch", "conditions"},
			}
			if test.wantBuildError {
				w.filters = []string{"badfilter"}
			}
			got, err := w.CheckWeather(test.latitude, test.longitude, test.date)
			if (err != nil) != test.wantErr {
				t.Errorf("Weatherman.CheckWeather() error = %v, wantErr = %v", err, test.wantErr)
			}
			require.Equal(t, got.Conditions, test.want.Conditions)
		})
	}
}

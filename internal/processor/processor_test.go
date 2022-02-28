package processor

import (
	"fmt"
	"testing"
	"time"

	"github.com/adrianos93/nomenclator/internal/locator"
	"github.com/adrianos93/nomenclator/internal/weatherman"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockLocator struct {
	mock.Mock
}

func (l *mockLocator) Locate(latitude, longitude float64) (locator.Location, error) {
	args := l.Called(latitude, longitude)
	return args.Get(0).(locator.Location), args.Error(1)
}

type mockWeatherman struct {
	mock.Mock
}

func (w *mockWeatherman) CheckWeather(latitude, longitude float64, date time.Time) (weatherman.Forecast, error) {
	args := w.Called(latitude, longitude, date)
	return args.Get(0).(weatherman.Forecast), args.Error(1)
}

func TestProcessor_New(t *testing.T) {
	locator := &mockLocator{}
	weatherman := &mockWeatherman{}
	got := New(locator, weatherman)
	require.IsType(t, &Processor{}, got)
}

func dateParser(date string) time.Time {
	datetime, _ := time.Parse("2006-01-02T15:04:05", date)
	return datetime
}
func TestProcessor_Process(t *testing.T) {
	for name, test := range map[string]struct {
		input [][]string

		locatorCalls    []mock.Call
		weathermanCalls []mock.Call

		expect string

		wantErr bool
	}{
		"foggy weekend": {
			input: [][]string{
				{"2020-03-30T14:12:19Z", "40.728808", "-73.996106"},
				{"2020-03-29T14:20:10Z", "40.728656", "-73.998790"},
				{"2020-03-28T14:32:02Z", "40.727160", "-73.996044"},
			},
			locatorCalls: []mock.Call{
				{Method: "Locate", Arguments: []interface{}{40.728808, -73.996106}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
				{Method: "Locate", Arguments: []interface{}{40.728656, -73.998790}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
				{Method: "Locate", Arguments: []interface{}{40.727160, -73.996044}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
			},
			weathermanCalls: []mock.Call{
				{Method: "CheckWeather", Arguments: []interface{}{40.728808, -73.996106, dateParser("2020-03-30T14:12:19")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Rain"}, nil}},
				{Method: "CheckWeather", Arguments: []interface{}{40.728656, -73.998790, dateParser("2020-03-29T14:20:10")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Overcast"}, nil}},
				{Method: "CheckWeather", Arguments: []interface{}{40.727160, -73.996044, dateParser("2020-03-28T14:32:02")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Fog"}, nil}},
			},

			expect: "A foggy weekend in New York",
		},
		"get snowy few days": {
			input: [][]string{
				{"2020-02-11T14:12:19Z", "40.728808", "-73.996106"},
				{"2020-02-12T14:20:10Z", "40.728656", "-73.998790"},
				{"2020-02-13T14:32:02Z", "40.727160", "-73.996044"},
			},
			locatorCalls: []mock.Call{
				{Method: "Locate", Arguments: []interface{}{40.728808, -73.996106}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
				{Method: "Locate", Arguments: []interface{}{40.728656, -73.998790}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
				{Method: "Locate", Arguments: []interface{}{40.727160, -73.996044}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
			},
			weathermanCalls: []mock.Call{
				{Method: "CheckWeather", Arguments: []interface{}{40.728808, -73.996106, dateParser("2020-02-11T14:12:19")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Snow"}, nil}},
				{Method: "CheckWeather", Arguments: []interface{}{40.728656, -73.998790, dateParser("2020-02-12T14:20:10")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Snow"}, nil}},
				{Method: "CheckWeather", Arguments: []interface{}{40.727160, -73.996044, dateParser("2020-02-13T14:32:02")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Drizzle"}, nil}},
			},

			expect: "A snowy few days in New York",
		},
		"get stormy day": {
			input: [][]string{
				{"2020-02-11T14:12:19Z", "40.728808", "-73.996106"},
			},
			locatorCalls: []mock.Call{
				{Method: "Locate", Arguments: []interface{}{40.728808, -73.996106}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
			},
			weathermanCalls: []mock.Call{
				{Method: "CheckWeather", Arguments: []interface{}{40.728808, -73.996106, dateParser("2020-02-11T14:12:19")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Tornado"}, nil}},
			},

			expect: "A stormy day in New York",
		},
		"get chilly day": {
			input: [][]string{
				{"2020-02-11T14:12:19Z", "40.728808", "-73.996106"},
			},
			locatorCalls: []mock.Call{
				{Method: "Locate", Arguments: []interface{}{40.728808, -73.996106}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
			},
			weathermanCalls: []mock.Call{
				{Method: "CheckWeather", Arguments: []interface{}{40.728808, -73.996106, dateParser("2020-02-11T14:12:19")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Ice"}, nil}},
			},

			expect: "A chilly day in New York",
		},
		"sunny week": {
			input: [][]string{
				{"2020-03-30T14:12:19Z", "40.728808", "-73.996106"},
				{"2020-03-29T14:20:10Z", "40.728656", "-73.998790"},
				{"2020-03-28T14:32:02Z", "40.727160", "-73.996044"},
				{"2020-03-31T14:32:02Z", "40.727160", "-73.996044"},
				{"2020-04-01T14:32:02Z", "40.727160", "-73.996044"},
			},
			locatorCalls: []mock.Call{
				{Method: "Locate", Arguments: []interface{}{40.728808, -73.996106}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
				{Method: "Locate", Arguments: []interface{}{40.728656, -73.998790}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
				{Method: "Locate", Arguments: []interface{}{40.727160, -73.996044}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
				{Method: "Locate", Arguments: []interface{}{40.727160, -73.996044}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
				{Method: "Locate", Arguments: []interface{}{40.727160, -73.996044}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
			},
			weathermanCalls: []mock.Call{
				{Method: "CheckWeather", Arguments: []interface{}{40.728808, -73.996106, dateParser("2020-03-30T14:12:19")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Clear"}, nil}},
				{Method: "CheckWeather", Arguments: []interface{}{40.728656, -73.998790, dateParser("2020-03-29T14:20:10")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Rain"}, nil}},
				{Method: "CheckWeather", Arguments: []interface{}{40.727160, -73.996044, dateParser("2020-03-28T14:32:02")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Clear"}, nil}},
				{Method: "CheckWeather", Arguments: []interface{}{40.727160, -73.996044, dateParser("2020-03-31T14:32:02")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Overcast"}, nil}},
				{Method: "CheckWeather", Arguments: []interface{}{40.727160, -73.996044, dateParser("2020-04-01T14:32:02")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Clear"}, nil}},
			},

			expect: "A sunny week in New York",
		},
		"invalid data": {
			input: [][]string{
				{"2020-03-30 14:12:19Z", "40.728808", "-73.996106"},
				{"2020-03-29T14:20:10Z", "40.728656", "-73.998790"},
			},
			locatorCalls: []mock.Call{
				{Method: "Locate", Arguments: []interface{}{40.728656, -73.998790}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
			},
			weathermanCalls: []mock.Call{
				{Method: "CheckWeather", Arguments: []interface{}{40.728656, -73.998790, dateParser("2020-03-29T14:20:10")}, ReturnArguments: []interface{}{weatherman.Forecast{Conditions: "Rain"}, nil}},
			},
			wantErr: true,
			expect:  "A rainy day in New York",
		},
		"coordinates fail": {
			input: [][]string{
				{"2020-03-30T14:12:19Z", "what", "-73.998790"},
				{"2020-03-29T14:20:10Z", "40.728656", "what"},
			},
			wantErr: true,
		},
		"locator/weatherman fails": {
			input: [][]string{
				{"2020-03-30T14:12:19Z", "40.728808", "-73.996106"},
				{"2020-03-29T14:20:10Z", "40.728656", "-73.998790"},
			},
			locatorCalls: []mock.Call{
				{Method: "Locate", Arguments: []interface{}{40.728808, -73.996106}, ReturnArguments: []interface{}{locator.Location{}, fmt.Errorf("argh")}},
				{Method: "Locate", Arguments: []interface{}{40.728656, -73.998790}, ReturnArguments: []interface{}{locator.Location{City: "New York", Country: "USA"}, nil}},
			},
			weathermanCalls: []mock.Call{
				{Method: "CheckWeather", Arguments: []interface{}{40.728656, -73.998790, dateParser("2020-03-29T14:20:10")}, ReturnArguments: []interface{}{weatherman.Forecast{}, fmt.Errorf("argh")}},
			},
			wantErr: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			locatorDouble := &mockLocator{}
			locatorDouble.Test(t)
			defer locatorDouble.AssertExpectations(t)

			weathermanDouble := &mockWeatherman{}
			weathermanDouble.Test(t)
			defer weathermanDouble.AssertExpectations(t)

			for _, call := range test.locatorCalls {
				locatorDouble.On(call.Method, call.Arguments...).Return(call.ReturnArguments...)
			}
			for _, call := range test.weathermanCalls {
				weathermanDouble.On(call.Method, call.Arguments...).Return(call.ReturnArguments...)
			}

			p := &Processor{
				locator:    locatorDouble,
				weatherman: weathermanDouble,
			}

			got, errs := p.Process(test.input)
			if len(errs) > 0 {
				require.Error(t, errs[0])
			}
			require.Equal(t, test.expect, got)
		})
	}
}

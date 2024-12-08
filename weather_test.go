package main

import (
	"testing"
)

func TestGetUnit(t *testing.T) {
	tests := []struct {
		name      string
		unit      string
		want      Unit
		expectErr bool
	}{
		{
			name:      "valid metric unit",
			unit:      "metric",
			want:      UnitMetric,
			expectErr: false,
		},
		{
			name:      "valid imperial unit",
			unit:      "imperial",
			want:      UnitImperial,
			expectErr: false,
		},
		{
			name:      "invalid unit",
			unit:      "unknown",
			want:      "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUnit(tt.unit)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetUnit() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUnit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCity(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		cityName    string
		want        City
		expectErr   bool
	}{
		{
			name:        "valid city",
			countryCode: "US",
			cityName:    "New York",
			want:        City{ID: 5128581, Name: "New York", Country: "US"},
			expectErr:   false,
		},
		{
			name:        "invalid city",
			countryCode: "US",
			cityName:    "Unknown City",
			want:        City{},
			expectErr:   true,
		},
		{
			name:        "invalid country code",
			countryCode: "XX",
			cityName:    "New York",
			want:        City{},
			expectErr:   true,
		},
	}

	// Mock the citiesList for testing
	citiesList = []byte(`[
		{"id": 5128581, "name": "New York", "country": "US"},
		{"id": 2643743, "name": "London", "country": "GB"}
	]`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCity(tt.countryCode, tt.cityName)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetCity() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCities_Get(t *testing.T) {
	tests := []struct {
		name        string
		cities      Cities
		countryCode string
		cityName    string
		want        City
	}{
		{
			name: "city found",
			cities: Cities{
				{ID: 5128581, Name: "New York", Country: "US"},
				{ID: 2643743, Name: "London", Country: "GB"},
			},
			countryCode: "US",
			cityName:    "New York",
			want:        City{ID: 5128581, Name: "New York", Country: "US"},
		},
		{
			name: "city not found",
			cities: Cities{
				{ID: 5128581, Name: "New York", Country: "US"},
				{ID: 2643743, Name: "London", Country: "GB"},
			},
			countryCode: "US",
			cityName:    "Unknown City",
			want:        City{},
		},
		{
			name: "country code not found",
			cities: Cities{
				{ID: 5128581, Name: "New York", Country: "US"},
				{ID: 2643743, Name: "London", Country: "GB"},
			},
			countryCode: "XX",
			cityName:    "New York",
			want:        City{},
		},
		{
			name: "case insensitive match",
			cities: Cities{
				{ID: 5128581, Name: "New York", Country: "US"},
				{ID: 2643743, Name: "London", Country: "GB"},
			},
			countryCode: "us",
			cityName:    "new york",
			want:        City{ID: 5128581, Name: "New York", Country: "US"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cities.Get(tt.countryCode, tt.cityName)
			if got != tt.want {
				t.Errorf("Cities.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeather_UnitIcon(t *testing.T) {
	tests := []struct {
		name    string
		weather Weather
		want    string
	}{
		{
			name:    "metric unit",
			weather: Weather{Unit: UnitMetric},
			want:    CelsiusIcon,
		},
		{
			name:    "imperial unit",
			weather: Weather{Unit: UnitImperial},
			want:    FahrenheitIcon,
		},
		{
			name:    "unknown unit",
			weather: Weather{Unit: "unknown"},
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.weather.UnitIcon(); got != tt.want {
				t.Errorf("Weather.UnitIcon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindDirections_Get(t *testing.T) {
	tests := []struct {
		name       string
		windDirs   WindDirections
		windDirDeg float32
		want       string
	}{
		{
			name:       "exact north",
			windDirs:   windDirections,
			windDirDeg: 0.0,
			want:       "N",
		},
		{
			name:       "north-northeast",
			windDirs:   windDirections,
			windDirDeg: 22.5,
			want:       "NNE",
		},
		{
			name:       "east",
			windDirs:   windDirections,
			windDirDeg: 90.0,
			want:       "E",
		},
		{
			name:       "south-southwest",
			windDirs:   windDirections,
			windDirDeg: 202.5,
			want:       "SSW",
		},
		{
			name:       "west-northwest",
			windDirs:   windDirections,
			windDirDeg: 292.5,
			want:       "WNW",
		},
		{
			name:       "exact north again",
			windDirs:   windDirections,
			windDirDeg: 360.0,
			want:       "N",
		},
		{
			name:       "out of range",
			windDirs:   windDirections,
			windDirDeg: 400.0,
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.windDirs.Get(tt.windDirDeg); got != tt.want {
				t.Errorf("WindDirections.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

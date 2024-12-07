package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

// Unit is a unit of measurement.
type Unit string

const (
	// UnitMetric is the metric unit. Used by most countries.
	UnitMetric Unit = "metric"
	// UnitImperial is the imperial unit. Used by the US.
	UnitImperial Unit = "imperial"
)

// GetUnit returns a unit by name or an error if not found.
func GetUnit(unit string) (Unit, error) {
	switch unit {
	case string(UnitMetric):
		return UnitMetric, nil
	case string(UnitImperial):
		return UnitImperial, nil
	default:
		return "", fmt.Errorf("invalid unit '%s': valid units are metric and imperial", unit)
	}
}

const (
	// WindSpeedUnit is the unit of wind speed meters per second.
	WindSpeedUnit = "m/s"
	// WindDeg is the increment of wind direction in degrees used for determining the wind direction.
	WindDeg = 11.25
	// ThermometerIcon is the thermometer icon.
	ThermometerIcon = "🌡"
	// WindIcon is the wind icon.
	WindIcon = "💨"
	// CelsiusIcon is the Celsius icon.
	CelsiusIcon = "°C"
	// FahrenheitIcon is the Fahrenheit icon.
	FahrenheitIcon = "°F"
)

var (
	//go:embed public/city.list.min.json
	citiesList []byte
)

// GetCity returns a city by country code and city name or an error if not found.
func GetCity(countryCode, cityName string) (City, error) {
	cities := Cities{}
	err := json.Unmarshal(citiesList, &cities)
	if err != nil {
		return City{}, fmt.Errorf("unable to unmarshal city file: %s", err)
	}

	city := cities.Get(countryCode, cityName)
	if city == (City{}) {
		return City{}, fmt.Errorf("country code '%s' and city name '%s' not found", countryCode, cityName)
	}

	return city, nil
}

// City is a city with an ID, name and country.
type City struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

// Cities is a list of cities.
type Cities []City

// Get returns a city by country code and city name or an empty city if not found.
func (c Cities) Get(countryCode, cityName string) City {
	for _, city := range c {
		if strings.EqualFold(city.Country, countryCode) && strings.EqualFold(city.Name, cityName) {
			return city
		}
	}
	return City{}
}

// Weather is a weather with a temperature, icon and wind.
type Weather struct {
	Temperature float32
	Unit        Unit
	Icon        string
	Wind        Wind
}

func (w Weather) UnitIcon() string {
	switch w.Unit {
	case UnitMetric:
		return CelsiusIcon
	case UnitImperial:
		return FahrenheitIcon
	default:
		return ""
	}
}

// Wind is a wind with a speed and direction.
type Wind struct {
	Speed     float32
	Direction string
}

// WindDirection is a wind direction with a degree and direction.
type WindDirection struct {
	Degree    float32
	Direction string
}

// windDirections is a list of wind directions.
var windDirections = WindDirections{
	{Degree: 0.0, Direction: "N"},
	{Degree: 360.0, Direction: "N"},
	{Degree: 22.5, Direction: "NNE"},
	{Degree: 45.0, Direction: "NE"},
	{Degree: 67.5, Direction: "ENE"},
	{Degree: 90.0, Direction: "E"},
	{Degree: 112.5, Direction: "ESE"},
	{Degree: 135.0, Direction: "SE"},
	{Degree: 157.5, Direction: "SSE"},
	{Degree: 180.0, Direction: "S"},
	{Degree: 202.5, Direction: "SSW"},
	{Degree: 225.0, Direction: "SW"},
	{Degree: 247.5, Direction: "WSW"},
	{Degree: 270.0, Direction: "W"},
	{Degree: 292.5, Direction: "WNW"},
	{Degree: 315.0, Direction: "NW"},
	{Degree: 337.5, Direction: "NNW"},
}

// WindDirections is a list of wind.
type WindDirections []WindDirection

// Get returns a wind direction by degree or an empty string if not found.
func (w WindDirections) Get(windDirDeg float32) string {
	for _, wind := range w {
		if wind.Degree == 0 {
			if (windDirDeg >= wind.Degree) && (windDirDeg < wind.Degree+WindDeg) {
				return wind.Direction
			}
		} else if wind.Degree == 360 {
			if (windDirDeg <= wind.Degree) && (windDirDeg > wind.Degree-WindDeg) {
				return wind.Direction
			}
		} else {
			if (windDirDeg >= wind.Degree) && (windDirDeg < wind.Degree+WindDeg) {
				return wind.Direction
			} else if (windDirDeg <= wind.Degree) && (windDirDeg > wind.Degree-WindDeg) {
				return wind.Direction
			}
		}
	}

	return ""
}

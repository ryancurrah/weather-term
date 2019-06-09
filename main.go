package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	_ "github.com/ryancurrah/weather-term/statik"
)

const (
	openWeatherMapAPIURL = "https://api.openweathermap.org/data/2.5/weather?id=%d&units=%s&APPID=%s"
	windSpeedUnit        = "m/s"
	windDeg              = 11.25
	thermometerIcon      = ""
	windIcon             = ""
	celsiusIcon          = ""
	fahrenheitIcon       = ""
)

var (
	unitIcon string
)

type Cities []City

type City struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}
type WeatherJSON struct {
	Weather []struct {
		ID int `json:"id"`
	} `json:"weather"`

	Main struct {
		Temp float32 `json:"temp"`
	} `json:"main"`

	Wind struct {
		Speed float32 `json:"speed"`
		Deg   float32 `json:"deg"`
	} `json:"wind"`
}

type WindDirections []WindDirection

type WindDirection struct {
	Degree    float32
	Direction string
}

// WeatherConditionIcons https://openweathermap.org/weather-conditions
var WeatherConditionIcons = map[int]string{
	200: "",
	201: "",
	202: "",
	210: "",
	211: "",
	212: "",
	221: "",
	230: "",
	231: "",
	232: "",
	300: "",
	301: "",
	302: "",
	310: "",
	312: "",
	313: "",
	314: "",
	321: "",
	500: "",
	501: "",
	502: "",
	503: "",
	504: "",
	511: "",
	520: "",
	521: "",
	522: "",
	531: "",
	600: "",
	601: "",
	602: "",
	611: "",
	612: "",
	613: "",
	615: "",
	616: "",
	620: "",
	621: "",
	622: "",
	701: "",
	711: "",
	721: "",
	731: "",
	741: "",
	751: "",
	761: "",
	762: "",
	771: "",
	781: "",
	800: "",
	801: "",
	802: "",
	803: "",
	804: "",
}

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

func main() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to determine home directory"))
	}

	countryCode := flag.String("country", "", "a country code eg: CA or US")
	cityName := flag.String("city", "", "a city name")
	openWeatherMapAPIKey := flag.String("key", "", "openweathermap.com api key")
	unit := flag.String("unit", "metric", "metric or imperial unit")
	sleepTime := flag.Int64("sleep", 300, "number of seconds to wait before updating weather")
	weatherFile := flag.String("file", fmt.Sprintf("%s/.weather", home), "number of seconds to wait before updating weather")
	flag.Parse()

	switch *unit {
	case "metric":
		unitIcon = celsiusIcon
	case "imperial":
		unitIcon = fahrenheitIcon
	default:
		log.Fatal(errors.Errorf("invalid unit '%s'. valid units are metric and imperial", *unit))
	}

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	cities := Cities{}
	file, err := statikFS.Open("/city.list.min.json")
	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to open city file"))
	}

	err = json.NewDecoder(file).Decode(&cities)
	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to unmarshall city file"))
	}

	city := cities.Get(*countryCode, *cityName)
	if city == (City{}) {
		log.Fatalf("country code '%s' and city name '%s' not found", *countryCode, *cityName)
	}

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-quitChan
		os.Exit(0)
	}()

	for {
		resp, err := http.Get(fmt.Sprintf(openWeatherMapAPIURL, city.ID, *unit, *openWeatherMapAPIKey))
		if err != nil {
			log.Fatal(errors.Wrap(err, "unable to get current weather"))
		}

		weatherJSON := WeatherJSON{}
		err = json.NewDecoder(resp.Body).Decode(&weatherJSON)
		resp.Body.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "unable to unmarshall weather"))
		}

		if reflect.DeepEqual(WeatherJSON{}, weatherJSON) {
			log.Fatal("unable to get current weather")
		}

		weatherIcon := WeatherConditionIcons[weatherJSON.Weather[0].ID]
		weather := fmt.Sprintf("%s %g%s %s", thermometerIcon, weatherJSON.Main.Temp, unitIcon, weatherIcon)

		windDirection := windDirections.Get(weatherJSON.Wind.Deg).Direction
		wind := fmt.Sprintf("%s %g%s %s", windIcon, weatherJSON.Wind.Speed, windSpeedUnit, windDirection)

		err = ioutil.WriteFile(*weatherFile, []byte(fmt.Sprintf("%s   %s", weather, wind)), 0644)
		if err != nil {
			log.Fatal(errors.Wrap(err, "unable to write weather to file"))
		}

		time.Sleep(time.Duration(*sleepTime) * time.Second)
	}
}

func (cities Cities) Get(countryCode, cityName string) City {
	for _, city := range cities {
		if strings.ToLower(city.Country) == strings.ToLower(countryCode) &&
			strings.ToLower(city.Name) == strings.ToLower(cityName) {
			return city
		}
	}
	return City{}
}

func (windDirections WindDirections) Get(windDirDeg float32) WindDirection {
	for _, windDirection := range windDirections {
		if windDirection.Degree == 0 {
			if (windDirDeg >= windDirection.Degree) && (windDirDeg < windDirection.Degree+windDeg) {
				return windDirection
			}
		} else if windDirection.Degree == 360 {
			if (windDirDeg <= windDirection.Degree) && (windDirDeg > windDirection.Degree-windDeg) {
				return windDirection
			}
		} else {
			if (windDirDeg >= windDirection.Degree) && (windDirDeg < windDirection.Degree+windDeg) {
				return windDirection
			} else if (windDirDeg <= windDirection.Degree) && (windDirDeg > windDirection.Degree-windDeg) {
				return windDirection
			}
		}
	}
	return WindDirection{}
}

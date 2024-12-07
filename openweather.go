package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// openWeatherConditionIcons https://openweathermap.org/weather-conditions.
var openWeatherConditionIcons = map[int]string{
	// Group 2xx: Thunderstorm
	200: "⛈️",
	201: "⛈️",
	202: "⛈️",
	210: "🌩️",
	211: "🌩️",
	212: "🌩️",
	221: "🌩️",
	230: "⛈️",
	231: "⛈️",
	232: "⛈️",

	// Group 3xx: Drizzle
	300: "🌦️",
	301: "🌦️",
	302: "🌦️",
	310: "🌦️",
	311: "🌦️",
	312: "🌦️",
	313: "🌦️",
	314: "🌦️",
	321: "🌦️",

	// Group 5xx: Rain
	500: "🌧️",
	501: "🌧️",
	502: "🌧️",
	503: "🌧️",
	504: "🌧️",
	511: "🌨️", // Freezing rain
	520: "🌧️",
	521: "🌧️",
	522: "🌧️",
	531: "🌧️",

	// Group 6xx: Snow
	600: "❄️",
	601: "❄️",
	602: "❄️",
	611: "🌨️", // Sleet
	612: "🌨️",
	613: "🌨️",
	615: "🌨️",
	616: "🌨️",
	620: "❄️",
	621: "❄️",
	622: "❄️",

	// Group 7xx: Atmosphere
	701: "🌫️", // Mist
	711: "💨",  // Smoke
	721: "🌫️", // Haze
	731: "💨",  // Dust/sand
	741: "🌫️", // Fog
	751: "💨",  // Sand
	761: "💨",  // Dust
	762: "🌋",  // Volcanic ash
	771: "🌬️", // Squalls
	781: "🌪️", // Tornado

	// Group 800: Clear
	800: "☀️", // Clear sky

	// Group 80x: Clouds
	801: "🌤️", // Few clouds
	802: "⛅",  // Scattered clouds
	803: "🌥️", // Broken clouds
	804: "☁️", // Overcast clouds
}

const openWeatherMapApiUrl = "https://api.openweathermap.org/data/2.5/weather?id=%d&units=%s&APPID=%s"

// OpenWeatherResp is the response from the open weather map API.
type OpenWeatherResp struct {
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

// OpenWeather is an open weather map client used for getting weather.
type OpenWeather struct {
	APIKey string
}

// Report returns a weather report for a city.
func (o OpenWeather) Report(city City, unit Unit) (weather Weather, err error) {
	resp, err := http.Get(fmt.Sprintf(openWeatherMapApiUrl, city.ID, unit, o.APIKey))
	if err != nil {
		return Weather{}, fmt.Errorf("unable to get current weather: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return Weather{}, fmt.Errorf("unable to get current weather: %s: %s", resp.Status, string(bodyBytes))
	}

	weatherResp := OpenWeatherResp{}
	err = json.NewDecoder(resp.Body).Decode(&weatherResp)
	resp.Body.Close()
	if err != nil {
		return Weather{}, fmt.Errorf("unable to unmarshal weather: %s", err)
	}

	weather = Weather{
		Temperature: weatherResp.Main.Temp,
		Unit:        unit,
		Icon:        openWeatherConditionIcons[weatherResp.Weather[0].ID],
		Wind: Wind{
			Speed:     weatherResp.Wind.Speed,
			Direction: windDirections.Get(weatherResp.Wind.Deg),
		},
	}

	return weather, nil
}

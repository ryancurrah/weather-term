package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mitchellh/go-homedir"
)

const failMsg = "⚠️ Unable to Get Weather"

func main() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal("unable to determine home directory:", err)
	}

	countryCode := flag.String("country", "", "a country code eg: CA or US")
	cityName := flag.String("city", "", "a city name")
	openWeatherMapAPIKey := flag.String("key", "", "openweathermap.com api key")
	unit := flag.String("unit", "metric", "metric or imperial unit")
	sleepTime := flag.Int64("sleep", 300, "number of seconds to wait before updating weather")
	weatherFile := flag.String("file", fmt.Sprintf("%s/.weatherterm", home), "file to write weather to, if empty writes to stdout, by default writes to ~/.weatherterm")
	flag.Parse()

	unitType, err := GetUnit(*unit)
	if err != nil {
		log.Fatal(err)
	}

	city, err := GetCity(*countryCode, *cityName)
	if err != nil {
		writeWeather(failMsg, *weatherFile)
		log.Fatal(err)
	}

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	openWeatherClient := OpenWeather{APIKey: *openWeatherMapAPIKey}

	reportWeather(*weatherFile, openWeatherClient, city, unitType)

	for {
		select {
		case <-quitChan:
			fmt.Println("received shutdown signal, exiting...")
			return
		case <-time.After(time.Duration(*sleepTime) * time.Second):
			reportWeather(*weatherFile, openWeatherClient, city, unitType)
		}
	}
}

func reportWeather(weatherFile string, openWeatherClient OpenWeather, city City, unit Unit) {
	weather, err := openWeatherClient.Report(city, unit)
	if err != nil {
		writeWeather(failMsg, weatherFile)
		log.Fatal(err)
	}

	weatherStr := fmt.Sprintf("%s %g%s %s", ThermometerIcon, weather.Temperature, weather.UnitIcon(), weather.Icon)
	windStr := fmt.Sprintf("%s %g%s %s", WindIcon, weather.Wind.Speed, WindSpeedUnit, weather.Wind.Direction)
	weatherMsg := fmt.Sprintf("%s   %s", weatherStr, windStr)

	writeWeather(weatherMsg, weatherFile)
}

// writeWeather writes the weather to a stdout or file.
func writeWeather(msg string, file string) {
	if strings.EqualFold(file, "") {
		fmt.Println(msg)
	} else {
		err := os.WriteFile(file, []byte(msg), 0644)
		if err != nil {
			log.Fatal("unable to write weather to file:", err)
		}
	}
}

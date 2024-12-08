package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

var failMsg = "⚠️ Unable to Get Weather"

var home = func() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal("unable to determine home directory:", err)
	}

	return home
}()

var cliFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "country",
		Usage: "a country code eg: CA or US",
	},
	cli.StringFlag{
		Name:  "city",
		Usage: "a city name",
	},
	cli.StringFlag{
		Name:  "key",
		Usage: "openweathermap.com api key",
	},
	cli.StringFlag{
		Name:  "unit",
		Value: "metric",
		Usage: "metric or imperial unit",
	},
	cli.Int64Flag{
		Name:  "sleep",
		Value: 300,
		Usage: "number of seconds to wait before updating weather",
	},
	cli.StringFlag{
		Name:  "file",
		Value: filepath.Join(home, ".weatherterm"),
		Usage: "file to write weather to, if empty writes to stdout, by default writes to ~/.weatherterm",
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "weatherterm"
	app.Usage = "A weather application for the terminal"
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Run the weatherterm application",
			Flags:  cliFlags,
			Action: run,
		},
		{
			Name:   "install",
			Usage:  "Install the weatherterm service",
			Flags:  cliFlags,
			Action: install,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// run executes the weatherterm application.
func run(c *cli.Context) error {
	unitType, err := GetUnit(c.String("unit"))
	if err != nil {
		writeWeatherReport(failMsg, c.String("file"))

		return err
	}

	city, err := GetCity(c.String("country"), c.String("city"))
	if err != nil {
		writeWeatherReport(failMsg, c.String("file"))

		return err
	}

	openWeatherClient := OpenWeather{APIKey: c.String("key")}

	weatherReport, err := reportWeather(openWeatherClient, city, unitType)
	if err != nil {
		writeWeatherReport(failMsg, c.String("file"))

		return err
	}

	writeWeatherReport(weatherReport, c.String("file"))

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-quitChan:
			log.Println("received shutdown signal, exiting...")

			return nil
		case <-time.After(time.Duration(c.Int64("sleep")) * time.Second):
			weatherReport, err := reportWeather(openWeatherClient, city, unitType)
			if err != nil {
				writeWeatherReport(failMsg, c.String("file"))

				return err
			}

			writeWeatherReport(weatherReport, c.String("file"))
		}
	}
}

// install installs the weatherterm service.
func install(c *cli.Context) error {
	err := installWeatherTermService(
		home,
		c.String("country"),
		c.String("city"),
		c.String("key"),
		c.String("unit"),
		c.Int64("sleep"),
		c.String("file"),
	)
	if err != nil {
		return err
	}

	log.Println("WeatherTerm service installed successfully")

	return nil
}

func reportWeather(openWeatherClient OpenWeather, city City, unit Unit) (string, error) {
	weather, err := openWeatherClient.Report(context.Background(), city, unit)
	if err != nil {
		return "", err
	}

	weatherStr := fmt.Sprintf("%s   %g%s %s", ThermometerIcon, weather.Temperature, weather.UnitIcon(), weather.Icon)
	windStr := fmt.Sprintf("%s   %g%s %s", WindIcon, weather.Wind.Speed, WindSpeedUnit, weather.Wind.Direction)
	weatherReport := fmt.Sprintf("%s  %s", weatherStr, windStr)

	return weatherReport, nil
}

// writeWeatherReport writes the weather to a stdout or file.
func writeWeatherReport(msg string, file string) {
	if strings.EqualFold(file, "") {
		log.Println(msg)
	} else {
		err := os.WriteFile(file, []byte(msg), 0644)
		if err != nil {
			log.Fatalf("unable to write weather to file: %s", err)
		}
	}
}

//go:embed com.weatherterm.plist
var weatherTermServiceTmpl string

func installWeatherTermService(home, countryCode, cityName, apiKey, unit string, sleepTime int64, file string) error {
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to determine binary path: %w", err)
	}

	data := map[string]interface{}{
		"BinaryPath":  binaryPath,
		"CountryCode": countryCode,
		"CityName":    cityName,
		"APIKey":      apiKey,
		"Unit":        unit,
		"SleepTime":   sleepTime,
		"File":        file,
	}

	tmpl, err := template.New("weatherTermService").Parse(weatherTermServiceTmpl)
	if err != nil {
		return fmt.Errorf("unable to parse template: %w", err)
	}

	// Create the LaunchAgents directory if it doesn't exist
	launchAgentsDir := filepath.Join(home, "Library/LaunchAgents")

	err = os.MkdirAll(launchAgentsDir, 0755)
	if err != nil {
		return fmt.Errorf("unable to create LaunchAgents directory: %w", err)
	}

	// Write the rendered template to the com.weatherterm.plist file
	plistPath := filepath.Join(launchAgentsDir, "com.weatherterm.plist")

	weatherFile, err := os.Create(plistPath)
	if err != nil {
		return fmt.Errorf("unable to create plist file: %w", err)
	}
	defer func() {
		_ = weatherFile.Close()
	}()

	err = tmpl.Execute(weatherFile, data)
	if err != nil {
		return fmt.Errorf("unable to render plist template: %w", err)
	}

	return nil
}

# weatherterm

Weather in your terminal! This binary updates a `.weatherterm` file in your home directory that can be used in your terminal.

# Requirements

- [OpenWeatherMap](https://openweathermap.org/) API key
- `weatherterm` binary running in the background preferably as a service

# Add to startup on mac

The install subcommand will create a plist service file. You will need to load and run it.

```shell
weatherterm run -country US -city Miami -key 0000000000000000000 -unit imperial

launchctl load -w ~/Library/LaunchAgents/com.weatherterm.plist

launchctl start -w ~/Library/LaunchAgents/com.weatherterm.plist
```

# Get it working with POWERLEVEL9k

In your .zshrc file add the following lines...

```shell
POWERLEVEL9K_CUSTOM_WEATHER="prompt_weather"
POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS=(custom_weather)

function prompt_weather() {
    echo "$(cat "${HOME}/.weatherterm")"
}
```

# Example

```shell
weatherterm run -country US -city Miami -key 0000000000000000000 -unit imperial
^C

cat ~/.weather
🌡 0.82 ⛅   💨 7.72m/s SSW
```

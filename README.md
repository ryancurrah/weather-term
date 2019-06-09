# weather-term

Weather in your terminal! This binary updates a `.weather` file in your home directory that can be used in your terminal.

# Requirements

- [Nerdfonts](https://nerdfonts.com/) as your terminals font
- [OpenWeatherMap](https://openweathermap.org/) API key
- `weather-term` binary running in the background preferably as a service

# Add to startup on mac

Fill out `com.weather-term.plist` with your values

```shell
cp ./com.weather-term.plist ~/Library/LaunchAgents/com.weather-term.plist

launchctl load -w ~/Library/LaunchAgents/com.weather-term.plist

launchctl start -w ~/Library/LaunchAgents/com.weather-term.plist
```

# Get it working with POWERLEVEL9k

In your .zshrc file add the following lines...

```shell
POWERLEVEL9K_MODE='nerdfont-complete'
POWERLEVEL9K_CUSTOM_WEATHER="prompt_weather"
POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS=(custom_weather)

function prompt_weather() {
    echo "$(cat "${HOME}/.weather")"
}
```

# Example

```shell
weather-term -country US -city Miami -key 0000000000000000000 -unit imperial
^C

cat ~/.weather
 77.41     5.82m/s SW
```

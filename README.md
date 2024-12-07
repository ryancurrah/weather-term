# weatherterm

Weather in your terminal! This binary updates a `.weatherterm` file in your home directory that can be used in your terminal.

# Requirements

- [OpenWeatherMap](https://openweathermap.org/) API key
- `weatherterm` binary running in the background preferably as a service

# Add to startup on mac

Fill out `com.weatherterm.plist` with your values

```shell
cp ./com.weatherterm.plist ~/Library/LaunchAgents/com.weatherterm.plist

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
weatherterm -country US -city Miami -key 0000000000000000000 -unit imperial
^C

cat ~/.weather
🌡 0.82 ⛅   💨 7.72m/s SSW
```

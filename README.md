# Weather API for MiWeather app

This simple API fetches weather data from open-meteo API and converts the data into a format usable by MiWeather app. Unless otherwise specified all data is supposedly measured at 80m above sea level. Returned data is a list containing one or more JSON objects.

## Endpoints

### Current forecast

```
/WeatherForecast/current

Params:
lat: Latitude (float, required)
lon: Longitude (float, required)
```


### Hourly forecast:
```
/WeatherForecast/hourly

Params:
lat: Latitude (float, required)
lon: Longitude (float, required)
hours: Number of hours for the forecast (int, optional, default = 5)
```

### Weather Forecast Data JSON
- Time: Time of the forecast in **[RFC3339](https://www.rfc-editor.org/rfc/rfc3339.html)** format with zero Z offset (GMT +0)
- Humidity: Relative humidity at 2m above sea level
- Temperature: Temperature in C
- WindData: See description beneath
- Description: Weather data description following the **[WMO code descriptions](https://open-meteo.com/en/docs#weather_variable_documentation)**. Only the ones specified in open-meteo api description are used.

### Wind data 
- Wind Direction: Wind direction in degrees
- Wind Speed: Wind speed in m/s
- SpeedUnit: Currently hardcoded to 'm/s'
- Description: Wind speed description as according to **[https://www.weather.gov/pqr/wind](https://www.weather.gov/pqr/wind)**

Example:
```
http://localhost:8082/WeatherForecast/hourly?lat=56.16&lon=10.20&hours=3

[
  {
    "Time": "2025-03-04T20:00:00Z",
    "Humidity": 89,
    "Description": "Partly Cloudy",
    "WindData": {
      "Direction": 255,
      "Speed": 10.89,
      "SpeedUnit": "m/s",
      "Description": ""
    },
    "Temperature": 6
  },
  {
    "Time": "2025-03-04T21:00:00Z",
    "Humidity": 85,
    "Description": "Partly Cloudy",
    "WindData": {
      "Direction": 249,
      "Speed": 10.83,
      "SpeedUnit": "m/s",
      "Description": ""
    },
    "Temperature": 6
  }
]
```

### Want to try it?

1. Install Go
2. Navigate to repo dir in a terminal program
3. `go run . `
4. The api now runs at "http://localhost:8082"
5. You can try it by navigating to "http://localhost:8082/WeatherForecast/hourly?lat=56.16&lon=10.20" in a browser, or running it in curl tool.


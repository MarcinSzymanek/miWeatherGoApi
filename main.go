package main

import (
	"fmt"
	"net/http"

	"github.com/MarcinSzymanek/miWeatherGoApi/wfetch"
	"github.com/gin-gonic/gin"
)

type CurrentQueryParam struct {
	Lat float64 `form:"lat" binding:"required"`
	Lon float64 `form:"lon" binding:"required"`
}

type HourlyQueryParam struct {
	Lat   float64 `form:"lat" binding:"required"`
	Lon   float64 `form:"lon" binding:"required"`
	Count int     `form:"count"`
}

// If a user does not enter a float as lat and lon queries, return status 422
func main() {
	router := gin.Default()

	router.GET("/WeatherForecast/current", func(c *gin.Context) {
		var queryParams CurrentQueryParam
		if err := c.ShouldBind(&queryParams); err == nil {
			result, err := getCurrentWeather(queryParams.Lat, queryParams.Lon)
			if err != nil {
				fmt.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get weather data"})
				return
			}
			c.JSON(http.StatusOK, result)
		} else {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"Parameter error": "Incorrect query parameters (lat, lon)"})
		}
	})

	router.GET("/WeatherForecast/hourly", func(c *gin.Context) {
		var queryParams HourlyQueryParam
		if err := c.ShouldBind(&queryParams); err == nil {
			if queryParams.Count <= 0 {
				queryParams.Count = 5
			}
			results, err := getHourlyWeather(queryParams.Lat, queryParams.Lon, queryParams.Count)
			if err != nil {
				// If we get this far, something's wrong with OM http request or conversion
				fmt.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get weather data"})
				return
			}
			c.JSON(http.StatusOK, results)
		} else {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"Parameter error": "Incorrect query parameters (lat, lon)"})
		}
	})

	router.GET("/WeatherForecast/daily", func(c *gin.Context) {
		var queryParams HourlyQueryParam
		if err := c.ShouldBind(&queryParams); err == nil {
			if queryParams.Count <= 0 {
				queryParams.Count = 5
			}
			results, err := getDailyWeather(queryParams.Lat, queryParams.Lon, queryParams.Count)
			if err != nil {
				// If we get this far, something's wrong with OM http request or conversion
				fmt.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get weather data"})
				return
			}
			c.JSON(http.StatusOK, results)
		} else {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"Parameter error": "Incorrect query parameters (lat, lon)"})
		}
	})
	router.Run("localhost:8082")
}

func getCurrentWeather(lat float64, lon float64) (WeatherDataModel, error) {
	fmt.Println("Received Get Current request")
	reply, err := wfetch.FetchCurrentForecast(lat, lon)
	if err != nil {
		fmt.Println("Error fetching current forecast")
		return WeatherDataModel{}, err
	}

	var result WeatherDataModel
	result.Description = wmoCodeToDescription(reply.WMO)
	result.Temperature = int(reply.Temperature)
	result.Humidity = reply.Humidity
	result.WindData.Direction = reply.WindDirection
	result.WindData.Speed = reply.WindSpeed
	result.WindData.SpeedUnit = "m/s"
	result.Time = reply.CTime.Time

	return result, nil
}

// Fetch open-meteo weather data, then convert it to []WeatherDataModel
func getHourlyWeather(lat float64, lon float64, count int) ([]WeatherDataModel, error) {
	fmt.Println("Received get hourly request")
	reply, err := wfetch.FetchHourlyForecast(lat, lon, count)
	if err != nil {
		fmt.Println("Error fetching hourly forecast")
		fmt.Println(err.Error())
		return []WeatherDataModel{}, err
	}

	results := make([]WeatherDataModel, count)

	for i := 0; i < count; i++ {
		results[i].Humidity = reply.Humidity[i]
		results[i].Temperature = int(reply.Temperature[i])
		results[i].WindData.Direction = reply.WindDirection[i]
		results[i].WindData.Speed = reply.WindSpeed[i]
		results[i].WindData.SpeedUnit = "m/s"
		results[i].Time = reply.CTime[i].Time
		results[i].Description = wmoCodeToDescription(reply.WMO[i])
	}

	return results, nil
}

func getDailyWeather(lat float64, lon float64, count int) ([]WeatherDataDailyModel, error) {
	fmt.Println("Received get daily request")
	reply, err := wfetch.FetchDailyForecast(lat, lon, count)
	if err != nil {
		fmt.Println("Error fetching daily forecast")
		fmt.Println(err.Error())
		return []WeatherDataDailyModel{}, err
	}

	results := make([]WeatherDataDailyModel, count)

	for i := 0; i < count; i++ {
		results[i].Date = reply.Date[i].Time
		results[i].TemperatureMax = int(reply.TemperatureMax[i])
		results[i].TemperatureMin = int(reply.TemperatureMin[i])
		results[i].Description = wmoCodeToDescription(reply.WMO[i])
	}

	return results, nil
}

// For a full list of WMO code weather descriptions see https://www.nodc.noaa.gov/archive/arc0021/0002199/1.1/data/0-data/HTML/WMO-CODE/WMO4677.HTM
func wmoCodeToDescription(wmo int) string {
	switch {
	case wmo >= 0 && wmo < 2:
		return "Clear Sky"
	case wmo >= 2 && wmo < 4:
		return "Partly Cloudy"
	case wmo == 45 || wmo == 48:
		return "Foggy"
	case wmo == 51 || wmo == 53 || wmo == 55:
		return "Drizzle"
	case wmo == 56 || wmo == 57:
		return "Freezing Drizzle"
	case wmo == 61:
		return "Light Rain"
	case wmo == 63:
		return "Moderate Rain"
	case wmo == 65:
		return "Heavy Rain"
	case wmo >= 66 && wmo < 60:
		return "Freezing rain"
	case wmo == 71:
		return "Light Snow"
	case wmo == 73:
		return "Moderate Snow"
	case wmo == 75:
		return "Heavy Snow"
	case wmo == 77:
		return "Snow Grains"
	case wmo >= 80 && wmo < 83:
		return "Rain Showers"
	case wmo >= 85 && wmo < 87:
		return "Snow Showers"
	case wmo >= 95 && wmo < 100:
		return "Thunderstorm"
	}
	return "Unknown, wmo code: " + string(wmo)
}

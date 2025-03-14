package wfetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type ReqType int

const (
	CURRENT ReqType = iota
	HOURLY
	DAILY
)

// "https://api.open-meteo.com/v1/forecast?latitude=56.1567&longitude=10.2108&daily=weather_code,temperature_2m_max,temperature_2m_min&wind_speed_unit=ms"
const (
	OM_API_STR      string = "https://api.open-meteo.com/v1/forecast?"
	OM_LAT_LON_STR  string = "latitude=%.2f&longitude=%.2f"
	OM_HOURLY       string = "&hourly="
	OM_CURRENT      string = "&current="
	OM_DAILY        string = "&daily="
	OM_PARAMS       string = "temperature_80m,relative_humidity_2m,weather_code,wind_speed_80m,wind_direction_80m&wind_speed_unit=ms"
	OM_DAILY_PARAMS string = "weather_code,temperature_2m_max,temperature_2m_min"
	OM_HOURLY_COUNT string = "&forecast_days=1&forecast_hours=%d"
	OM_DAILY_COUNT  string = "&forecast_days=%d"
)

type OpenMeteoCurrentReply struct {
	Current CurrentForecast `json:"current"`
}

type OpenMeteoHourlyReply struct {
	Hourly HourlyForecast `json:"hourly"`
}

type OpenMeteoDailyReply struct {
	Daily DailyForecast `json:"daily"`
}

type CurrentForecast struct {
	Humidity      int     `json:"relative_humidity_2m"`
	WindDirection float32 `json:"wind_direction_80m"`
	WindSpeed     float32 `json:"wind_speed_80m"`
	Temperature   float32 `json:"temperature_80m"`
	CTime         CTime   `json:"time"`
	WMO           int     `json:"weather_code"`
}

type HourlyForecast struct {
	Humidity      []int     `json:"relative_humidity_2m"`
	WindDirection []float32 `json:"wind_direction_80m"`
	WindSpeed     []float32 `json:"wind_speed_80m"`
	Temperature   []float32 `json:"temperature_80m"`
	CTime         []CTime   `json:"time"`
	WMO           []int     `json:"weather_code"`
}

type DailyForecast struct {
	Date           []DateTime `json:"time"`
	TemperatureMax []float32  `json:"temperature_2m_max"`
	TemperatureMin []float32  `json:"temperature_2m_min"`
	WMO            []int      `json:"weather_code"`
}

type CTime struct {
	time.Time
}

type DateTime struct {
	time.Time
}

const (
	// Match 1-2 digits followed by '"' at the end of the string, capture digits
	TIME_END = `([\d]{1,2})$`
)

// Custom Unmarshall is necessary since
// Time.Unmarshall does not accept Open Meteo format (missing seconds and Z)
func (ct *CTime) UnmarshalJSON(bytes []byte) (err error) {
	exp := regexp.MustCompile(TIME_END)
	trimmedString := strings.Trim(string(bytes), "\"")
	validTime := exp.ReplaceAllString(trimmedString, `$1:00Z`)
	ct.Time, err = time.Parse(time.RFC3339, validTime)
	return err
}

func (dt *DateTime) UnmarshalJSON(bytes []byte) (err error) {
	dt.Time, err = time.Parse(time.DateOnly, strings.Trim(string(bytes), "\""))
	return err
}

func FetchCurrentForecast(lat float64, lon float64) (CurrentForecast, error) {
	resp, err := fetchFromOpenMeteo(lat, lon, 0, CURRENT)
	if err != nil {
		return CurrentForecast{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CurrentForecast{}, err
	}

	var omreply OpenMeteoCurrentReply
	err = json.Unmarshal(body, &omreply)
	return omreply.Current, err
}

func FetchHourlyForecast(lat float64, lon float64, count int) (HourlyForecast, error) {
	resp, err := fetchFromOpenMeteo(lat, lon, count, HOURLY)
	if err != nil {
		return HourlyForecast{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return HourlyForecast{}, err
	}

	var omreply OpenMeteoHourlyReply
	err = json.Unmarshal(body, &omreply)
	return omreply.Hourly, err
}

func FetchDailyForecast(lat float64, lon float64, count int) (DailyForecast, error) {
	resp, err := fetchFromOpenMeteo(lat, lon, count, DAILY)
	if err != nil {
		return DailyForecast{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DailyForecast{}, err
	}
	var omreply OpenMeteoDailyReply
	err = json.Unmarshal(body, &omreply)
	return omreply.Daily, err
}

// Get http response from open-meteo forecast api at given latitude and longitude
// and number of hours specified
func fetchFromOpenMeteo(lat float64, lon float64, count int, t ReqType) (*http.Response, error) {
	var b strings.Builder
	_, err := b.WriteString(OM_API_STR)
	if err != nil {
		return new(http.Response), err
	}
	fmt.Fprintf(&b, OM_LAT_LON_STR, lat, lon)

	switch {
	case t == CURRENT:
		fmt.Fprint(&b, OM_CURRENT)
		fmt.Fprint(&b, OM_PARAMS)
	case t == HOURLY:
		fmt.Fprint(&b, OM_HOURLY)
		fmt.Fprint(&b, OM_PARAMS)
	case t == DAILY:
		fmt.Fprint(&b, OM_DAILY)
		fmt.Fprint(&b, OM_DAILY_PARAMS)
	}

	if t == HOURLY {
		fmt.Fprintf(&b, OM_HOURLY_COUNT, count)
	} else if t == DAILY {
		fmt.Fprintf(&b, OM_DAILY_COUNT, count)
	}

	fmt.Println("Sent http request: ")
	fmt.Println(b.String())
	resp, err := http.Get(b.String())
	if err != nil {
		return resp, err
	}
	return resp, nil
}

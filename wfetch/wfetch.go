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

const (
	OM_API_STR      string = "https://api.open-meteo.com/v1/forecast?"
	OM_LAT_LON_STR  string = "latitude=%.2f&longitude=%.2f"
	OM_HOURLY       string = "&hourly="
	OM_CURRENT      string = "&current="
	OM_PARAMS       string = "temperature_80m,relative_humidity_2m,weather_code,wind_speed_80m,wind_direction_80m&wind_speed_unit=ms"
	OM_HOURLY_COUNT string = "&forecast_days=1&forecast_hours=%d"
	OM_DAILY_COUNT  string = "&forecast_days=%d"
)

type OpenMeteoHourlyReply struct {
	Hourly HourlyForecast `json:"hourly"`
}

type OpenMeteoCurrentReply struct {
	Current CurrentForecast `json:"current"`
}

type HourlyForecast struct {
	Humidity      []int     `json:"relative_humidity_2m"`
	WindDirection []float32 `json:"wind_direction_80m"`
	WindSpeed     []float32 `json:"wind_speed_80m"`
	Temperature   []float32 `json:"temperature_80m"`
	CTime         []CTime   `json:"time"`
	WMO           []int     `json:"weather_code"`
}

type CurrentForecast struct {
	Humidity      int     `json:"relative_humidity_2m"`
	WindDirection float32 `json:"wind_direction_80m"`
	WindSpeed     float32 `json:"wind_speed_80m"`
	Temperature   float32 `json:"temperature_80m"`
	CTime         CTime   `json:"time"`
	WMO           int     `json:"weather_code"`
}

type CTime struct {
	time.Time
}

const (
	TIME_END = `([\d]{1,2})\"$`
)

// Custom Unmarshall is necessary since
// Time.Unmarshall does not accept Open Meteo format (missing seconds and Z)
func (ct *CTime) UnmarshalJSON(bytes []byte) (err error) {
	exp := regexp.MustCompile(TIME_END)
	validTime := exp.ReplaceAllString(string(bytes), `$1:00Z"`)
	ct.Time, err = time.Parse(`"`+time.RFC3339+`"`, validTime)
	if err != nil {
		return err
	}
	return nil
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
	case t == HOURLY:
		fmt.Fprint(&b, OM_HOURLY)

	}

	fmt.Fprint(&b, OM_PARAMS)

	if t == HOURLY {
		fmt.Fprintf(&b, OM_HOURLY_COUNT, count)
	} else if t == DAILY {
		fmt.Fprintf(&b, OM_DAILY_COUNT, count)
	}

	resp, err := http.Get(b.String())
	if err != nil {
		return resp, err
	}
	return resp, nil
}

package wfetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	OM_API_STR     string = "https://api.open-meteo.com/v1/forecast?"
	OM_LAT_LON_STR string = "latitude=%.2f&longitude=%.2f"
	OM_PARAMS      string = "&hourly=temperature_80m,relative_humidity_2m,weather_code,wind_speed_80m,wind_direction_80m&wind_speed_unit=ms&forecast_days=1&forecast_hours=%d"
)

type OpenMeteoReply struct {
	Hourly HourlyForecast `json:"hourly"`
}

type HourlyForecast struct {
	Humidity      []int     `json:"relative_humidity_2m"`
	WindDirection []float32 `json:"wind_direction_80m"`
	WindSpeed     []float32 `json:"wind_speed_80m"`
	Temperature   []float32 `json:"temperature_80m"`
	CTime         []CTime   `json:"time"`
	WMO           []int     `json:"weather_code"`
}

type CTime struct {
	time.Time
}

// Custom Unmarshall is necessary since
// Time.Unmarshall does not accept Open Meteo format (missing seconds and Z)
func (ct *CTime) UnmarshalJSON(bytes []byte) (err error) {
	validTime := strings.Replace(string(bytes), "00\"", "00:00Z\"", 1)
	ct.Time, err = time.Parse(`"`+time.RFC3339+`"`, validTime)
	if err != nil {
		return err
	}
	return nil
}

func FetchHourlyForecast(lat float64, lon float64, count int) (HourlyForecast, error) {
	resp, err := fetchFromOpenMeteo(lat, lon, count)
	if err != nil {
		return HourlyForecast{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return HourlyForecast{}, err
	}

	var omreply OpenMeteoReply
	err = json.Unmarshal(body, &omreply)
	return omreply.Hourly, err
}

// Get http response from open-meteo forecast api at given latitude and longitude
// and number of hours specified
func fetchFromOpenMeteo(lat float64, lon float64, count int) (*http.Response, error) {
	var b strings.Builder
	_, err := b.WriteString(OM_API_STR)
	if err != nil {
		return new(http.Response), err
	}
	fmt.Fprintf(&b, OM_LAT_LON_STR, lat, lon)
	fmt.Fprintf(&b, OM_PARAMS, count)

	resp, err := http.Get(b.String())
	if err != nil {
		return resp, err
	}
	return resp, nil
}

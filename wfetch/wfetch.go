package wfetch

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	OM_API_STR     string = "https://api.open-meteo.com/v1/forecast?"
	OM_LAT_LON_STR string = "latitude=%.2f&longitude=%.2f"
	OM_PARAMS      string = "&hourly=temperature_80m,relative_humidity_2m,weather_code,wind_speed_80m,wind_direction_80m&wind_speed_unit=ms&forecast_days=1&forecast_hours=%d"
)

func FetchHourlyForecast(lat float32, lon float32, count int) string {
	var b strings.Builder
	_, err := b.WriteString(OM_API_STR)
	if err != nil {
		return ""
	}
	fmt.Fprintf(&b, OM_LAT_LON_STR, lat, lon)
	fmt.Fprintf(&b, OM_PARAMS, count)

	resp, err := http.Get(b.String())
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body[:])
}

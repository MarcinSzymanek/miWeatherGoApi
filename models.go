package main

import (
	"time"
)

type WeatherDataModel struct {
	Time        time.Time
	Humidity    int
	Description string
	WindData    WindData
	Temperature int
}

type WeatherDataDailyModel struct {
	Date           time.Time
	TemperatureMax int
	TemperatureMin int
	Description    string
}

type WindData struct {
	Direction   float32 `json: wind_direction_80m`
	Speed       float32 `json: wind_speed_80m`
	SpeedUnit   string
	Description string
}

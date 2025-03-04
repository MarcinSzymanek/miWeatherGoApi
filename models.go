package main

import (
	"time"
)

type WeatherDataModel struct {
	Time        time.Time `json:`
	Humidity    int       `json`
	Description string
	WindData    WindData
	Temperature int `json:`
}

type WindData struct {
	Direction   float32 `json: wind_direction_80m`
	Speed       float32 `json: wind_speed_80m`
	SpeedUnit   string
	Description string
}

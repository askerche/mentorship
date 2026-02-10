package models

type PrayerTimes struct {
	Fajr    string `json:"fajr"`
	Sunrise string `json:"sunrise"`
	Dhuhr   string `json:"dhuhr"`
	Asr     string `json:"asr"`
	Maghrib string `json:"maghrib"`
	Isha    string `json:"isha"`
}

type Response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		Timings PrayerTimes `json:"timings"`
	} `json:"data"`
}

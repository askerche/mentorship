package models

type PrayerTimes struct {
	Fajr    string `json:"Fajr"`
	Sunrise string `json:"Sunrise"`
	Dhuhr   string `json:"Dhuhr"`
	Asr     string `json:"Asr"`
	Maghrib string `json:"Maghrib"`
	Isha    string `json:"Isha"`
}

type Response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		Timings PrayerTimes `json:"timings"`
		Date    Date        `json:"date"`
	} `json:"data"`
}

type Date struct {
	Readable string `json:"readable"`
	Hijri    Hijri  `json:"hijri"`
}

type Hijri struct {
	Day   string     `json:"day"`
	Month HijriMonth `json:"month"`
	Year  string     `json:"year"`
}

type HijriMonth struct {
	En string `json:"en"`
}

type PrayerItem struct {
	Name string
	Time string
}

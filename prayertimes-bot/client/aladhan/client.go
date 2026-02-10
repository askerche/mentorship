package aladhan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"prayertimes/models"
)

type Client struct {
}

func GetTodayPrayerTimesByCity(ctx context.Context, city string) (models.PrayerTimes, error) {
	host := "https://api.aladhan.com/v1/timingsByCity"
	u, err := url.Parse(host)
	if err != nil {
		return models.PrayerTimes{}, fmt.Errorf("Invalid URL: %w", err)
	}
	q := u.Query()
	q.Set("city", city)
	q.Set("country", "Russia")
	q.Set("method", "3")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return models.PrayerTimes{}, fmt.Errorf("Failed to create request: %w", err)
	}

	var apiResp models.Response
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.PrayerTimes{}, fmt.Errorf("Failed to do request: %w", err)
	}
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return models.PrayerTimes{}, fmt.Errorf("Failed to unmarshall json: %w", err)
	}
	return apiResp.Data.Timings, nil
}

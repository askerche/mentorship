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
	httpclient http.Client
	host       string
}

func New() *Client {
	return &Client{
		httpclient: http.Client{},
		host:       "https://api.aladhan.com/",
	}
}

func (c *Client) GetTodayPrayerTimesByCity(ctx context.Context, city string) (models.PrayerTimes, models.Hijri, error) {
	u, err := url.Parse(c.host + "v1/timingsByCity")
	if err != nil {
		return models.PrayerTimes{}, models.Hijri{}, fmt.Errorf("Invalid URL: %w", err)
	}
	q := u.Query()
	q.Set("city", city)
	q.Set("country", "RU")
	q.Set("method", "3")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return models.PrayerTimes{}, models.Hijri{}, fmt.Errorf("Failed to create request: %w", err)
	}
	var apiResp models.Response
	resp, err := c.httpclient.Do(req)
	if err != nil {
		return models.PrayerTimes{}, models.Hijri{}, fmt.Errorf("Failed to do request: %w", err)
	}
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return models.PrayerTimes{}, models.Hijri{}, fmt.Errorf("Failed to unmarshall json: %w", err)
	}
	return apiResp.Data.Timings, apiResp.Data.Date.Hijri, nil
}

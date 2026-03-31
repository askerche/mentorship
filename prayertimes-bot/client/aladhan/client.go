package aladhan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"prayertimes/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	httpclient http.Client
	host       string
	rdb        *redis.Client
}

func New(rdb *redis.Client) *Client {
	return &Client{
		httpclient: http.Client{},
		host:       "https://api.aladhan.com/",
		rdb:        rdb,
	}
}

func (c *Client) GetTodayPrayerTimesByCity(ctx context.Context, city string) (models.PrayerTimes, models.Hijri, error) {
	today := time.Now().Format("02.01.2006")
	cachedKey := fmt.Sprintf("prayers_%s_%s", city, today)
	prayerTimesCachedJSON, err := c.rdb.Get(ctx, cachedKey).Result()
	if err != nil && err == redis.Nil {
		fmt.Println("error get cached prayer times:", err)
	}
	if prayerTimesCachedJSON != "" {
		var prayerTimesCached models.Response
		err = json.Unmarshal([]byte(prayerTimesCachedJSON), &prayerTimesCached)
		if err != nil {
			fmt.Println("error unmarshal prayer times cached", err)
		}
		return prayerTimesCached.Data.Timings, prayerTimesCached.Data.Date.Hijri, nil
	}

	u, err := url.Parse(c.host + "v1/timingsByCity")
	if err != nil {
		return models.PrayerTimes{}, models.Hijri{}, fmt.Errorf("Invalid URL: %w", err)
	}
	q := u.Query()
	q.Set("city", city)
	q.Set("country", "RU")
	q.Set("method", "2")
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
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return models.PrayerTimes{}, models.Hijri{}, fmt.Errorf("Failed to unmarshall json: %w", err)
	}
	prayerTimesJSON, err := json.Marshal(apiResp)
	if err != nil {
		fmt.Println("error marshal prayer times for cache: %w", err)
	}

	cmd := c.rdb.Set(ctx, cachedKey, prayerTimesJSON, 0)
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
	}
	return apiResp.Data.Timings, apiResp.Data.Date.Hijri, nil
}

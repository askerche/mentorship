package models

import (
	"time"
)

type ShortenLinkReq struct {
	Link string `json:"link"`
}

type ShortenLinkResp struct {
	ShortLink string `json:"short_link"`
}

type AnalyticsResp struct {
	RedirectsCount int                 `json: "redirects_count"`
	Redirects      []RedirectAnalytics `json: "redirects"`
}

type RedirectAnalytics struct {
	UserAgent string    `json: "user_agent"`
	CreatedAt time.Time `json: "created_at"`
}

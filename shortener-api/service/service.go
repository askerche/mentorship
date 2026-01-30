package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"shortener/models"
	"shortener/repo"
)

type Service struct {
	Repo repo.Repo
}

func (s *Service) ShortenLink(ctx context.Context, longLink string) (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	shortLink := base64.RawURLEncoding.EncodeToString(b)
	fmt.Println("will insert:", "long =", longLink, "short =", shortLink)

	err = s.Repo.StoreShortLink(ctx, longLink, shortLink)
	if err != nil {
		return "", err
	}
	return shortLink, nil
}

func (s *Service) ShortLink(ctx context.Context, shortLink string, userAgent string) (string, error) {
	longLink, err := s.Repo.GetLongLink(ctx, shortLink)
	if err != nil {
		return "", err
	}
	err = s.Repo.StoreRedirect(ctx, shortLink, userAgent)
	if err != nil {
		return "", err
	}
	return longLink, nil
}

func (s *Service) LinkAnalytics(ctx context.Context, shortLink string) (models.AnalyticsResp, error) {
	res, err := s.Repo.GetLinkAnalytics(ctx, shortLink)
	if err != nil {
		return models.AnalyticsResp{}, err
	}
	return res, nil
}

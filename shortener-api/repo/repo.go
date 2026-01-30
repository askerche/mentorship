package repo

import (
	"context"
	"errors"
	"fmt"
	"shortener/models"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	DB *pgx.Conn
}

func (r *Repo) StoreShortLink(ctx context.Context, longLink string, shortLink string) error {
	_, err := r.DB.Exec(ctx, "INSERT INTO links (long_url, short_url) VALUES ($1, $2)",
		longLink, shortLink)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetLongLink(ctx context.Context, shortLink string) (string, error) {
	var longLink string
	row := r.DB.QueryRow(ctx, "SELECT long_url FROM links WHERE short_url = $1", shortLink)
	err := row.Scan(&longLink)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
		return "", err
	}
	return longLink, nil
}

func (r *Repo) StoreRedirect(ctx context.Context, shortLink string, userAgent string) error {
	_, err := r.DB.Exec(ctx, "INSERT INTO redirects (short_link, user_agent) VALUES ($1, $2)", shortLink, userAgent)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetLinkAnalytics(ctx context.Context, shortLink string) (models.AnalyticsResp, error) {
	rows, err := r.DB.Query(ctx, "SELECT user_agent, created_at FROM redirects WHERE short_link = $1", shortLink)
	if err != nil {
		fmt.Printf("Ошибка при получении переходов: %v", err)
		return models.AnalyticsResp{}, err
	}
	var res models.AnalyticsResp

	for rows.Next() {
		var r models.RedirectAnalytics
		err = rows.Scan(&r.UserAgent, &r.CreatedAt)
		if err != nil {
			fmt.Printf("Ошибка при сканировании переходов: %v", err)
			return models.AnalyticsResp{}, err
		}
		res.Redirects = append(res.Redirects, r)
	}
	res.RedirectsCount = len(res.Redirects)
	return res, nil
}

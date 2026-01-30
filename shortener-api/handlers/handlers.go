package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

type Handler struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) ShortenLinkHandler(c *gin.Context) {
	var req ShortenLinkReq
	c.ShouldBindJSON(&req)
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("failed to generate short link", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate short link"})
		return
	}

	shortLink := base64.RawURLEncoding.EncodeToString(b)
	fmt.Println("will insert:", "long =", req.Link, "short =", shortLink)

	_, err = s.db.Exec(c, "INSERT INTO links (long_url, short_url) VALUES ($1, $2)", req.Link, shortLink)
	if err != nil {
		fmt.Println("error db.Exec: ,", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate short link"})
		return
	}

	c.JSON(http.StatusOK, ShortenLinkResp{
		ShortLink: shortLink,
	})

}

func (h *Handler) ShortLinkHandler(c *gin.Context) {
	shortLink := c.Param("link")
	var longLink string
	row := s.db.QueryRow(c, "SELECT long_url FROM links WHERE short_url = $1", shortLink)
	err := row.Scan(&longLink)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
			return
		}
		fmt.Println("error row.Scan: ,", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "try again later"})
		return
	}

	_, err = s.db.Exec(c, "INSERT INTO redirects (short_link, user_agent) VALUES ($1, $2)", shortLink, c.Request.UserAgent())
	if err != nil {
		fmt.Printf("Ошибка в сохранении аналитики: %v", err)
	}

	c.Redirect(http.StatusTemporaryRedirect, longLink)

}

func (h *Handler) LinkAnalyticsHandler(c *gin.Context) {
	shortLink := c.Param("link")
	var res AnalyticsResp
	rows, err := s.db.Query(c, "SELECT user_agent, created_at FROM redirects WHERE short_link = $1", shortLink)
	if err != nil {
		fmt.Printf("Ошибка при получении переходов: %v", err)
		c.String(http.StatusInternalServerError, "Попробуйте позже")
		return
	}

	for rows.Next() {
		var r RedirectAnalytics
		err = rows.Scan(&r.UserAgent, &r.CreatedAt)
		if err != nil {
			fmt.Printf("Ошибка при сканировании переходов: %v", err)
			c.String(http.StatusInternalServerError, "Попробуйте позже")
			return
		}
		res.Redirects = append(res.Redirects, r)
	}

	res.RedirectsCount = len(res.Redirects)

	c.JSON(http.StatusOK, res)
}

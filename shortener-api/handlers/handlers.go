package handlers

import (
	"net/http"
	"shortener/models"
	"shortener/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Svc service.Service
}

func (h *Handler) ShortenLinkHandler(c *gin.Context) {
	var req models.ShortenLinkReq
	c.ShouldBindJSON(&req)
	shortLink, err := h.Svc.ShortenLink(c, req.Link)
	if err != nil {
		c.String(http.StatusInternalServerError, "Попробуйте позже")
		return
	}
	c.JSON(http.StatusOK, models.ShortenLinkResp{
		ShortLink: shortLink,
	})

}

func (h *Handler) ShortLinkHandler(c *gin.Context) {
	shortLink := c.Param("link")
	longLink, err := h.Svc.ShortLink(c, shortLink, c.Request.UserAgent())
	if err != nil {
		c.String(http.StatusInternalServerError, "Попробуйте позже")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, longLink)
}

func (h *Handler) LinkAnalyticsHandler(c *gin.Context) {
	shortLink := c.Param("link")
	res, err := h.Svc.LinkAnalytics(c, shortLink)
	if err != nil {
		c.String(http.StatusInternalServerError, "Попробуйте позже")
		return
	}
	c.JSON(http.StatusOK, res)
}

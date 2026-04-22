package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *ApiServer) GetCategoriesHandler(c *gin.Context) {
	categories, err := s.repo.GetAllCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении списка категорий из базы данных",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count": len(categories),
		"data":  categories,
	})
}

package api

import (
	"errors"
	"fmt"
	"net/http"
	"perfume-bot/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (s *ApiServer) GetBrandsHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		limit = 0
	}

	brands, err := s.repo.GetBrandsPage(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении списка брендов",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Status": "Succes",
		"count":  len(brands),
		"data":   brands,
	})
}

func (s *ApiServer) CreateBrandHandler(c *gin.Context) {
	var req models.CreateBrandRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных",
		})
		return
	}

	id, err := s.repo.CreateBrand(c, req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при сохранении бренда",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"Status": "Бренд добавлен!",
		"id":     id,
	})
}

func (s *ApiServer) UpdateBrandHandler(c *gin.Context) {
	brandIdStr := c.Param("id")
	brandID, err := strconv.Atoi(brandIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var brand models.CreateBrandRequest
	err = c.ShouldBindJSON(&brand)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err = s.repo.UpdateBrand(c, brandID, brand.Title, brand.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update brand"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Brand successfully updated",
	})
}

func (s *ApiServer) DeleteBrandHandler(c *gin.Context) {
	brandIdStr := c.Param("id")
	brandID, err := strconv.Atoi(brandIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	err = s.repo.DeleteBrand(c, brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete brand"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "Brand deleted."})
}

func (s *ApiServer) GetBrandHandler(c *gin.Context) {
	brandIdStr := c.Param("id")
	brandID, err := strconv.Atoi(brandIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	brand, err := s.repo.GetBrand(c, brandID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "brand not found"})
		} else {
			fmt.Println("db error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
		}
		return
	}
	c.JSON(http.StatusOK, brand)
}

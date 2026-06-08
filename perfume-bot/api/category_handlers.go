package api

import (
	"net/http"
	"perfume-bot/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *ApiServer) GetCategoriesHandler(c *gin.Context) {
	categories, err := s.repo.GetAllCategories(c)
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

func (s *ApiServer) CreateCategoryHandler(c *gin.Context) {
	var req models.CreateCategoryRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных",
		})
		return
	}

	id, err := s.repo.CreateCategory(c, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при сохранении категории",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"Status": "Категория добавлена!",
		"id":     id,
	})
}

func (s *ApiServer) UpdateCategoryHandler(c *gin.Context) {
	categoryIdStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var category models.CreateCategoryRequest
	err = c.ShouldBindJSON(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err = s.repo.UpdateCategory(c, categoryID, category.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Category successfully updated",
	})
}

func (s *ApiServer) DeleteCategoryHandler(c *gin.Context) {
	categoryIdStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	err = s.repo.DeleteCategory(c, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "Category deleted."})
}

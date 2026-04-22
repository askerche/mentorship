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

func (s *ApiServer) CreateProductHandler(c *gin.Context) {
	var req models.CreateProductRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных или пропущены обязательные поля",
		})
		return
	}

	id, err := s.repo.CreateProduct(c.Request.Context(), req.BrandID, req.Title, req.Price, req.Description, req.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при сохраненния товара в БД",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"Status": "Товар добавлен!",
		"id":     id,
	})
}

func (s *ApiServer) GetProductsHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		limit = 0
	}

	products, err := s.repo.GetAllProductsPage(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении списка товаров из базы данных",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Status": "Succes",
		"count":  len(products),
		"data":   products,
	})
}

func (s *ApiServer) GetProductHandler(c *gin.Context) {
	productIdStr := c.Param("id")
	productID, err := strconv.Atoi(productIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	product, err := s.repo.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "failed to find product"})
		} else {
			fmt.Println("db error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
		}
		return
	}
	c.JSON(http.StatusOK, product)
}

func (s *ApiServer) UpdateProductHandler(c *gin.Context) {
	productIdStr := c.Param("id")
	productID, err := strconv.Atoi(productIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var product models.CreateProductRequest
	err = c.ShouldBindJSON(&product)
	if err != nil {
		fmt.Println("error unmarhall update product: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err = s.repo.UpdateProduct(c.Request.Context(), productID, product.BrandID, product.Title, product.Price, product.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product successfully updated",
	})
}

func (s *ApiServer) DeleteProductHandler(c *gin.Context) {
	productIdStr := c.Param("id")
	productID, err := strconv.Atoi(productIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	err = s.repo.DeleteProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "Product deleted."})
}

package api

import (
	"net/http"
	"perfume-bot/repository"

	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *ApiServer {
	return &ApiServer{
		repo: repo,
	}
}

func (s *ApiServer) Run() error {
	router := gin.Default()
	router.GET("/api/health", s.HealthCheckHandler)
	router.POST("/api/products", s.CreateProductHandler)
	router.GET("/api/products", s.GetProductsHandler)
	router.GET("/api/products/:id", s.GetProductHandler)
	router.PUT("/api/products/:id", s.UpdateProductHandler)
	router.DELETE("/api/products/:id", s.DeleteProductHandler)

	router.GET("/api/brands", s.GetBrandsHandler)
	router.POST("/api/brands", s.CreateBrandHandler)
	router.GET("/api/brands/:id", s.GetBrandHandler)
	router.PUT("/api/brands/:id", s.UpdateBrandHandler)
	router.DELETE("/api/brands/:id", s.DeleteBrandHandler)

	router.GET("/api/categories", s.GetCategoriesHandler)

	router.StaticFile("/admin", "./static/admin.html")

	return router.Run(":8181")
}

func (s *ApiServer) HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "API is running",
	})
}

package api

import (
	"log"
	"net/http"
	"perfume-bot/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *ApiServer) GetOrdersHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	orders, totalOrders, err := s.repo.GetAdminOrders(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении списка заказов из базы данных",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   orders,
		"meta": gin.H{
			"limit":       limit,
			"offset":      offset,
			"count":       len(orders),
			"total_count": totalOrders,
		},
	})
}

func (s *ApiServer) GetOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат ID",
		})
		return
	}

	order, err := s.repo.GetAdminOrder(c, orderID)
	if err != nil {
		log.Printf("Error to get order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка БД",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"order_id": orderID,
		"order":    order,
	})
}

func (s *ApiServer) UpdateOrderStatusHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат ID",
		})
		return
	}

	var status models.StatusRequest

	err = c.ShouldBindJSON(&status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный запрос",
		})
		return
	}

	validStatus := map[string]bool{
		"новый":       true,
		"в обработке": true,
		"отправлен":   true,
		"доставлен":   true,
		"отменен":     true,
	}

	if validStatus[status.Status] == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Не валидный статус заказа",
		})
		return
	}

	err = s.repo.UpdateOrderStatus(c, orderID, status.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при обновлении статуса заказа",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"order_id":   orderID,
		"new_status": status.Status,
	})
}

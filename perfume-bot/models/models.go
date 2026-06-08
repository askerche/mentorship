package models

import (
	"time"
)

type Product struct {
	ID          int
	Title       string
	Brand       Brand
	Price       int
	Description string
	Created_at  time.Time
	ImageFileID string
	CategoryIDs []int64 `json:"category_ids"`
}

type Brand struct {
	ID          int
	Title       string
	Description string
}

type Category struct {
	ID    int
	Title string
}

type CartItem struct {
	ProductID int
	Title     string
	BrandName string
	Price     int
	Quantity  int
}

type CreateProductRequest struct {
	BrandID     int    `json:"brand_id"`
	Title       string `json:"title"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	CategoryIDs []int  `json:"category_ids"`
	ImageFileID string `json:"image_file_id"`
}

type CreateBrandRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateCategoryRequest struct {
	Title string `json:"title"`
}

type OrderProduct struct {
	Id       int64
	Price    int
	Quantity int
}

type Order struct {
	ID         int       `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	TotalPrice int       `json:"total_price"`
	Status     string    `json:"status"`
	Created_at time.Time `json:"created_at"`
}

type OrderItemDetail struct {
	ProductID       int    `json:"product_id"`
	Title           string `json:"title"`
	BrandName       string `json:"brand_name"`
	PriceAtPurchase int    `json:"price_at_purchase"`
	Quantity        int    `json:"quantity"`
}

type StatusRequest struct {
	Status string `json:"status"`
}

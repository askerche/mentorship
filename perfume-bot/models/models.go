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
	Category    string `json:"category"`
}

type CreateBrandRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

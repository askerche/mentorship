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

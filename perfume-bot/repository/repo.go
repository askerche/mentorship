package repository

import (
	"context"
	"fmt"
	"perfume-bot/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, `
	SELECT p.id, p.title, p.price, p.description, p.created_at, b.id, b.title, b.description
	FROM products p 
	JOIN brands b
	ON p.brand_id = b.id;
	`)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err = rows.Scan(&p.ID, &p.Title, &p.Price, &p.Description,
			&p.Created_at, &p.Brand.ID, &p.Brand.Title, &p.Brand.Description)
		if err != nil {
			return nil, fmt.Errorf("scanning product row: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title FROM categories ORDER BY id;`)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		err = rows.Scan(&c.ID, &c.Title)
		if err != nil {
			return nil, fmt.Errorf("scanning category row: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *Repository) GetAllBrands(ctx context.Context) ([]models.Brand, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, description FROM brands ORDER BY id;`)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()

	brands := make([]models.Brand, 0)

	for rows.Next() {
		var b models.Brand
		err = rows.Scan(&b.ID, &b.Title, &b.Description)
		if err != nil {
			return nil, fmt.Errorf("scanning brand row: %w", err)
		}
		brands = append(brands, b)
	}
	return brands, nil
}

func (r *Repository) GetProductsByCategoryID(ctx context.Context, categoryID int) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, `
	SELECT  p.id, p.title, p.price, p.description, p.created_at, b.id, b.title, b.description
	FROM products p
	JOIN products_categories pc
	ON p.id = pc.product_id
	JOIN brands b 
	ON b.id = p.brand_id
	WHERE pc.category_id = $1
	`, categoryID)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()
	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err = rows.Scan(
			&p.ID, &p.Title, &p.Price, &p.Description, &p.Created_at,
			&p.Brand.ID, &p.Brand.Title, &p.Brand.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning product by category row: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repository) GetProductsByBrandID(ctx context.Context, brandID int) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, `
	SELECT  p.id, p.title, p.price, p.description, p.created_at, b.id, b.title, b.description
	FROM products p
	JOIN brands b 
	ON b.id = p.brand_id
	WHERE b.id = $1
	`, brandID)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()
	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err = rows.Scan(
			&p.ID, &p.Title, &p.Price, &p.Description, &p.Created_at,
			&p.Brand.ID, &p.Brand.Title, &p.Brand.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning product by brand row: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

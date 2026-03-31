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
		return nil, fmt.Errorf("error GetAllProducts db.Query: %w", err)
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err = rows.Scan(&p.ID, &p.Title, &p.Price, &p.Description,
			&p.Created_at, &p.Brand.ID, &p.Brand.Title, &p.Brand.Description)
		if err != nil {
			return nil, fmt.Errorf("error GetAllProducts rows.Scan: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title FROM categories ORDER BY id;`)
	if err != nil {
		return nil, fmt.Errorf("error GetAllCategories db.Query: %w", err)
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		err = rows.Scan(&c.ID, &c.Title)
		if err != nil {
			return nil, fmt.Errorf("error GetAllCategories rows.Scan: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *Repository) GetProductsByCategoryID(ctx context.Context, categoryID int) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, `
	SELECT  p.id, p.title, p.price, p.description, p.created_at, b.id, b.title, b.description
	FROM products p
	JOIN brands b
	ON p.brand_id = b.id
	JOIN products_categories pc 
	ON p.id = pc.product_id
	WHERE pc.category_id = $1;
	`)
	if err != nil {
		return nil, fmt.Errorf("error GetProductsByCategoryID db.Query: %w", err)
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
			return nil, fmt.Errorf("error GetProductsByCategoryID rows.Scan: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

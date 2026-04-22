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

func (r *Repository) CreateUserIfNotExists(ctx context.Context, telegram_id int64) error {
	query := `INSERT INTO users (telegram_id) VALUES($1) ON CONFLICT (telegram_id) DO NOTHING`
	_, err := r.db.Exec(ctx, query, telegram_id)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *Repository) AddToCart(ctx context.Context, telegramID int64, productID int) error {
	query := `
			INSERT INTO cart_items (telegram_id, product_id, quantity)
			VALUES($1, $2, 1)
			ON CONFLICT (telegram_id, product_id)
			DO UPDATE SET quantity = cart_items.quantity + 1
`
	_, err := r.db.Exec(ctx, query, telegramID, productID)
	if err != nil {
		return fmt.Errorf("failed to add item to cart: %w", err)
	}
	return nil
}

func (r *Repository) GetCartByTelegramID(ctx context.Context, telegramID int64) ([]models.CartItem, error) {
	query := `
				SELECT p.id, p.title, b.title, p.price, c.quantity
				FROM cart_items c
				JOIN products p ON p.id = c.product_id
				JOIN brands b ON b.id = p.brand_id
				WHERE c.telegram_id = $1
				ORDER BY c.created_at ASC
				`
	rows, err := r.db.Query(ctx, query, telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(&item.ProductID, &item.Title, &item.BrandName, &item.Price, &item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) GetAllProductsPage(ctx context.Context, limit, offset int) ([]models.Product, error) {
	query := `
			SELECT
					p.id,
					p.title,
					p.price, 
					COALESCE(p.description, ''), 
					p.created_at, 
					COALESCE(p.image_file_id, ''), 
					b.id, 
					b.title, 
					COALESCE(b.description, '')
					COALESCE(ARRAY_AGG(pc.category_id) FILTER (WHERE pc.category_id IS NOT NULL), '{}') as category_ids
			FROM products p 
			JOIN brands b ON p.brand_id = b.id
			LEFT JOIN products_categories pc ON p.id = pc.product_id
			GROUP BY p.id, b.id
			ORDER BY p.id DESC
			LIMIT $1
			OFFSET $2;
			`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err = rows.Scan(&p.ID, &p.Title, &p.Price, &p.Description, &p.Created_at, &p.ImageFileID,
			&p.Brand.ID, &p.Brand.Title, &p.Brand.Description, &p.CategoryIDs)
		if err != nil {
			return nil, fmt.Errorf("scanning product row: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repository) GetAllProductsCounts(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM products`
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting brands: %w", err)
	}
	return count, nil
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

func (r *Repository) GetProductsPageByBrandID(ctx context.Context, brandID int, limit, offset int) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, `
	SELECT p.id, p.title, p.price, COALESCE(p.description, ''), p.created_at, COALESCE(p.image_file_id, ''),
	b.id, b.title, COALESCE(b.description, '')
	FROM products p
	JOIN brands b 
	ON b.id = p.brand_id
	WHERE b.id = $1
	ORDER BY p.id DESC
	LIMIT $2
	OFFSET $3
	`, brandID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()
	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err = rows.Scan(
			&p.ID, &p.Title, &p.Price, &p.Description, &p.Created_at, &p.ImageFileID,
			&p.Brand.ID, &p.Brand.Title, &p.Brand.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning product by brand row: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repository) GetProductsPageByCategoryID(ctx context.Context, categoryID int, limit, offset int) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, `
	SELECT  p.id, p.title, p.price, p.description, p.created_at, COALESCE(p.image_file_id, ''), b.id, b.title, b.description
	FROM products p
	JOIN products_categories pc
	ON p.id = pc.product_id
	JOIN brands b 
	ON b.id = p.brand_id
	WHERE pc.category_id = $1
	ORDER BY p.id DESC
	LIMIT $2
	OFFSET $3
	`, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()
	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err = rows.Scan(
			&p.ID, &p.Title, &p.Price, &p.Description, &p.Created_at, &p.ImageFileID,
			&p.Brand.ID, &p.Brand.Title, &p.Brand.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning product by category row: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repository) GetProductsCountByCategoryID(ctx context.Context, categoryID int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM products p JOIN products_categories pc ON p.id = pc.product_id WHERE pc.category_id = $1`, categoryID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting products for category: %w", err)
	}
	return count, nil
}

func (r *Repository) GetProductsCountByBrandID(ctx context.Context, brandID int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT (*) FROM products p 
							   JOIN brands b
							   ON b.id = p.brand_id 
							   WHERE b.id = $1`,
		brandID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting products for brand: %w", err)
	}
	return count, nil
}

func (r *Repository) CreateProduct(ctx context.Context, brandID int, title string, price int, description string, category string) (int, error) {
	query := `
			  INSERT INTO products (brand_id, title, price, description, category)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id			
	`
	var newID int

	err := r.db.QueryRow(ctx, query, brandID, title, price, description, category).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new product: %w", err)
	}
	return newID, nil
}

func (r *Repository) GetProductByID(ctx context.Context, productID int) (*models.Product, error) {
	query := `
		SELECT 
			p.id, 
			p.title, 
			p.price, 
			COALESCE(p.description, ''), 
			p.created_at,
			COALESCE(p.image_file_id, ''),
			b.id, 
			b.title, 
			COALESCE(b.description, '')
		FROM products p
		JOIN brands b ON p.brand_id = b.id
		WHERE p.id = $1
	`
	var p models.Product
	err := r.db.QueryRow(ctx, query, productID).
		Scan(&p.ID,
			&p.Title,
			&p.Price,
			&p.Description,
			&p.Created_at,
			&p.ImageFileID,
			&p.Brand.ID,
			&p.Brand.Title,
			&p.Brand.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to scan product: %w", err)
	}
	return &p, nil
}

func (r *Repository) UpdateProduct(ctx context.Context, id int, brandID int, title string, price int, description string) error {
	query := `UPDATE products SET brand_id = $1, title = $2, price = $3, description = $4 WHERE id = $5`
	res, err := r.db.Exec(ctx, query, brandID, title, price, description, id)
	if err != nil {
		return fmt.Errorf("failed to update product %d: %w", id, err)
	}

	rowsAffected := res.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("failed to find product id")
	}
	return nil
}

func (r *Repository) DeleteProduct(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product %d: %w", id, err)
	}
	rowsAffected := res.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}
	return nil
}

func (r *Repository) CreateBrand(ctx context.Context, title string, description string) (int, error) {
	query := `
		INSERT INTO brands (title, description)
		VALUES ($1, $2)
		RETURNING id
	`
	var newID int
	err := r.db.QueryRow(ctx, query, title, description).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new brand: %w", err)
	}
	return newID, nil
}

func (r *Repository) GetBrand(ctx context.Context, id int) (models.Brand, error) {
	var b models.Brand
	query := `SELECT id, title, description FROM brands WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&b.ID, &b.Title, &b.Description)
	if err != nil {
		return b, fmt.Errorf("failed to get brand: %w", err)
	}
	return b, nil
}

func (r *Repository) GetBrandsPage(ctx context.Context, limit, offset int) ([]models.Brand, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, description FROM brands ORDER BY id LIMIT $1 OFFSET $2;`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("executing get brands page: %w", err)
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

func (r *Repository) GetBrandsCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM brands`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting brands: %w", err)
	}
	return count, nil
}

func (r *Repository) UpdateBrand(ctx context.Context, id int, title string, description string) error {
	query := `UPDATE brands SET title = $1, description = $2 WHERE id = $3`
	res, err := r.db.Exec(ctx, query, title, description, id)
	if err != nil {
		return fmt.Errorf("failed to update brand %d: %w", id, err)
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("brand with id %d not found", id)
	}
	return nil
}

func (r *Repository) DeleteBrand(ctx context.Context, id int) error {
	query := `DELETE FROM brands WHERE id = $1`
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete brand %d: %w", id, err)
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("brand with id %d not found", id)
	}
	return nil
}

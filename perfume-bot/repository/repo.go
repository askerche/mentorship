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
	_, err := r.db.Exec(ctx, queryCreateUserIfNotExists, telegram_id)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *Repository) AddToCart(
	ctx context.Context,
	telegramID int64,
	productID int,
) error {
	_, err := r.db.Exec(ctx, queryAddToCart, telegramID, productID)
	if err != nil {
		return fmt.Errorf("failed to add item to cart: %w", err)
	}
	return nil
}

func (r *Repository) GetCartByTelegramID(
	ctx context.Context,
	telegramID int64,
) ([]models.CartItem, error) {
	rows, err := r.db.Query(
		ctx, queryGetCartByTelegramID, telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	var items []models.CartItem

	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(
			&item.ProductID,
			&item.Title,
			&item.BrandName,
			&item.Price,
			&item.Quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) ClearCart(ctx context.Context, telegram_id int64) error {
	_, err := r.db.Exec(ctx, queryClearCart, telegram_id)
	if err != nil {
		return fmt.Errorf("failed to clear cart for user %d: %w", telegram_id, err)
	}
	return nil
}

func (r *Repository) CreateOrderFromCart(ctx context.Context, telegramID int, username string) (int, error) {
	rows, err := r.db.Query(ctx, querySelectItemsFromCart, telegramID)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	var order []models.OrderProduct
	totalPrice := 0

	for rows.Next() {
		var p models.OrderProduct

		err := rows.Scan(&p.Id, &p.Price, &p.Quantity)
		if err != nil {
			return 0, fmt.Errorf("failed to scan cart row: %w", err)
		}
		order = append(order, p)
		totalPrice += p.Price * p.Quantity
	}
	if len(order) == 0 {
		return 0, fmt.Errorf("cart is empty")
	}

	var orderID int
	err = r.db.QueryRow(ctx, queryCreateOrder, telegramID, username, totalPrice).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert into orders: %w", err)
	}

	for _, s := range order {
		_, err = r.db.Exec(ctx, queryInsertOrderItems, orderID, s.Id, s.Price, s.Quantity)
		if err != nil {
			return 0, fmt.Errorf("failed to insert order_item: %w", err)
		}
	}
	_, err = r.db.Exec(ctx, queryClearCart, telegramID)
	if err != nil {
		return 0, fmt.Errorf("failed to clear user cart: %w", err)
	}
	return orderID, nil
}

func (r *Repository) GetAllProductsPage(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.Product, error) {
	rows, err := r.db.Query(
		ctx, queryGetAllProductsPage, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Price,
			&p.Description,
			&p.Created_at,
			&p.ImageFileID,
			&p.Brand.ID,
			&p.Brand.Title,
			&p.CategoryIDs,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scanning product row: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repository) GetAllProductsCounts(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(
		ctx, queryGetAllProductsCounts).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting brands: %w", err)
	}
	return count, nil
}

func (r *Repository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := r.db.Query(ctx, queryGetAllCategories)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		err = rows.Scan(&c.ID, &c.Title)
		if err != nil {
			return nil, fmt.Errorf(
				"scanning category row: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *Repository) GetProductsPageByBrandID(
	ctx context.Context,
	brandID int,
	limit int,
	offset int,
) ([]models.Product, error) {
	rows, err := r.db.Query(
		ctx,
		queryGetProductsPageByBrandID,
		brandID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()
	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Price,
			&p.Description,
			&p.Created_at,
			&p.ImageFileID,
			&p.Brand.ID,
			&p.Brand.Title,
			&p.Brand.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning product by brand row: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repository) GetProductsPageByCategoryID(
	ctx context.Context,
	categoryID int,
	limit int,
	offset int,
) ([]models.Product, error) {
	rows, err := r.db.Query(
		ctx,
		queryGetProductsPageByCategoryID,
		categoryID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()
	products := make([]models.Product, 0)

	for rows.Next() {
		var p models.Product
		p.CategoryIDs = []int64{int64(categoryID)}
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Price,
			&p.Description,
			&p.Created_at,
			&p.ImageFileID,
			&p.Brand.ID,
			&p.Brand.Title,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning product by category row: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repository) GetProductsCountByCategoryID(
	ctx context.Context,
	categoryID int,
) (int, error) {
	var count int
	err := r.db.QueryRow(
		ctx, queryGetProductsCountByCategoryID, categoryID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf(
			"counting products for category: %w", err)
	}
	return count, nil
}

func (r *Repository) GetProductsCountByBrandID(
	ctx context.Context, brandID int,
) (int, error) {
	var count int
	err := r.db.QueryRow(
		ctx, queryGetProductsCountByBrandID, brandID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting products for brand: %w", err)
	}
	return count, nil
}

func (r *Repository) CreateProduct(
	ctx context.Context,
	brandID int,
	title string,
	price int,
	description string,
	image_file_id string,
) (int, error) {
	var newID int
	err := r.db.QueryRow(
		ctx,
		queryCreateProduct,
		brandID,
		title,
		price,
		description,
		image_file_id,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new product: %w", err)
	}
	return newID, nil
}

func (r *Repository) UpdateProductCategories(
	ctx context.Context,
	productID int,
	categoriesIDs []int,
) error {
	_, err := r.db.Exec(ctx, queryDeleteProductCategories, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product categories before update: %w", err)
	}
	_, err = r.db.Exec(ctx, queryInsertProductCategories, productID, categoriesIDs)
	if err != nil {
		return fmt.Errorf("failed to insert product categories: %w", err)
	}
	return nil
}

func (r *Repository) GetProductByID(
	ctx context.Context, productID int,
) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(ctx, queryGetProductByID, productID).Scan(
		&p.ID,
		&p.Title,
		&p.Price,
		&p.Description,
		&p.Created_at,
		&p.ImageFileID,
		&p.Brand.ID,
		&p.Brand.Title,
		&p.CategoryIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan product: %w", err)
	}
	return &p, nil
}

func (r *Repository) UpdateProduct(
	ctx context.Context,
	id int,
	brandID int,
	title string,
	price int,
	description string,
	imageFileID string,
) error {
	res, err := r.db.Exec(
		ctx, queryUpdateProduct, brandID, title, price, description, imageFileID, id)
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
	_, err := r.db.Exec(ctx, queryDeleteProductCategories, id)
	if err != nil {
		return fmt.Errorf("failed to delete product categories: %w", err)
	}
	res, err := r.db.Exec(ctx, queryDeleteProduct, id)
	if err != nil {
		return fmt.Errorf("failed to delete product %d: %w", id, err)
	}
	rowsAffected := res.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}
	return nil
}

func (r *Repository) CreateBrand(
	ctx context.Context,
	title string,
	description string,
) (int, error) {
	var newID int
	err := r.db.QueryRow(ctx, queryCreateBrand, title, description).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new brand: %w", err)
	}
	return newID, nil
}

func (r *Repository) GetBrand(ctx context.Context, id int) (models.Brand, error) {
	var b models.Brand
	err := r.db.QueryRow(ctx, queryGetBrand, id).Scan(
		&b.ID, &b.Title, &b.Description)
	if err != nil {
		return b, fmt.Errorf("failed to get brand: %w", err)
	}
	return b, nil
}

func (r *Repository) GetBrandsPage(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.Brand, error) {
	rows, err := r.db.Query(ctx, queryGetBrandsPage, limit, offset)
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
	err := r.db.QueryRow(ctx, queryGetBrandsCount).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting brands: %w", err)
	}
	return count, nil
}

func (r *Repository) UpdateBrand(
	ctx context.Context,
	id int,
	title string,
	description string,
) error {
	res, err := r.db.Exec(ctx, queryUpdateBrand, title, description, id)
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
	res, err := r.db.Exec(ctx, queryDeleteBrand, id)
	if err != nil {
		return fmt.Errorf("failed to delete brand %d: %w", id, err)
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("brand with id %d not found", id)
	}
	return nil
}

func (r *Repository) CreateCategory(
	ctx context.Context,
	title string,
) (int, error) {
	var newID int
	err := r.db.QueryRow(ctx, queryCreateCategory, title).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new category: %w", err)
	}
	return newID, nil
}

func (r *Repository) UpdateCategory(
	ctx context.Context,
	id int,
	title string,
) error {
	res, err := r.db.Exec(ctx, queryUpdateCategory, title, id)
	if err != nil {
		return fmt.Errorf("failed to update category %d: %w", id, err)
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("category with id %d not found", id)
	}
	return nil
}

func (r *Repository) DeleteCategory(ctx context.Context, id int) error {
	res, err := r.db.Exec(ctx, queryDeleteCategory, id)
	if err != nil {
		return fmt.Errorf("failed to delete category %d: %w", id, err)
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("category with id %d not found", id)
	}
	return nil
}

func (r *Repository) GetAdminOrders(ctx context.Context, limit int, offset int) ([]models.Order, int, error) {
	rows, err := r.db.Query(ctx, queryGetAdminOrders, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get orders: %w", err)
	}
	orders := make([]models.Order, 0)
	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.ID, &o.TelegramID, &o.Username, &o.TotalPrice, &o.Status, &o.Created_at)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order row: %w", err)
		}
		orders = append(orders, o)
	}
	var totalOrders int
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders`).Scan(&totalOrders)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}
	return orders, totalOrders, nil
}

func (r *Repository) GetAdminOrder(
	ctx context.Context, orderID int,
) ([]models.OrderItemDetail, error) {
	rows, err := r.db.Query(ctx, queryGetAdminOrder, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders items: %w", err)
	}
	items := make([]models.OrderItemDetail, 0)
	for rows.Next() {
		var item models.OrderItemDetail
		err := rows.Scan(
			&item.ProductID,
			&item.Title,
			&item.BrandName,
			&item.PriceAtPurchase,
			&item.Quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan items: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) UpdateOrderStatus(
	ctx context.Context,
	orderID int,
	status string,
) error {
	res, err := r.db.Exec(ctx, queryUpdateOrderStatus, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status %d: %w", orderID, err)
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("order with id %d not found", orderID)
	}
	return nil
}

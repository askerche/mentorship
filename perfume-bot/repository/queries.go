package repository

const queryCreateUserIfNotExists = `INSERT INTO users (telegram_id) VALUES($1) ON CONFLICT (telegram_id) DO NOTHING`
const queryAddToCart = `
			INSERT INTO cart_items (telegram_id, product_id, quantity)
			VALUES($1, $2, 1)
			ON CONFLICT (telegram_id, product_id)
			DO UPDATE SET quantity = cart_items.quantity + 1
`
const queryGetCartByTelegramID = `
				SELECT p.id, p.title, b.title, p.price, c.quantity
				FROM cart_items c
				JOIN products p ON p.id = c.product_id
				JOIN brands b ON b.id = p.brand_id
				WHERE c.telegram_id = $1
				ORDER BY c.created_at ASC
				`
const queryGetAllProductsPage = `
			SELECT
			        p.id,
					p.title,
					p.price, 
					COALESCE(p.description, ''), 
					p.created_at, 
					COALESCE(p.image_file_id, ''),
                    b.id, 
					b.title,
					COALESCE (ARRAY_AGG(pc.category_id) FILTER (WHERE pc.category_id IS NOT NULL), ARRAY[]::bigint[]) AS category_ids 
			FROM products p 
			JOIN brands b ON p.brand_id = b.id
			LEFT JOIN products_categories pc ON p.id = pc.product_id
            GROUP BY p.id, b.id
			ORDER BY p.id DESC
			LIMIT $1
			OFFSET $2;
			`
const queryGetAllProductsCounts = `SELECT COUNT(*) FROM products`
const queryGetAllCategories = `SELECT id, title FROM categories ORDER BY id;`
const queryGetProductsPageByBrandID = `
	SELECT p.id, p.title, p.price, COALESCE(p.description, ''), p.created_at, COALESCE(p.image_file_id, ''),
	b.id, b.title, COALESCE(b.description, '')
	FROM products p
	JOIN brands b 
	ON b.id = p.brand_id
	WHERE b.id = $1
	ORDER BY p.id DESC
	LIMIT $2
	OFFSET $3
	`
const queryGetProductsPageByCategoryID = `
	SELECT  p.id, p.title, p.price, COALESCE(p.description, ''), p.created_at, COALESCE(p.image_file_id, ''), b.id, b.title
	FROM products p
	JOIN products_categories pc
	ON p.id = pc.product_id
	JOIN brands b 
	ON b.id = p.brand_id
	WHERE pc.category_id = $1
	ORDER BY p.id DESC
	LIMIT $2
	OFFSET $3
	`

const queryGetProductsCountByCategoryID = `SELECT COUNT(*) FROM products p JOIN products_categories pc ON p.id = pc.product_id WHERE pc.category_id = $1`
const queryGetProductsCountByBrandID = `SELECT COUNT (*) FROM products p 
							   JOIN brands b
							   ON b.id = p.brand_id 
							   WHERE b.id = $1`
const queryCreateProduct = `
			  INSERT INTO products (brand_id, title, price, description, image_file_id)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id			
	`
const queryGetProductByID = `
		SELECT 
			p.id, 
			p.title, 
			p.price, 
			COALESCE(p.description, ''), 
			p.created_at,
			COALESCE(p.image_file_id, ''),
			b.id, 
			b.title,
		    COALESCE (ARRAY_AGG(pc.category_id) FILTER (WHERE pc.category_id IS NOT NULL), ARRAY[]::bigint[]) AS category_ids 
		
		FROM products p
		JOIN brands b ON p.brand_id = b.id
		LEFT JOIN products_categories pc ON p.id = pc.product_id
		WHERE p.id = $1
		GROUP BY p.id, b.id
	`

const queryUpdateProduct = `UPDATE products SET brand_id = $1, title = $2, price = $3, description = $4, image_file_id = $5 WHERE id = $6`
const queryDeleteProduct = `DELETE FROM products WHERE id = $1`
const queryGetBrand = `SELECT id, title, description FROM brands WHERE id = $1`
const queryGetBrandsPage = `SELECT id, title, description FROM brands ORDER BY id LIMIT $1 OFFSET $2;`
const queryGetBrandsCount = `SELECT COUNT(*) FROM brands`
const queryUpdateBrand = `UPDATE brands SET title = $1, description = $2 WHERE id = $3`
const queryDeleteBrand = `DELETE FROM brands WHERE id = $1`
const queryCreateBrand = `INSERT INTO brands (title, description)VALUES ($1, $2)RETURNING id`
const queryDeleteProductCategories = `DELETE FROM products_categories WHERE product_id = $1`
const queryInsertProductCategories = `INSERT INTO products_categories (product_id, category_id) SELECT $1, unnest($2::bigint[])`
const queryCreateCategory = `INSERT INTO categories (title)VALUES ($1) RETURNING id`
const queryDeleteCategory = `DELETE FROM category WHERE id = $1`
const queryUpdateCategory = `UPDATE categories SET title = $1 WHERE id = $2`
const queryClearCart = `DELETE FROM cart_items WHERE telegram_id = $1`

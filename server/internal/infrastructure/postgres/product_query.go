package postgres

const (
	createProductQuery = `
		INSERT INTO products (id, organization_id, category_id, sku, name, unit,net_unit,net_amount, product_type, show_on_web, cost_price, sell_price, stock_quantity, min_stock, is_active)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,$14, $15)
	`
	listProductsByOrgIDQuery = `
		SELECT id, organization_id, category_id, sku, name, unit,net_unit,net_amount, product_type, show_on_web,
		       cost_price, sell_price, stock_quantity, min_stock, is_active, created_at, updated_at
		FROM products
		WHERE organization_id = $1
	`
	searchProductsFilter = `
		AND ($2 = '' OR name ILIKE $2 OR sku ILIKE $2 OR product_type ILIKE $2)
	`
	filterProductsStatus = `
		AND ($3 = '' 
		OR ($3 = 'active' AND is_active = TRUE)
		OR ($3 = 'low_stock' AND stock_quantity <= min_stock AND is_active = TRUE)
		OR ($3 = 'out_of_stock' AND stock_quantity = 0 AND is_active = TRUE))
	`
	
	sortProductsFilter = `
		ORDER BY 
			CASE WHEN $4 = 'created_at_desc' THEN created_at END DESC,
			CASE WHEN $4 = 'price_asc' THEN sell_price END ASC,
			CASE WHEN $4 = 'price_desc' THEN sell_price END DESC,
			CASE WHEN $4 = 'stock_asc' THEN stock_quantity END ASC,
			created_at DESC
	`
	countProductsByOrgIDQuery = `
		SELECT COUNT(*)
		FROM products
		WHERE organization_id = $1
	`
	updateProductQuery = `
		UPDATE products
		SET category_id = $3, sku = $4, name = $5, unit = $6, product_type = $7, show_on_web = $8,
		    cost_price = $9, sell_price = $10, stock_quantity = $11, min_stock = $12, is_active = $13
		WHERE id = $1 AND organization_id = $2
	`
	deleteProductQuery = `
		DELETE FROM products
		WHERE id = $1 AND organization_id = $2
	`
	getProductByIdQuery =`
		SELECT id, sku, name, unit,net_unit,net_amount, cost_price, sell_price
		FROM products
		WHERE id = $1 AND organization_id = $2
	`
)

package postgres

const (
	createOrderQuery = `
		INSERT INTO orders (id, organization_id, code, customer_id, appointment_id, subtotal_amount, discount_amount, total_amount, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	bulkInsertOrderItemsQuery = `
		INSERT INTO order_items (id, organization_id, order_id, item_type, service_id, product_id, staff_id, quantity, unit_price, line_total)
		VALUES
	`
	bulkInsertOrderItemMaterialsQuery = `
		INSERT INTO order_item_materials (id, order_item_id, product_id, quantity, unit_cost_at_time)
		VALUES
	`
	// Trừ kho có điều kiện chống âm, dùng cho bán hàng và vật tư tiêu hao.
	deductStockQuery = `
		UPDATE products
		SET stock_quantity = stock_quantity - $3
		WHERE id = $1 AND organization_id = $2 AND stock_quantity >= $3
	`
	// Giữ chỗ mềm lúc tạo đơn: tăng reserved khi tồn khả dụng (stock - reserved) còn đủ.
	reserveStockQuery = `
		UPDATE products
		SET reserved_quantity = reserved_quantity + $3
		WHERE id = $1 AND organization_id = $2 AND stock_quantity - reserved_quantity >= $3
	`
	// Nhả giữ chỗ khi hủy đơn chưa phục vụ (idempotent, không bao giờ âm).
	releaseReserveQuery = `
		UPDATE products
		SET reserved_quantity = GREATEST(reserved_quantity - $3, 0)
		WHERE id = $1 AND organization_id = $2
	`
	// Trừ kho thật lúc phục vụ: trừ cả tồn và giữ chỗ, chỉ cần tồn đủ (giữ chỗ có thể đã hết hạn).
	serveDeductStockQuery = `
		UPDATE products
		SET stock_quantity = stock_quantity - $3,
		    reserved_quantity = GREATEST(reserved_quantity - $3, 0)
		WHERE id = $1 AND organization_id = $2 AND stock_quantity >= $3
	`
	// Chuyển paid -> serving, guard chống serve 2 lần / serve đơn chưa thu tiền.
	flipToServingQuery = `
		UPDATE orders
		SET status = 'serving'
		WHERE organization_id = $1 AND id = $2 AND status = 'paid'
	`
	// Vô hiệu phiếu xuất của đơn (no-op nếu chưa từng có phiếu).
	voidSaleSlipQuery = `
		UPDATE inventory_documents
		SET status = 'cancelled'
		WHERE organization_id = $1 AND order_id = $2 AND status = 'completed'
	`
	// Khóa dòng đơn để serialize các payment/serve đồng thời.
	lockOrderRowQuery = `
		SELECT status
		FROM orders
		WHERE organization_id = $1 AND id = $2
		FOR UPDATE
	`
	// Hoàn kho luôn thành công (cộng lại, không cần điều kiện).
	restoreStockQuery = `
		UPDATE products
		SET stock_quantity = stock_quantity + $3
		WHERE id = $1 AND organization_id = $2
	`
	insertInventoryLedgerQuery = `
		INSERT INTO inventory_transactions (id, organization_id, product_id, type, quantity, reference_type, reference_id, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	orderColumns = `id, organization_id, code, customer_id, appointment_id, subtotal_amount, discount_amount, total_amount, status, created_by, created_at, updated_at`
	getOrderByIDQuery = `
		SELECT ` + orderColumns + `
		FROM orders
		WHERE organization_id = $1 AND id = $2
	`
	updateOrderStatusQuery = `
		UPDATE orders
		SET status = $3
		WHERE organization_id = $1 AND id = $2
	`
	customerExistsInOrgQuery = `
		SELECT EXISTS(SELECT 1 FROM customers WHERE organization_id = $1 AND id = $2)
	`
	appointmentExistsInOrgQuery = `
		SELECT EXISTS(SELECT 1 FROM appointments WHERE organization_id = $1 AND id = $2)
	`
	productPricingQuery = `
		SELECT id, sell_price, cost_price
		FROM products
		WHERE organization_id = $1 AND id = ANY($2)
	`
	serviceMaterialsWithCostQuery = `
		SELECT sm.service_id, sm.product_id, sm.quantity, p.cost_price
		FROM service_materials sm
		JOIN products p ON p.id = sm.product_id AND p.organization_id = sm.organization_id
		WHERE sm.organization_id = $1 AND sm.service_id = ANY($2)
	`
	servicePriceQuery = `
		SELECT price
		FROM services
		WHERE organization_id = $1 AND id = $2
	`
	findPaymentByProviderTxnQuery = `
		SELECT id, organization_id, order_id, method, provider, provider_txn_id, amount, status, paid_at, created_at
		FROM payments
		WHERE provider = $1 AND provider_txn_id = $2
	`
	createPaymentQuery = `
		INSERT INTO payments (id, organization_id, order_id, method, provider, provider_txn_id, amount, status, paid_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	totalPaidByOrderQuery = `
		SELECT COALESCE(SUM(amount), 0)
		FROM payments
		WHERE organization_id = $1 AND order_id = $2 AND status = 'succeeded'
	`
	listOrderPaymentsQuery = `
		SELECT id, organization_id, order_id, method, provider, provider_txn_id, amount, status, paid_at, created_at
		FROM payments
		WHERE organization_id = $1 AND order_id = $2
		ORDER BY created_at ASC
	`
	listOrderItemsQuery = `
		SELECT id, organization_id, order_id, item_type, service_id, product_id, staff_id, quantity, unit_price, line_total
		FROM order_items
		WHERE organization_id = $1 AND order_id = $2
	`
	listOrderMaterialsQuery = `
		SELECT m.id, m.order_item_id, m.product_id, m.quantity, m.unit_cost_at_time
		FROM order_item_materials m
		JOIN order_items i ON i.id = m.order_item_id
		WHERE i.organization_id = $1 AND i.order_id = $2
	`

    listItemPosQuery = `
        SELECT * FROM (
            SELECT id, organization_id, category_id, 'product' AS type, name, NULL AS description, NULL AS duration_minutes, unit, sell_price AS price, stock_quantity, is_active
            FROM products p
            WHERE p.organization_id = $1 AND p.is_active = true 
            AND (product_type = $2 OR product_type = $3)
            UNION ALL
            SELECT id, organization_id, category_id, 'service' AS type, name, description, duration_minutes, NULL AS unit, price, NULL AS stock_quantity, is_active
            FROM services s
            WHERE s.organization_id = $1 AND s.is_active = true
        ) AS pos_items
        WHERE 1 = 1
    `

    searchItemPosQuery = `
        AND ($4 = '' OR unaccent(name) ILIKE unaccent($4))
    `
    
    filterItemPosQuery = `
        AND ($5 = '' OR category_id::text = $5) 
    `
    
    // Sắp xếp lại thứ tự con trỏ khớp với ORDER BY: (type, name, id)
    // NULLIF + cast ::uuid để so sánh đúng kiểu UUID khi phân trang;
    // nhánh $6 = '' short-circuit nên không bao giờ cast chuỗi rỗng.
    paginateItemPosQuery = `
        AND ($6 = '' OR (type, name, id) > ($6, $7, NULLIF($8, '')::uuid))
        ORDER BY type ASC, name ASC, id ASC
        LIMIT $9
    `
)

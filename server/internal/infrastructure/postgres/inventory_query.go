package postgres

const (
	createInventoryDocumentQuery =`
		INSERT INTO inventory_documents (id,organization_id,supplier_id,order_id,document_code,type,total_amount,reason,note,created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10)
	`
	createInventoryDocumentItemQuery = `
		INSERT INTO inventory_document_items (id,document_id,product_id,product_name,sku,unit, net_unit, net_amount, quantity, unit_price, total_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	getProductsByIDsQuery = `
		SELECT id, sku, name, unit, net_unit, net_amount
		FROM products
		WHERE organization_id = $1 AND id = ANY($2)
	`
	createInventoryTransactionQuery = `
		INSERT INTO inventory_transactions (id, organization_id, product_id, type, quantity, reference_type, reference_id, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	// Trừ/cộng kho có điều kiện chống âm, trả về tồn mới để usecase không cần đọc lại.
	applyStockChangeQuery = `
		UPDATE products
		SET stock_quantity = stock_quantity + $3
		WHERE id = $1 AND organization_id = $2 AND stock_quantity + $3 >= 0
		RETURNING stock_quantity
	`
	listInventoryTransactionsQuery = `
		SELECT id, organization_id, product_id, type, quantity, reference_type, reference_id, created_by, created_at
		FROM inventory_transactions
		WHERE organization_id = $1 AND ($2 = '' OR product_id::text = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	listInventoryDocumentsQuery = `
		SELECT
			id.id AS id,
			id.organization_id AS organization_id,
			id.document_code AS document_code,
			id.type AS type,
			id.total_amount AS total_amount,
			id.note AS note,
			id.reason AS reason,
			s.id AS supplier_id,
			s.name AS supplier_name,

			o.id AS order_id,
			o.code AS order_code,

			u.id AS created_by,
			u.full_name AS created_by_name,

			id.created_at AS created_at,

			(
				SELECT COUNT(*)
				FROM inventory_document_items idi
				WHERE idi.document_id = id.id
			) AS items_count
		FROM inventory_documents id
		LEFT JOIN suppliers s ON s.id = id.supplier_id
		LEFT JOIN orders o ON o.id = id.order_id
		LEFT JOIN users u ON u.id = id.created_by
		WHERE id.organization_id = $1
	`

	// 2. Phần lọc tìm kiếm (Search trên nhiều cột)
	searchInventoryFilter = `
		AND ($2 = '' OR 
			id.document_code ILIKE $2 OR 
			s.name ILIKE $2 OR 
			o.code ILIKE $2 OR 
			u.full_name ILIKE $2 OR 
			u.phone ILIKE $2
		)
	`

	// 3. Phần lọc trạng thái / loại phiếu (Filter theo type)
	filterInventoryType = `
		AND ($3 = '' OR id.type = $3)
	`

	// 4. Phần sắp xếp (Sort theo yêu cầu)
	sortInventoryFilter = `
		ORDER BY 
			CASE WHEN $4 = 'total_amount_asc' THEN id.total_amount END ASC,
			CASE WHEN $4 = 'total_amount_desc' THEN id.total_amount END DESC,
			CASE WHEN $4 = 'created_at_asc' THEN id.created_at END ASC,
			CASE WHEN $4 = 'created_at_desc' THEN id.created_at END DESC,
			id.created_at DESC
	`
	countInventoryDocumentsQuery = `
		SELECT COUNT(*)
		FROM inventory_documents
		WHERE organization_id = $1
	`
	countInventoryTransactionsQuery = `
		SELECT COUNT(*)
		FROM inventory_transactions
		WHERE organization_id = $1 AND ($2 = '' OR product_id::text = $2)
	`
	listLowStockProductsQuery = `
		SELECT id AS product_id, sku, name, stock_quantity, min_stock
		FROM products
		WHERE organization_id = $1 AND stock_quantity <= min_stock
		ORDER BY stock_quantity ASC
		LIMIT $2 OFFSET $3
	`
	countLowStockProductsQuery = `
		SELECT COUNT(*)
		FROM products
		WHERE organization_id = $1 AND stock_quantity <= min_stock
	`
	InventoryDocumentDetailsQuery = `
			SELECT 
				id.id, 
				id.organization_id, 
				id.document_code, 
				id.type, 
				id.total_amount, 
				id.note, 
				id.reason,            -- Đã bổ sung
				id.created_by,        -- Đã bổ sung
				u.full_name AS created_by_name,
				id.created_at, 
				o.id AS order_id, 
				o.code AS order_code, 
				s.id AS supplier_id, 
				s.name AS supplier_name, 
				s.phone AS supplier_phone,
				s.address AS supplier_address, 
				COALESCE(
					jsonb_agg(
					jsonb_build_object(
						'id', idi.id,
						'document_id', idi.document_id,
						'product_id', idi.product_id,
						'product_name', idi.product_name,
						'sku', idi.sku,
						'unit', idi.unit,
						'net_unit', idi.net_unit,
						'net_amount', idi.net_amount,
						'quantity', idi.quantity,
						'unit_price', idi.unit_price,
						'total_price', idi.total_price,
						-- Quy đổi hiển thị: có net đầy đủ và đủ 1 đơn vị lớn thì show theo unit lớn,
						-- ngược lại giữ nguyên theo đơn vị nhỏ (số đã lưu luôn ở đơn vị nhỏ nhất)
						'display_quantity', CASE
							WHEN idi.net_unit IS NOT NULL AND idi.net_unit <> ''
								AND COALESCE(idi.net_amount, 0) > 0
								AND idi.quantity >= idi.net_amount
							THEN idi.quantity / idi.net_amount
							ELSE idi.quantity
						END,
						'display_unit', CASE
							WHEN idi.net_unit IS NOT NULL AND idi.net_unit <> ''
								AND COALESCE(idi.net_amount, 0) > 0
								AND idi.quantity >= idi.net_amount
							THEN idi.unit
							ELSE COALESCE(NULLIF(idi.net_unit, ''), idi.unit)
						END,
						'display_unit_price', CASE
							WHEN idi.net_unit IS NOT NULL AND idi.net_unit <> ''
								AND COALESCE(idi.net_amount, 0) > 0
								AND idi.quantity >= idi.net_amount
							THEN idi.unit_price * idi.net_amount
							ELSE idi.unit_price
						END
					)
					) FILTER (WHERE idi.id IS NOT NULL), 
					'[]'::jsonb
				) AS items
		FROM inventory_documents id
		LEFT JOIN orders o ON o.id = id.order_id
		LEFT JOIN suppliers s ON s.id = id.supplier_id
		LEFT JOIN inventory_document_items idi ON idi.document_id = id.id
		LEFT JOIN users u ON u.id = id.created_by
		WHERE id.organization_id =$1 
		AND id.id =$2
		GROUP BY id.id, o.id, s.id, u.full_name;
	`
	listInventoryTransactionsByProductIDQuery = `
		SELECT
			-- ===== 1. KHỐI THÔNG TIN CHÍNH (Bảng transactions) =====
			it.id,
			it.organization_id,
			it.type,
			it.quantity,
			it.reference_type,
			it.reference_id,
			it.created_at,

			it.product_id,
			id.id AS document_id, 
			id.document_code, 

			o.id AS order_id, 
			o.code AS order_code, 
			s.id AS supplier_id, 
			s.name AS supplier_name

		FROM inventory_transactions it
		LEFT JOIN inventory_documents id ON id.id = it.reference_id
		LEFT JOIN products p ON p.id = it.product_id
		LEFT JOIN orders o ON o.id = id.order_id
		LEFT JOIN suppliers s ON s.id = id.supplier_id
		WHERE it.organization_id = $1 AND it.product_id = $2
	`
	FilterInventoryTransactionsByType = `
		AND ($3 = '' OR it.type = $3)
		AND ($4 = '' OR DATE(it.created_at) = NULLIF($4, '')::DATE)
	`
	sortInventoryTransactionsFilter = `
		ORDER BY 
		CASE WHEN $5 = 'created_desc' THEN it.created_at END DESC, 
		CASE WHEN $5 = 'created_asc' THEN it.created_at END ASC,
		it.created_at DESC
	`
	countInventoryTransactionsByProductIDQuery = `
		SELECT COUNT(*)
		FROM inventory_transactions it
		WHERE it.organization_id = $1 AND it.product_id = $2
	`
)

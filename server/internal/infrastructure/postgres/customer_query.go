package postgres

const (
	createCustomerQuery = `
		INSERT INTO customers (id, organization_id, full_name, phone, gender, birth_date, note)
		VALUES($1, $2, $3, $4, $5, $6, $7)
	`
	listCustomersByOrgIDQuery = `
		SELECT id, organization_id, full_name, phone, gender, birth_date, note,
		       total_spent, total_visits, created_at, updated_at, deleted_at
		FROM customers
		WHERE organization_id = $1 AND deleted_at IS NULL
	`
	countCustomersByOrgIDQuery = `
		SELECT COUNT(*)
		FROM customers
		WHERE organization_id = $1 AND deleted_at IS NULL
	`
	// searchCustomersFilter lọc theo tên hoặc phone (ILIKE, không phân biệt hoa thường).
	// $2 là từ khóa đã bọc %...% ở tầng repo.
	searchCustomersFilter = `
		AND (full_name ILIKE $2 OR phone ILIKE $2)
	`
	getCustomerByIDQuery = `
		SELECT id, organization_id, full_name, phone, gender, birth_date, note,
		       total_spent, total_visits, created_at, updated_at, deleted_at
		FROM customers
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`
	updateCustomerQuery = `
		UPDATE customers
		SET full_name = $3, phone = $4, gender = $5, birth_date = $6, note = $7
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`
	deleteCustomerQuery = `
		UPDATE customers
		SET deleted_at = now()
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`
	findCustomerByPhoneQuery = `
		SELECT *
		FROM customers
		WHERE organization_id = $1 AND phone = $2 AND deleted_at IS NULL
		LIMIT 1
	`
)

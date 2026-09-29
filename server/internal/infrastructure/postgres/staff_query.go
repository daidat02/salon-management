package postgres

const (
	createStaffQuery = `
		INSERT INTO staff (id, organization_id, code, full_name, phone, position, commission_rate, status, hire_date)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	getStaffsByOrgIDQuery = `
		SELECT *
		FROM staff
		WHERE organization_id = $1 AND deleted_at IS NULL
	`
	searchStaffsFilter = `
		AND ($2 = '' OR full_name ILIKE $2 OR code ILIKE $2 OR phone ILIKE $2 OR position ILIKE $2)
	`
	filterStaffsStatus = `
		AND ($3 = '' OR status = $3)
	`
	sortStaffsFilter = `
		ORDER BY
			CASE WHEN $4 = 'name_asc' THEN full_name END ASC,
			CASE WHEN $4 = 'name_desc' THEN full_name END DESC,
			CASE WHEN $4 = 'hire_date_asc' THEN hire_date END ASC,
			CASE WHEN $4 = 'hire_date_desc' THEN hire_date END DESC,
			created_at DESC
	`
	countStaffsByOrgIDQuery = `
		SELECT COUNT(*)
		FROM staff
		WHERE organization_id = $1 AND deleted_at IS NULL`
	updateStaffQuery = `
		UPDATE staff
		SET code = $3, full_name = $4, phone = $5, position = $6,
		    commission_rate = $7, status = $8, hire_date = $9
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`
	deleteStaffQuery = `
		UPDATE staff
		SET deleted_at = now()
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`
)
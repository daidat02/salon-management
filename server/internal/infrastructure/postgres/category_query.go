package postgres

const (
	createCategoryQuery = `
		INSERT INTO categories (id, organization_id, type, name, description, is_active)
		VALUES($1, $2, $3, $4, $5, $6)
	`
	listCategoriesByOrgIDQuery = `
		SELECT id, organization_id, type, name, description, is_active, created_at, updated_at
		FROM categories
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	countCategoriesByOrgIDQuery = `
		SELECT COUNT(*)
		FROM categories
		WHERE organization_id = $1
	`
	updateCategoryQuery = `
		UPDATE categories
		SET type = $3, name = $4, description = $5, is_active = $6
		WHERE id = $1 AND organization_id = $2
	`
	deleteCategoryQuery = `
		DELETE FROM categories
		WHERE id = $1 AND organization_id = $2
	`
)

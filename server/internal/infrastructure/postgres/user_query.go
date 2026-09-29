package postgres

const (

	createUserQuery = `
		INSERT INTO users (id,organization_id, email, password_hash, full_name, phone, role, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	findUserByPhoneNumberQuery = `
		SELECT id, organization_id, email, password_hash, full_name, phone, role, status,last_login_at,created_at, updated_at
		FROM users
		WHERE phone = $1
		Limit 1
	`
	
	
)
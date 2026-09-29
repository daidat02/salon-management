package postgres

const (
	createAppointmentQuery = `
		INSERT INTO appointments (id, organization_id, code, customer_id, staff_id, start_time, end_time, status, source, total_amount, note, device_id, client_ip, required_deposit_amount, deposit_status)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	// appointment_items không có unique nên bulk INSERT thường (không upsert).
	bulkInsertAppointmentItemsQuery = `
		INSERT INTO appointment_items (id, organization_id, appointment_id, item_type, service_id, product_id, quantity, unit_price, line_total)
		VALUES
	`
	appointmentColumns = `
		id, organization_id, code, customer_id, staff_id, start_time, end_time,
		status, source, total_amount, note, cancelled_reason, device_id, client_ip,
		required_deposit_amount, deposit_status, created_at, updated_at
	`
	getAppointmentByIDQuery = `
		SELECT ` + appointmentColumns + `
		FROM appointments
		WHERE organization_id = $1 AND id = $2
	`
	updateAppointmentTimesQuery = `
		UPDATE appointments
		SET start_time = $3, end_time = $4
		WHERE organization_id = $1 AND id = $2
	`
	updateAppointmentStatusQuery = `
		UPDATE appointments
		SET status = $3, cancelled_reason = NULLIF($4, '')
		WHERE organization_id = $1 AND id = $2
	`
	deleteAppointmentQuery = `
		DELETE FROM appointments
		WHERE organization_id = $1 AND id = $2
	`
	listAppointmentItemsQuery = `
		SELECT id, organization_id, appointment_id, item_type, service_id, product_id,
		       quantity, unit_price, line_total
		FROM appointment_items
		WHERE organization_id = $1 AND appointment_id = $2
		ORDER BY id
	`
)

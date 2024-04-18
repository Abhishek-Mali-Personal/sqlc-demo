-- name: CreateLookup :one
INSERT INTO lookups(
	lookup_id, table_name, display_order, display_text, is_active, parent_id, internal_key, concurrency_key, create_date, create_user_id, update_date, update_user_id, value_text)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING *;

-- name: NewLookupWithConcurrencyKey :one
INSERT INTO lookups(
	table_name, display_order, display_text, is_active, internal_key, concurrency_key, create_date, create_user_id, value_text)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: NewLookup :one
INSERT INTO lookups(
	table_name, display_order, display_text, is_active, create_date, create_user_id, value_text)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetLookup :one
SELECT * FROM lookups WHERE lookup_id = $1;

-- name: GetLookupByID :one
SELECT table_name, display_order, display_text, is_active FROM lookups WHERE lookup_id = $1;

-- name: ListLookupsByDisplayText :many
SELECT table_name, display_order, display_text, is_active FROM lookups WHERE display_text ILIKE $1;

-- name: ListLookupsByTableName :many
SELECT table_name, display_order, display_text, is_active FROM lookups WHERE table_name ILIKE $1;

-- name: ListLookups :many
SELECT table_name, display_order, display_text, is_active FROM lookups;


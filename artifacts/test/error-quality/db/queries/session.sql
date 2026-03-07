-- name: Create :one
INSERT INTO sessions (project_id) VALUES ($1) RETURNING *;

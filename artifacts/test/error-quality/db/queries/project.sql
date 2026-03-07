-- name: FindByID :one
SELECT * FROM projects WHERE id = $1;

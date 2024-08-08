-- name: GerRoom :one
SELECT * FROM rooms
WHERE id = $1 LIMIT 1;
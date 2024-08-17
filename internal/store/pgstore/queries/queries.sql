-- name: GetRoomByID :one
SELECT * FROM rooms
WHERE id = $1 LIMIT 1;

-- name: GetRoomByTheme :one
SELECT * FROM rooms
WHERE theme = $1 LIMIT 1;

-- name: GetRoomByName :one
SELECT * FROM rooms
WHERE name = $1 LIMIT 1;

-- name: ListRooms :many
SELECT * FROM rooms
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateRoom :one
INSERT INTO rooms (theme, name)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE id = $1;

-- name: UpdateRoom :one
UPDATE rooms
SET theme = $2, name = $3
WHERE id = $1
RETURNING *;

-- name: UpdateRoomName :one
UPDATE rooms
SET name = $2
WHERE id = $1
RETURNING *;

-- name: UpdateRoomTheme :one
UPDATE rooms
SET theme = $2
WHERE id = $1
RETURNING *;

-- name: AddReactionFromMessage :one
UPDATE messages
SET reaction_count = reaction_count + 1
WHERE id = $1
RETURNING *;

-- name: RemoveReactionFromMessage :one
UPDATE messages
SET reaction_count = reaction_count - 1
WHERE id = $1
RETURNING *;

-- name: AnswerMessage :one
UPDATE messages
SET answered = true
WHERE id = $1
RETURNING *;

-- name: UnAnswerMessage :one
UPDATE messages
SET answered = false
WHERE id = $1
RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1 LIMIT 1;
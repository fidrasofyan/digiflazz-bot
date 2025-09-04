-- name: IsChatExists :one
SELECT EXISTS(SELECT 1 FROM chats WHERE id = ? LIMIT 1);

-- name: GetChat :one
SELECT * FROM chats WHERE id = ? LIMIT 1;

-- name: CreateChat :one
INSERT INTO chats (id, command, step, data) 
VALUES (?, ?, ?, ?) 
RETURNING *;

-- name: UpdateChat :one
UPDATE chats
SET 
  command = ?, 
  step = ?, 
  data = ?
WHERE id = ?
RETURNING *;

-- name: UpdateReplyMarkup1 :exec
UPDATE chats
SET 
  reply_markup_1 = ?
WHERE id = ?;

-- name: UpdateReplyMarkup2 :exec
UPDATE chats
SET 
  reply_markup_2 = ?
WHERE id = ?;

-- name: UpdateReplyMarkup3 :exec
UPDATE chats
SET 
  reply_markup_3 = ?
WHERE id = ?;

-- name: UpdateReplyMarkup4 :exec
UPDATE chats
SET 
  reply_markup_4 = ?
WHERE id = ?;

-- name: DeleteChat :exec
DELETE FROM chats WHERE id = ?;
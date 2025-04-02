-- name: AddNewChrip :one
INSERT INTO chrips (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING * ;

-- name: GetChrips :many
SELECT * FROM chrips
ORDER BY created_at ASC;

-- name: GetChripByID :one
SELECT * FROM chrips
WHERE id = $1;

-- name: DeleteChripsByID :exec
DELETE FROM chrips
WHERE id = $1;

-- name: GetChirpsByAuthorID :many
SELECT * FROM chrips
WHERE user_id = $1
ORDER BY created_at ASC
LIMIT 4;
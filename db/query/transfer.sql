-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1
LIMIT 1;

-- name: ListAllTransfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListAccountTransfers :many
SELECT * FROM transfers
WHERE from_account_id = sqlc.arg(account_id) 
OR to_account_id = sqlc.arg(account_id)
ORDER BY id
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;
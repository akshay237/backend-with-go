-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount
) values (
    $1, $2, $3
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
where id = $1
LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
where   
    from_account_id = $1 or
    to_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;
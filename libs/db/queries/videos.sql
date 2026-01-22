-- name: CreateVideo :one
INSERT INTO videos (
    id,
    user_id,
    title,
    status,
    duration_sec
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetVideoByID :one
SELECT *
FROM videos
WHERE id = $1;

-- name: ListVideosByUser :many
SELECT *
FROM videos
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListVideosByUserPaginated :many
SELECT *
FROM videos
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchVideo :many
SELECT *
FROM videos
WHERE
    user_id = COALESCE(sqlc.narg(user_id), user_id)
AND status  = COALESCE(sqlc.narg(status), status)
AND (
    sqlc.narg(title) IS NULL
    OR title ILIKE '%' || sqlc.narg(title) || '%'
)
ORDER BY created_at DESC
LIMIT COALESCE(sqlc.narg(size), 30)
OFFSET COALESCE(sqlc.narg(page), 0) * COALESCE(sqlc.narg(size), 30);


-- name: UpdateVideoStatus :one
UPDATE videos
SET status = $2
WHERE id = $1
RETURNING *;

-- name: UpdateVideoDuration :one
UPDATE videos
SET duration_sec = $2
WHERE id = $1
RETURNING *;

-- name: UpdateVideoTitle :one
UPDATE videos
SET title = $2
WHERE id = $1
RETURNING *;

-- name: DeleteVideo :exec
DELETE FROM videos
WHERE id = $1;

-- name: ListVideosByStatus :many
SELECT *
FROM videos
WHERE status = $1
ORDER BY created_at ASC;

-- name: ListStaleProcessingVideos :many
SELECT *
FROM videos
WHERE status = 'PROCESSING'
  AND created_at < now() - interval '30 minutes'
ORDER BY created_at ASC;

-- name: PatchVideos :exec
UPDATE videos
SET 
  title = COALESCE(sqlc.narg(title), title),
  status = COALESCE(sqlc.narg(status), status),
  duration_sec = COALESCE(sqlc.narg(duration_sec), duration_sec)
WHERE
  id = @id AND user_id = @user_id;
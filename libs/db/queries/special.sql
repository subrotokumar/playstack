-- name: ListVideosWithUsers :many
SELECT
    v.*,
    u.email,
    u.cognito_sub
FROM videos v
JOIN users u ON u.id = v.user_id
ORDER BY v.created_at DESC;

-- name: GetVideoWithUser :one
SELECT
    v.*,
    u.email,
    u.cognito_sub
FROM videos v
JOIN users u ON u.id = v.user_id
WHERE v.id = $1;

-- name: CountVideosByUser :one
SELECT COUNT(*) AS video_count
FROM videos
WHERE user_id = $1;

-- name: CountVideosByStatus :many
SELECT status, COUNT(*) AS count
FROM videos
GROUP BY status;

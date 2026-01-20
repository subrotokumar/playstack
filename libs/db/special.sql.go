// source: special.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countVideosByStatus = `-- name: CountVideosByStatus :many
SELECT status, COUNT(*) AS count
FROM videos
GROUP BY status
`

type CountVideosByStatusRow struct {
	Status VideoStatus `json:"status"`
	Count  int64       `json:"count"`
}

func (q *Queries) CountVideosByStatus(ctx context.Context) ([]CountVideosByStatusRow, error) {
	rows, err := q.db.Query(ctx, countVideosByStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []CountVideosByStatusRow{}
	for rows.Next() {
		var i CountVideosByStatusRow
		if err := rows.Scan(&i.Status, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const countVideosByUser = `-- name: CountVideosByUser :one
SELECT COUNT(*) AS video_count
FROM videos
WHERE user_id = $1
`

func (q *Queries) CountVideosByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countVideosByUser, userID)
	var video_count int64
	err := row.Scan(&video_count)
	return video_count, err
}

const getVideoWithUser = `-- name: GetVideoWithUser :one
SELECT
    v.id, v.user_id, v.title, v.status, v.original_s3_key, v.duration_sec, v.created_at,
    u.email
FROM videos v
JOIN users u ON u.id = v.user_id
WHERE v.id = $1
`

type GetVideoWithUserRow struct {
	ID            uuid.UUID        `json:"id"`
	UserID        uuid.UUID        `json:"user_id"`
	Title         string           `json:"title"`
	Status        VideoStatus      `json:"status"`
	OriginalS3Key string           `json:"original_s3_key"`
	DurationSec   pgtype.Int4      `json:"duration_sec"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	Email         pgtype.Text      `json:"email"`
}

func (q *Queries) GetVideoWithUser(ctx context.Context, id uuid.UUID) (GetVideoWithUserRow, error) {
	row := q.db.QueryRow(ctx, getVideoWithUser, id)
	var i GetVideoWithUserRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Status,
		&i.OriginalS3Key,
		&i.DurationSec,
		&i.CreatedAt,
		&i.Email,
	)
	return i, err
}

const listVideosWithUsers = `-- name: ListVideosWithUsers :many
SELECT
    v.id, v.user_id, v.title, v.status, v.original_s3_key, v.duration_sec, v.created_at,
    u.email
FROM videos v
JOIN users u ON u.id = v.user_id
ORDER BY v.created_at DESC
`

type ListVideosWithUsersRow struct {
	ID            uuid.UUID        `json:"id"`
	UserID        uuid.UUID        `json:"user_id"`
	Title         string           `json:"title"`
	Status        VideoStatus      `json:"status"`
	OriginalS3Key string           `json:"original_s3_key"`
	DurationSec   pgtype.Int4      `json:"duration_sec"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	Email         pgtype.Text      `json:"email"`
}

func (q *Queries) ListVideosWithUsers(ctx context.Context) ([]ListVideosWithUsersRow, error) {
	rows, err := q.db.Query(ctx, listVideosWithUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListVideosWithUsersRow{}
	for rows.Next() {
		var i ListVideosWithUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Status,
			&i.OriginalS3Key,
			&i.DurationSec,
			&i.CreatedAt,
			&i.Email,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

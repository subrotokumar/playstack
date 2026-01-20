package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CountVideosByStatus(ctx context.Context) ([]CountVideosByStatusRow, error)
	CountVideosByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateVideo(ctx context.Context, arg CreateVideoParams) (Video, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	DeleteVideo(ctx context.Context, id uuid.UUID) error
	GetUserByEmail(ctx context.Context, email pgtype.Text) (User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetVideoByID(ctx context.Context, id uuid.UUID) (Video, error)
	GetVideoWithUser(ctx context.Context, id uuid.UUID) (GetVideoWithUserRow, error)
	ListStaleProcessingVideos(ctx context.Context) ([]Video, error)
	ListVideosByStatus(ctx context.Context, status VideoStatus) ([]Video, error)
	ListVideosByUser(ctx context.Context, userID uuid.UUID) ([]Video, error)
	ListVideosByUserPaginated(ctx context.Context, arg ListVideosByUserPaginatedParams) ([]Video, error)
	ListVideosWithUsers(ctx context.Context) ([]ListVideosWithUsersRow, error)
	UpdateVideoDuration(ctx context.Context, arg UpdateVideoDurationParams) (Video, error)
	UpdateVideoStatus(ctx context.Context, arg UpdateVideoStatusParams) (Video, error)
	UpdateVideoTitle(ctx context.Context, arg UpdateVideoTitleParams) (Video, error)
}

var _ Querier = (*Queries)(nil)

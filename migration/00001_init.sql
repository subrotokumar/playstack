-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TYPE video_status AS ENUM (
    'PREUPLOAD',
    'UPLOADED',
    'PROCESSING',
    'READY',
    'FAILED'
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    status video_status NOT NULL,
    duration_sec INT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_videos_user_id ON videos(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS idx_videos_user_id;

DROP TABLE IF EXISTS videos;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS video_status;
-- +goose StatementEnd

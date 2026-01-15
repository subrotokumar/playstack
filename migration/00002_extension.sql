-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TYPE  IF NOT EXISTS job_status AS ENUM (
    'PENDING',
    'RUNNING',
    'SUCCESS',
    'FAILED'
);

CREATE TYPE  IF NOT EXISTS video_resolution AS ENUM (
    '240p',
    '360p',
    '480p',
    '720p',
    '1080p'
);

CREATE TABLE IF NOT EXISTS video_renditions (
    id UUID PRIMARY KEY,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    resolution video_resolution NOT NULL,
    bitrate_kbps INT NOT NULL,
    s3_key TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_renditions_video_id ON video_renditions(video_id);

CREATE TABLE IF NOT EXISTS manifests (
    id UUID PRIMARY KEY, 
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    s3_key TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (video_id)
);

CREATE TABLE IF NOT EXISTS transcoding_jobs (
    id UUID PRIMARY KEY,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    status job_status NOT NULL,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX  IF NOT EXISTS ON idx_jobs_video_id ON transcoding_jobs(video_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS idx_jobs_video_id;
DROP INDEX IF EXISTS idx_renditions_video_id;

DROP TABLE IF EXISTS transcoding_jobs;
DROP TABLE IF EXISTS dash_manifests;
DROP TABLE IF EXISTS video_renditions;

DROP TYPE IF EXISTS job_status;
-- +goose StatementEnd

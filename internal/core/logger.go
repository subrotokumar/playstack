package core

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

func NewLogger(env Env, service string, level *string) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: env != EnvDevelopment,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				a.Key = "time"
				a.Value = slog.StringValue(
					time.Now().Local().Format("2006/01/02 15:04:05"),
				)

			case slog.LevelKey:
				a.Key = "level"

			case slog.SourceKey:
				a.Key = "caller"
			}
			return a
		},
	}

	var handler slog.Handler

	switch env {
	case EnvDevelopment:
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)

	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	svc := fmt.Sprintf("%s:%s", service, env)
	return slog.New(handler).With(
		slog.String("app", svc),
	)
}

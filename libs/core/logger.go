package core

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(env Env, service string, level *string) *Logger {
	slogLevel := slog.LevelInfo
	if level != nil {
		slogLevel = parseLevel(*level)
	}
	opts := &slog.HandlerOptions{
		Level:     slogLevel,
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
	slog := slog.New(handler).With(
		slog.String("app", svc),
	)
	return &Logger{slog}
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.Logger.Error(msg, args...)
	os.Exit(1)
}

package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
)

const (
	LevelTrace slog.Level = slog.Level(-8)
)

var _ Logger = (*Slog)(nil)

type (
	Slog struct {
		slogger *slog.Logger
	}

	SlogOpts struct {
		// Component enriches each log line with a componenent key/value.
		// Useful for aggregating/filtering with your log collector.
		Component string
		// Handler allows overriding of the defult Logfmt handler.
		Handler slog.Handler
		// Minimal level to log. Defaults to Info.
		// No effect when passing a custom handler.
		LogLevel slog.Level
		// Add source location to each log line. Defaults to false.
		// No effect when passing a custom handler.
		IncludeSource bool
	}
)

func NewSlog(o SlogOpts) *Slog {
	if o.Handler == nil {
		o.Handler = buildDefaultHandler(os.Stderr, o.LogLevel, o.IncludeSource)
	}
	if o.Component == "" {
		o.Component = "vise"
	}

	return &Slog{
		slogger: slog.New(o.Handler).With("component", o.Component),
	}
}

func buildDefaultHandler(w io.Writer, level slog.Level, includeSource bool) slog.Handler {
	return slog.NewTextHandler(w, &slog.HandlerOptions{
		AddSource: includeSource,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				switch a.Value.Any().(slog.Level) {
				case LevelTrace:
					return slog.String(slog.LevelKey, "TRACE")
				}
			}
			return a
		},
	})
}

// WithDomain sets the domain for the logger.
// Returns a new instance of Slog with the domain set.
func (s *Slog) WithDomain(domain string) *Slog {
	return &Slog{slogger: s.slogger.With("domain", domain)}
}

func (s *Slog) Writef(w io.Writer, level int, msg string, args ...any) {
	s.slogger.Warn("Writef not implemented")
}

func (s *Slog) WriteCtxf(ctx context.Context, w io.Writer, level int, msg string, args ...any) {
	s.slogger.Warn("WriteCtxf not implemented")
}

func (s *Slog) Printf(level int, msg string, args ...any) {
	s.slogger.Warn("Printf not implemented")
}

func (s *Slog) PrintCtxf(ctx context.Context, level int, msg string, args ...any) {
	s.slogger.Warn("PrintCtxf not implemented")
}

func (s *Slog) Tracef(msg string, args ...any) {
	s.slogger.Log(nil, LevelTrace, msg, args...)
}

func (s *Slog) TraceCtxf(ctx context.Context, msg string, args ...any) {
	s.slogger.Log(ctx, LevelTrace, msg, args...)
}

func (s *Slog) Debugf(msg string, args ...any) {
	s.slogger.Debug(msg, args...)
}

func (s *Slog) DebugCtxf(ctx context.Context, msg string, args ...any) {
	s.slogger.DebugContext(ctx, msg, args...)
}

func (s *Slog) Infof(msg string, args ...any) {
	s.slogger.Info(msg, args...)
}

func (s *Slog) InfoCtxf(ctx context.Context, msg string, args ...any) {
	s.slogger.InfoContext(ctx, msg, args...)
}

func (s *Slog) Warnf(msg string, args ...any) {
	s.slogger.Warn(msg, args...)
}

func (s *Slog) WarnCtxf(ctx context.Context, msg string, args ...any) {
	s.slogger.WarnContext(ctx, msg, args...)
}

func (s *Slog) Errorf(msg string, args ...any) {
	s.slogger.Error(msg, args...)
}

func (s *Slog) ErrorCtxf(ctx context.Context, msg string, args ...any) {
	s.slogger.ErrorContext(ctx, msg, args...)
}

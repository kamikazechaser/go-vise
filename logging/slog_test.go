package logging

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func newTestHandler(w io.Writer, o SlogOpts) slog.Handler {
	return slog.NewTextHandler(w, &slog.HandlerOptions{
		AddSource: o.IncludeSource,
		Level:     o.LogLevel,
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

func TestNewSlogOutput(t *testing.T) {
	var buf bytes.Buffer

	logger := NewSlog(SlogOpts{
		Handler: newTestHandler(&buf, SlogOpts{
			LogLevel:      LevelTrace,
			IncludeSource: false,
		}),
	})

	logger.Tracef("trace test")
	logger.Debugf("debug test")
	logger.Infof("info test")
	logger.Warnf("warn test")
	logger.Errorf("error test")

	logOutput := buf.String()
	t.Logf("Log output: %s", logOutput)
	if !strings.Contains(logOutput, "TRACE") || !strings.Contains(logOutput, "trace test") {
		t.Errorf("expected TRACE message in log output: %s", logOutput)
	}
}

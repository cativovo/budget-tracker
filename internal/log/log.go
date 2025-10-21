package log

import "log/slog"

// Error returns an slog.Attr for error with `error` key
func Error(err error) slog.Attr {
	return slog.Any("error", err)
}

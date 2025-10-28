package log

import "log/slog"

// ErrAttr returns an slog.Attr for error with `error` key
func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

package testutil

import "log/slog"

// DisableLogs disables the logs
func DisableLogs() {
	slog.SetDefault(slog.New(slog.DiscardHandler))
}

package conman

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"context"
	"fmt"
	"log/slog"
)

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, keyLogger, logger)
}

func GetLogger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(keyLogger).(*slog.Logger)
	if !ok {
		return slog.With("internal_error", "context error", "context_key", keyLogger, "context_value", ctx.Value(keyLogger), "context_type", fmt.Sprintf("%T", ctx.Value(keyLogger)))
	}
	return logger
}

package main

import (
	"context"
	"io"
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func setupLogger() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.New(slog.NewTextHandler(os.Stdout, nil)).Error("could not open log file", "error", err)
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
		return
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	handler := slog.NewTextHandler(multiWriter, nil)
	logger = slog.New(handler)
}

const traceIdKey string = "traceID"

func Log(ctx context.Context) *slog.Logger {
	id, _ := ctx.Value(traceIdKey).(string)
	return logger.With("traceID", id)
}

func WithNewTraceId() context.Context {
	return context.WithValue(context.Background(), traceIdKey, GenerateTraceID())
}

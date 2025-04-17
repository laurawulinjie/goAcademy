package main

import (
	"context"
	"io"
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func setupLogger() {
	logFile, err := os.OpenFile("./data/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
	if id, ok := ctx.Value(traceIdKey).(string); ok {
		return logger.With("traceID", id)
	}
	return logger.With("traceID", "no-trace-id")

}

func WithNewTraceId() context.Context {
	return context.WithValue(context.Background(), traceIdKey, GenerateTraceID())
}

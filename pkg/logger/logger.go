package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/laurawulinjie/goAcademy/pkg/utils"
)

const TraceIdKey string = "traceID"

type traceIdHandler struct {
	slog.Handler
}

func (h *traceIdHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID, ok := ctx.Value(TraceIdKey).(string); ok {
		r.AddAttrs(slog.String(TraceIdKey, traceID))
	}
	return h.Handler.Handle(ctx, r)
}

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func SetupLogger() {
	logFile, err := os.OpenFile("./data/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.New(slog.NewTextHandler(os.Stdout, nil)).Error("could not open log file", "error", err)
		return
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	baseHandler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{AddSource: true})
	handler := &traceIdHandler{Handler: baseHandler}
	logger = slog.New(handler)
	slog.SetDefault(logger)
}

func WithNewTraceId() context.Context {
	return context.WithValue(context.Background(), TraceIdKey, utils.GenerateTraceID())
}

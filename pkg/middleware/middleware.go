package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/laurawulinjie/goAcademy/pkg/logger"
	"github.com/laurawulinjie/goAcademy/pkg/utils"
)

func TraceIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId := utils.GenerateTraceID()
		ctx := context.WithValue(r.Context(), logger.TraceIdKey, traceId)
		slog.InfoContext(ctx, "Received request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

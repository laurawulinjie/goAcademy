package main

import (
	"context"
	"net/http"
)

func TraceIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId := GenerateTraceID()
		ctx := context.WithValue(r.Context(), traceIdKey, traceId)
		Log(ctx).Info("Received request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

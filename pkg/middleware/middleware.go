package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

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

func UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := 0
		cookie, err := r.Cookie("user_id")

		if err != nil {
			slog.WarnContext(ctx, "No user_id cookie found", "error", err)
		} else {
			if id, err := strconv.Atoi(cookie.Value); err == nil {
				userID = id
				ctx = context.WithValue(ctx, "user_id", userID)
			} else {
				slog.WarnContext(ctx, "Invalid user_id cookie", "error", err)
			}
		}

		publicPaths := map[string]bool{
			"/register": true,
			"/login":    true,
			"/about":    true,
		}

		if !publicPaths[r.URL.Path] && userID == 0 {
			slog.ErrorContext(ctx, "Unauthorized access attempt", "path", r.URL.Path)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

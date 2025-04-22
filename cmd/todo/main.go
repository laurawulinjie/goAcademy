package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	setupLogger()
	ctx := WithNewTraceId()
	ctx, ctxDone := context.WithCancel(ctx)

	if err := LoadTodos(ctx); err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	if err := setupDynamicPages(); err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", HomePageHandler)
	mux.HandleFunc("/list", ListPageHandler)
	mux.HandleFunc("/about", AboutPageHandler)
	mux.HandleFunc("/create", CreateHandler)
	mux.HandleFunc("/get", GetHandler)
	mux.HandleFunc("/update", UpdateHandler)
	mux.HandleFunc("/delete", DeleteHandler)

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: TraceIdMiddleware(mux),
	}

	slog.InfoContext(ctx, "Server is running on http://localhost:8080/")

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic("ListenAndServe: " + err.Error())
		}
	}()

	go func() {
		defer ctxDone()
		close := make(chan os.Signal, 1)
		signal.Notify(close, os.Interrupt)
		sig := <-close
		slog.InfoContext(ctx, "got signal: ["+sig.String()+"] now closing")
	}()

	<-ctx.Done()
	slog.InfoContext(ctx, "shutdown application")
}

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/laurawulinjie/goAcademy/pkg/handlers"
	"github.com/laurawulinjie/goAcademy/pkg/logger"
	"github.com/laurawulinjie/goAcademy/pkg/middleware"
	"github.com/laurawulinjie/goAcademy/pkg/todo"
)

func main() {
	logger.SetupLogger()
	ctx := logger.WithNewTraceId()
	ctx, ctxDone := context.WithCancel(ctx)

	if err := todo.LoadTodos(ctx); err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	if err := SetupDynamicPages(); err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", ServeHomePage)
	mux.HandleFunc("/list", ServeListPage)
	mux.HandleFunc("/about", ServeAboutPage)
	mux.HandleFunc("/create", handlers.CreateHandler)
	mux.HandleFunc("/get", handlers.GetHandler)
	mux.HandleFunc("/update", handlers.UpdateHandler)
	mux.HandleFunc("/delete", handlers.DeleteHandler)

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: middleware.TraceIdMiddleware(mux),
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

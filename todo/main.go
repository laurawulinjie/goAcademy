package main

import (
	"net/http"
	"os"
	"os/signal"
)

func main() {
	setupLogger()
	ctx := WithNewTraceId()

	if err := LoadTodos(ctx); err != nil {
		Log(ctx).Error(err.Error())
	}

	if err := setupDynamicPages(); err != nil {
		Log(ctx).Error(err.Error())
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

	close := make(chan os.Signal, 1)
	signal.Notify(close, os.Interrupt)

	go func() {
		Log(ctx).Info("Server is running on http://localhost:8080/")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Log(ctx).Error("Server error", "error", err)
		}
	}()

	<-close
	Log(ctx).Info("shutdown application")
}

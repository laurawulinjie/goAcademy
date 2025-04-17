package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	setupLogger()
	ctx := WithNewTraceId()
	ctx, ctxDone := context.WithCancel(ctx)

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

	Log(ctx).Info("Server is running on http://localhost:8080/")

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
		Log(ctx).Info("got signal: [" + sig.String() + "] now closing")
	}()

	<-ctx.Done()
	Log(ctx).Info("shutdown application")
}

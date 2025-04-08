package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {
	slog.Info("hello, world")

	slog.Info("hello, world", "user", os.Getenv("USER"))

	logger1 := slog.Default()
	logger1.Info("hello, world", "user", os.Getenv("USER"))

	logger2 := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger2.Info("hello, world", "user", os.Getenv("USER"))

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("hello, world", "user", os.Getenv("USER"))

	slog.LogAttrs(context.Background(), slog.LevelInfo, "hello, world",
		slog.String("user", os.Getenv("USER")))

	value1 := 1
	value2 := "dummy value 2"
	slog.Info("message", slog.Int("key1", value1), slog.String("key2", value2))

}

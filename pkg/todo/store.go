package todo

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
)

const dataFile = "./data/todos.json"

func LoadTodos(ctx context.Context) error {
	file, err := os.Open(dataFile)

	if err != nil {
		if os.IsNotExist(err) {
			slog.ErrorContext(ctx, "todos file not found, starting with emtpy list", "file", dataFile)
			return nil
		}
		slog.ErrorContext(ctx, "failed to open file")
		return err
	}

	defer file.Close()
	var loaded map[int]Todo
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&loaded); err != nil {
		slog.ErrorContext(ctx, "failed to decode todos")
		return err
	}

	todos = loaded
	maxId := 0

	for id := range todos {
		if id > maxId {
			maxId = id
		}
	}

	nextId = maxId + 1
	slog.InfoContext(ctx, "todos loaded", "todoLength", len(todos))
	return nil
}

func SaveTodos(ctx context.Context) error {
	file, err := os.Create(dataFile)

	if err != nil {
		slog.ErrorContext(ctx, "failed to create file")
		return err
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	")

	if err := encoder.Encode(todos); err != nil {
		slog.ErrorContext(ctx, "failed to encode todos")
		return err
	}

	slog.InfoContext(ctx, "todos saved", "todoLength", len(todos))
	return nil
}

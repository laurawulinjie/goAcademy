package main

import (
	"context"
	"encoding/json"
	"os"
)

const dataFile = "todo.json"

func LoadTodos(ctx context.Context) error {
	file, err := os.Open(dataFile)

	if err != nil {
		if os.IsNotExist(err) {
			Log(ctx).Info("todos file not found, starting with emtpy list", "file", dataFile)
			return nil
		}
		Log(ctx).Error("failed to open file")
		return err
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&todos); err != nil {
		Log(ctx).Error("failed to decode todos")
		return err
	}

	if len(todos) > 0 {
		nextId = todos[len(todos)-1].ID + 1
	}

	Log(ctx).Info("todos loaded", "todoLength", len(todos))
	return nil
}

func SaveTodos(ctx context.Context) error {
	file, err := os.Create(dataFile)

	if err != nil {
		Log(ctx).Error("failed to create file")
		return err
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	")

	if err := encoder.Encode(todos); err != nil {
		Log(ctx).Error("failed to encode todos")
		return err
	}

	Log(ctx).Info("todos saved", "todoLength", len(todos))
	return nil
}

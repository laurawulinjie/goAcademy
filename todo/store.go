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
	var loaded map[int]Todo
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&loaded); err != nil {
		Log(ctx).Error("failed to decode todos")
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

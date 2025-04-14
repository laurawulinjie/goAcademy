package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
)

var todos []Todo
var nextId = 1

func GetAllTodos() []Todo {
	return todos
}

func CreateTodo(ctx context.Context, task string) (Todo, error) {
	if task == "" {
		return Todo{}, errors.New("task is required")
	}

	newTodo := Todo{
		ID:     nextId,
		Task:   task,
		Status: NotStarted,
	}

	todos = append(todos, newTodo)
	nextId++
	return newTodo, nil
}

func UpdateTodo(ctx context.Context, id int, task string, status string) error {
	for i, todo := range todos {
		if todo.ID == id {
			if task != "" {
				todos[i].Task = task
			}

			if status != "" {
				todos[i].Status = status
			}

			Log(ctx).Info("Updated todo", "id", id, "task", todos[i].Task, "status", todos[i].Status)
			return nil
		}
	}

	Log(ctx).Error("Todo not found", "id", id)
	return fmt.Errorf("todo not found")
}

func DeleteTodo(ctx context.Context, id int) error {
	for i, todo := range todos {
		if todo.ID == id {
			todos = slices.Delete(todos, i, i+1)
			Log(ctx).Info("Deleted todo", "id", id)
			return nil
		}
	}

	Log(ctx).Error("Todo not found", "id", id)
	return errors.New("todo not found")
}

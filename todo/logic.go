package main

import (
	"context"
	"errors"
	"log/slog"
)

var todos = make(map[int]Todo)

var nextId = 1

func GetAllTodos(ctx context.Context) map[int]Todo {
	if len(todos) == 0 {
		slog.InfoContext(ctx, "No todos found")
	}

	return todos
}

func CreateTodo(ctx context.Context, task string) (Todo, error) {
	if task == "" {
		slog.ErrorContext(ctx, "Task cannot be empty")
		return Todo{}, errors.New("task cannot be empty")
	}

	newTodo := Todo{
		Task:   task,
		Status: NotStarted,
	}

	todos[nextId] = newTodo
	slog.InfoContext(ctx, "Created new todo", "id", nextId, "task", newTodo.Task, "status", newTodo.Status)
	nextId++
	return newTodo, nil
}

func UpdateTodo(ctx context.Context, id int, task string, status string) error {
	todo, exists := todos[id]
	if !exists {
		slog.ErrorContext(ctx, "Todo not found", "id", id)
		return errors.New("todo not found")
	}

	if task != "" {
		todo.Task = task
	}

	if status != "" {
		todo.Status = status
	}

	todos[id] = todo

	slog.InfoContext(ctx, "Updated todo", "id", id, "task", todos[id].Task, "status", todos[id].Status)
	return nil
}

func DeleteTodo(ctx context.Context, id int) error {
	_, exists := todos[id]
	if !exists {
		slog.ErrorContext(ctx, "Todo not found", "id", id)
		return errors.New("todo not found")
	}

	delete(todos, id)
	slog.InfoContext(ctx, "Deleted todo", "id", id)
	return nil
}

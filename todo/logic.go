package main

import (
	"context"
	"errors"
)

var todos = make(map[int]Todo)

var nextId = 1

func GetAllTodos(ctx context.Context) map[int]Todo {
	if len(todos) == 0 {
		Log(ctx).Info("No todos found")
	}

	return todos
}

func CreateTodo(ctx context.Context, task string) (Todo, error) {
	if task == "" {
		Log(ctx).Error("Task cannot be empty")
		return Todo{}, errors.New("task cannot be empty")
	}

	newTodo := Todo{
		ID:     nextId,
		Task:   task,
		Status: NotStarted,
	}

	todos[nextId] = newTodo
	nextId++
	Log(ctx).Info("Created new todo", "id", newTodo.ID, "task", newTodo.Task, "status", newTodo.Status)
	return newTodo, nil
}

func UpdateTodo(ctx context.Context, id int, task string, status string) error {
	todo, exists := todos[id]
	if !exists {
		Log(ctx).Error("Todo not found", "id", id)
		return errors.New("todo not found")
	}

	if task != "" {
		todo.Task = task
	}

	if status != "" {
		todo.Status = status
	}

	todos[id] = todo

	Log(ctx).Info("Updated todo", "id", id, "task", todos[id].Task, "status", todos[id].Status)
	return nil
}

func DeleteTodo(ctx context.Context, id int) error {
	_, exists := todos[id]
	if !exists {
		Log(ctx).Error("Todo not found", "id", id)
		return errors.New("todo not found")
	}

	delete(todos, id)
	Log(ctx).Info("Deleted todo", "id", id)
	return nil
}

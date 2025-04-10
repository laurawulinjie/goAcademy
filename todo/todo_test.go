package main

import (
	"context"
	"testing"
)

func resetTodos() context.Context {
	todos = []Todo{}
	nextID = 1
	setupLogger()
	return withTraceId(context.Background())
}

func TestAddTodo(t *testing.T) {
	t.Run("add todo with description", func(t *testing.T) {
		ctx := resetTodos()
		err := addTodo(ctx, "Dummy Task")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(todos) != 1 {
			t.Errorf("Expected 1 todo, got %d", len(todos))
		}

		if todos[0].Task != "Dummy Task" {
			t.Errorf("Expected task 'Dummy Task', got '%s'", todos[0].Task)
		}

		if todos[0].Status != NotStarted {
			t.Errorf("Expected status 'not started', got '%s'", todos[0].Status)
		}
	})

	t.Run("add todo with no description", func(t *testing.T) {
		ctx := resetTodos()
		err := addTodo(ctx, "")

		if err == nil {
			t.Errorf("Expected error for empty description, got nil")
		}
	})
}

func TestUpdateTodo(t *testing.T) {
	t.Run("update existing task", func(t *testing.T) {
		ctx := resetTodos()
		addTodo(ctx, "Initial Task")
		id := todos[0].ID
		err := updateTodo(ctx, id, "Updated task", Started)

		if err != nil {
			t.Errorf("Expected no err, got '%v'", err)
		}

		if todos[0].Task != "Updated task" {
			t.Errorf("Expected task 'Updated task', got '%s'", todos[0].Task)
		}

		if todos[0].Status != Started {
			t.Errorf("Expected status 'started', got '%s'", todos[0].Status)
		}
	})

	t.Run("update nonexistent task", func(t *testing.T) {
		ctx := resetTodos()
		err := updateTodo(ctx, 1, "Nonexistent", Started)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestDeleteTodo(t *testing.T) {
	t.Run("delete existing todo", func(t *testing.T) {
		ctx := resetTodos()
		addTodo(ctx, "Todo to delete")
		id := todos[0].ID
		err := deleteTodo(ctx, id)

		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		if len(todos) != 0 {
			t.Errorf("Expected 0 todos after deletion, got %d", len(todos))
		}
	})

	t.Run("delete nonexistent todo", func(t *testing.T) {
		ctx := resetTodos()
		err := deleteTodo(ctx, 1)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

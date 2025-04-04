package main

import "testing"

func resetTodos() {
	todos = []Todo{}
	nextID = 1
}

func TestAddTodo(t *testing.T) {
	resetTodos()
	addTodo("Dummy Task")

	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}

	if todos[0].Task != "Dummy Task" {
		t.Errorf("Expected task 'Dummy Task', got '%s'", todos[0].Task)
	}

	if todos[0].Status != NotStarted {
		t.Errorf("Expected status 'not started', got '%s'", todos[0].Status)
	}
}

func TestUpdateTodo(t *testing.T) {
	t.Run("update existing task", func(t *testing.T) {
		resetTodos()
		addTodo("Initial Task")
		id := todos[0].ID
		updateTodo(id, "Updated task", Started)

		if todos[0].Task != "Updated task" {
			t.Errorf("Expected task 'Updated task', got '%s'", todos[0].Task)
		}

		if todos[0].Status != Started {
			t.Errorf("Expected status 'started', got '%s'", todos[0].Status)
		}
	})

	t.Run("update nonexistent task", func(t *testing.T) {
		resetTodos()
		updateTodo(1, "Nonexistent", Started)

	})
}

func TestDeleteTodo(t *testing.T) {
	t.Run("delete existing todo", func(t *testing.T) {
		resetTodos()
		addTodo("Todo to delete")
		id := todos[0].ID
		deleteTodo(id)

		if len(todos) != 0 {
			t.Errorf("Expected 0 todos after deletion, got %d", len(todos))
		}
	})

	t.Run("delete nonexistent todo", func(t *testing.T) {
		resetTodos()
		deleteTodo(1)
	})
}

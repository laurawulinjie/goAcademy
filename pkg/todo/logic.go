package todo

import (
	"context"
	"errors"
	"log/slog"
)

func GetAllTodos(ctx context.Context, userID int) (map[int]Todo, error) {
	rows, err := DB.QueryContext(ctx, "SELECT id, task, status FROM todos WHERE user_id = $1", userID)

	if err != nil {
		slog.ErrorContext(ctx, "failed to query todos", "error", err)
		return nil, errors.New("failed to query todos")
	}

	defer rows.Close()
	todos := make(map[int]Todo)

	for rows.Next() {
		var todo Todo
		var id int
		if err := rows.Scan(&id, &todo.Task, &todo.Status); err != nil {
			slog.ErrorContext(ctx, "failed to scan todo", "id", id, "error", err)
			continue
		}
		todos[id] = todo
	}

	return todos, nil
}

func CreateTodo(ctx context.Context, task string, userID int) (Todo, error) {
	if task == "" {
		slog.ErrorContext(ctx, "Task cannot be empty")
		return Todo{}, errors.New("task cannot be empty")
	}

	newTodo := Todo{
		Task:   task,
		Status: NotStarted,
		UserID: userID,
	}

	err := DB.QueryRowContext(ctx,
		"INSERT INTO todos (task, status, user_id) VALUES ($1, $2, $3) RETURNING id",
		newTodo.Task,
		newTodo.Status,
		newTodo.UserID).Scan(&newTodo.ID)

	if err != nil {
		slog.ErrorContext(ctx, "failed to insert todo", "error", err)
		return Todo{}, err
	}

	slog.InfoContext(ctx, "Created new todo", "id", newTodo.ID, "task", newTodo.Task, "status", newTodo.Status)
	return newTodo, nil
}

func UpdateTodo(ctx context.Context, id int, task string, status string, userID int) error {
	res, err := DB.ExecContext(ctx,
		"UPDATE todos SET task = $1, status = $2 WHERE id = $3 AND user_id = $4",
		task, status, id, userID)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to update todo", "error", err)
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		slog.ErrorContext(ctx, "Todo not found for updated", "id", id)
		return errors.New("todo not found")
	}

	slog.InfoContext(ctx, "Updated todo", "id", id, "task", task, "status", status)
	return nil
}

func DeleteTodo(ctx context.Context, id int, userID int) error {
	res, err := DB.ExecContext(ctx,
		"DELETE FROM todos WHERE id = $1 AND user_id = $2", id, userID)

	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete todo", "error", err)
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		slog.ErrorContext(ctx, "Todo not found for deletion", "id", id)
		return errors.New("todo not found")
	}
	slog.InfoContext(ctx, "Deleted todo", "id", id)
	return nil
}

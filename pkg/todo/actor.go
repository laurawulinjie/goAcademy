package todo

import (
	"context"
	"log/slog"
)

type Request struct {
	Ctx      context.Context
	Action   string
	Payload  any
	Response chan any
}

var RequestQueue = make(chan Request, 100)

func StartTodoActor(ctx context.Context) {
	slog.InfoContext(ctx, "Starting Todo Actor")

	go func() {
		for {
			select {
			case req := <-RequestQueue:
				switch req.Action {
				case "create":
					payload := req.Payload.(struct {
						Task   string
						UserID int
					})

					todo, err := CreateTodo(req.Ctx, payload.Task, payload.UserID)

					req.Response <- struct {
						Todo Todo
						Err  error
					}{todo, err}

				case "update":
					payload := req.Payload.(struct {
						ID     int
						Task   string
						Status string
						UserID int
					})

					err := UpdateTodo(req.Ctx, payload.ID, payload.Task, payload.Status, payload.UserID)

					req.Response <- struct {
						Err error
					}{Err: err}

				case "delete":
					payload := req.Payload.(struct {
						ID     int
						UserID int
					})

					err := DeleteTodo(req.Ctx, payload.ID, payload.UserID)

					req.Response <- struct {
						Err error
					}{Err: err}

				case "getAll":
					userID := req.Payload.(int)
					todos, err := GetAllTodos(req.Ctx, userID)
					req.Response <- struct {
						Todos map[int]Todo
						Err   error
					}{Todos: todos, Err: err}

				case "register":
					payload := req.Payload.(struct {
						Username string
						Password string
					})

					_, err := DB.ExecContext(req.Ctx,
						"INSERT INTO users (username, password) VALUES ($1, $2)",
						payload.Username, payload.Password)

					req.Response <- struct {
						Err error
					}{Err: err}
				}
			case <-ctx.Done():
				slog.InfoContext(ctx, "shutting down todo actor")
				return
			}
		}
	}()
}

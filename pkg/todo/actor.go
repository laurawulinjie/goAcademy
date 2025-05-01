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
					todo, err := CreateTodo(req.Ctx, req.Payload.(string))

					req.Response <- struct {
						Todo Todo
						Err  error
					}{todo, err}

				case "update":
					payload := req.Payload.(struct {
						ID     int
						Task   string
						Status string
					})

					err := UpdateTodo(req.Ctx, payload.ID, payload.Task, payload.Status)

					req.Response <- struct {
						Err error
					}{Err: err}

				case "delete":
					err := DeleteTodo(req.Ctx, req.Payload.(int))

					req.Response <- struct {
						Err error
					}{Err: err}

				case "getAll":
					todos, err := GetAllTodos(req.Ctx)
					req.Response <- struct {
						Todos map[int]Todo
						Err   error
					}{Todos: todos, Err: err}
				}
			case <-ctx.Done():
				slog.InfoContext(ctx, "shutting down todo actor")
				return
			}
		}
	}()
}

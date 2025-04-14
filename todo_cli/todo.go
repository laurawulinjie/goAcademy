package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"text/tabwriter"

	"github.com/google/uuid"
)

type Status string

const (
	NotStarted Status = "not started"
	Started    Status = "started"
	Completed  Status = "completed"
)

type Todo struct {
	ID     int    `json:"id"`
	Task   string `json:"task"`
	Status Status `json:"status"`
}

var todos []Todo
var nextID = 1

const dataFile = "todos.json"

var logger *slog.Logger

func setupLogger() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

// not really type safe, get rid of it
type contextKey string

const traceIdKey contextKey = "traceID"

func generateTraceId() string {
	return uuid.New().String()
}

func getTraceId(ctx context.Context) string {
	if traceId, ok := ctx.Value(traceIdKey).(string); ok {
		return traceId
	}
	return "no-trace-id"
}

func logWithTraceId(ctx context.Context) *slog.Logger {
	return logger.With("traceId", getTraceId(ctx))
}

func loadTodos(ctx context.Context) error {
	file, err := os.Open(dataFile)

	if err != nil {
		if os.IsNotExist(err) {
			logWithTraceId(ctx).Info("todos file not found, starting with emtpy list", "file", dataFile)
			return nil
		}
		logWithTraceId(ctx).Error("failed to open file")
		return err
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&todos); err != nil {
		logWithTraceId(ctx).Error("failed to decode todos")
		return err
	}

	if len(todos) > 0 {
		nextID = todos[len(todos)-1].ID + 1
	}

	logWithTraceId(ctx).Info("todos loaded", "todoLength", len(todos))
	return nil
}

func saveTodos(ctx context.Context) error {
	file, err := os.Create(dataFile)

	if err != nil {
		logWithTraceId(ctx).Error("failed to create file")
		return err
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	")

	if err := encoder.Encode(todos); err != nil {
		logWithTraceId(ctx).Error("failed to encode todos")
		return err
	}

	logWithTraceId(ctx).Info("todos saved", "todoLength", len(todos))
	return nil
}

func listTodos(ctx context.Context) {
	if len(todos) == 0 {
		logWithTraceId(ctx).Info("No todos found")
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "ID\tTask\tStatus")

	for _, todo := range todos {
		fmt.Fprintf(writer, "%d\t%s\t%s\n", todo.ID, todo.Task, todo.Status)
	}

	writer.Flush()
}

func addTodo(ctx context.Context, description string) error {
	if description == "" {
		return errors.New("description cannot be empty")
	}

	newTodo := Todo{
		ID:     nextID,
		Task:   description,
		Status: NotStarted,
	}

	todos = append(todos, newTodo)
	nextID++
	logWithTraceId(ctx).Info("Created new todo", "id", newTodo.ID, "task", newTodo.Task, "status", newTodo.Status)
	return nil
}

func updateTodo(ctx context.Context, id int, description string, status Status) error {
	for i, todo := range todos {
		if todo.ID == id {
			if description != "" {
				todos[i].Task = description
			}

			if status != "" {
				todos[i].Status = status
			}

			logWithTraceId(ctx).Info("Updated todo", "id", id, "task", todos[i].Task, "status", todos[i].Status)
			return nil
		}
	}

	logWithTraceId(ctx).Error("Todo not found", "id", id)
	return fmt.Errorf("todo not found")
}

func deleteTodo(ctx context.Context, id int) error {
	for i, todo := range todos {
		if todo.ID == id {
			todos = slices.Delete(todos, i, i+1)
			logWithTraceId(ctx).Info("Deleted todo", "id", id)
			return nil
		}
	}

	logWithTraceId(ctx).Error("Todo not found", "id", id)
	return errors.New("todo not found")
}

func validateStatus(status string) (Status, error) {
	switch status {
	case "not started":
		return NotStarted, nil
	case "started":
		return Started, nil
	case "completed":
		return Completed, nil
	case "":
		return "", nil
	default:
		return "", fmt.Errorf("invalid status. Status must be one of: not started, started, completed")
	}
}

func main() {
	setupLogger()
	ctx, ctxDone := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, traceIdKey, generateTraceId())
	logWithTraceId(ctx).InfoContext(ctx, "start application")

	go func() {
		defer ctxDone()
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		s := <-c
		logWithTraceId(ctx).InfoContext(ctx, "got signal: ["+s.String()+"] now closing")
	}()

	listFlag := flag.Bool("list", false, "List all todos")
	addFlag := flag.Bool("add", false, "Add a new todo")
	updateFlag := flag.Bool("update", false, "Update a todo")
	deleteFlag := flag.Bool("delete", false, "Delete a todo")
	description := flag.String("description", "", "Task description")
	id := flag.Int("id", 0, "Todo ID")
	status := flag.String("status", "", "Task status (not started, started, completed)")

	flag.Parse()

	if err := loadTodos(ctx); err != nil {
		logWithTraceId(ctx).Error("Error loading todos", "error", err)
		return
	}

	switch {
	case *listFlag:
		listTodos(ctx)
	case *addFlag:
		if *description == "" {
			logWithTraceId(ctx).Error("--description is required for --add")
			break
		}
		addTodo(ctx, *description)
		saveTodos(ctx)
	case *updateFlag:
		if *id == 0 {
			logWithTraceId(ctx).Error("--id is required for --update")
			break
		}
		if *description == "" && *status == "" {
			logWithTraceId(ctx).Error("At least one of --description or --status must be provided for --update")
			break
		}
		validStatus, err := validateStatus(*status)
		if err != nil {
			logWithTraceId(ctx).Error(err.Error())
			break
		}
		err = updateTodo(ctx, *id, *description, validStatus)
		if err != nil {
			break
		}

		saveTodos(ctx)
	case *deleteFlag:
		if *id == 0 {
			logWithTraceId(ctx).Error("--id is required for --delete")
			break
		}
		err := deleteTodo(ctx, *id)
		if err != nil {
			break
		}
		saveTodos(ctx)
	default:
		logWithTraceId(ctx).Info("No valid flag provided")
	}

	fmt.Println("Application is running. Press Ctrl + C to exit...")
	<-ctx.Done()
	logWithTraceId(ctx).Info("shutdown application")
}

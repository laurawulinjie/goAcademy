package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"text/tabwriter"
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

func loadTodos() error {
	file, err := os.Open(dataFile)

	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("todos file not found, starting with emtpy list", "file", dataFile)
			return nil
		}
		return fmt.Errorf("failed to open file: %v", err)
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&todos); err != nil {
		return fmt.Errorf("failed to decode todos: %v", err)
	}

	if len(todos) > 0 {
		nextID = todos[len(todos)-1].ID + 1
	}

	logger.Info("todos loaded", "todoLength", len(todos))
	return nil
}

func saveTodos() error {
	file, err := os.Create(dataFile)

	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	")

	if err := encoder.Encode(todos); err != nil {
		return fmt.Errorf("failed to encode todos: %v", err)
	}

	logger.Info("todos saved", "todoLength", len(todos))
	return nil
}

func listTodos() {
	if len(todos) == 0 {
		logger.Info("No todos found")
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "ID\tTask\tStatus")

	for _, todo := range todos {
		fmt.Fprintf(writer, "%d\t%s\t%s\n", todo.ID, todo.Task, todo.Status)
	}

	writer.Flush()
}

func addTodo(description string) error {
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
	logger.Info("Created new todo", "id", newTodo.ID, "task", newTodo.Task, "status", newTodo.Status)
	return nil
}

func updateTodo(id int, description string, status Status) error {
	for i, todo := range todos {
		if todo.ID == id {
			if description != "" {
				todos[i].Task = description
			}

			if status != "" {
				todos[i].Status = status
			}

			logger.Info("Updated todo", "id", id, "task", todos[i].Task, "status", todos[i].Status)
			return nil
		}
	}

	logger.Error("Todo not found", "id", id)
	return fmt.Errorf("todo not found")
}

func deleteTodo(id int) error {
	for i, todo := range todos {
		if todo.ID == id {
			todos = slices.Delete(todos, i, i+1)
			fmt.Printf("Deleted todo ID %d\n", id)
			return nil
		}
	}

	logger.Error("Todo not found", "id", id)
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

	listFlag := flag.Bool("list", false, "List all todos")
	addFlag := flag.Bool("add", false, "Add a new todo")
	updateFlag := flag.Bool("update", false, "Update a todo")
	deleteFlag := flag.Bool("delete", false, "Delete a todo")
	description := flag.String("description", "", "Task description")
	id := flag.Int("id", 0, "Todo ID")
	status := flag.String("status", "", "Task status (not started, started, completed)")

	flag.Parse()

	if err := loadTodos(); err != nil {
		logger.Error("Error loading todos", "error", err)
		return
	}

	switch {
	case *listFlag:
		listTodos()
	case *addFlag:
		if *description == "" {
			logger.Error("Error: --description is required for --add")
			return
		}
		addTodo(*description)
		saveTodos()
	case *updateFlag:
		if *id == 0 {
			logger.Error("Error: --id is required for --update")
			return
		}
		if *description == "" && *status == "" {
			logger.Error("Error: At least one of --description or --status must be provided for --update")
			return
		}
		validStatus, err := validateStatus(*status)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		updateTodo(*id, *description, validStatus)
		saveTodos()
	case *deleteFlag:
		if *id == 0 {
			fmt.Println("Error: --id is required for --delete")
			return
		}
		deleteTodo(*id)
		saveTodos()
	}
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

func loadTodos() error {
	file, err := os.Open(dataFile)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open file: %v", err)
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.Decode(&todos)

	if len(todos) > 0 {
		nextID = todos[len(todos)-1].ID + 1
	}

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

	return nil
}

func listTodos() {
	if len(todos) == 0 {
		fmt.Println("No todos found")
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "ID\tTask\tStatus")

	for _, todo := range todos {
		fmt.Fprintf(writer, "%d\t%s\t%s\n", todo.ID, todo.Task, todo.Status)
	}

	writer.Flush()
}

func addTodo(description string) {
	newTodo := Todo{
		ID:     nextID,
		Task:   description,
		Status: NotStarted,
	}

	todos = append(todos, newTodo)
	nextID++
	fmt.Printf("Created todo with ID %d: %s [%s]\n", newTodo.ID, newTodo.Task, newTodo.Status)
}

func updateTodo(id int, description string, status Status) {
	for i, todo := range todos {
		if todo.ID == id {
			if description != "" {
				todos[i].Task = description
			}

			if status != "" {
				todos[i].Status = status
			}

			fmt.Printf("Updated todo ID %d: %s [%s]\n", id, todos[i].Task, todos[i].Status)
			return
		}
	}

	fmt.Printf("Error: Todo with ID %d not found\n", id)
	os.Exit(1)
}

func deleteTodo(id int) {
	for i, todo := range todos {
		if todo.ID == id {
			todos = slices.Delete(todos, i, i+1)
			fmt.Printf("Deleted todo ID %d\n", id)
			return
		}
	}

	fmt.Printf("Error: Todo with ID %d not found\n", id)
	os.Exit(1)
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
	listFlag := flag.Bool("list", false, "List all todos")
	addFlag := flag.Bool("add", false, "Add a new todo")
	updateFlag := flag.Bool("update", false, "Update a todo")
	deleteFlag := flag.Bool("delete", false, "Delete a todo")
	description := flag.String("description", "", "Task description")
	id := flag.Int("id", 0, "Todo ID")
	status := flag.String("status", "", "Task status (not started, started, completed)")

	flag.Parse()

	if err := loadTodos(); err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		os.Exit(1)
	}

	switch {
	case *listFlag:
		listTodos()
	case *addFlag:
		if *description == "" {
			fmt.Println("Error: --description is required for --add")
			os.Exit(1)
		}
		addTodo(*description)
		saveTodos()
	case *updateFlag:
		if *id == 0 {
			fmt.Println("Error: --id is required for --update")
			os.Exit(1)
		}
		if *description == "" && *status == "" {
			fmt.Println("Error: At least one of --description or --status must be provided for --update")
			os.Exit(1)
		}
		validStatus, err := validateStatus(*status)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		updateTodo(*id, *description, validStatus)
		saveTodos()
	case *deleteFlag:
		if *id == 0 {
			fmt.Println("Error: --id is required for --delete")
			os.Exit(1)
		}
		deleteTodo(*id)
		saveTodos()
	}
}

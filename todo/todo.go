package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type Status string

const (
	NotStarted Status = "not started"
	Started    Status = "started"
	Completed  Status = "completed"
)

type Todo struct {
	ID     int
	Task   string
	Status Status
}

var todos []Todo
var nextID = 1

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

func createTodo(description string, status Status) {
	newTodo := Todo{
		ID:     nextID,
		Task:   description,
		Status: status,
	}
	todos = append(todos, newTodo)
	nextID++
	fmt.Printf("Created todo with ID %d: %s [%s]\n", newTodo.ID, newTodo.Task, newTodo.Status)
}

func updateTodo(id int, description string, status Status) {
	var found bool
	for i, todo := range todos {
		if todo.ID == id {
			found = true
			if description != "" {
				todos[i].Task = description
			}
			if status != "" {
				todos[i].Status = status
			}
			fmt.Printf("Updated todo ID %d: %s [%s]\n", id, todos[i].Task, todos[i].Status)
			break
		}
	}
	if !found {
		fmt.Printf("Error: Todo with ID %d not found\n", id)
		os.Exit(1)
	}
}

func deleteTodo(id int) {
	var found bool
	for i, todo := range todos {
		if todo.ID == id {
			found = true
			todos = append(todos[:i], todos[i+1:]...)
			fmt.Printf("Deleted todo ID %d\n", id)
			break
		}
	}
	if !found {
		fmt.Printf("Error: Todo with ID %d not found\n", id)
		os.Exit(1)
	}
}

func main() {

}

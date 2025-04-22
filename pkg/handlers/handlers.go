package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/laurawulinjie/goAcademy/pkg/todo"
	"github.com/laurawulinjie/goAcademy/pkg/utils"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Task string `json:"task"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	newTodo, err := todo.CreateTodo(ctx, req.Task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.SaveTodos(ctx)
	json.NewEncoder(w).Encode(newTodo)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "Returning todos")
	json.NewEncoder(w).Encode(todo.GetAllTodos(ctx))
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		ID     int    `json:"id"`
		Task   string `json:"task"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	status, err := utils.ValidateStatus(req.Status)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := todo.UpdateTodo(ctx, req.ID, req.Task, status); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	todo.SaveTodos(ctx)
	w.WriteHeader(http.StatusOK)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := todo.DeleteTodo(ctx, req.ID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	todo.SaveTodos(ctx)
	w.WriteHeader(http.StatusOK)
}

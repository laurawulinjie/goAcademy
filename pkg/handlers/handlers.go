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
		return
	}

	response := make(chan any)
	todo.RequestQueue <- todo.Request{
		Ctx:      ctx,
		Action:   "create",
		Payload:  req.Task,
		Response: response,
	}

	res := (<-response).(struct {
		Todo todo.Todo
		Err  error
	})

	if res.Err != nil {
		http.Error(w, res.Err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(res.Todo)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := make(chan any)

	todo.RequestQueue <- todo.Request{
		Ctx:      ctx,
		Action:   "getAll",
		Payload:  nil,
		Response: response,
	}

	slog.InfoContext(ctx, "Returning todos")

	res := (<-response).(struct {
		Todos map[int]todo.Todo
		Err   error
	})

	if res.Err != nil {
		http.Error(w, res.Err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res.Todos)
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

	response := make(chan any)
	todo.RequestQueue <- todo.Request{
		Ctx:    ctx,
		Action: "update",
		Payload: struct {
			ID     int
			Task   string
			Status string
		}{req.ID, req.Task, status},
		Response: response,
	}

	res := (<-response).(struct {
		Err error
	})

	if res.Err != nil {
		http.Error(w, res.Err.Error(), http.StatusNoContent)
		return
	}

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

	response := make(chan any)
	todo.RequestQueue <- todo.Request{
		Ctx:      ctx,
		Action:   "delete",
		Payload:  req.ID,
		Response: response,
	}

	res := (<-response).(struct {
		Err error
	})
	if res.Err != nil {
		http.Error(w, res.Err.Error(), http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
}

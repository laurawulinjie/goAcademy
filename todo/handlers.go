package main

import (
	"encoding/json"
	"net/http"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Task string `json:"task"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	todo, err := CreateTodo(ctx, req.Task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SaveTodos(ctx)
	json.NewEncoder(w).Encode(todo)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	Log(ctx).Info("Returning todos")
	json.NewEncoder(w).Encode(GetAllTodos(ctx))
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

	status, err := ValidateStatus(req.Status)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := UpdateTodo(ctx, req.ID, req.Task, status); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	SaveTodos(ctx)
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

	if err := DeleteTodo(ctx, req.ID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	SaveTodos(ctx)
	w.WriteHeader(http.StatusOK)
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := PageData{
		Todos: GetAllTodos(ctx),
	}

	if err := tmpl.ExecuteTemplate(w, "home.html", data); err != nil {
		Log(ctx).Error("render error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ListPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := PageData{
		Todos: GetAllTodos(ctx),
	}

	if err := tmpl.ExecuteTemplate(w, "list.html", data); err != nil {
		Log(ctx).Error("render error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func AboutPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	staticFS, err := getStaticFS()
	if err != nil {
		Log(ctx).Error("failed to get static FS", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	r = r.Clone(r.Context())
	r.URL.Path = "about.html"

	http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
}

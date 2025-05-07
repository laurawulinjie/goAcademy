package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/laurawulinjie/goAcademy/pkg/todo"
	"github.com/laurawulinjie/goAcademy/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func getUserID(ctx context.Context) (int, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserID(ctx)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Task string `json:"task"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := make(chan any)
	todo.RequestQueue <- todo.Request{
		Ctx:    ctx,
		Action: "create",
		Payload: struct {
			Task   string
			UserID int
		}{req.Task, userID},
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
	userID, err := getUserID(ctx)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := make(chan any)

	todo.RequestQueue <- todo.Request{
		Ctx:      ctx,
		Action:   "getAll",
		Payload:  userID,
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
	userID, err := getUserID(ctx)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
			UserID int
		}{req.ID, req.Task, status, userID},
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
	userID, err := getUserID(ctx)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
		Payload:  struct{ ID, UserID int }{req.ID, userID},
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = todo.DB.Exec(
		`INSERT INTO users (username, password) VALUES ($1, $2)`,
		req.Username, hashedPassword)

	if err != nil {
		slog.Error("Failed to register user", "username", req.Username, "error", err)
		http.Error(w, "Failed to register user (possibly duplicate username)", http.StatusInternalServerError)
		return
	}

	slog.Info("User registered successfully", "username", req.Username)
	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	var user todo.User

	err := todo.DB.QueryRow(
		`SELECT id, password FROM users WHERE username = $1`,
		req.Username).Scan(&user.ID, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		slog.Error("Failed to query user", "username", req.Username, "error", err)
		http.Error(w, "Failed to query user", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: fmt.Sprint(user.ID),
		Path:  "/",
	})

	slog.Info("User logged in successfully", "user_id", user.ID, "username", req.Username)
	w.WriteHeader(http.StatusOK)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

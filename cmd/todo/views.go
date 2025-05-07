package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/laurawulinjie/goAcademy/pkg/todo"
)

var (
	//go:embed pages/static
	staticFS embed.FS
	//go:embed pages/dynamic
	dynamicFS embed.FS
	tmpl      *template.Template
)

type PageData struct {
	Todos    map[int]todo.Todo
	Username string
}

func SetupDynamicPages() error {
	subFS, err := fs.Sub(dynamicFS, "pages/dynamic")
	if err != nil {
		return err
	}
	tmpl, err = template.ParseFS(subFS, "*.html")
	return err
}

func ServeHomePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int)

	if !ok {
		slog.ErrorContext(ctx, "user ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	todos, _ := todo.GetAllTodos(ctx, userID)
	var username string
	err := todo.DB.QueryRowContext(ctx, "SELECT username FROM users WHERE id = $1", userID).Scan(&username)

	if err != nil {
		slog.ErrorContext(ctx, "failed to get username", "error", err)
	}

	data := PageData{
		Todos:    todos,
		Username: username,
	}

	if err := tmpl.ExecuteTemplate(w, "home.html", data); err != nil {
		slog.ErrorContext(ctx, "render error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ServeListPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int)

	if !ok {
		slog.ErrorContext(ctx, "user ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	todos, _ := todo.GetAllTodos(ctx, userID)

	data := PageData{
		Todos: todos,
	}

	if err := tmpl.ExecuteTemplate(w, "list.html", data); err != nil {
		slog.ErrorContext(ctx, "render error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := tmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		slog.ErrorContext(ctx, "render error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ServeRegisterPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := tmpl.ExecuteTemplate(w, "register.html", nil); err != nil {
		slog.ErrorContext(ctx, "render error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ServeAboutPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	staticFS, err := fs.Sub(staticFS, "pages/static")
	if err != nil {
		slog.ErrorContext(ctx, "failed to get static FS", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	r = r.Clone(ctx)
	r.URL.Path = "about.html"

	http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
}

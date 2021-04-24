package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/imgabe/todo/pkg/api/web/controllers"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/tasks", controllers.NewTasksController())

	return r
}

func NewServer(port string, handler http.Handler) *http.Server {
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      handler,
		Addr:         "0.0.0.0:" + port,
	}

	return srv
}

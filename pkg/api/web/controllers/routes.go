package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/imgabe/todo/pkg/errors"
	"github.com/imgabe/todo/pkg/models"
	"github.com/imgabe/todo/pkg/store"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// ContextKey ...
type ContextKey string

var ctx context.Context

var (
	// TaskStoreContextKey ...
	TaskStoreContextKey ContextKey = "taskstore"
	// DBContextKey ...
	DatabaseContextKey ContextKey = "db"
	// TaskContexKey ...
	TaskContexKey ContextKey = "task"
)

func init() {
	db := OpenDatabase("database")
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	ctx = context.WithValue(context.Background(), DatabaseContextKey, db)
	ts := store.TaskStore{DB: db}
	ctx = context.WithValue(context.Background(), TaskStoreContextKey, ts)
}

func OpenDatabase(path string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("error opening database: %s", err)
	}

	db.MustExec(store.CreateDatabaseStatement)

	return db
}

func NewTasksController() *chi.Mux {
	r := chi.NewRouter()
	tc := TasksController{}

	r.Get("/", tc.All)   // GET /tasks - read a list of tasks
	r.Post("/", tc.Post) // POST /tasks - create a new task and persist it

	r.Route("/{taskID}", func(r chi.Router) {
		r.Use(TaskCtx)

		r.Get("/", tc.Get)       // GET /tasks/{taskID} - read a single task by :taskID
		r.Put("/", tc.Put)       // PUT /tasks/{taskID} - update a single task by :taskID
		r.Delete("/", tc.Delete) // DELETE /tasks/{taskID} - delete a single task by :taskID
	})

	return r
}

type TasksController struct{}

func TaskCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var task models.Task

		if strTaskID := chi.URLParam(r, "taskID"); strTaskID != "" {
			taskID, err := strconv.ParseInt(strTaskID, 10, 64)
			if err != nil {
				log.Fatal("Fail convert to int64")
				return
			}

			if ts, ok := ctx.Value(TaskStoreContextKey).(store.TaskStore); ok {
				task, err = ts.Select(models.Task{ID: taskID})
				if err != nil {
					render.Render(w, r, errors.ErrNotFound)
					return
				}
			}
		}

		ctx := context.WithValue(r.Context(), TaskContexKey, task)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (t TasksController) All(w http.ResponseWriter, r *http.Request) {
	ts := ctx.Value(TaskStoreContextKey).(store.TaskStore)
	tasks, err := ts.SelectAll(true)
	if err != nil {
		render.Render(w, r, errors.ErrNotFound)
	}

	render.JSON(w, r, tasks)
}

func (t TasksController) Post(w http.ResponseWriter, r *http.Request) {
	data := &models.Task{}
	ts := ctx.Value(TaskStoreContextKey).(store.TaskStore)

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, errors.ErrInvalidRequest(err))
		return
	}

	_, err := ts.Insert(*data)
	if err != nil {
		render.Render(w, r, errors.ErrNotFound)
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, data)
}

func (t TasksController) Get(w http.ResponseWriter, r *http.Request) {
	if task, ok := r.Context().Value(TaskContexKey).(models.Task); ok {
		render.JSON(w, r, task)
	}
}

func (t TasksController) Put(w http.ResponseWriter, r *http.Request) {
	if task, ok := r.Context().Value(TaskContexKey).(models.Task); ok {
		ts := ctx.Value(TaskStoreContextKey).(store.TaskStore)
		data := &models.Task{}

		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}

		newTask := &models.Task{ID: task.ID, Description: data.Description, Done: data.Done}
		updateTask, err := ts.Update(*newTask)
		if err != nil {
			render.Render(w, r, &updateTask)
		}

		render.Render(w, r, &updateTask)
	}
}

func (t TasksController) Delete(w http.ResponseWriter, r *http.Request) {
	if task, ok := r.Context().Value(TaskContexKey).(*models.Task); ok {
		ts := ctx.Value(TaskStoreContextKey).(store.TaskStore)
		taskID := int64(task.ID)

		err := ts.Delete(taskID)
		if err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, task)
	}
}

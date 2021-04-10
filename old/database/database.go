package database

import (
	"log"

	"github.com/imgabe/todo/database/models"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type todo struct {
	db *sqlx.DB
}

func StartDB(path string) (*todo, error) {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db.MustExec(`
	CREATE TABLE IF NOT EXISTS todo (
		id   INTEGER NOT NULL PRIMARY KEY,
		task TEXT    NOT NULL,
		done BOOL    NOT NULL
	)`)

	return &todo{db}, nil
}

func (t *todo) Add(task string) (models.TodoModel, error) {
	tx := t.db.MustBegin()
	todo := models.TodoModel{Task: task, Done: false}
	result, err := tx.NamedExec("INSERT INTO todo (task, done) VALUES (:task, :done)", todo)
	if err != nil {
		log.Fatal(err)
	}

	todoID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	todo.ID = int(todoID)
	err = tx.Commit()

	return todo, err
}

func (t *todo) Check(id int) (int, error) {
	tx := t.db.MustBegin()
	tx.MustExec("UPDATE todo SET done = true WHERE ID = $1", id)

	err := tx.Commit()

	return id, err
}

func (t *todo) Clear() error {
	tx := t.db.MustBegin()
	tx.MustExec("DELETE FROM todo")
	err := tx.Commit()

	return err
}

func (t *todo) List() ([]models.TodoModel, error) {
	todos := []models.TodoModel{}
	tx := t.db.MustBegin()
	err := tx.Select(&todos, "SELECT * from todo")

	return todos, err
}

func (t *todo) Delete(id int) (int, error) {
	tx := t.db.MustBegin()
	tx.MustExec("DELETE FROM todo WHERE ID = $1", id)
	err := tx.Commit()

	return id, err
}

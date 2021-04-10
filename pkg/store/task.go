package store

import (
	"database/sql"

	"github.com/imgabe/todo/pkg/models"
	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	DB *sqlx.DB
}

func (s TaskStore) Insert(task models.Task) (models.Task, error) {
	insertStmt := `
		INSERT INTO task (description, done)
		VALUES (:description, :done);
	`

	selectStmt := `
		SELECT *
		FROM task
		WHERE id = $1
	`

	var received models.Task

	result, err := s.DB.NamedExec(insertStmt, &task)
	if err != nil {
		return received, err
	}

	lastId, _ := result.LastInsertId()
	err = s.DB.Get(&received, selectStmt, lastId)
	if err != nil {
		return received, err
	}

	return received, nil
}

func (s TaskStore) Update(task models.Task) (models.Task, error) {
	stmt := `
		UPDATE task
		SET description = :description,
			done = :done
		WHERE id = :id
	`

	selectStmt := `
		SELECT *
		FROM task
		WHERE id = $1
	`

	var received models.Task

	result, err := s.DB.NamedExec(stmt, &task)
	if err != nil {
		return received, err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return received, sql.ErrNoRows
	}

	err = s.DB.Get(&received, selectStmt, task.ID)
	if err != nil {
		return received, err
	}

	return received, nil
}

func (s TaskStore) Delete(task models.Task) error {
	stmt := `
		DELETE FROM task
		WHERE id = :id
	`

	result, err := s.DB.NamedExec(stmt, &task)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s TaskStore) Select(task models.Task) (models.Task, error) {
	stmt := `
		SELECT *
		FROM task
		WHERE id = :id
	`

	var received models.Task
	err := s.DB.Get(&received, stmt, task.ID)
	if err != nil {
		return received, err
	}

	return received, err
}

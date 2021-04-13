package store

import (
	"database/sql"

	"github.com/imgabe/todo/pkg/models"
	"github.com/jmoiron/sqlx"
)

// TaskStore is responsible for all database actions related to tasks
type TaskStore struct {
	DB *sqlx.DB
}

// Insert inserts a new task on the database
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

// Update updates an existent task on the database
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

// Delete deletes a task on the database
func (s TaskStore) Delete(taskID int64) error {
	stmt := `
		DELETE FROM task
		WHERE id = $1
	`

	result, err := s.DB.Exec(stmt, &taskID)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Select retrieves a task from the database
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

// SelectAll retrieves all tasks from the database
func (s TaskStore) SelectAll(done bool) ([]models.Task, error) {
	stmt := `
		SELECT *
		FROM task
		WHERE done = false OR done = :done
	`

	var tasks []models.Task

	err := s.DB.Select(&tasks, stmt, done)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// Check checks a task on the database
func (s TaskStore) Check(taskID int64) error {
	stmt := `
	UPDATE task
	SET done = True
	WHERE id = :id
	`

	result, err := s.DB.Exec(stmt, &taskID)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

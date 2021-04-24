package models

import (
	"errors"
	"net/http"
)

type Task struct {
	ID          int64  `db:"id" json:"id"`
	Description string `db:"description" json:"description"`
	Done        bool   `db:"done" json:"done"`
}

func (t *Task) Bind(r *http.Request) error {
	if t.Description == "" {
		return errors.New("missing required Description fields")
	}

	return nil

}

func (t *Task) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

package models

type Task struct {
	ID          int64  `db:"id" json:"id"`
	Description string `db:"description" json:"description"`
	Done        bool   `db:"done" json:"done"`
}

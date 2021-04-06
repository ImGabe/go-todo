package models

type Todo interface {
	Add(task string) (*TodoModel, error)
	Check(id int) (int, error)
	Clear() error
	List() (*[]TodoModel, error)
	Delete(id int) (int, error)
}

type TodoModel struct {
	ID   int    `db:"id"`
	Task string `db:"task"`
	Done bool   `db:"done"`
}

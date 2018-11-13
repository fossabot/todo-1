package store

import "fmt"

// ErrNoResults is a generic error of sql.ErrNoRows
var ErrNoResults = fmt.Errorf("no results returned")

// Todo holds information about a Todo
type Todo struct {
	ID          int64  `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
	IsCompleted bool   `json:"isCompleted" db:"is_completed"`
}

// Service provides methods for interacting with a store
type Service interface {
	CreateTodo(todo *Todo) error
	GetTodos() ([]Todo, error)
	GetTodo(id int64) (Todo, error)
	UpdateTodo(todo Todo) error
	DeleteTodo(id int64) error
	Close() error
}

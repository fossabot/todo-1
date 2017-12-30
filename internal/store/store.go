package store

import "fmt"

// ErrNoResults is a generic error of sql.ErrNoRows
var ErrNoResults = fmt.Errorf("no results returned")

// Todo holds information about a Todo
type Todo struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	IsCompleted bool   `json:"isCompleted"`
}

// NullableTodo is a Todo with all nullable fields
type NullableTodo struct {
	ID          *int64  `json:"id"`
	Description *string `json:"description"`
	IsCompleted *bool   `json:"isCompleted"`
}

// Service provides methods for interacting with a store
type Service interface {
	CreateTodo(todo Todo) (int64, error)
	GetTodo(id int64) (Todo, error)
	GetTodos() ([]Todo, error)
	UpdateTodo(todo Todo) error
	PatchTodo(nt NullableTodo) error
	DeleteTodo(id int64) error
	Close() error
}

package store

import "fmt"

// ErrNoResults is a generic error of sql.ErrNoRows
var ErrNoResults = fmt.Errorf("no results returned")

// Todo holds information about a Todo
type Todo struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"createdAt"`
	CompletedAt *int64 `json:"completedAt,omitempty"` // nullable
}

// NullableTodo is a Todo with all nullable fields
type NullableTodo struct {
	ID          *int64  `json:"id"`
	Description *string `json:"description"`
	CreatedAt   *int64  `json:"createdAt"`
	CompletedAt *int64  `json:"completedAt"`
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

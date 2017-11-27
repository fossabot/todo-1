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
	DeleteTodo(id int64) error
	Close() error
}

// Populate populates an existing Todo with the not null fields of a NullableTodo, and
// returns the updated Todo
func Populate(t Todo, nt NullableTodo) Todo {
	if nt.ID != nil {
		t.ID = *nt.ID
	}
	if nt.Description != nil {
		t.Description = *nt.Description
	}
	if nt.CreatedAt != nil {
		t.CreatedAt = *nt.CreatedAt
	}
	if nt.CompletedAt != nil {
		t.CompletedAt = nt.CompletedAt
	}

	return t
}

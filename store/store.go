package store

import "fmt"

// ErrNoResults is a generic error of sql.ErrNoRows
var ErrNoResults = fmt.Errorf("no results returned")

type Todo struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"createdAt"`
	CompletedAt *int64 `json:"completedAt,omitempty"` // nullable
}

type NullableTodo struct {
	ID          *int64  `json:"id"`
	Description *string `json:"description"`
	CreatedAt   *int64  `json:"createdAt"`
	CompletedAt *int64  `json:"completedAt"`
}

type Service interface {
	CreateTodo(todo Todo) (int64, error)
	GetTodo(id int64) (Todo, error)
	GetTodos() ([]Todo, error)
	UpdateTodo(todo Todo) error
	DeleteTodo(id int64) error
	Close() error
}

func Populate(nt NullableTodo, t Todo) Todo {
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

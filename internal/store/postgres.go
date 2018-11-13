package store

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // pq is used as the postgres sql driver
)

type postgres struct {
	db *sqlx.DB
}

// ErrNotFound is returned when a todo is not found.
type ErrNotFound error

// NewPostgres returns a new postgresql store from the given postgres-specific
// data source name.
func NewPostgres(dsn string) (Service, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return postgres{db}, db.Ping()
}

func (p postgres) CreateTodo(todo *Todo) error {
	err := p.db.Get(&todo.ID, `
	INSERT INTO
		todos (description, is_completed)
		VALUES ($1, $2)
		RETURNING id
	`, todo.Description, todo.IsCompleted)
	return err
}

func (p postgres) GetTodo(id int64) (Todo, error) {
	var todo Todo

	err := p.db.Get(&todo, "SELECT * FROM todos WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return todo, ErrNotFound(err)
	}

	return todo, err
}

func (p postgres) GetTodos() ([]Todo, error) {
	var todos []Todo
	err := p.db.Select(&todos, "SELECT * FROM todos")
	return todos, err
}

func (p postgres) UpdateTodo(todo Todo) error {
	_, err := p.db.NamedExec(`
	UPDATE todos
		SET
			description = :description,
			is_completed = :is_completed
		WHERE
			id = :id
	`, todo)
	return err
}

func (p postgres) DeleteTodo(id int64) error {
	res, err := p.db.Exec("DELETE FROM todos WHERE id = $1", id)
	if rows, err := res.RowsAffected(); err != nil && rows == 0 {
		return ErrNotFound(err)
	}

	return err
}

func (p postgres) Close() error {
	return p.db.Close()
}

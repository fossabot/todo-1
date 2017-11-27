package postgres

import (
	"database/sql"
	"fmt"

	"github.com/fharding1/todo/store"

	// for the postgres sql driver
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type service struct {
	db *sql.DB
}

// Options holds information for connecting to a postgresql server
type Options struct {
	User, Pass string
	Host       string
	Port       int
	DBName     string
	SSLMode    string
}

func (o Options) connectionInfo() string {
	return fmt.Sprintf("host='%s' port='%d' user='%s' password='%s' dbname='%s' sslmode='%s'",
		o.Host, o.Port, o.User, o.Pass, o.DBName, o.SSLMode)
}

const todoTableCreationQuery = `
CREATE TABLE IF NOT EXISTS todos (
	id          SERIAL PRIMARY KEY,
	description varchar(256),
	createdAt   bigint,
	completedAt bigint
)`

// New connects to a postgres server with specified options and returns a store.Service
func New(options Options) (store.Service, error) {
	db, err := sql.Open("postgres", options.connectionInfo())
	if err != nil {
		return nil, errors.Wrap(err, "connecting to postgres database")
	}

	_, err = db.Exec(todoTableCreationQuery)
	if err != nil {
		return nil, errors.Wrap(err, "creating todos table")
	}

	return &service{db: db}, nil
}

func (s *service) CreateTodo(todo store.Todo) (id int64, err error) {
	err = s.db.QueryRow(
		"INSERT INTO todos (description, createdAt, completedAt) VALUES ($1, $2, $3) RETURNING id",
		todo.Description, todo.CreatedAt, todo.CompletedAt).Scan(&id)
	return
}

func (s *service) GetTodo(id int64) (store.Todo, error) {
	todo := store.Todo{ID: id}
	err := s.db.QueryRow("SELECT description, createdAt, completedAt FROM todos WHERE id = $1", id).Scan(
		&todo.Description, &todo.CreatedAt, &todo.CompletedAt)
	if err == sql.ErrNoRows {
		err = store.ErrNoResults
	}
	return todo, err
}

func (s *service) GetTodos() ([]store.Todo, error) {
	rows, err := s.db.Query("SELECT id, description, createdAt, completedAt FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []store.Todo{}

	for rows.Next() {
		var todo store.Todo
		if err := rows.Scan(&todo.ID, &todo.Description, &todo.CreatedAt, &todo.CompletedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (s *service) UpdateTodo(todo store.Todo) error {
	_, err := s.db.Exec("UPDATE todos SET description = $1, createdAt = $2, completedAt = $3 WHERE id = $4",
		todo.Description, todo.CreatedAt, todo.CompletedAt, todo.ID)
	return err
}

func (s *service) PatchTodo(nt store.NullableTodo) error {
	_, err := s.db.Exec(`
		UPDATE todos SET
		description = COALESCE($1, description),
		createdAt = COALESCE($2, createdAt),
		completedAt = COALESCE($3, completedAt)
		WHERE id = $4
		`, nt.Description, nt.CreatedAt, nt.CompletedAt, nt.ID)
	return err
}

func (s *service) DeleteTodo(id int64) error {
	_, err := s.db.Exec("DELETE FROM todos WHERE id = $1", id)
	return err
}

func (s *service) Close() error {
	return s.db.Close()
}

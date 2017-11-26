package postgres

import (
	"database/sql"
	"fmt"

	gstore "github.com/fharding1/todo/store"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type store struct {
	db *sql.DB
}

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

func New(options Options) (gstore.Service, error) {
	db, err := sql.Open("postgres", options.connectionInfo())
	if err != nil {
		return nil, errors.Wrap(err, "connecting to postgres database")
	}

	_, err = db.Exec(todoTableCreationQuery)
	if err != nil {
		return nil, errors.Wrap(err, "creating todos table")
	}

	return &store{db: db}, nil
}

func (s *store) CreateTodo(todo gstore.Todo) (id int64, err error) {
	err = s.db.QueryRow(
		"INSERT INTO todos (description, createdAt, completedAt) VALUES ($1, $2, $3) RETURNING id",
		todo.Description, todo.CreatedAt, todo.CompletedAt).Scan(&id)
	return
}

func (s *store) GetTodo(id int64) (gstore.Todo, error) {
	todo := gstore.Todo{ID: id}
	err := s.db.QueryRow("SELECT description, createdAt, completedAt FROM todos WHERE id = $1", id).Scan(
		&todo.Description, &todo.CreatedAt, &todo.CompletedAt)
	if err == sql.ErrNoRows {
		err = gstore.ErrNoResults
	}
	return todo, err
}

func (s *store) GetTodos() ([]gstore.Todo, error) {
	rows, err := s.db.Query("SELECT description, createdAt, completedAt FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []gstore.Todo{}

	for rows.Next() {
		var todo gstore.Todo
		if err := rows.Scan(&todo.Description, &todo.CreatedAt, &todo.CompletedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (s *store) UpdateTodo(todo gstore.Todo) error {
	_, err := s.db.Exec("UPDATE todos SET description = $1, createdAt = $2, completedAt = $3 WHERE id = $4",
		todo.Description, todo.CreatedAt, todo.CompletedAt, todo.ID)
	return err
}

func (s *store) DeleteTodo(id int64) error {
	_, err := s.db.Exec("DELETE FROM todos WHERE id = $1", id)
	return err
}

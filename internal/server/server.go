package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/fharding1/todo/internal/respond"
	"github.com/fharding1/todo/internal/store"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type server struct {
	sto     store.Service
	handler http.Handler
}

// New creates a new server from a store and populates the handler
func New(sto store.Service) *server {
	s := &server{sto: sto}

	router := mux.NewRouter()

	router.Handle("/todo", allowedMethods(
		[]string{"OPTIONS", "GET", "POST"},
		handlers.MethodHandler{
			"GET":  http.HandlerFunc(s.getTodos),
			"POST": http.HandlerFunc(s.createTodo),
		}))

	router.Handle("/todo/{id}", idMiddleware(allowedMethods(
		[]string{"OPTIONS", "GET", "PUT", "PATCH", "DELETE"},
		handlers.MethodHandler{
			"GET":    http.HandlerFunc(s.getTodo),
			"PUT":    http.HandlerFunc(s.putTodo),
			"PATCH":  http.HandlerFunc(s.patchTodo),
			"DELETE": http.HandlerFunc(s.deleteTodo),
		})))

	s.handler = limitBody(defaultHeaders(router))

	return s
}

// Run starts the server listening on what address is specified
func (s *server) Run(addr string) error {
	return http.ListenAndServe(addr, s.handler)
}

func allowedMethods(methods []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))

		next.ServeHTTP(w, r)
	})
}

func (s *server) createTodo(w http.ResponseWriter, r *http.Request) {
	var todo store.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := s.sto.CreateTodo(todo)
	if err != nil {
		respond.JSON(w, err)
		return
	}

	respond.JSON(w, map[string]int64{"id": id})
}

func (s *server) getTodo(w http.ResponseWriter, r *http.Request) {
	todo, err := s.sto.GetTodo(r.Context().Value(keyIDCtx).(int64))
	if err != nil {
		if err == store.ErrNoResults {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			respond.JSON(w, err)
		}
		return
	}

	respond.JSON(w, todo)
}

func (s *server) getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := s.sto.GetTodos()
	if err != nil {
		if err == store.ErrNoResults {
			todos = []store.Todo{}
		} else {
			respond.JSON(w, err)
			return
		}
	}

	respond.JSON(w, todos)
}

func (s *server) putTodo(w http.ResponseWriter, r *http.Request) {
	var todo store.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	todo.ID = r.Context().Value(keyIDCtx).(int64)

	if err := s.sto.UpdateTodo(todo); err != nil {
		respond.JSON(w, err)
		return
	}
}

func (s *server) patchTodo(w http.ResponseWriter, r *http.Request) {
	var todo store.NullableTodo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id := r.Context().Value(keyIDCtx).(int64)
	todo.ID = &id

	if err := s.sto.PatchTodo(todo); err != nil {
		respond.JSON(w, err)
		return
	}
}

func (s *server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	if err := s.sto.DeleteTodo(r.Context().Value(keyIDCtx).(int64)); err != nil {
		respond.JSON(w, err)
		return
	}
}

type key int

const (
	keyIDCtx key = iota
)

func idMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawID := mux.Vars(r)["id"]

		id, err := strconv.ParseInt(rawID, 10, 64)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), keyIDCtx, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func defaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		next.ServeHTTP(w, r)
	})
}

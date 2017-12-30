package server

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	router.Handle("/todo/{id}", allowedMethods(
		[]string{"OPTIONS", "GET", "PUT", "PATCH", "DELETE"},
		handlers.MethodHandler{
			"GET":    http.HandlerFunc(s.getTodo),
			"PUT":    http.HandlerFunc(s.putTodo),
			"PATCH":  http.HandlerFunc(s.patchTodo),
			"DELETE": http.HandlerFunc(s.deleteTodo),
		}))

	s.handler = limitBody(defaultHeaders(router))

	return s
}

// Run starts the server listening on what address is specified
func (s *server) Run(addr string) error {
	return http.ListenAndServe(addr, s.handler)
}

func commaify(ss []string) (out string) {
	for i, s := range ss {
		out += s
		if i != len(ss)-1 {
			out += ","
		}
	}
	return
}

func allowedMethods(methods []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", commaify(methods))

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
	rawID := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	todo, err := s.sto.GetTodo(id)
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
	rawID := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var todo store.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	todo.ID = id

	if err := s.sto.UpdateTodo(todo); err != nil {
		respond.JSON(w, err)
		return
	}
}

func (s *server) patchTodo(w http.ResponseWriter, r *http.Request) {
	rawID := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var todo store.NullableTodo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	todo.ID = &id

	if err := s.sto.PatchTodo(todo); err != nil {
		respond.JSON(w, err)
		return
	}
}

func (s *server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	rawID := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := s.sto.DeleteTodo(id); err != nil {
		respond.JSON(w, err)
		return
	}
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

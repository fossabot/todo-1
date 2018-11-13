package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/fharding1/todo/internal/respond"
	"github.com/fharding1/todo/internal/store"
	"github.com/gorilla/mux"
)

type Server struct {
	sto     store.Service
	handler http.Handler
}

// New creates a new server from a store and populates the handler
func New(sto store.Service) *Server {
	s := &Server{sto: sto}

	router := mux.NewRouter()

	router.HandleFunc("/todos", s.getTodos).Methods("GET")
	router.HandleFunc("/todos", s.createTodo).Methods("POST")
	router.Handle("/todos/{id}", idMiddleware(http.HandlerFunc(s.getTodo))).Methods("GET")
	router.Handle("/todos/{id}", idMiddleware(http.HandlerFunc(s.putTodo))).Methods("PUT")
	router.Handle("/todos/{id}", idMiddleware(http.HandlerFunc(s.patchTodo))).Methods("PATCH")
	router.Handle("/todos/{id}", idMiddleware(http.HandlerFunc(s.deleteTodo))).Methods("DELETE")

	s.handler = limitBody(defaultHeaders(router))

	return s
}

// Run starts the server listening on what address is specified
func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.handler)
}

func (s *Server) createTodo(w http.ResponseWriter, r *http.Request) {
	var todo store.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respond.JSON(w, nil, err, http.StatusUnprocessableEntity)
		return
	}

	err := s.sto.CreateTodo(&todo)
	if err != nil {
		respond.JSON(w, nil, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, todo, nil, http.StatusCreated)
}

func (s *Server) getTodo(w http.ResponseWriter, r *http.Request) {
	todo, err := s.sto.GetTodo(r.Context().Value(keyIDCtx).(int64))
	if err != nil {
		if err == store.ErrNoResults {
			respond.JSON(w, nil, err, http.StatusNotFound)
			return
		}
		respond.JSON(w, nil, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, todo, nil, http.StatusOK)
}

func (s *Server) getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := s.sto.GetTodos()
	if err != nil {
		respond.JSON(w, nil, err, http.StatusInternalServerError)
		return
	}

	if todos == nil {
		todos = make([]store.Todo, 0)
	}

	respond.JSON(w, todos, nil, http.StatusOK)
}

func (s *Server) putTodo(w http.ResponseWriter, r *http.Request) {
	var todo store.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respond.JSON(w, nil, err, http.StatusUnprocessableEntity)
		return
	}

	todo.ID = r.Context().Value(keyIDCtx).(int64)

	if err := s.sto.UpdateTodo(todo); err != nil {
		if err, ok := err.(store.ErrNotFound); ok {
			respond.JSON(w, nil, err, http.StatusNotFound)
			return
		}
		respond.JSON(w, nil, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

const jsonPatchContentType = "application/json-patch+json"

func (s *Server) patchTodo(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != jsonPatchContentType {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var patch jsonpatch.Patch
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		respond.JSON(w, nil, err, http.StatusUnprocessableEntity)
		return
	}

	todo, err := s.sto.GetTodo(r.Context().Value(keyIDCtx).(int64))
	if err != nil {
		if err, ok := err.(store.ErrNotFound); ok {
			respond.JSON(w, nil, err, http.StatusNotFound)
			return
		}
		respond.JSON(w, nil, err, http.StatusInternalServerError)
		return
	}

	originalDocument, err := json.Marshal(todo)
	if err != nil {
		respond.JSON(w, nil, err, http.StatusInternalServerError)
		return
	}

	modifiedDocument, err := patch.Apply(originalDocument)
	if err != nil {
		respond.JSON(w, nil, err, http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(modifiedDocument, &todo); err != nil {
		respond.JSON(w, nil, err, http.StatusInternalServerError)
		return
	}

	if err := s.sto.UpdateTodo(todo); err != nil {
		if err, ok := err.(store.ErrNotFound); ok {
			respond.JSON(w, nil, err, http.StatusNotFound)
			return
		}
		respond.JSON(w, nil, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	err := s.sto.DeleteTodo(r.Context().Value(keyIDCtx).(int64))
	if err != nil {
		if err, ok := err.(store.ErrNotFound); ok {
			respond.JSON(w, nil, err, http.StatusNotFound)
			return
		}
		respond.JSON(w, nil, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

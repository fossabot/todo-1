package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/fharding1/todo/store"
	"github.com/gorilla/mux"
)

// Server holds information necessary for a Todo server, such as the store
// and the http handler
type Server struct {
	sto     store.Service
	handler http.Handler
}

// New creates a new Server from a store and populates the handler
func New(sto store.Service) *Server {
	router := mux.NewRouter()

	s := &Server{sto: sto}

	router.HandleFunc("/todo", s.createTodo).Methods("POST")
	router.HandleFunc("/todo", s.getTodos).Methods("GET")
	router.HandleFunc("/todo/{id}", s.getTodo).Methods("GET")
	router.HandleFunc("/todo/{id}", s.putTodo).Methods("PUT")
	router.HandleFunc("/todo/{id}", s.patchTodo).Methods("PATCH")
	router.HandleFunc("/todo/{id}", s.deleteTodo).Methods("DELETE")

	s.handler = limitBody(router)

	return s
}

// Run starts the server listening on what address is specified
func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.handler)
}

func (s *Server) createTodo(w http.ResponseWriter, r *http.Request) {
	var todo store.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	todo.CreatedAt = time.Now().Unix()

	id, err := s.sto.CreateTodo(todo)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func (s *Server) getTodo(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func (s *Server) getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := s.sto.GetTodos()
	if err != nil {
		if err == store.ErrNoResults {
			todos = []store.Todo{}
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(todos)
}

func (s *Server) putTodo(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) patchTodo(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	rawID := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := s.sto.DeleteTodo(id); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		next.ServeHTTP(w, r)
	})
}

package server

import (
	"net/http"

	"github.com/fharding1/todo/store"
	"github.com/gorilla/mux"
)

type Server struct {
	sto     store.Service
	handler http.Handler
}

func New(sto store.Service) *Server {
	router := mux.NewRouter()

	s := &Server{sto: sto}

	router.HandleFunc("/todo", s.createTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", s.getTodo).Methods("GET")
	router.HandleFunc("/todo", s.getTodos).Methods("GET")
	router.HandleFunc("/todo", s.updateTodo).Methods("PUT")
	router.HandleFunc("/todo/{id}", s.deleteTodo).Methods("DELETE")

	s.handler = limitBody(router)

	return s
}

func (s *Server) Run() error {
	return http.ListenAndServe(":8080", s.handler)
}

func (s *Server) createTodo(w http.ResponseWriter, r *http.Request) {}
func (s *Server) getTodo(w http.ResponseWriter, r *http.Request)    {}
func (s *Server) getTodos(w http.ResponseWriter, r *http.Request)   {}
func (s *Server) updateTodo(w http.ResponseWriter, r *http.Request) {}
func (s *Server) deleteTodo(w http.ResponseWriter, r *http.Request) {}

func limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.Body = http.MaxBytesReader(w, r, 8192)

		next.ServeHTTP(w, r)
	})
}

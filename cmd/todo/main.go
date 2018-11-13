package main

import (
	"log"
	"os"

	"github.com/fharding1/todo/internal/server"
	"github.com/fharding1/todo/internal/store"
)

func main() {
	sto, err := store.NewPostgres(os.Getenv("TODO_POSTGRES_DSN"))
	if err != nil {
		log.Fatalf("connecting to postgres database: %v\n", err)
	}

	s := server.New(sto)

	addr := os.Getenv("TODO_ADDR")

	if err := s.Run(addr); err != nil {
		log.Fatalf("running server: %v\n", err)
	}

	sto.Close()
}

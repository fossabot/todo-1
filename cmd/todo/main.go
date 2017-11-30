package main

import (
	"log"
	"os"
	"strconv"

	"github.com/fharding1/todo/internal/server"
	"github.com/fharding1/todo/internal/store/postgres"
)

func main() {
	portString := os.Getenv("POSTGRES_PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		log.Fatalf("invalid port: %s\n", portString)
	}

	sto, err := postgres.New(postgres.Options{
		User:    os.Getenv("POSTGRES_USER"),
		Pass:    os.Getenv("POSTGRES_PASS"),
		Host:    os.Getenv("POSTGRES_HOST"),
		Port:    port,
		DBName:  os.Getenv("POSTGRES_DB_NAME"),
		SSLMode: os.Getenv("POSTGRES_SSL_MODE"),
	})

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

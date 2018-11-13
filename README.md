# Todo

A clean todo web backend written in Go.

## Setup

### Dependencies

- [Install PostgreSQL](https://www.postgresql.org/docs/9.2/tutorial-install.html)
- [Install migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cli)

### Migrations

    migrate -path ./internal/store/migrations -database "postgres://user:pass@addr?sslmode=disable&dbname=dbname" up

### Build

    cd cmd/todo
    go build

### Configure

    export TODO_POSTGRES_DSN="sslmode=disable user=postgres dbname=todos"
    export TODO_ADDR=":8080"

### Run

    ./todo

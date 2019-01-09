# Todo
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ffharding1%2Ftodo.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Ffharding1%2Ftodo?ref=badge_shield)


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

    export TODO_POSTGRES_DSN="sslmode=disable user=user dbname=dbname"
    export TODO_ADDR=":8080"

### Run

    ./todo


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ffharding1%2Ftodo.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Ffharding1%2Ftodo?ref=badge_large)
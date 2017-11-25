package store

import "fmt"

// ErrNoResults is a generic error of sql.ErrNoRows
var ErrNoResults = fmt.Errorf("no results returned")

type Service interface {
}

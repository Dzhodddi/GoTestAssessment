package graph

import (
	"WorkAssigment/internal/store"
	"database/sql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	catService *store.CatStore
}

func NewResolver(db *sql.DB) *Resolver {
	return &Resolver{
		catService: store.NewCatStore(db),
	}
}

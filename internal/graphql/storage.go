package graphql

import (
	"FIDOtestBackendApp/internal/store"
	"database/sql"
	"github.com/graphql-go/graphql"
)

type GPQLStorage struct {
	Cat interface {
		GetListOfCats() *graphql.Result
	}
}

func NewGPQLStorage(conn *sql.DB) *GPQLStorage {
	return &GPQLStorage{
		Cat: &Cat{
			resolver: &Resolver{
				catService: store.NewCatStore(conn),
			},
		},
	}
}

package graphql

import (
	"FIDOtestBackendApp/internal/store"
	"database/sql"
	"github.com/graphql-go/graphql"
)

type Resolver struct {
	catService *store.CatStore
}

func NewResolver(db *sql.DB) *Resolver {
	return &Resolver{
		catService: store.NewCatStore(db),
	}
}
func populate() []SpyCatInfo {
	return []SpyCatInfo{
		{
			ID:               1,
			Name:             "Whisker Shadow",
			YearOfExperience: 5,
			Breed:            "Siberian",
			Salary:           50000,
			Mission: Mission{
				ID:       101,
				Complete: false,
				Targets: []Target{
					{
						ID:       201,
						Name:     "Dr. Meow",
						Country:  "Germany",
						Notes:    "Has a secret lab in Berlin",
						Complete: false,
					},
				},
			},
		},
		{
			ID:               2,
			Name:             "Agent Purr",
			YearOfExperience: 3,
			Breed:            "British Shorthair",
			Salary:           45000,
			Mission: Mission{
				ID:       102,
				Complete: true,
				Targets: []Target{
					{
						ID:       202,
						Name:     "The Dogfather",
						Country:  "Italy",
						Notes:    "Operates from Naples",
						Complete: true,
					},
					{
						ID:       203,
						Name:     "Ruffian",
						Country:  "France",
						Notes:    "Allied with The Dogfather",
						Complete: true,
					},
				},
			},
		},
	}
}

func (r *Resolver) getListOfCats() (interface{}, error) {
	return populate(), nil
}

func (r *Resolver) getOneCat(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if ok {
		return r.catService.GetByID(p.Context, int64(id))
	}
	return nil, nil
}

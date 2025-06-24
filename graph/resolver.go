package graph

import (
	"WorkAssigment/graph/model"
	"WorkAssigment/internal/store"
	"database/sql"
	"errors"
	"github.com/graphql-go/graphql"
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

func (r *Resolver) getListOfCats(p graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) getOneCat(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, errors.New("id is required")
	}
	catId := int64(id)

	smt, err := r.catService.GetCatWithMissionAndTargets(p.Context, catId)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			return nil, nil
		default:
			return nil, err
		}
	}
	var targets []*model.Target
	for _, t := range smt.Target {
		targets = append(targets, &model.Target{
			ID:       int(t.ID),
			Name:     t.Name,
			Country:  t.Country,
			Notes:    t.Notes,
			Complete: t.Completed,
		})
	}
	var mission *model.Mission
	if smt.Mission != nil {
		mission = &model.Mission{
			ID:       int(smt.Mission.ID),
			Complete: smt.Mission.Completed,
			Targets:  targets,
		}
	}

	return &model.SpyCatInfo{
		ID:               int(smt.Cat.ID),
		Name:             smt.Cat.Name,
		YearOfExperience: int32(smt.Cat.Experience),
		Breed:            smt.Cat.Breed,
		Mission:          mission,
	}, nil
}

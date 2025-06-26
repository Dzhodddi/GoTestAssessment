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

//func (r *Resolver) getOneCat(ctx context.Context, id int64) (interface{}, error) {
//
//	smt, err := r.catService.GetCatWithMissionAndTargets(ctx, id)
//
//	if err != nil {
//		switch err {
//		case store.ErrNotFound:
//			return nil, nil
//		default:
//			return nil, err
//		}
//	}
//	var targets []*model.Target
//	for _, t := range smt.Target {
//		targets = append(targets, &model.Target{
//			ID:       int(t.ID),
//			Name:     t.Name,
//			Country:  t.Country,
//			Notes:    t.Notes,
//			Complete: t.Completed,
//		})
//	}
//	var mission *model.Mission
//	if smt.Mission != nil {
//		mission = &model.Mission{
//			ID:       int(smt.Mission.ID),
//			Complete: smt.Mission.Completed,
//			Targets:  targets,
//		}
//	}
//
//	return &model.SpyCatInfo{
//		ID:               int(smt.Cat.ID),
//		Name:             smt.Cat.Name,
//		YearOfExperience: int32(smt.Cat.Experience),
//		Breed:            smt.Cat.Breed,
//		Mission:          mission,
//	}, nil
//}

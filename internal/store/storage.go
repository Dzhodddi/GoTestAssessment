package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	QueryTimeOut      = 5 * time.Second
	ErrNotFound       = errors.New("record not found")
	MissionedAssigned = errors.New("mission assigned")
)

type Storage struct {
	Cat interface {
		CreateSpyCat(ctx context.Context, spyCat *Cat) error
		DeleteSpyCat(ctx context.Context, id int64) error
		GetByID(ctx context.Context, id int64) (*Cat, error)
		UpdateSpyCat(ctx context.Context, spyCat *Cat) error
		GetPaginatedSpyCatList(ctx context.Context, paginatedQuery PaginatedQuery) ([]*Cat, error)
	}
	Mission interface {
		CreateMission(ctx context.Context, mission *MissionWithTargets) error
		DeleteMission(ctx context.Context, id int64) error
		UpdateMissionStatus(ctx context.Context, mission *UpdatedMission) error
	}
	Target interface {
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Cat:     &CatStore{db},
		Mission: &MissionStore{db},
		Target:  &TargetStore{db},
	}
}

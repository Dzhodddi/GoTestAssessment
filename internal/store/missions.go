package store

import (
	"context"
	"database/sql"
	"time"
)

type Mission struct {
	ID        int64  `json:"id"`
	CatID     *int64 `json:"cat_id"`
	Completed bool   `json:"completed"`
}

type MissionWithTargets struct {
	Mission Mission
	Targets []Target
}

type MissionStore struct {
	db *sql.DB
}

type UpdatedMission struct {
	ID     int64 `json:"id"`
	Status bool  `json:"status"`
}

func (s *MissionStore) CreateMission(ctx context.Context, mission *MissionWithTargets) error {
	const queryAddMission = `INSERT INTO missions (cat_id, completed) VALUES ($1, $2) RETURNING id`
	const queryAddTargets = `INSERT INTO targets (mission_id, name, country, notes, completed) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, _ := s.db.BeginTx(ctx, nil)

	var missionID int64
	err := tx.QueryRowContext(ctx, queryAddMission, mission.Mission.CatID, mission.Mission.Completed).Scan(&missionID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for i := range mission.Targets {
		target := &mission.Targets[i]
		target.MissionID = missionID
		err = tx.QueryRowContext(ctx, queryAddTargets, missionID, target.Name, target.Country, target.Notes, target.Completed).Scan(&target.ID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	_ = tx.Commit()

	return nil
}

func (s *MissionStore) DeleteMission(ctx context.Context, id int64) error {
	checkIdQuery := `SELECT cat_id FROM missions WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var catID *int64
	err := s.db.QueryRowContext(ctx, checkIdQuery, id).Scan(&catID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	if catID != nil {
		return MissionedAssigned
	}

	res, err := s.db.ExecContext(ctx, "DELETE FROM missions WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *MissionStore) UpdateMissionStatus(ctx context.Context, missionState *UpdatedMission) error {
	query := `UPDATE missions SET completed = $1 WHERE id = $2`

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, missionState.Status, missionState.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

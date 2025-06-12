package store

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"log"
	"time"
)

type Target struct {
	ID        int64  `json:"id"`
	MissionID int64  `json:"mission_id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Notes     string `json:"notes"`
	Completed bool   `json:"completed"`
}

type UpdateTargetNote struct {
	ID        int64  `json:"id"`
	MissionID int64  `json:"mission_id"`
	Note      string `json:"notes"`
}

type UpdateTargetStatus struct {
	ID        int64 `json:"id"`
	MissionID int64 `json:"mission_id"`
	Status    bool  `json:"status"`
}

type TargetStore struct {
	db *sql.DB
}

func (s *TargetStore) UpdateTargetNote(ctx context.Context, updateNote *UpdateTargetNote) error {
	query := `
	UPDATE targets t 
	SET notes = $1
	FROM missions m
	WHERE t.id = $2 AND t.mission_id = m.id AND m.id = $3 AND t.completed = false AND m.completed = false
	RETURNING t.id;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	var id int64
	err := s.db.QueryRowContext(ctx, query, updateNote.Note, updateNote.ID, updateNote.MissionID).Scan(&id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *TargetStore) UpdateTargetStatus(ctx context.Context, updateTargetStatus *UpdateTargetStatus) error {
	updateQuery := `UPDATE targets SET completed = $1 WHERE id = $2 AND mission_id = $3 RETURNING id;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	tx, _ := s.db.BeginTx(ctx, nil)
	var id int64
	err := tx.QueryRowContext(ctx, updateQuery, updateTargetStatus.Status, updateTargetStatus.ID, updateTargetStatus.MissionID).Scan(&id)
	if err != nil {
		_ = tx.Rollback()
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	if updateTargetStatus.Status {
		allTargetsCompleted := `SELECT NOT EXISTS (
    	SELECT 1 FROM targets WHERE mission_id = $1 AND completed = false)`
		var completed bool
		err = tx.QueryRowContext(ctx, allTargetsCompleted, updateTargetStatus.ID).Scan(&completed)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if completed {
			updateMissionStatus := `UPDATE missions SET completed = true WHERE id = $1`
			_, err = tx.ExecContext(ctx, updateMissionStatus, updateTargetStatus.MissionID)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (s *TargetStore) DeleteTarget(ctx context.Context, missionID, targetID int64) error {
	query := `DELETE FROM targets WHERE id = $1 AND mission_id = $2 AND completed = false`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, targetID, missionID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *TargetStore) AddTarget(ctx context.Context, target *Target) error {
	query := `
	SELECT
	  m.completed,
	  (SELECT COUNT(*) FROM targets WHERE mission_id = $1) AS existing_count
	FROM missions m
	WHERE m.id = $1;
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	var completed bool
	var count int64
	err := s.db.QueryRowContext(ctx, query, target.MissionID).Scan(&completed, &count)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}
	log.Printf("count %v", count)
	if completed || count >= 3 {
		return TargetAmountError
	}
	insertQuery := `INSERT INTO targets (mission_id, name, country, notes, completed) VALUES ($1, $2, $3, $4, $5)`
	_, err = s.db.ExecContext(ctx, insertQuery, target.MissionID, target.Name, target.Country, target.Notes, completed)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return ViolatePK
			}
		}
		return err
	}
	return nil
}

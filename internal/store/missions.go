package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
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
type MissionWithMetadata struct {
	Mission Mission
	Cat     *Cat
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
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
			if pgErr, ok := err.(*pq.Error); ok {
				if pgErr.Code == "23505" {
					return ViolatePK
				}
			}
			_ = tx.Rollback()
			return err
		}
	}

	_ = tx.Commit()

	return nil
}

func (s *MissionStore) DeleteMission(ctx context.Context, id int64) error {
	checkIdQuery := `SELECT cat_id FROM missions WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
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

func (s *MissionStore) AddCatToMission(ctx context.Context, catID, missionID int64) error {
	var completed bool
	err := s.db.QueryRowContext(ctx, `
		SELECT completed FROM missions WHERE id = $1
	`, missionID).Scan(&completed)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	if completed {
		return MissionCompleted
	}

	var exists bool
	err = s.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM spycat WHERE id = $1)`, catID).Scan(&exists)

	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}

	result, err := s.db.ExecContext(ctx, `
		UPDATE missions
		SET cat_id = $1
		WHERE id = $2
	`, catID, missionID)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return ViolatePK
			}
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *MissionStore) GetMissionList(ctx context.Context) ([]*MissionWithMetadata, error) {
	rows, err := s.db.QueryContext(ctx, `
	SELECT 
		m.id,
		m.completed,
		m.cat_id,
		c.id,
		c.name,
		c.years,
		c.breed,
		c.salary
	FROM missions m
	LEFT JOIN spycat c ON m.cat_id = c.id
	ORDER BY m.id ASC`)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	var missions []*MissionWithMetadata
	for rows.Next() {
		m := &MissionWithMetadata{}
		var catID sql.NullInt64
		var catName sql.NullString
		var catYears sql.NullInt64
		var catBreed sql.NullString
		var catSalary sql.NullInt64

		err = rows.Scan(
			&m.Mission.ID,
			&m.Mission.Completed,
			&m.Mission.CatID,
			&catID,
			&catName,
			&catYears,
			&catBreed,
			&catSalary,
		)
		if err != nil {
			return nil, err
		}

		if catID.Valid {
			m.Cat = &Cat{
				ID:         catID.Int64,
				Name:       catName.String,
				Experience: int(catYears.Int64),
				Breed:      catBreed.String,
				Salary:     int(catSalary.Int64),
			}
		}

		missions = append(missions, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return missions, nil
}

func (s *MissionStore) GetOneMission(ctx context.Context, id int64) (*MissionWithMetadata, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT 
			m.id,
			m.completed,
			m.cat_id,
			c.id,
			c.name,
			c.years,
			c.breed,
			c.salary
		FROM missions m
		LEFT JOIN spycat c ON m.cat_id = c.id
		WHERE m.id = $1`, id)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	m := &MissionWithMetadata{}
	var catID sql.NullInt64
	var catName sql.NullString
	var catYears sql.NullInt64
	var catBreed sql.NullString
	var catSalary sql.NullInt64

	err := row.Scan(
		&m.Mission.ID,
		&m.Mission.Completed,
		&m.Mission.CatID,
		&catID,
		&catName,
		&catYears,
		&catBreed,
		&catSalary,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if catID.Valid {
		m.Cat = &Cat{
			ID:         catID.Int64,
			Name:       catName.String,
			Experience: int(catYears.Int64),
			Breed:      catBreed.String,
			Salary:     int(catSalary.Int64),
		}
	}

	return m, nil
}

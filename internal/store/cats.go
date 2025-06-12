package store

import (
	"context"
	"database/sql"
	"errors"
)

type Cat struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Experience int    `json:"year_of_experience"`
	Breed      string `json:"breed"`
	Salary     int    `json:"salary"`
}

type CatStore struct {
	db *sql.DB
}

func (s *CatStore) CreateSpyCat(ctx context.Context, cat *Cat) error {
	query := `INSERT INTO spycat (name, years, breed, salary) VALUES ($1, $2, $3, $4) RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, cat.Name, cat.Experience, cat.Breed, cat.Salary).Scan(&cat.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *CatStore) DeleteSpyCat(ctx context.Context, id int64) error {
	query := `DELETE FROM spycat WHERE ID =  $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
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

func (s *CatStore) GetByID(ctx context.Context, id int64) (*Cat, error) {
	query := `SELECT id, name, years, breed, salary FROM spycat WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	cat := &Cat{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&cat.ID, &cat.Name, &cat.Experience, &cat.Breed, &cat.Salary)
	if err != nil {
		return nil, ErrNotFound
	}
	return cat, nil
}

func (s *CatStore) UpdateSpyCat(ctx context.Context, cat *Cat) error {
	query := `UPDATE spycat SET salary = $1 WHERE ID = $2 RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, cat.Salary, cat.ID).Scan(&cat.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *CatStore) GetPaginatedSpyCatList(ctx context.Context, paginatedQuery PaginatedQuery) ([]*Cat, error) {
	query := "SELECT id, name, years, breed, salary FROM spycat ORDER BY id LIMIT $1 OFFSET $2;"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, paginatedQuery.Limit, paginatedQuery.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var cats []*Cat
	for rows.Next() {
		var cat Cat
		err = rows.Scan(&cat.ID, &cat.Name, &cat.Experience, &cat.Breed, &cat.Salary)
		if err != nil {
			return nil, err
		}
		cats = append(cats, &cat)
	}
	return cats, nil
}

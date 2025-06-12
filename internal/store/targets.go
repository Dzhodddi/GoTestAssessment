package store

import "database/sql"

type Target struct {
	ID        int64  `json:"id"`
	MissionID int64  `json:"mission_id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Notes     string `json:"notes"`
	Completed bool   `json:"completed"`
}

type TargetStore struct {
	db *sql.DB
}

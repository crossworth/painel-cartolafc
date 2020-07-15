package database

import (
	"database/sql"
)

type DatabaseService struct {
	db *sql.DB
}

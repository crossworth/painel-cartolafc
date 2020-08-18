package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	db *sql.DB
}

func NewPostgreSQL(dsn string) (*PostgreSQL, error) {
	psql := &PostgreSQL{}

	var err error
	psql.db, err = sql.Open("postgres", dsn)
	if err != nil {
		return psql, err
	}

	err = psql.db.Ping()
	if err != nil {
		return psql, err
	}

	return psql, nil
}

func (p *PostgreSQL) GetDB() *sql.DB {
	return p.db
}

package database

import (
	"context"

	"github.com/travelaudience/go-sx"
)

func (p *PostgreSQL) SettingByName(context context.Context, name string) (string, error) {
	var value string

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT value FROM settings WHERE name = $1`, name).MustScan(&value)
	})

	return value, err
}

func (p *PostgreSQL) UpdateSetting(context context.Context, name string, value string) error {
	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustExecContext(context, `UPDATE settings SET value = $1 WHERE name = $2`, value, name)
	})

	return err
}

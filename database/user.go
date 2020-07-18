package database

import (
	"github.com/travelaudience/go-sx"

	"github.com/crossworth/cartola-web-admin/model"
)

func (d *PostgreSQL) FindProfileByID(id int) (model.Profile, error) {
	var profile model.Profile

	err := sx.Do(d.db, func(tx *sx.Tx) {
		tx.MustQueryRow(`SELECT * FROM profiles WHERE id = $1`, id).MustScans(&profile)
	})

	return profile, err
}

func (d *PostgreSQL) FindTopicByUser(id int, before int64, after int64, limit int) ([]model.Topic, int64, error) {
	var topics []model.Topic
	var total int64

	err := sx.Do(d.db, func(tx *sx.Tx) {
		tx.MustQueryRow(`SELECT COUNT(*) FROM topics WHERE created_by = $1`, id).MustScan(&total)

		query := `SELECT * FROM topics WHERE created_by = $1 AND created_at < $2 ORDER BY created_at DESC LIMIT $3`
		cond := before

		// TODO(Pedro): Does this work?
		if cond == 0 {
			query = `SELECT * FROM topics WHERE created_by = $1 AND created_at > $2 ORDER BY created_at DESC LIMIT $3`
			cond = after
		}

		tx.MustQuery(query, id, cond, limit).Each(func(r *sx.Rows) {
			var o model.Topic
			r.MustScans(&o)
			topics = append(topics, o)
		})
	})

	return topics, total, err
}

func (d *PostgreSQL) FindTopicCount(id int) (uint64, error) {
	return 0, nil
}

func (d *PostgreSQL) FindCommentsByUser(id int) ([]model.Comment, error) {
	return []model.Comment{}, nil
}

func (d *PostgreSQL) FindCommentCount(id int) (uint64, error) {
	return 0, nil
}

func (d *PostgreSQL) FindProfileHistory(id int) ([]model.ProfileNames, error) {
	return []model.ProfileNames{}, nil
}

package database

import (
	"github.com/travelaudience/go-sx"

	"github.com/crossworth/cartola-web-admin/model"
)

func (d *PostgreSQL) ProfileByUserID(id int) (model.Profile, error) {
	var profile model.Profile

	err := sx.Do(d.db, func(tx *sx.Tx) {
		tx.MustQueryRow(`SELECT * FROM profiles WHERE id = $1`, id).MustScans(&profile)
	})

	return profile, err
}

func (d *PostgreSQL) TopicsByUserID(id int, before int, limit int) ([]model.Topic, error) {
	var topics []model.Topic

	err := sx.Do(d.db, func(tx *sx.Tx) {
		query := `SELECT * FROM topics WHERE created_by = $1 AND created_at < $2 ORDER BY created_at DESC LIMIT $3`

		tx.MustQuery(query, id, before, limit).Each(func(r *sx.Rows) {
			var t model.Topic
			r.MustScans(&t)
			topics = append(topics, t)
		})
	})

	return topics, err
}

func (d *PostgreSQL) TopicsCountByUserID(id int) (int, error) {
	var total int

	err := sx.Do(d.db, func(tx *sx.Tx) {
		tx.MustQueryRow(`SELECT COUNT(*) FROM topics WHERE created_by = $1`, id).MustScan(&total)
	})

	return total, err
}

func (d *PostgreSQL) CommentsByUserID(id int, before int, limit int) ([]model.Comment, error) {
	var comments []model.Comment

	err := sx.Do(d.db, func(tx *sx.Tx) {
		query := `SELECT * FROM comments WHERE from_id = $1 AND date < $2 ORDER BY date DESC LIMIT $3`

		tx.MustQuery(query, id, before, limit).Each(func(r *sx.Rows) {
			var c model.Comment
			r.MustScans(&c)
			comments = append(comments, c)
		})
	})

	return comments, err
}

func (d *PostgreSQL) CommentsCountByUserID(id int) (int, error) {
	var total int

	err := sx.Do(d.db, func(tx *sx.Tx) {
		tx.MustQueryRow(`SELECT COUNT(*) FROM comments WHERE from_id = $1`, id).MustScan(&total)
	})

	return total, err
}

func (d *PostgreSQL) ProfileHistoryByUserID(id int) ([]model.ProfileNames, error) {
	var profileHistory []model.ProfileNames

	err := sx.Do(d.db, func(tx *sx.Tx) {
		query := `SELECT * FROM profile_names WHERE profile_id = $1 ORDER BY date DESC`

		tx.MustQuery(query, id).Each(func(r *sx.Rows) {
			var p model.ProfileNames
			r.MustScans(&p)
			profileHistory = append(profileHistory, p)
		})
	})

	return profileHistory, err
}

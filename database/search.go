package database

import (
	"github.com/travelaudience/go-sx"

	"github.com/crossworth/cartola-web-admin/model"
)

func (d *PostgreSQL) SearchTopics(text string, before int, limit int) ([]model.Topic, error) {
	var topics []model.Topic

	err := sx.Do(d.db, func(tx *sx.Tx) {
		query := `SELECT * FROM topics WHERE title LIKE '%$1% AND created_at < $2 ORDER BY created_at DESC LIMIT $3`

		tx.MustQuery(query, text, before, limit).Each(func(r *sx.Rows) {
			var t model.Topic
			r.MustScans(&t)
			topics = append(topics, t)
		})
	})

	return topics, err
}

func (d *PostgreSQL) SearchComments(text string, before int, limit int) ([]model.Comment, error) {
	var comments []model.Comment

	err := sx.Do(d.db, func(tx *sx.Tx) {
		query := `SELECT * FROM comments WHERE text LIKE '%$1%' AND date < $2 ORDER BY date DESC LIMIT $3`

		tx.MustQuery(query, text, before, limit).Each(func(r *sx.Rows) {
			var c model.Comment
			r.MustScans(&c)
			comments = append(comments, c)
		})
	})

	return comments, err
}

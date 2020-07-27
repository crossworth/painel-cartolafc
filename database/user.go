package database

import (
	"context"
	"strconv"
	"strings"

	"github.com/travelaudience/go-sx"

	"github.com/crossworth/cartola-web-admin/model"
)

func (d *PostgreSQL) ProfileByUserID(context context.Context, id int) (model.Profile, error) {
	var profile model.Profile

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT * FROM profiles WHERE id = $1`, id).MustScans(&profile)
	})

	return profile, err
}

func (d *PostgreSQL) TopicsByUserID(context context.Context, id int, before int, limit int) ([]model.Topic, error) {
	var topics []model.Topic

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT * FROM topics WHERE created_by = $1 AND created_at < $2 ORDER BY created_at DESC LIMIT $3`

		tx.MustQueryContext(context, query, id, before, limit).Each(func(r *sx.Rows) {
			var t model.Topic
			r.MustScans(&t)
			topics = append(topics, t)
		})
	})

	return topics, err
}

func (d *PostgreSQL) PrevTopicTimestampByUserID(context context.Context, id int, before int, limit int) (int, error) {
	var prev int

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT * FROM topics WHERE created_by = $1 AND created_at > $2 ORDER BY created_at ASC limit $3`

		var topics []model.Topic
		tx.MustQueryContext(context, query, id, before, limit).Each(func(r *sx.Rows) {
			var t model.Topic
			r.MustScans(&t)
			topics = append(topics, t)
		})

		if len(topics) > 0 {
			prev = topics[len(topics)-1].CreatedAt + 1 // NOTE(Pedro): para incluir o tópico já que usamos before < e não <=
		}
	})

	return prev, err
}

func (d *PostgreSQL) TopicsCountByUserID(context context.Context, id int) (int, error) {
	var total int

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM topics WHERE created_by = $1`, id).MustScan(&total)
	})

	return total, err
}

func (d *PostgreSQL) CommentsByUserID(context context.Context, id int, before int, limit int) ([]model.Comment, error) {
	var comments []model.Comment

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT * FROM comments WHERE from_id = $1 AND date < $2 ORDER BY date DESC LIMIT $3`

		tx.MustQueryContext(context, query, id, before, limit).Each(func(r *sx.Rows) {
			var c model.Comment
			r.MustScans(&c)
			comments = append(comments, c)
		})
	})

	return comments, err
}

func (d *PostgreSQL) CommentsCountByUserID(context context.Context, id int) (int, error) {
	var total int

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM comments WHERE from_id = $1`, id).MustScan(&total)
	})

	return total, err
}

func (d *PostgreSQL) PrevCommentTimestampByUserID(context context.Context, id int, before int, limit int) (int, error) {
	var prev int

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT * FROM comments WHERE from_id = $1 AND date > $2 ORDER BY date ASC limit $3`

		var comments []model.Comment
		tx.MustQueryContext(context, query, id, before, limit).Each(func(r *sx.Rows) {
			var t model.Comment
			r.MustScans(&t)
			comments = append(comments, t)
		})

		if len(comments) > 0 {
			prev = comments[len(comments)-1].Date + 1 // NOTE(Pedro): para incluir o tópico já que usamos before < e não <=
		}
	})

	return prev, err
}

func (d *PostgreSQL) ProfileHistoryByUserID(context context.Context, id int) ([]model.ProfileNames, error) {
	var profileHistory []model.ProfileNames

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT * FROM profile_names WHERE profile_id = $1 ORDER BY date DESC`

		tx.MustQueryContext(context, query, id).Each(func(r *sx.Rows) {
			var p model.ProfileNames
			r.MustScans(&p)
			profileHistory = append(profileHistory, p)
		})
	})

	return profileHistory, err
}

func (d *PostgreSQL) SearchProfileName(context context.Context, text string) ([]model.Profile, error) {
	var profiles []model.Profile

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT * FROM profiles WHERE LOWER(first_name) LIKE '%' || $1 || '%' OR LOWER(last_name) LIKE '%' || $1 || '%' OR LOWER(screen_name) LIKE '%' || $1 || '%' OR LOWER(CONCAT(first_name, ' ', last_name)) LIKE '%' || $1 || '%'`

		tx.MustQueryContext(context, query, strings.ToLower(text)).Each(func(r *sx.Rows) {
			var p model.Profile
			r.MustScans(&p)

			// NOTE(Pedro): when the user dont have an screen name we normalize to the conical name
			if p.ScreenName == "" {
				p.ScreenName = "id" + strconv.Itoa(p.ID)
			}

			profiles = append(profiles, p)
		})
	})

	return profiles, err
}

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
		query := `SELECT * FROM topics WHERE created_by = $1 AND created_at <= $2 ORDER BY created_at DESC LIMIT $3`

		tx.MustQueryContext(context, query, id, before, limit).Each(func(r *sx.Rows) {
			var t model.Topic
			r.MustScans(&t)
			topics = append(topics, t)
		})
	})

	return topics, err
}

func (d *PostgreSQL) PaginationTimestampTopicByUserID(context context.Context, id int, before int, limit int) (PaginationTimestamps, error) {
	var timestamps PaginationTimestamps

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT 
					COALESCE((SELECT created_at FROM topics WHERE created_by = $1 AND created_at <= $2 ORDER BY created_at DESC OFFSET $3 LIMIT 1), 0) as next,
					COALESCE((SELECT created_at FROM topics WHERE created_by = $1 AND created_at >= $2 ORDER BY created_at ASC OFFSET $3 LIMIT 1), 0) as prev
				`
		tx.MustQueryRowContext(context, query, id, before, limit).MustScan(&timestamps.Next, &timestamps.Prev)
	})

	return timestamps, err
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
		query := `SELECT * FROM comments WHERE from_id = $1 AND date <= $2 ORDER BY date DESC LIMIT $3`

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

func (d *PostgreSQL) PaginationTimestampCommentByUserID(context context.Context, id int, before int, limit int) (PaginationTimestamps, error) {
	var timestamps PaginationTimestamps

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT 
					COALESCE((SELECT date FROM comments WHERE from_id = $1 AND date <= $2 ORDER BY date DESC OFFSET $3 LIMIT 1), 0) as next,
					COALESCE((SELECT date FROM comments WHERE from_id = $1 AND date >= $2 ORDER BY date ASC OFFSET $3 LIMIT 1), 0) as prev
				`
		tx.MustQueryRowContext(context, query, id, before, limit).MustScan(&timestamps.Next, &timestamps.Prev)
	})

	return timestamps, err
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

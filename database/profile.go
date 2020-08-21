package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/travelaudience/go-sx"

	"github.com/crossworth/cartola-web-admin/model"
)

var (
	ErrInvalidMemberOrderBy = errors.New("order by de membros inválido")
)

func (p *PostgreSQL) ProfileByID(context context.Context, id int) (model.Profile, error) {
	var profile model.Profile

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT * FROM profiles WHERE id = $1`, id).MustScans(&profile)
	})

	return profile, err
}

func (p *PostgreSQL) AdminProfiles(context context.Context) ([]model.Profile, error) {
	var profiles []model.Profile

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryContext(context, `SELECT id, first_name, last_name, screen_name, photo
FROM profiles
WHERE id IN (SELECT id FROM administrators)`).Each(func(rows *sx.Rows) {
			var profile model.Profile
			rows.MustScan(
				&profile.ID,
				&profile.FirstName,
				&profile.LastName,
				&profile.ScreenName,
				&profile.Photo,
			)
			profiles = append(profiles, profile)
		})
	})

	return profiles, err
}

func (p *PostgreSQL) TopicsByProfileID(context context.Context, id int, before int, limit int) ([]model.Topic, error) {
	var topics []model.Topic

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		query := `SELECT * FROM topics WHERE created_by = $1 AND created_at <= $2 ORDER BY created_at DESC LIMIT $3`

		tx.MustQueryContext(context, query, id, before, limit).Each(func(r *sx.Rows) {
			var t model.Topic
			r.MustScans(&t)
			topics = append(topics, t)
		})
	})

	return topics, err
}

func (p *PostgreSQL) TopicsPaginationTimestampByProfileID(context context.Context, id int, before int, limit int) (PaginationTimestamps, error) {
	var timestamps PaginationTimestamps

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		query := `SELECT 
					COALESCE((SELECT created_at FROM topics WHERE created_by = $1 AND created_at <= $2 ORDER BY created_at DESC OFFSET $3 LIMIT 1), 0) as next,
					COALESCE((SELECT created_at FROM topics WHERE created_by = $1 AND created_at >= $2 ORDER BY created_at ASC OFFSET $3 LIMIT 1), 0) as prev
				`
		tx.MustQueryRowContext(context, query, id, before, limit).MustScan(&timestamps.Next, &timestamps.Prev)
	})

	return timestamps, err
}

func (p *PostgreSQL) TopicsCountByProfileID(context context.Context, id int) (int, error) {
	var total int

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM topics WHERE created_by = $1`, id).MustScan(&total)
	})

	return total, err
}

func (p *PostgreSQL) CommentsByProfileID(context context.Context, id int, before int, limit int) ([]model.Comment, error) {
	var comments []model.Comment

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		query := `SELECT * FROM comments WHERE from_id = $1 AND date <= $2 ORDER BY date DESC LIMIT $3`

		tx.MustQueryContext(context, query, id, before, limit).Each(func(r *sx.Rows) {
			var c model.Comment
			r.MustScans(&c)
			comments = append(comments, c)
		})
	})

	return comments, err
}

func (p *PostgreSQL) CommentsCountByProfileID(context context.Context, id int) (int, error) {
	var total int

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM comments WHERE from_id = $1`, id).MustScan(&total)
	})

	return total, err
}

func (p *PostgreSQL) CommentsPaginationTimestampByProfileID(context context.Context, id int, before int, limit int) (PaginationTimestamps, error) {
	var timestamps PaginationTimestamps

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		query := `SELECT 
					COALESCE((SELECT date FROM comments WHERE from_id = $1 AND date <= $2 ORDER BY date DESC OFFSET $3 LIMIT 1), 0) as next,
					COALESCE((SELECT date FROM comments WHERE from_id = $1 AND date >= $2 ORDER BY date ASC OFFSET $3 LIMIT 1), 0) as prev
				`
		tx.MustQueryRowContext(context, query, id, before, limit).MustScan(&timestamps.Next, &timestamps.Prev)
	})

	return timestamps, err
}

func (p *PostgreSQL) ProfileHistoryByProfileID(context context.Context, id int) ([]model.ProfileNames, error) {
	var profileHistory []model.ProfileNames

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		query := `SELECT * FROM profile_names WHERE profile_id = $1 ORDER BY date DESC`

		tx.MustQueryContext(context, query, id).Each(func(r *sx.Rows) {
			var p model.ProfileNames
			r.MustScans(&p)
			profileHistory = append(profileHistory, p)
		})
	})

	return profileHistory, err
}

func (p *PostgreSQL) SearchProfileName(context context.Context, text string) ([]model.Profile, error) {
	var profiles []model.Profile

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		query := `SELECT * FROM profiles WHERE LOWER(first_name) LIKE '%' || $1 || '%' OR LOWER(last_name) LIKE '%' || $1 || '%' OR LOWER(screen_name) LIKE '%' || $1 || '%' OR LOWER(CONCAT(first_name, ' ', last_name)) LIKE '%' || $1 || '%'`

		tx.MustQueryContext(context, query, strings.ToLower(text)).Each(func(r *sx.Rows) {
			var p model.Profile
			r.MustScans(&p)

			// NOTE(Pedro): when the profile dont have an screen name we normalize to the conical name
			if p.ScreenName == "" {
				p.ScreenName = "id" + strconv.Itoa(p.ID)
			}

			profiles = append(profiles, p)
		})
	})

	return profiles, err
}

func (p *PostgreSQL) ProfileStatsByProfileID(context context.Context, id int) (ProfileWithStats, error) {
	var profile ProfileWithStats

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		query := `SELECT
    p.*,
    (SELECT COUNT(t.*) FROM topics t WHERE t.created_by = p.id) as topics,
	COALESCE(c.comments, 0)::INTEGER as comments,
    COALESCE(c.likes, 0)::INTEGER as likes
FROM profiles p
LEFT JOIN (
        SELECT from_id, COUNT(from_id) as comments, SUM(likes) as likes FROM comments GROUP BY from_id
    ) as c ON c.from_id = p.id
WHERE p.id = $1`

		tx.MustQueryRowContext(context, query, id).MustScan(
			&profile.ID,
			&profile.FirstName,
			&profile.LastName,
			&profile.ScreenName,
			&profile.Photo,
			&profile.Topics,
			&profile.Comments,
			&profile.Likes,
		)
	})

	return profile, err
}

func (p *PostgreSQL) ProfileWithStats(context context.Context, orderBy string, orderDirection OrderByDirection, period Period, page int, limit int) ([]ProfileWithStats, error) {
	var profiles []ProfileWithStats

	if orderBy != "topics" && orderBy != "comments" && orderBy != "likes" && orderBy != "topics_comments" {
		return profiles, ErrInvalidMemberOrderBy
	}

	periodTopics := ""
	periodComments := ""

	if period != PeriodAll {
		periodTopics = fmt.Sprintf("WHERE created_at >= EXTRACT(epoch FROM current_date - INTERVAL '1 %s') AND created_at < EXTRACT(epoch FROM current_date)", period.Stringer())
		periodComments = fmt.Sprintf("WHERE date >= EXTRACT(epoch FROM current_date - INTERVAL '1 %s') AND date < EXTRACT(epoch FROM current_date)", period.Stringer())
	}

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		query := `SELECT
    p.*,
    COALESCE(t.total, 0)::INTEGER as topics,
    COALESCE(c.comments, 0)::INTEGER as comments,
    COALESCE(c.likes, 0)::INTEGER as likes,
	(COALESCE(t.total, 0)::INTEGER + COALESCE(c.comments, 0)::INTEGER) as topics_comments
FROM profiles p
    LEFT JOIN (
        SELECT created_by, COUNT(created_by) as total FROM topics ` + periodTopics + ` GROUP BY created_by
    ) as t ON t.created_by = p.id
     LEFT JOIN (
        SELECT from_id, COUNT(from_id) as comments, SUM(likes) as likes FROM comments ` + periodComments + ` GROUP BY from_id
    ) as c ON c.from_id = p.id
GROUP BY p.id, t.total, c.comments, c.likes ORDER BY ` + orderBy + ` ` + orderDirection.Stringer() + ` OFFSET $1 LIMIT $2`

		i := 1
		tx.MustQueryContext(context, query, (page-1)*limit, limit).Each(func(r *sx.Rows) {
			var p ProfileWithStats
			r.MustScan(
				&p.ID,
				&p.FirstName,
				&p.LastName,
				&p.ScreenName,
				&p.Photo,
				&p.Topics,
				&p.Comments,
				&p.Likes,
				&p.TopicsPlusComments,
			)
			p.Position = ((page - 1) * limit) + i
			profiles = append(profiles, p)
			i++
		})
	})

	return profiles, err
}

func (p *PostgreSQL) ProfileCount(context context.Context) (int, error) {
	var total int

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM profiles`).MustScan(&total)
	})

	return total, err
}

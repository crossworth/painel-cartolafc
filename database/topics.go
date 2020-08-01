package database

import (
	"context"

	"github.com/travelaudience/go-sx"

	"github.com/crossworth/cartola-web-admin/model"
)

func (d *PostgreSQL) Topics(context context.Context, before int, limit int) ([]TopicWithPollAndCommentsCount, error) {
	var topics []TopicWithPollAndCommentsCount

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		var t TopicWithPollAndCommentsCount
		tx.MustQueryContext(context, `SELECT * FROM topics WHERE created_at <= $1 ORDER BY created_at DESC LIMIT $2`, before, limit).Each(func(rows *sx.Rows) {
			rows.MustScans(&t.Topic)

			tx.MustQueryRowContext(context, `SELECT * FROM polls WHERE topic_id = $1`, t.Topic.ID).MustScans(&t.Poll)
			tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM comments WHERE topic_id = $1`, t.Topic.ID).MustScans(&t.CommentsCount)

			tx.MustQueryContext(context, `SELECT * FROM poll_answers WHERE poll_id = $1`, t.Poll.ID).Each(func(rows *sx.Rows) {
				var pa model.PollAnswer
				rows.MustScans(&pa)
				t.PollWithAnswers.Answers = append(t.PollWithAnswers.Answers, pa)
			})

			topics = append(topics, t)
		})
	})

	return topics, err
}

func (d *PostgreSQL) TopicsPaginationTimestamp(context context.Context, before int, limit int) (PaginationTimestamps, error) {
	var timestamps PaginationTimestamps

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT 
					COALESCE((SELECT created_at FROM topics WHERE created_at <= $1 ORDER BY created_at DESC OFFSET $2 LIMIT 1), 0) as next,
					COALESCE((SELECT created_at FROM topics WHERE created_at >= $1 ORDER BY created_at ASC OFFSET $2 LIMIT 1), 0) as prev
				`
		tx.MustQueryRowContext(context, query, before, limit).MustScan(&timestamps.Next, &timestamps.Prev)
	})

	return timestamps, err
}

func (d *PostgreSQL) TopicsCount(context context.Context) (int, error) {
	var total int

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM topics`).MustScan(&total)
	})

	return total, err
}

func (d *PostgreSQL) TopicByID(context context.Context, id int) (TopicWithPoll, error) {
	var topic TopicWithPoll

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT * FROM topics WHERE id = $1`, id).MustScans(&topic.Topic)
		tx.MustQueryRowContext(context, `SELECT * FROM polls WHERE topic_id = $1`, topic.Topic.ID).MustScans(&topic.Poll)

		tx.MustQueryContext(context, `SELECT * FROM poll_answers WHERE poll_id = $1`, topic.Poll.ID).Each(func(rows *sx.Rows) {
			var pa model.PollAnswer
			rows.MustScans(&pa)
			topic.PollWithAnswers.Answers = append(topic.PollWithAnswers.Answers, pa)
		})
	})

	return topic, err
}

func (d *PostgreSQL) CommentsByTopicID(context context.Context, id int, before int, limit int) ([]CommentWithProfileAndAttachment, error) {
	var comments []CommentWithProfileAndAttachment

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		queryComment := `SELECT * FROM comments WHERE topic_id = $1 AND date <= $2 ORDER BY date DESC LIMIT $3`

		// TODO(Pedro): there is a faster way to to this, we could use a where in after we get all the
		// comments and populate in code
		queryProfile := `SELECT * FROM profiles FROM id = $1`
		queryAttachments := `SELECT * FROM attachments FROM comment_id = $1`

		tx.MustQueryContext(context, queryComment, id, before, limit).Each(func(r *sx.Rows) {
			var c model.Comment
			r.MustScans(&c)

			var p model.Profile
			tx.MustQueryRowContext(context, queryProfile, c.FromID).MustScans(&p)

			var attachments []model.Attachment

			tx.MustQueryContext(context, queryAttachments, c.ID).Each(func(rows *sx.Rows) {
				var a model.Attachment
				r.MustScans(&a)
				attachments = append(attachments, a)
			})

			comments = append(comments, CommentWithProfileAndAttachment{
				Comment:     c,
				Profile:     p,
				Attachments: attachments,
			})
		})
	})

	return comments, err
}

func (d *PostgreSQL) CommentsCountByTopicID(context context.Context, id int) (int, error) {
	var total int

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM comments WHERE topic_id = $1`, id).MustScan(&total)
	})

	return total, err
}

func (d *PostgreSQL) CommentsPaginationTimestampByTopicID(context context.Context, id int, before int, limit int) (PaginationTimestamps, error) {
	var timestamps PaginationTimestamps

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT 
					COALESCE((SELECT date FROM comments WHERE topic_id = $1 AND date <= $2 ORDER BY date DESC OFFSET $3 LIMIT 1), 0) as next,
					COALESCE((SELECT date FROM comments WHERE topic_id = $1 AND date >= $2 ORDER BY date ASC OFFSET $3 LIMIT 1), 0) as prev
				`
		tx.MustQueryRowContext(context, query, id, before, limit).MustScan(&timestamps.Next, &timestamps.Prev)
	})

	return timestamps, err
}

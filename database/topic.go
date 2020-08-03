package database

import (
	"context"
	"strconv"
	"strings"

	"github.com/travelaudience/go-sx"

	"github.com/crossworth/cartola-web-admin/model"
)

func (d *PostgreSQL) Topics(context context.Context, before int, limit int) ([]TopicWithPollAndCommentsCount, error) {
	var topics []TopicWithPollAndCommentsCount

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		var t TopicWithPollAndCommentsCount
		tx.MustQueryContext(context, `SELECT * FROM topics WHERE created_at <= $1 ORDER BY created_at DESC LIMIT $2`, before, limit).Each(func(rows *sx.Rows) {
			rows.MustScans(&t.Topic)
			topics = append(topics, t)
		})

		var topidsIds []int

		for i := range topics {
			topidsIds = append(topidsIds, topics[i].ID)
			// TODO(Pedro): Slow count, fix it
			tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM comments WHERE topic_id = $1`, topics[i].ID).MustScan(&topics[i].CommentsCount)
		}

		// TODO(Pedro): Add poll to the result struct using where in
		// tx.MustQueryRowContext(context, `SELECT * FROM polls WHERE topic_id = $1`, t.Topic.ID).MustScan(
		// 	&t.Poll.ID,
		// 	&t.Poll.Question,
		// 	&t.Poll.Votes,
		// 	&t.Poll.Multiple,
		// 	&t.Poll.EndDate,
		// 	&t.Poll.Closed,
		// 	&t.Poll.TopicID,
		// )

		// //
		// // tx.MustQueryContext(context, `SELECT * FROM poll_answers WHERE poll_id = $1`, t.Poll.ID).Each(func(rows *sx.Rows) {
		// // 	var pa model.PollAnswer
		// // 	rows.MustScans(&pa)
		// // 	t.PollWithAnswers.Answers = append(t.PollWithAnswers.Answers, pa)
		// // })
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

func (d *PostgreSQL) TopicByID(context context.Context, id int) (TopicWithPollAndCommentsCount, error) {
	var topic TopicWithPollAndCommentsCount

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT * FROM topics WHERE id = $1`, id).MustScans(&topic.Topic)
		tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM comments WHERE topic_id = $1`, topic.Topic.ID).MustScan(&topic.CommentsCount)

		// TODO(Pedro): fix this
		// var p model.Poll
		// tx.MustQueryRowContext(context, `SELECT * FROM polls WHERE topic_id = $1`, topic.Topic.ID).MustScans(&p)
		// if p.ID != 0 {
		// 	topic.Poll = new(PollWithAnswers)
		// 	topic.Poll.Poll = p
		// }
		//
		// tx.MustQueryContext(context, `SELECT * FROM poll_answers WHERE poll_id = $1`, topic.Poll.ID).Each(func(rows *sx.Rows) {
		// 	var pa model.PollAnswer
		// 	rows.MustScans(&pa)
		// 	topic.Poll.Answers = append(topic.Poll.Answers, pa)
		// })
	})

	return topic, err
}

func (d *PostgreSQL) CreatedAtByTopic(context context.Context, id int) (int, error) {
	var date int
	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT created_at FROM topics WHERE id = $1`, id).MustScan(&date)
	})

	return date, err
}

func (d *PostgreSQL) CommentsByTopicID(context context.Context, id int, after int, limit int) ([]CommentWithProfileAndAttachment, error) {
	var comments []CommentWithProfileAndAttachment

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		queryComment := `SELECT * FROM comments WHERE topic_id = $1 AND date >= $2 ORDER BY date ASC LIMIT $3`
		tx.MustQueryContext(context, queryComment, id, after, limit).Each(func(r *sx.Rows) {
			var c model.Comment
			r.MustScans(&c)

			comments = append(comments, CommentWithProfileAndAttachment{
				Comment: c,
			})
		})

		if len(comments) == 0 {
			return
		}

		var commentsIDs []string
		for i := range comments {
			commentsIDs = append(commentsIDs, strconv.Itoa(comments[i].Comment.ID))
		}

		queryProfile := `SELECT * FROM profiles WHERE id IN (` + strings.Join(commentsIDs, ",") + `)`
		tx.MustQueryContext(context, queryProfile).Each(func(rows *sx.Rows) {
			var p model.Profile
			rows.MustScans(&p)

			for i := range comments {
				if comments[i].FromID == p.ID {
					comments[i].Profile = p
				}
			}
		})

		queryAttachments := `SELECT * FROM attachments WHERE comment_id IN (` + strings.Join(commentsIDs, ",") + `)`
		tx.MustQueryContext(context, queryAttachments).Each(func(rows *sx.Rows) {
			var a model.Attachment
			rows.MustScans(&a)

			for i := range comments {
				if comments[i].ID == a.CommentID {
					comments[i].Attachments = append(comments[i].Attachments, a)
				}
			}
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

func (d *PostgreSQL) CommentsPaginationTimestampByTopicID(context context.Context, id int, after int, limit int) (PaginationTimestamps, error) {
	var timestamps PaginationTimestamps

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT 
					COALESCE((SELECT date FROM comments WHERE topic_id = $1 AND date >= $2 ORDER BY date ASC OFFSET $3 LIMIT 1), 0) as next,
					COALESCE((SELECT date FROM comments WHERE topic_id = $1 AND date <= $2 ORDER BY date DESC OFFSET $3 LIMIT 1), 0) as prev
				`
		tx.MustQueryRowContext(context, query, id, after, limit).MustScan(&timestamps.Next, &timestamps.Prev)
	})

	return timestamps, err
}

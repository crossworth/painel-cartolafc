package database

import (
	"context"
	fmt "fmt"
	"strconv"
	"strings"

	"github.com/travelaudience/go-sx"

	"github.com/crossworth/cartola-web-admin/model"
)

func (d *PostgreSQL) Topics(context context.Context, before int, limit int, orderBy OrderBy) ([]TopicWithPollAndCommentsCount, error) {
	var topics []TopicWithPollAndCommentsCount

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		var t TopicWithPollAndCommentsCount
		tx.MustQueryContext(context, `SELECT * FROM topics WHERE `+orderBy.Stringer()+` <= $1 ORDER BY `+orderBy.Stringer()+` DESC LIMIT $2`, before, limit).Each(func(rows *sx.Rows) {
			rows.MustScans(&t.Topic)
			topics = append(topics, t)
		})

		var topicIDs []string

		for i := range topics {
			topicIDs = append(topicIDs, strconv.Itoa(topics[i].ID))
			tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM comments WHERE topic_id = $1`, topics[i].ID).MustScan(&topics[i].CommentsCount)
		}

		if len(topicIDs) == 0 {
			return
		}

		var pollsIDs []string

		pollsQuery := `SELECT * FROM polls WHERE topic_id IN (` + strings.Join(topicIDs, ",") + `)`
		tx.MustQueryContext(context, pollsQuery).Each(func(rows *sx.Rows) {
			var p model.Poll
			rows.MustScans(&p)

			pollsIDs = append(pollsIDs, strconv.Itoa(p.ID))

			for i := range topics {
				if topics[i].ID == p.TopicID {
					topics[i].Poll = &PollWithAnswers{
						Poll: p,
					}
				}
			}
		})

		if len(pollsIDs) == 0 {
			return
		}

		pollAnswersQuery := `SELECT * FROM poll_answers WHERE poll_id IN (` + strings.Join(pollsIDs, ",") + `)`
		tx.MustQueryContext(context, pollAnswersQuery).Each(func(rows *sx.Rows) {
			var pa model.PollAnswer
			rows.MustScans(&pa)

			for i := range topics {
				if topics[i].Poll != nil && topics[i].Poll.ID == pa.PollID {
					topics[i].Poll.Answers = append(topics[i].Poll.Answers, pa)
				}
			}
		})
	})

	return topics, err
}

func (d *PostgreSQL) TopicsPaginationTimestamp(context context.Context, before int, limit int, orderBy OrderBy) (PaginationTimestamps, error) {
	var timestamps PaginationTimestamps

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT 
					COALESCE((SELECT ` + orderBy.Stringer() + ` FROM topics WHERE ` + orderBy.Stringer() + ` <= $1 ORDER BY ` + orderBy.Stringer() + ` DESC OFFSET $2 LIMIT 1), 0) as next,
					COALESCE((SELECT ` + orderBy.Stringer() + ` FROM topics WHERE ` + orderBy.Stringer() + ` >= $1 ORDER BY ` + orderBy.Stringer() + ` ASC OFFSET $2 LIMIT 1), 0) as prev
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
		tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM comments WHERE topic_id = $1`, id).MustScan(&topic.CommentsCount)

		pollsQuery := `SELECT * FROM polls WHERE topic_id = $1`
		tx.MustQueryContext(context, pollsQuery, id).Each(func(rows *sx.Rows) {
			var p model.Poll
			rows.MustScans(&p)

			topic.Poll = &PollWithAnswers{
				Poll: p,
			}
		})

		if topic.Poll == nil {
			return
		}

		pollAnswersQuery := `SELECT * FROM poll_answers WHERE poll_id = $1`
		tx.MustQueryContext(context, pollAnswersQuery, topic.Poll.ID).Each(func(rows *sx.Rows) {
			var pa model.PollAnswer
			rows.MustScans(&pa)
			topic.Poll.Answers = append(topic.Poll.Answers, pa)
		})
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

		var profileIDs []string
		for i := range comments {
			profileIDs = append(profileIDs, strconv.Itoa(comments[i].Comment.FromID))
		}

		queryProfile := `SELECT * FROM profiles WHERE id IN (` + strings.Join(profileIDs, ",") + `)`
		tx.MustQueryContext(context, queryProfile).Each(func(rows *sx.Rows) {
			var p model.Profile
			rows.MustScans(&p)

			for i := range comments {
				if comments[i].FromID == p.ID {
					comments[i].Profile = p
				}
			}
		})

		var commentsIDs []string
		for i := range comments {
			commentsIDs = append(commentsIDs, strconv.Itoa(comments[i].Comment.ID))
		}

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

func (d *PostgreSQL) TopicWithStats(context context.Context, orderBy string, orderDirection OrderByDirection, period Period, showOlderTopics bool, page int, limit int) ([]TopicsWithStats, error) {
	var topics []TopicsWithStats

	if orderBy != "comments" && orderBy != "likes" {
		return topics, ErrInvalidMemberOrderBy
	}

	periodTopics := ""
	periodComments := ""


	if period != PeriodAll {
		periodComments = fmt.Sprintf("WHERE date >= EXTRACT(epoch FROM date_trunc('%[1]s', current_date - INTERVAL '1 %[1]s')) AND date < EXTRACT(epoch FROM date_trunc('%[1]s', current_date))", period.Stringer())
	}

	if !showOlderTopics && period != PeriodAll {
		periodTopics = fmt.Sprintf("HAVING t.created_at >= EXTRACT(epoch FROM date_trunc('%[1]s', current_date - INTERVAL '1 %[1]s')) AND t.created_at < EXTRACT(epoch FROM date_trunc('%[1]s', current_date))", period.Stringer())
	}

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		query := `SELECT
    t.*,
    COALESCE(c.total, 0)::INTEGER as comments,
    COALESCE(l.total, 0)::INTEGER as likes
FROM topics t
    LEFT JOIN (
        SELECT topic_id, COUNT(topic_id) as total FROM comments ` + periodComments + ` GROUP BY topic_id
    ) as c ON c.topic_id = t.id
    LEFT JOIN (
        SELECT topic_id, SUM(likes) as total FROM comments ` + periodComments + ` GROUP BY topic_id
    ) as l ON l.topic_id = t.id
GROUP BY t.id, c.total, l.total ` + periodTopics + ` ORDER BY ` + orderBy + ` ` + orderDirection.Stringer() + ` OFFSET $1 LIMIT $2`

		i := 1
		tx.MustQueryContext(context, query, (page-1)*limit, limit).Each(func(r *sx.Rows) {
			var t TopicsWithStats
			r.MustScan(
				&t.ID,
				&t.Title,
				&t.IsClosed,
				&t.IsFixed,
				&t.CreatedAt,
				&t.UpdatedAt,
				&t.CreatedBy,
				&t.UpdatedBy,
				&t.Deleted,
				&t.Comments,
				&t.Likes,
			)
			t.Position = ((page - 1) * limit) + i
			topics = append(topics, t)
			i++
		})
	})

	return topics, err
}

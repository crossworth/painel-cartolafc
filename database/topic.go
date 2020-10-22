package database

import (
	"context"
	fmt "fmt"
	"strconv"
	"strings"

	"github.com/travelaudience/go-sx"

	"github.com/crossworth/painel-cartolafc/model"
)

func (p *PostgreSQL) Topics(context context.Context, before int, limit int, orderBy OrderBy) ([]TopicWithPollAndCommentsCount, error) {
	var topics []TopicWithPollAndCommentsCount

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		var t TopicWithPollAndCommentsCount
		tx.MustQueryContext(context, `SELECT t.*, (SELECT COUNT(*) FROM comments c WHERE c.topic_id = t.id) as comment_count FROM topics t WHERE t.`+orderBy.Stringer()+` <= $1 ORDER BY t.`+orderBy.Stringer()+` DESC LIMIT $2`, before, limit).Each(func(rows *sx.Rows) {
			rows.MustScan(
				&t.Topic.ID,
				&t.Topic.Title,
				&t.Topic.IsClosed,
				&t.Topic.IsFixed,
				&t.Topic.CreatedAt,
				&t.Topic.UpdatedAt,
				&t.Topic.CreatedBy,
				&t.Topic.UpdatedBy,
				&t.Topic.Deleted,
				&t.CommentsCount,
			)
			topics = append(topics, t)
		})

		var topicIDs []string

		for i := range topics {
			topicIDs = append(topicIDs, strconv.Itoa(topics[i].ID))
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

func (p *PostgreSQL) TopicsPaginationTimestamp(context context.Context, before int, limit int, orderBy OrderBy) (PaginationTimestamps, error) {
	var timestamps PaginationTimestamps

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		query := `SELECT 
					COALESCE((SELECT ` + orderBy.Stringer() + ` FROM topics WHERE ` + orderBy.Stringer() + ` <= $1 ORDER BY ` + orderBy.Stringer() + ` DESC OFFSET $2 LIMIT 1), 0) as next,
					COALESCE((SELECT ` + orderBy.Stringer() + ` FROM topics WHERE ` + orderBy.Stringer() + ` >= $1 ORDER BY ` + orderBy.Stringer() + ` ASC OFFSET $2 LIMIT 1), 0) as prev
				`
		tx.MustQueryRowContext(context, query, before, limit).MustScan(&timestamps.Next, &timestamps.Prev)
	})

	return timestamps, err
}

func (p *PostgreSQL) TopicsCount(context context.Context) (int, error) {
	var total int

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT reltuples::INTEGER FROM pg_class WHERE oid = 'topics'::regclass`).MustScan(&total)
	})

	return total, err
}

func (p *PostgreSQL) TopicByID(context context.Context, id int) (TopicWithPollAndCommentsCount, error) {
	var topic TopicWithPollAndCommentsCount

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT t.*, (SELECT COUNT(*) FROM comments c WHERE c.topic_id = t.id) as comment_count FROM topics t WHERE t.id = $1`, id).
			MustScan(
				&topic.ID,
				&topic.Title,
				&topic.IsClosed,
				&topic.IsFixed,
				&topic.CreatedAt,
				&topic.UpdatedAt,
				&topic.CreatedBy,
				&topic.UpdatedBy,
				&topic.Deleted,
				&topic.CommentsCount,
			)

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

func (p *PostgreSQL) CreatedAtByTopic(context context.Context, id int) (int, error) {
	var date int
	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT created_at FROM topics WHERE id = $1`, id).MustScan(&date)
	})

	return date, err
}

func (p *PostgreSQL) CommentsByTopicID(context context.Context, id int, after int, limit int) ([]CommentWithProfileAndAttachment, error) {
	var comments []CommentWithProfileAndAttachment

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
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

func (p *PostgreSQL) CommentsCountByTopicID(context context.Context, id int) (int, error) {
	var total int

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM comments WHERE topic_id = $1`, id).MustScan(&total)
	})

	return total, err
}

func (p *PostgreSQL) CommentsPaginationTimestampByTopicID(context context.Context, id int, after int, limit int) (PaginationTimestamps, error) {
	var timestamps PaginationTimestamps

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		query := `SELECT 
					COALESCE((SELECT date FROM comments WHERE topic_id = $1 AND date >= $2 ORDER BY date ASC OFFSET $3 LIMIT 1), 0) as next,
					COALESCE((SELECT date FROM comments WHERE topic_id = $1 AND date <= $2 ORDER BY date DESC OFFSET $3 LIMIT 1), 0) as prev
				`
		tx.MustQueryRowContext(context, query, id, after, limit).MustScan(&timestamps.Next, &timestamps.Prev)
	})

	return timestamps, err
}

func (p *PostgreSQL) TopicWithStats(context context.Context, orderBy string, orderDirection OrderByDirection, period Period, showOlderTopics bool, excludePseudoFixed bool, page int, limit int) ([]TopicsWithStats, error) {
	var topics []TopicsWithStats

	if orderBy != "comments" && orderBy != "likes" {
		return topics, ErrInvalidMemberOrderBy
	}

	periodTopics := ""
	periodComments := ""

	if period != PeriodAll {
		periodComments = fmt.Sprintf("WHERE date >= EXTRACT(epoch FROM current_date - INTERVAL '1 %s') AND date < EXTRACT(epoch FROM current_date)", period.Stringer())
	}

	hasHaving := false

	if !showOlderTopics && period != PeriodAll {
		hasHaving = true
		periodTopics = fmt.Sprintf("HAVING t.created_at >= EXTRACT(epoch FROM current_date - INTERVAL '1 %s') AND t.created_at < EXTRACT(epoch FROM current_date)", period.Stringer())
	}

	excludePseudoFixedComparison := `t.title NOT ILIKE 'FIXO%' AND t.title NOT ILIKE 'CARTOLA%' AND t.title NOT ILIKE '###%' AND t.title NOT ILIKE '%A/D/D%' AND t.title NOT ILIKE '%Antes/Durante/Depois%'`

	if excludePseudoFixed && !hasHaving {
		periodTopics = fmt.Sprintf("HAVING %s", excludePseudoFixedComparison)
	}

	if excludePseudoFixed && hasHaving {
		periodTopics = fmt.Sprintf("%s AND %s", periodTopics, excludePseudoFixedComparison)
	}

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
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

func (p *PostgreSQL) TopicWithStatsCount(context context.Context, period Period, showOlderTopics bool, excludePseudoFixed bool) (int, error) {
	var total int

	periodTopics := ""
	hasWhere := false

	if !showOlderTopics && period != PeriodAll {
		hasWhere = true
		periodTopics = fmt.Sprintf("WHERE t.created_at >= EXTRACT(epoch FROM current_date - INTERVAL '1 %s') AND t.created_at < EXTRACT(epoch FROM current_date)", period.Stringer())
	}

	excludePseudoFixedComparison := `t.title NOT ILIKE 'FIXO%' AND t.title NOT ILIKE 'CARTOLA%' AND t.title NOT ILIKE '###%' AND t.title NOT ILIKE '%A/D/D%' AND t.title NOT ILIKE '%Antes/Durante/Depois%'`

	if excludePseudoFixed && !hasWhere {
		periodTopics = fmt.Sprintf("WHERE %s", excludePseudoFixedComparison)
	}

	if excludePseudoFixed && hasWhere {
		periodTopics = fmt.Sprintf("%s AND %s", periodTopics, excludePseudoFixedComparison)
	}

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		query := `SELECT
    COUNT(*)
FROM topics t ` + periodTopics
		tx.MustQueryRowContext(context, query).MustScan(&total)
	})

	return total, err
}

func (p *PostgreSQL) TopicsWithMoreCommentsByID(context context.Context, id int, limit int) ([]model.TopicWithComments, error) {
	var result []model.TopicWithComments

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryContext(context, `SELECT t.*,
       COALESCE((SELECT COUNT(*) FROM "comments" c WHERE c.topic_id = t.id), 0) AS comments_count
FROM topics t
WHERE t.created_by = $1
ORDER BY comments_count DESC LIMIT $2`, id, limit).Each(func(rows *sx.Rows) {
			var t model.TopicWithComments
			rows.MustScan(
				&t.ID,
				&t.Title,
				&t.IsClosed,
				&t.IsFixed,
				&t.CreatedAt,
				&t.UpdatedAt,
				&t.CreatedBy,
				&t.UpdatedBy,
				&t.Deleted,
				&t.CommentsCount,
			)
			result = append(result, t)
		})
	})

	return result, err
}

func (p *PostgreSQL) TopicsWithMoreLikesByID(context context.Context, id int, limit int) ([]model.TopicWithLikes, error) {
	var result []model.TopicWithLikes

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryContext(context, `SELECT t.*,
       COALESCE((SELECT SUM(c.likes) FROM "comments" c WHERE c.topic_id = t.id), 0) AS likes_count
FROM topics t
WHERE t.created_by = $1
ORDER BY likes_count DESC LIMIT $2`, id, limit).Each(func(rows *sx.Rows) {
			var t model.TopicWithLikes
			rows.MustScan(
				&t.ID,
				&t.Title,
				&t.IsClosed,
				&t.IsFixed,
				&t.CreatedAt,
				&t.UpdatedAt,
				&t.CreatedBy,
				&t.UpdatedBy,
				&t.Deleted,
				&t.LikesCount,
			)
			result = append(result, t)
		})
	})

	return result, err
}

func (p *PostgreSQL) CommentsWithMoreLikes(context context.Context, id int, limit int) ([]model.Comment, error) {
	var result []model.Comment

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryContext(context, `SELECT *
FROM "comments" c
WHERE from_id = $1
ORDER BY likes DESC LIMIT $2`, id, limit).Each(func(rows *sx.Rows) {
			var c model.Comment
			rows.MustScans(&c)
			result = append(result, c)
		})
	})

	return result, err
}

func (p *PostgreSQL) TopicsNumberByIDGraph(context context.Context, id int) ([]GraphValue, error) {
	var results []GraphValue

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryContext(context, `SELECT TO_DATE(d.date, 'YYYY/MM/DD'),
       count(t.*)
FROM (
         SELECT TO_CHAR(DATE_TRUNC('day', (current_date - offs)), 'YYYY-MM-DD') AS date
         FROM generate_series(0, 365, 1) AS offs) d
         LEFT OUTER JOIN topics t ON
            d.date = TO_CHAR(DATE_TRUNC('day', TO_TIMESTAMP(t.created_at)), 'YYYY-MM-DD') AND
            t.created_by = $1
GROUP BY d.date
ORDER BY d.date;`, id).Each(func(rows *sx.Rows) {
			var t GraphValue
			rows.MustScans(&t)
			results = append(results, t)
		})
	})

	return results, err
}

func (p *PostgreSQL) CommentsNumberByIDGraph(context context.Context, id int) ([]GraphValue, error) {
	var results []GraphValue

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryContext(context, `SELECT TO_DATE(d.date, 'YYYY/MM/DD'),
       count(c.*)
FROM (
         SELECT TO_CHAR(DATE_TRUNC('day', (current_date - offs)), 'YYYY-MM-DD') AS date
         FROM generate_series(0, 365, 1) AS offs) d
         LEFT OUTER JOIN comments c ON
            d.date = TO_CHAR(DATE_TRUNC('day', TO_TIMESTAMP(c.date)), 'YYYY-MM-DD') AND
            c.from_id = $1
GROUP BY d.date
ORDER BY d.date;`, id).Each(func(rows *sx.Rows) {
			var t GraphValue
			rows.MustScans(&t)
			results = append(results, t)
		})
	})

	return results, err
}

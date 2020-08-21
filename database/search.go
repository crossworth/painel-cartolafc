package database

import (
	"context"

	"github.com/travelaudience/go-sx"
)

func (p *PostgreSQL) SearchTopics(context context.Context, term string, page int, limit int) ([]Search, error) {
	var result []Search

	query := `SELECT $1                                                                    as term,
       ts_headline((select title from topics t where t.id = ft.topic_id), q) as headline,
       ft.topic_id                                                           as topic_id,
       ft.date                                                               as date,
       ts_rank(tsv, q)                                                       as rank
FROM full_text_search_topic ft,
     plainto_tsquery($1) q
WHERE tsv @@ q
ORDER BY rank ASC, date DESC
OFFSET $2 LIMIT $3`

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		// NOTE(Pedro): to get the correct number of results we have to divide by two
		// since its two queries
		tx.MustQueryContext(context, query, term, (page-1)*limit, limit).Each(func(r *sx.Rows) {
			var search Search
			r.MustScan(
				&search.Term,
				&search.Headline,
				&search.TopicID,
				&search.Date,
				&search.Rank,
			)
			search.Type = SearchTypeTopic
			result = append(result, search)
		})
	})

	return result, err
}

func (p *PostgreSQL) SearchTopicsCount(context context.Context, term string) (int, error) {
	var total int

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*)
FROM full_text_search_topic ft,
     plainto_tsquery($1) q
WHERE tsv @@ q`, term).MustScan(&total)
	})

	return total, err
}

func (p *PostgreSQL) SearchComments(context context.Context, term string, page int, limit int) ([]Search, error) {
	var result []Search

	query := `SELECT $1                                                                       as term,
       ts_headline((select text from comments c where c.id = fc.comment_id), q) as headline,
       fc.topic_id                                                              as topic_id,
       fc.comment_id                                                            as comment_id,
       fc.date                                                                  as date,
       ts_rank(tsv, q)                                                          as rank
FROM full_text_search_comment fc,
     plainto_tsquery($1) q
WHERE tsv @@ q
ORDER BY rank ASC, date DESC
OFFSET $2 LIMIT $3`

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		// NOTE(Pedro): to get the correct number of results we have to divide by two
		// since its two queries
		tx.MustQueryContext(context, query, term, (page-1)*limit, limit).Each(func(r *sx.Rows) {
			var search Search
			r.MustScan(
				&search.Term,
				&search.Headline,
				&search.TopicID,
				&search.CommentID,
				&search.Date,
				&search.Rank,
			)
			search.Type = SearchTypeComment
			result = append(result, search)
		})
	})

	return result, err
}

func (p *PostgreSQL) SearchCommentsCount(context context.Context, term string) (int, error) {
	var total int

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*)
FROM full_text_search_comment ft,
     plainto_tsquery($1) q
WHERE tsv @@ q`, term).MustScan(&total)
	})

	return total, err
}

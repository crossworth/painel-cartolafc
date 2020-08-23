package database

import (
	"context"
	"strings"

	"github.com/travelaudience/go-sx"
)

func (p *PostgreSQL) SearchTopics(context context.Context, term string, page int, limit int, fullText bool) ([]Search, error) {
	if fullText {
		return p.SearchTopicsFullText(context, term, page, limit)
	} else {
		return p.SearchTopicsILike(context, term, page, limit)
	}
}

func (p *PostgreSQL) SearchTopicsFullText(context context.Context, term string, page int, limit int) ([]Search, error) {
	var result []Search

	query := `SELECT $1                                                                    as term,
       ts_headline((select title from topics t where t.id = ft.topic_id), q) as headline,
       ft.topic_id                                                           as topic_id,
       ft.date                                                               as date,
       ts_rank(tsv, q)                                                       as rank,
       (SELECT COUNT(c.*) FROM comments c WHERE c.topic_id = ft.topic_id)    as comments_count
FROM full_text_search_topic ft,
     plainto_tsquery($1) q
WHERE tsv @@ q
ORDER BY rank ASC, date DESC
OFFSET $2 LIMIT $3`

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryContext(context, query, term, (page-1)*limit, limit).Each(func(r *sx.Rows) {
			var search Search
			r.MustScan(
				&search.Term,
				&search.Headline,
				&search.TopicID,
				&search.Date,
				&search.Rank,
				&search.CommentsCount,
			)
			search.Type = SearchTypeTopic
			result = append(result, search)
		})
	})

	return result, err
}

func (p *PostgreSQL) SearchTopicsILike(context context.Context, term string, page int, limit int) ([]Search, error) {
	var result []Search

	query := `SELECT $1                                                          as term,
       t.title                                                     as headline,
       t.id                                                        as topic_id,
       t.created_at                                                as date,
       (SELECT COUNT(c.*) FROM comments c WHERE c.topic_id = t.id) as comments_count
FROM topics t
WHERE t.title ILIKE '%' || $1 || '%'
ORDER BY t.created_at DESC
OFFSET $2 LIMIT $3`
	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryContext(context, query, term, (page-1)*limit, limit).Each(func(r *sx.Rows) {
			var search Search
			r.MustScan(
				&search.Term,
				&search.Headline,
				&search.TopicID,
				&search.Date,
				&search.CommentsCount,
			)
			search.Type = SearchTypeTopic
			start := strings.Index(strings.ToLower(search.Headline), strings.ToLower(term))
			headline := search.Headline[0:start] + "<b>" + search.Headline[start:start+len(term)] + "</b>" + search.Headline[start+len(term):]
			search.Headline = headline
			result = append(result, search)
		})
	})

	return result, err
}

func (p *PostgreSQL) SearchTopicsCount(context context.Context, term string, fullText bool) (int, error) {
	if fullText {
		return p.SearchTopicsCountFullText(context, term)
	} else {
		return p.SearchTopicsCountILike(context, term)
	}
}

func (p *PostgreSQL) SearchTopicsCountFullText(context context.Context, term string) (int, error) {
	var total int

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*)
FROM full_text_search_topic ft,
     plainto_tsquery($1) q
WHERE tsv @@ q`, term).MustScan(&total)
	})

	return total, err
}

func (p *PostgreSQL) SearchTopicsCountILike(context context.Context, term string) (int, error) {
	var total int

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*)
FROM topics t
WHERE title ILIKE '%' || $1 || '%'`, term).MustScan(&total)
	})

	return total, err
}

func (p *PostgreSQL) SearchComments(context context.Context, term string, page int, limit int, fullText bool) ([]Search, error) {
	if fullText {
		return p.SearchCommentsFullText(context, term, page, limit)
	} else {
		return p.SearchCommentsILike(context, term, page, limit)
	}
}

func (p *PostgreSQL) SearchCommentsFullText(context context.Context, term string, page int, limit int) ([]Search, error) {
	var result []Search

	query := `SELECT $1                                                                                as term,
       ts_headline((SELECT noquote(text) FROM comments c WHERE c.id = fc.comment_id), q) as headline,
       fc.topic_id                                                                       as topic_id,
       fc.comment_id                                                                     as comment_id,
       fc.date                                                                           as date,
       ts_rank(tsv, q)                                                                   as rank,
       (SELECT c.likes FROM comments c WHERE c.id = fc.comment_id)                       as likes_count
FROM full_text_search_comment fc,
     plainto_tsquery($1) q
WHERE tsv @@ q
ORDER BY rank ASC, date DESC
OFFSET $2 LIMIT $3`

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryContext(context, query, term, (page-1)*limit, limit).Each(func(r *sx.Rows) {
			var search Search
			r.MustScan(
				&search.Term,
				&search.Headline,
				&search.TopicID,
				&search.CommentID,
				&search.Date,
				&search.Rank,
				&search.LikesCount,
			)
			search.Type = SearchTypeComment
			result = append(result, search)
		})
	})

	return result, err
}

func (p *PostgreSQL) SearchCommentsILike(context context.Context, term string, page int, limit int) ([]Search, error) {
	var result []Search

	query := `SELECT $1            as term,
       noquote(text) as headline,
       c.topic_id    as topic_id,
       c.id          as comment_id,
       c.date        as date,
       c.likes       as likes_count
FROM comments c
WHERE noquote(text) ILIKE '%' || $1 || '%'
ORDER BY date DESC
OFFSET $2 LIMIT $3`

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryContext(context, query, term, (page-1)*limit, limit).Each(func(r *sx.Rows) {
			var search Search
			r.MustScan(
				&search.Term,
				&search.Headline,
				&search.TopicID,
				&search.CommentID,
				&search.Date,
				&search.LikesCount,
			)
			search.Type = SearchTypeComment
			start := strings.Index(strings.ToLower(search.Headline), strings.ToLower(term))
			headline := search.Headline[0:start] + "<b>" + search.Headline[start:start+len(term)] + "</b>" + search.Headline[start+len(term):]
			search.Headline = headline
			result = append(result, search)
		})
	})

	return result, err
}

func (p *PostgreSQL) SearchCommentsCount(context context.Context, term string, fullText bool) (int, error) {
	if fullText {
		return p.SearchCommentsCountFullText(context, term)
	} else {
		return p.SearchCommentsCountILike(context, term)
	}
}

func (p *PostgreSQL) SearchCommentsCountFullText(context context.Context, term string) (int, error) {
	var total int

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*)
FROM full_text_search_comment ft,
     plainto_tsquery($1) q
WHERE tsv @@ q`, term).MustScan(&total)
	})

	return total, err
}

func (p *PostgreSQL) SearchCommentsCountILike(context context.Context, term string) (int, error) {
	var total int

	err := sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, `SELECT COUNT(*)
FROM comments c
WHERE
     noquote(text) ILIKE '%' || $1 || '%'
`, term).MustScan(&total)
	})

	return total, err
}

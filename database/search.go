package database

import (
	"context"
	"fmt"
	"strconv"
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

	query := `SELECT $1                                                                    AS term,
       ts_headline((SELECT title FROM topics t WHERE t.id = ft.topic_id), q) AS headline,
       ft.topic_id                                                           AS topic_id,
       ft.date                                                               AS date,
       ts_rank(tsv, q)                                                       AS rank,
       (SELECT COUNT(c.*) FROM comments c WHERE c.topic_id = ft.topic_id)    AS comments_count,
       (SELECT created_by FROM topics t WHERE t.id = ft.topic_id)            AS from_id
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
				&search.FromID,
			)
			search.Type = SearchTypeTopic
			result = append(result, search)
		})
	})

	return p.populateProfiles(context, result, err)
}

func (p *PostgreSQL) populateProfiles(context context.Context, input []Search, err error) ([]Search, error) {
	if err != nil {
		return input, err
	}

	if len(input) == 0 {
		return input, err
	}

	var ids []string

	for _, s := range input {
		ids = append(ids, strconv.Itoa(s.FromID))
	}

	type profileResult struct {
		ID         int
		Name       string
		ScreenName string
		Photo      string
	}

	var profiles []profileResult
	query := fmt.Sprintf(`SELECT p.id, p.first_name, p.screen_name, p.photo FROM profiles p WHERE id in (%s)`, strings.Join(ids, ", "))
	err = sx.DoContext(context, p.db, func(tx *sx.Tx) {
		tx.MustQueryContext(context, query).Each(func(r *sx.Rows) {
			var profile profileResult
			r.MustScan(
				&profile.ID,
				&profile.Name,
				&profile.ScreenName,
				&profile.Photo,
			)
			profiles = append(profiles, profile)
		})
	})

	for i := range input {
		for _, p := range profiles {
			if p.ID == input[i].FromID {
				input[i].FromName = p.Name
				input[i].FromPhoto = p.Photo
				input[i].FromScreenName = p.ScreenName
			}
		}
	}

	return input, err
}

func (p *PostgreSQL) SearchTopicsILike(context context.Context, term string, page int, limit int) ([]Search, error) {
	var result []Search

	query := `SELECT $1                                                          AS term,
       t.title                                                     AS headline,
       t.id                                                        AS topic_id,
       t.created_at                                                AS date,
       (SELECT COUNT(c.*) FROM comments c WHERE c.topic_id = t.id) AS comments_count,
	   t.created_by                                                AS from_id,
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
				&search.FromID,
			)
			search.Type = SearchTypeTopic
			start := strings.Index(strings.ToLower(search.Headline), strings.ToLower(term))
			headline := search.Headline[0:start] + "<b>" + search.Headline[start:start+len(term)] + "</b>" + search.Headline[start+len(term):]
			search.Headline = headline
			result = append(result, search)
		})
	})

	return p.populateProfiles(context, result, err)
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

	query := `SELECT $1                                                                                AS term,
       ts_headline((SELECT noquote(c.text) FROM comments c WHERE c.id = fc.comment_id), q) AS headline,
       fc.topic_id                                                                       AS topic_id,
       fc.comment_id                                                                     AS comment_id,
       fc.date                                                                           AS date,
       ts_rank(tsv, q)                                                                   AS rank,
       (SELECT c.likes FROM comments c WHERE c.id = fc.comment_id)                       AS likes_count,
       (SELECT c.from_id FROM comments c WHERE c.id = fc.comment_id)                     AS from_id
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
				&search.FromID,
			)
			search.Type = SearchTypeComment
			result = append(result, search)
		})
	})

	return p.populateProfiles(context, result, err)
}

func (p *PostgreSQL) SearchCommentsILike(context context.Context, term string, page int, limit int) ([]Search, error) {
	var result []Search

	query := `SELECT $1            AS term,
       noquote(text) AS headline,
       c.topic_id    AS topic_id,
       c.id          AS comment_id,
       c.date        AS date,
       c.likes       AS likes_count,
       c.from_id     AS from_id
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
				&search.FromID,
			)
			search.Type = SearchTypeComment
			start := strings.Index(strings.ToLower(search.Headline), strings.ToLower(term))
			headline := search.Headline[0:start] + "<b>" + search.Headline[start:start+len(term)] + "</b>" + search.Headline[start+len(term):]
			search.Headline = headline
			result = append(result, search)
		})
	})

	return p.populateProfiles(context, result, err)
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

package database

import (
	"bytes"
	"context"
	"strconv"
	"strings"

	"github.com/travelaudience/go-sx"
)

func (d *PostgreSQL) Search(context context.Context, text string, page int, limit int, fromID int, createdAfter int, createdBefore int) ([]Search, error) {
	var result []Search

	query := bytes.Buffer{}

	query.WriteString(`(SELECT title AS text, created_at AS date, created_by AS from_id, 'topic' AS type FROM topics t WHERE t.text LIKE '%' || $1 || '%'`)

	if fromID != 0 {
		query.WriteString(` AND t.from_id = ` + strconv.Itoa(fromID))
	}

	if createdAfter != 0 {
		query.WriteString(` AND t.date >= ` + strconv.Itoa(createdAfter))
	}

	if createdBefore != 0 {
		query.WriteString(` AND t.date <= ` + strconv.Itoa(createdBefore) + ` `)
	}

	query.WriteString(`OFFSET $2 LIMIT $3)`)
	query.WriteString(`UNION`)
	query.WriteString(`(SELECT text AS text, date AS date, from_id AS from_id, 'comment' AS type FROM comments c WHERE c.text LIKE '%' || $1 || '%'`)

	if fromID != 0 {
		query.WriteString(` AND c.from_id = ` + strconv.Itoa(fromID))
	}

	if createdAfter != 0 {
		query.WriteString(` AND c.date >= ` + strconv.Itoa(createdAfter))
	}

	if createdBefore != 0 {
		query.WriteString(` AND c.date <= ` + strconv.Itoa(createdBefore) + ` `)
	}

	query.WriteString(`OFFSET $2 LIMIT $3)`)

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		// NOTE(Pedro): to get the correct number of results we have to divide by two
		// since its two queries
		tx.MustQueryContext(context, query.String(), text, (page-1)*limit/2, limit/2).Each(func(r *sx.Rows) {
			var search Search
			var typeString string
			search.Term = text
			r.MustScan(&search.Text, &search.Date, &search.FromID, &typeString)

			search.Type = SearchTypeTopic

			if typeString == "comment" {
				search.Type = SearchTypeComment
			}

			search.HighlightedPart = strings.ReplaceAll(search.Text, text, "<strong>"+text+"</strong>")
			result = append(result, search)
		})
	})

	return result, err
}

func (d *PostgreSQL) SearchCount(context context.Context, text string, page int, limit int, fromID int, createdAfter int, createdBefore int) (int, error) {
	var total int

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		// TODO(Pedro): Fix this
		// tx.MustQueryRowContext(context, `SELECT COUNT(*) FROM comments WHERE topic_id = $1`, id).MustScan(&total)
	})

	return total, err
}

func (d *PostgreSQL) SearchPaginationTimestamp(context context.Context, text string, page int, limit int, fromID int, createdAfter int, createdBefore int) (PaginationTimestamps, error) {
	var timestamps PaginationTimestamps

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		// TODO(Pedro): Fix this queries
		// query := `SELECT
		// 			COALESCE((SELECT date FROM comments WHERE topic_id = $1 AND date >= $2 ORDER BY date ASC OFFSET $3 LIMIT 1), 0) as next,
		// 			COALESCE((SELECT date FROM comments WHERE topic_id = $1 AND date <= $2 ORDER BY date DESC OFFSET $3 LIMIT 1), 0) as prev
		// 		`
		// tx.MustQueryRowContext(context, query, id, after, limit).MustScan(&timestamps.Next, &timestamps.Prev)
	})

	return timestamps, err
}

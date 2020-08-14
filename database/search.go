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

	searchTerm := `'% ` + text +` %'`

	query := bytes.Buffer{}

	query.WriteString(`(SELECT title AS text, created_at AS date, created_by AS from_id, 'topic' AS type, id, 0 FROM topics WHERE title LIKE ` + searchTerm + ``)

	if fromID != 0 {
		query.WriteString(` AND created_by = ` + strconv.Itoa(fromID))
	}

	if createdAfter != 0 {
		query.WriteString(` AND created_at >= ` + strconv.Itoa(createdAfter))
	}

	if createdBefore != 0 {
		query.WriteString(` AND created_at <= ` + strconv.Itoa(createdBefore) + ` `)
	}

	query.WriteString(`ORDER BY title OFFSET $1 LIMIT $2)`)
	query.WriteString(`UNION`)
	query.WriteString(`(SELECT text AS text, date AS date, from_id AS from_id, 'comment' AS type, topic_id, id FROM comments WHERE text LIKE ` + searchTerm + ``)

	if fromID != 0 {
		query.WriteString(` AND from_id = ` + strconv.Itoa(fromID))
	}

	if createdAfter != 0 {
		query.WriteString(` AND date >= ` + strconv.Itoa(createdAfter))
	}

	if createdBefore != 0 {
		query.WriteString(` AND date <= ` + strconv.Itoa(createdBefore) + ` `)
	}

	query.WriteString(`ORDER BY text OFFSET $1 LIMIT $2)`)

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		// NOTE(Pedro): to get the correct number of results we have to divide by two
		// since its two queries
		tx.MustQueryContext(context, query.String(), (page-1)*limit/2, limit/2).Each(func(r *sx.Rows) {
			var search Search
			var typeString string
			search.Term = text
			r.MustScan(&search.Text, &search.Date, &search.FromID, &typeString, &search.TopicID, &search.CommentID)

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

	searchTerm := `'% ` + text +` %'`

	query := bytes.Buffer{}

	query.WriteString(`SELECT COUNT(*) FROM ((SELECT title AS text, created_at AS date, created_by AS from_id, 'topic' AS type, id, 0 FROM topics WHERE title LIKE ` + searchTerm + ``)

	if fromID != 0 {
		query.WriteString(` AND created_by = ` + strconv.Itoa(fromID))
	}

	if createdAfter != 0 {
		query.WriteString(` AND created_at >= ` + strconv.Itoa(createdAfter))
	}

	if createdBefore != 0 {
		query.WriteString(` AND created_at <= ` + strconv.Itoa(createdBefore) + ` `)
	}

	query.WriteString(`ORDER BY title)`)
	query.WriteString(`UNION`)
	query.WriteString(`(SELECT text AS text, date AS date, from_id AS from_id, 'comment' AS type, topic_id, id FROM comments WHERE text LIKE ` + searchTerm + ``)

	if fromID != 0 {
		query.WriteString(` AND from_id = ` + strconv.Itoa(fromID))
	}

	if createdAfter != 0 {
		query.WriteString(` AND date >= ` + strconv.Itoa(createdAfter))
	}

	if createdBefore != 0 {
		query.WriteString(` AND date <= ` + strconv.Itoa(createdBefore) + ` `)
	}

	query.WriteString(`ORDER BY text)) q`)

	err := sx.DoContext(context, d.db, func(tx *sx.Tx) {
		tx.MustQueryRowContext(context, query.String()).MustScan(&total)
	})

	return total, err
}

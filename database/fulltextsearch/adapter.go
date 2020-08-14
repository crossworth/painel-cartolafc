package fulltextsearch

import (
	"bytes"
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/travelaudience/go-sx"

	"github.com/crossworth/cartola-web-admin/database"
)

// https://www.infoq.com/br/articles/postgresql-fts/
// https://pt.slideshare.net/spjuliano/fts-26392077 pg 28
// https://medium.com/@gabrielfgularte/caracteres-especiais-e-full-text-search-no-postgre-c345b1da5b7b

type Adapter struct {
	db *sql.DB
}

func NewAdapter(db *sql.DB) *Adapter {
	return &Adapter{
		db: db,
	}
}

func (d *Adapter) Search(context context.Context, text string, page int, limit int, fromID int, createdAfter int, createdBefore int) ([]database.Search, error) {
	var result []database.Search

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
			var search database.Search
			var typeString string
			search.Term = text
			r.MustScan(&search.Text, &search.Date, &search.FromID, &typeString, &search.TopicID, &search.CommentID)

			search.Type = database.SearchTypeTopic

			if typeString == "comment" {
				search.Type = database.SearchTypeComment
			}

			search.HighlightedPart = strings.ReplaceAll(search.Text, text, "<strong>"+text+"</strong>")
			result = append(result, search)
		})
	})

	return result, err
}

func (d *Adapter) SearchCount(context context.Context, text string, page int, limit int, fromID int, createdAfter int, createdBefore int) (int, error) {
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

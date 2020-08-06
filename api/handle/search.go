package handle

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/util"
)

type SearchProvider interface {
	Search(context context.Context, text string, matchPartial bool, page int, limit int, fromID int, createdAfter int, createdBefore int) ([]database.Search, error)
	SearchCount(context context.Context, text string, matchPartial bool, page int, limit int, fromID int, createdAfter int, createdBefore int) (int, error)
}

func Search(provider SearchProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		term := util.StringWithDefault(r.URL.Query().Get("term"), "")
		matchPartial := util.BoolWithDefault(r.URL.Query().Get("includePartials"), false)
		fromID := util.ToIntWithDefaultMin(r.URL.Query().Get("fromID"), 0)
		createdAfter := util.ToIntWithDefaultMin(r.URL.Query().Get("createdAfter"), 0)
		createdBefore := util.ToIntWithDefaultMin(r.URL.Query().Get("createdBefore"), 0)
		page := util.ToIntWithDefaultMin(r.URL.Query().Get("page"), 1)
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)

		if term == "" {
			errorCode(w, NewError("nenhum termo fornecido"), 400)
			return
		}

		results, err := provider.Search(r.Context(), term, matchPartial, page, limit, fromID, createdAfter, createdBefore)
		if err != nil {
			databaseError(w, err)
			return
		}

		total, err := provider.SearchCount(r.Context(), term, matchPartial, page, limit, fromID, createdAfter, createdBefore)
		if err != nil {
			databaseError(w, err)
			return
		}

		next := ""
		prev := ""

		if page != 1 {
			prev = fmt.Sprintf("%s/search?limit=%d&page=%d&term=%s&includePartials=%t&fromID=%d&createdAfter=%d&createdBefore=%d", os.Getenv("APP_API_URL"), limit, page-1, term, matchPartial, fromID, createdAfter, createdBefore)
		}

		if page*limit < total {
			next = fmt.Sprintf("%s/search?limit=%d&page=%d&term=%s&includePartials=%t&fromID=%d&createdAfter=%d&createdBefore=%d", os.Getenv("APP_API_URL"), limit, page+1, term, matchPartial, fromID, createdAfter, createdBefore)
		}

		pagination(w, results, 200, PaginationMeta{
			Prev:    prev,
			Current: fmt.Sprintf("%s/search?limit=%d&page=%d&term=%s&includePartials=%t&fromID=%d&createdAfter=%d&createdBefore=%d", os.Getenv("APP_API_URL"), limit, page,  term, matchPartial, fromID, createdAfter, createdBefore),
			Next:    next,
			Total:   total,
		})
	}
}

package handle

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/crossworth/cartola-web-admin/cache"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/util"
)

type SearchProvider interface {
	Search(context context.Context, text string, page int, limit int, fromID int, createdAfter int, createdBefore int) ([]database.Search, error)
	SearchCount(context context.Context, text string, page int, limit int, fromID int, createdAfter int, createdBefore int) (int, error)
}

type SearchCache struct {
	Results   []database.Search
	Total     int
	CreatedAt time.Time
}

func Search(provider SearchProvider, cache *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		term := util.StringWithDefault(r.URL.Query().Get("term"), "")
		fromID := util.ToIntWithDefaultMin(r.URL.Query().Get("fromID"), 0)
		createdAfter := util.ToIntWithDefaultMin(r.URL.Query().Get("createdAfter"), 0)
		createdBefore := util.ToIntWithDefaultMin(r.URL.Query().Get("createdBefore"), 0)
		page := util.ToIntWithDefaultMin(r.URL.Query().Get("page"), 1)
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)

		if term == "" {
			errorCode(w, NewError("nenhum termo fornecido"), 400)
			return
		}

		cacheKey := fmt.Sprintf("search_%s_%d_%d_%d_%d_%d", term, fromID, createdAfter, createdBefore, page, limit)
		searchCache := cache.GetShortCache(cacheKey, func() interface{} {
			results, err := provider.Search(r.Context(), term, page, limit, fromID, createdAfter, createdBefore)
			if err != nil {
				return err
			}

			total, err := provider.SearchCount(r.Context(), term, page, limit, fromID, createdAfter, createdBefore)
			if err != nil {
				return err
			}

			searchCache := SearchCache{
				Results:   results,
				Total:     total,
				CreatedAt: time.Now(),
			}

			return searchCache
		})

		data, castOK := searchCache.(SearchCache)
		if !castOK {
			databaseError(w, searchCache.(error))
			return
		}

		next := ""
		prev := ""

		if page != 1 {
			prev = fmt.Sprintf("%s/search?limit=%d&page=%d&term=%s&fromID=%d&createdAfter=%d&createdBefore=%d", os.Getenv("APP_API_URL"), limit, page-1, term, fromID, createdAfter, createdBefore)
		}

		if page*limit < data.Total {
			next = fmt.Sprintf("%s/search?limit=%d&page=%d&term=%s&fromID=%d&createdAfter=%d&createdBefore=%d", os.Getenv("APP_API_URL"), limit, page+1, term, fromID, createdAfter, createdBefore)
		}

		pagination(w, data.Results, 200, PaginationMeta{
			Prev:    prev,
			Current: fmt.Sprintf("%s/search?limit=%d&page=%d&term=%s&fromID=%d&createdAfter=%d&createdBefore=%d", os.Getenv("APP_API_URL"), limit, page, term, fromID, createdAfter, createdBefore),
			Next:    next,
			Total:   data.Total,
		})
	}
}

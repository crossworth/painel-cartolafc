package handle

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/crossworth/painel-cartolafc/cache"
	"github.com/crossworth/painel-cartolafc/database"
	"github.com/crossworth/painel-cartolafc/httputil"
	"github.com/crossworth/painel-cartolafc/util"
)

type SearchProvider interface {
	SearchTopics(context context.Context, term string, page int, limit int, fullText bool) ([]database.Search, error)
	SearchTopicsCount(context context.Context, term string, fullText bool) (int, error)
	SearchComments(context context.Context, term string, page int, limit int, fullText bool) ([]database.Search, error)
	SearchCommentsCount(context context.Context, term string, fullText bool) (int, error)
}

type SearchCache struct {
	Results   []database.Search `json:"results"`
	Total     int               `json:"total"`
	CreatedAt time.Time         `json:"created_at"`
}

func Search(provider SearchProvider, cache *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		term := util.StringWithDefault(r.URL.Query().Get("term"), "")
		searchType := util.StringWithDefault(r.URL.Query().Get("searchType"), "title")
		fullText := util.BoolWithDefault(r.URL.Query().Get("fullText"), true)
		page := util.ToIntWithDefaultMin(r.URL.Query().Get("page"), 1)
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)

		if term == "" {
			httputil.SendErrorCode(w, httputil.NewError("nenhum termo fornecido"), 400)
			return
		}

		if searchType != "title" && searchType != "text" {
			searchType = "title"
		}

		cacheKey := fmt.Sprintf("search_%s_%s_%t_%d_%d", term, searchType, fullText, page, limit)
		searchCache := cache.GetShortCache(cacheKey, func() interface{} {
			var results []database.Search
			var err error

			if searchType == "title" {
				results, err = provider.SearchTopics(r.Context(), term, page, limit, fullText)
				if err != nil {
					return err
				}
			} else {
				results, err = provider.SearchComments(r.Context(), term, page, limit, fullText)
				if err != nil {
					return err
				}
			}

			var total int
			if searchType == "title" {
				total, err = provider.SearchTopicsCount(r.Context(), term, fullText)
				if err != nil {
					return err
				}
			} else {
				total, err = provider.SearchCommentsCount(r.Context(), term, fullText)
				if err != nil {
					return err
				}
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
			httputil.SendDatabaseError(w, searchCache.(error))
			return
		}

		next := ""
		prev := ""

		if page != 1 {
			prev = fmt.Sprintf("%s/search?limit=%d&page=%d&term=%s&searchType=%s", os.Getenv("APP_API_URL"), limit, page-1, term, searchType)
		}

		if page*limit < data.Total {
			next = fmt.Sprintf("%s/search?limit=%d&page=%d&term=%s&searchType=%s", os.Getenv("APP_API_URL"), limit, page+1, term, searchType)
		}

		httputil.SendPagination(w, data.Results, 200, httputil.PaginationMeta{
			Prev:    prev,
			Current: fmt.Sprintf("%s/search?limit=%d&page=%d&term=%s&searchType=%s", os.Getenv("APP_API_URL"), limit, page, term, searchType),
			Next:    next,
			Total:   data.Total,
		})
	}
}

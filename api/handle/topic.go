package handle

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/util"
)

type TopicsProvider interface {
	Topics(context context.Context, before int, limit int) ([]database.TopicWithPollAndCommentsCount, error)
	TopicsPaginationTimestamp(context context.Context, before int, limit int) (database.PaginationTimestamps, error)
	TopicsCount(context context.Context) (int, error)
}

func Topics(provider TopicsProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		before := util.ToIntWithDefault(r.URL.Query().Get("before"), int(time.Now().Unix()))
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)

		total, err := provider.TopicsCount(r.Context())
		if err != nil {
			databaseError(w, err)
			return
		}

		topics, err := provider.Topics(r.Context(), before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		paginationTimestamps, err := provider.TopicsPaginationTimestamp(r.Context(), before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		next := ""
		prev := ""

		if paginationTimestamps.Next != 0 {
			next = fmt.Sprintf("%s/topics?limit=%d&before=%d", os.Getenv("APP_API_URL"), limit, paginationTimestamps.Next)
		}
		if paginationTimestamps.Prev != 0 {
			prev = fmt.Sprintf("%s/topics?limit=%d&before=%d", os.Getenv("APP_API_URL"), limit, paginationTimestamps.Prev)
		}

		if len(topics) != 0 {
			before = topics[0].CreatedAt
		}

		pagination(w, topics, 200, PaginationMeta{
			Prev:    prev,
			Current: fmt.Sprintf("%s/topics?limit=%d&before=%d", os.Getenv("APP_API_URL"), limit, before),
			Next:    next,
			Total:   total,
		})
	}
}

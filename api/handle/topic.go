package handle

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"

	"github.com/crossworth/painel-cartolafc/cache"
	"github.com/crossworth/painel-cartolafc/database"
	"github.com/crossworth/painel-cartolafc/httputil"
	"github.com/crossworth/painel-cartolafc/util"
)

type TopicsProvider interface {
	Topics(context context.Context, before int, limit int, orderBy database.OrderBy) ([]database.TopicWithPollAndCommentsCount, error)
	TopicsPaginationTimestamp(context context.Context, before int, limit int, orderBy database.OrderBy) (database.PaginationTimestamps, error)
	TopicsCount(context context.Context) (int, error)
}

func Topics(provider TopicsProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		before := util.ToIntWithDefault(r.URL.Query().Get("before"), int(time.Now().Unix()))
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)
		orderByStr := util.StringWithDefault(r.URL.Query().Get("orderBy"), "updated_at")

		total, err := provider.TopicsCount(r.Context())
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		orderBy := database.OrderByUpdatedAt

		if strings.ToLower(orderByStr) == "created_at" {
			orderBy = database.OrderByCreatedAt
		}

		topics, err := provider.Topics(r.Context(), before, limit, orderBy)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		paginationTimestamps, err := provider.TopicsPaginationTimestamp(r.Context(), before, limit, orderBy)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		next := ""
		prev := ""

		if paginationTimestamps.Next != 0 {
			next = fmt.Sprintf("%s/topics?limit=%d&before=%d&orderBy=%s", os.Getenv("APP_API_URL"), limit, paginationTimestamps.Next, orderBy.Stringer())
		}
		if paginationTimestamps.Prev != 0 {
			prev = fmt.Sprintf("%s/topics?limit=%d&before=%d&orderBy=%s", os.Getenv("APP_API_URL"), limit, paginationTimestamps.Prev, orderBy.Stringer())
		}

		if len(topics) != 0 {
			before = topics[0].CreatedAt
		}

		httputil.SendPagination(w, topics, 200, httputil.PaginationMeta{
			Prev:    prev,
			Current: fmt.Sprintf("%s/topics?limit=%d&before=%d&orderBy=%s", os.Getenv("APP_API_URL"), limit, before, orderBy.Stringer()),
			Next:    next,
			Total:   total,
		})
	}
}

type TopicProvider interface {
	TopicByID(context context.Context, id int) (database.TopicWithPollAndCommentsCount, error)
}

func TopicByID(provider TopicProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "topic"))

		topic, err := provider.TopicByID(r.Context(), id)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		httputil.SendJSON(w, topic, 200)
	}
}

type TopicCommentsProvider interface {
	CreatedAtByTopic(context context.Context, id int) (int, error)
	CommentsByTopicID(context context.Context, id int, after int, limit int) ([]database.CommentWithProfileAndAttachment, error)
	CommentsCountByTopicID(context context.Context, id int) (int, error)
	CommentsPaginationTimestampByTopicID(context context.Context, id int, after int, limit int) (database.PaginationTimestamps, error)
}

func CommentFromTopicByID(provider TopicCommentsProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "topic"))
		after := util.ToIntWithDefault(r.URL.Query().Get("after"), 0)
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)

		total, err := provider.CommentsCountByTopicID(r.Context(), id)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		if after == 0 {
			createdAt, err := provider.CreatedAtByTopic(r.Context(), id)
			if err != nil {
				httputil.SendDatabaseError(w, err)
				return
			}

			after = createdAt
		}

		comments, err := provider.CommentsByTopicID(r.Context(), id, after, limit)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		paginationTimestamps, err := provider.CommentsPaginationTimestampByTopicID(r.Context(), id, after, limit)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		next := ""
		prev := ""

		if paginationTimestamps.Next != 0 {
			next = fmt.Sprintf("%s/topics/%d/comments?limit=%d&after=%d", os.Getenv("APP_API_URL"), id, limit, paginationTimestamps.Next)
		}
		if paginationTimestamps.Prev != 0 {
			prev = fmt.Sprintf("%s/topics/%d/comments?limit=%d&after=%d", os.Getenv("APP_API_URL"), id, limit, paginationTimestamps.Prev)
		}

		if len(comments) != 0 {
			after = comments[0].Date
		}

		httputil.SendPagination(w, comments, 200, httputil.PaginationMeta{
			Prev:    prev,
			Current: fmt.Sprintf("%s/topics/%d/comments?limit=%d&after=%d", os.Getenv("APP_API_URL"), id, limit, after),
			Next:    next,
			Total:   total,
		})
	}
}

type TopicsWithStatsProvider interface {
	TopicWithStats(context context.Context, orderBy string, orderDirection database.OrderByDirection, period database.Period, showOlderTopics bool, excludePseudoFixed bool, page int, limit int) ([]database.TopicsWithStats, error)
	TopicWithStatsCount(context context.Context, period database.Period, showOlderTopics bool, excludePseudoFixed bool) (int, error)
}

type TopicsWithStatsCache struct {
	Topics    []database.TopicsWithStats `json:"topics"`
	Total     int                        `json:"total"`
	CreatedAt time.Time                  `json:"created_at"`
}

func TopicsWithStats(provider TopicsWithStatsProvider, cache *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orderBy := util.StringWithDefault(r.URL.Query().Get("orderBy"), "comments")
		orderDirStr := util.StringWithDefault(r.URL.Query().Get("orderDir"), "DESC")
		periodStr := util.StringWithDefault(r.URL.Query().Get("period"), "all")
		page := util.ToIntWithDefaultMin(r.URL.Query().Get("page"), 1)
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)
		showOlderTopics := util.BoolWithDefault(r.URL.Query().Get("showOlderTopics"), true)
		excludePseudoFixed := util.BoolWithDefault(r.URL.Query().Get("excludePseudoFixed"), false)

		orderDir := database.OrderByASC

		if strings.ToLower(orderDirStr) == "desc" {
			orderDir = database.OrderByDESC
		}

		period := database.PeriodFromString(periodStr)

		// this avoids creating a duplicated cache record
		if period == database.PeriodAll {
			showOlderTopics = true
		}

		cacheKey := fmt.Sprintf("topics_with_stats_%s_%s_%d_%d_%s_%t_%t", orderBy, orderDir.Stringer(), page, limit, period.Stringer(), showOlderTopics, excludePseudoFixed)
		topicsCache := cache.Get(cacheKey, func() interface{} {
			topics, err := provider.TopicWithStats(context.TODO(), orderBy, orderDir, period, showOlderTopics, excludePseudoFixed, page, limit)
			if err != nil {
				return err
			}

			total, err := provider.TopicWithStatsCount(context.TODO(), period, showOlderTopics, excludePseudoFixed)
			if err != nil {
				return err
			}

			topicsCache := TopicsWithStatsCache{
				Topics:    topics,
				Total:     total,
				CreatedAt: time.Now(),
			}

			return topicsCache
		})

		data, castOK := topicsCache.(TopicsWithStatsCache)
		if !castOK {
			httputil.SendDatabaseError(w, topicsCache.(error))
			return
		}

		next := ""
		prev := ""

		if page != 1 {
			prev = fmt.Sprintf("%s/topics-ranking?limit=%d&page=%d&orderBy=%s&orderDir=%s&period=%s&showOlderTopics=%t", os.Getenv("APP_API_URL"), limit, page-1, orderBy, orderDir.Stringer(), period.URLString(), showOlderTopics)
		}

		if page*limit < data.Total {
			next = fmt.Sprintf("%s/topics-ranking?limit=%d&page=%d&orderBy=%s&orderDir=%s&period=%s&showOlderTopics=%t", os.Getenv("APP_API_URL"), limit, page+1, orderBy, orderDir.Stringer(), period.URLString(), showOlderTopics)
		}

		// NOTE(Pedro): To calculate the correct position for the records
		// when in asc order
		if orderDir == database.OrderByASC {
			for i := range data.Topics {
				data.Topics[i].Position = (data.Total + 1) - data.Topics[i].Position
			}
		}

		httputil.SendPagination(w, data.Topics, 200, httputil.PaginationMeta{
			Prev:     prev,
			Current:  fmt.Sprintf("%s/topics-ranking?limit=%d&page=%d&orderBy=%s&orderDir=%s&period=%s&showOlderTopics=%t", os.Getenv("APP_API_URL"), limit, page, orderBy, orderDir.Stringer(), period.URLString(), showOlderTopics),
			Next:     next,
			Total:    data.Total,
			CachedAt: &data.CreatedAt,
		})
	}
}

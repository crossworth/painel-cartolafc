package handle

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"

	"github.com/crossworth/cartola-web-admin/cache"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/model"
	"github.com/crossworth/cartola-web-admin/util"
	"github.com/crossworth/cartola-web-admin/vk"
)

type ScreeNameProvider interface {
	ScreenNameToProfileID(context context.Context, screenNameOrID string) (int, string, error)
	GroupScreenNameToProfileID(context context.Context, screenNameOrID string) (int, string, error)
}

type ProfileLinkResponse struct {
	ID                   int    `json:"id"`
	ScreenName           string `json:"screen_name"`
	ProfileLink          string `json:"profile_link"`
	CanonicalProfileLink string `json:"canonical_profile_link"`
}

func ProfileLinkToID(provider ScreeNameProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		link := r.URL.Query().Get("link")

		if link == "" {
			json(w, NewError("link de perfil não informado"), 400)
			return
		}

		screenNameOrID, err := vk.ProfileScreenNameOrIDFromURL(link)
		if err != nil {
			errorCode(w, err, 400)
			return
		}

		id, screenName, err := provider.ScreenNameToProfileID(r.Context(), screenNameOrID)
		if err != nil {

			// NOTE(Pedro): Maybe group?
			id, screenName, err = provider.GroupScreenNameToProfileID(r.Context(), screenNameOrID)
			if err != nil {
				errorCode(w, err, 400)
				return
			}
		}

		canonicalProfileLink := "https://vk.com/id" + util.ToString(id)

		if id < 0 {
			canonicalProfileLink = "https://vk.com/club" + util.ToString(int(math.Abs(float64(id))))
		}

		json(w, ProfileLinkResponse{
			ID:                   id,
			ScreenName:           screenName,
			ProfileLink:          "https://vk.com/" + screenName,
			CanonicalProfileLink: canonicalProfileLink,
		}, 200)
	}
}

type ProfileByIDProvider interface {
	ProfileByID(context context.Context, id int) (model.Profile, error)
}

func ProfileByID(provider ProfileByIDProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		profile, err := provider.ProfileByID(r.Context(), id)
		if err != nil {
			databaseError(w, err)
			return
		}

		json(w, profile, 200)
	}
}

type ProfileHistoryProvider interface {
	ProfileHistoryByProfileID(context context.Context, id int) ([]model.ProfileNames, error)
}

func ProfileHistoryByID(provider ProfileHistoryProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		profileHistory, err := provider.ProfileHistoryByProfileID(r.Context(), id)
		if err != nil {
			databaseError(w, err)
			return
		}

		json(w, profileHistory, 200)
	}
}

type ProfileTopicsProvider interface {
	TopicsByProfileID(context context.Context, id int, before int, limit int) ([]model.Topic, error)
	TopicsCountByProfileID(context context.Context, id int) (int, error)
	TopicsPaginationTimestampByProfileID(context context.Context, id int, before int, limit int) (database.PaginationTimestamps, error)
}

func TopicsByProfileID(provider ProfileTopicsProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))
		before := util.ToIntWithDefault(r.URL.Query().Get("before"), int(time.Now().Unix()))
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		total, err := provider.TopicsCountByProfileID(r.Context(), id)
		if err != nil {
			databaseError(w, err)
			return
		}

		topics, err := provider.TopicsByProfileID(r.Context(), id, before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		paginationTimestamps, err := provider.TopicsPaginationTimestampByProfileID(r.Context(), id, before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		next := ""
		prev := ""

		if paginationTimestamps.Next != 0 {
			next = fmt.Sprintf("%s/profiles/%d/topics?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, paginationTimestamps.Next)
		}
		if paginationTimestamps.Prev != 0 {
			prev = fmt.Sprintf("%s/profiles/%d/topics?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, paginationTimestamps.Prev)
		}

		if len(topics) != 0 {
			before = topics[0].CreatedAt
		}

		pagination(w, topics, 200, PaginationMeta{
			Prev:    prev,
			Current: fmt.Sprintf("%s/profiles/%d/topics?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, before),
			Next:    next,
			Total:   total,
		})
	}
}

type ProfileCommentsProvider interface {
	CommentsByProfileID(context context.Context, id int, before int, limit int) ([]model.Comment, error)
	CommentsCountByProfileID(context context.Context, id int) (int, error)
	CommentsPaginationTimestampByProfileID(context context.Context, id int, before int, limit int) (database.PaginationTimestamps, error)
}

func CommentsByProfileID(provider ProfileCommentsProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))
		before := util.ToIntWithDefault(r.URL.Query().Get("before"), int(time.Now().Unix()))
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		total, err := provider.CommentsCountByProfileID(r.Context(), id)
		if err != nil {
			databaseError(w, err)
			return
		}

		comments, err := provider.CommentsByProfileID(r.Context(), id, before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		paginationTimestamps, err := provider.CommentsPaginationTimestampByProfileID(r.Context(), id, before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		next := ""
		prev := ""

		if paginationTimestamps.Next != 0 {
			next = fmt.Sprintf("%s/profiles/%d/comments?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, paginationTimestamps.Next)
		}
		if paginationTimestamps.Prev != 0 {
			prev = fmt.Sprintf("%s/profiles/%d/comments?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, paginationTimestamps.Prev)
		}

		if len(comments) != 0 {
			before = comments[0].Date
		}

		pagination(w, comments, 200, PaginationMeta{
			Prev:    prev,
			Current: fmt.Sprintf("%s/profiles/%d/comments?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, before),
			Next:    next,
			Total:   total,
		})
	}
}

type ProfileStatsAndHistoryProvider interface {
	ProfileStatsByProfileID(context context.Context, id int) (database.ProfileWithStats, error)
	ProfileHistoryProvider
}

type ProfileStatsResponse struct {
	TotalTopics         int `json:"total_topics"`
	TotalComments       int `json:"total_comments"`
	TotalLikes          int `json:"total_likes"`
	TotalProfileChanges int `json:"total_profile_changes"`
}

func ProfileStatsByID(provider ProfileStatsAndHistoryProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		stats, err := provider.ProfileStatsByProfileID(r.Context(), id)
		if err != nil {
			databaseError(w, err)
			return
		}

		totalProfileChanges, err := provider.ProfileHistoryByProfileID(r.Context(), id)
		if err != nil {
			databaseError(w, err)
			return
		}

		json(w, ProfileStatsResponse{
			TotalTopics:         stats.Topics,
			TotalComments:       stats.Comments,
			TotalLikes:          stats.Likes,
			TotalProfileChanges: len(totalProfileChanges),
		}, 200)
	}
}

type ProfileNameProvider interface {
	SearchProfileName(context context.Context, text string) ([]model.Profile, error)
}

func AutoCompleteProfileName(provider ProfileNameProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profile := chi.URLParam(r, "profile")

		if profile == "" {
			json(w, []model.Profile{}, 200)
			return
		}

		profiles, err := provider.SearchProfileName(r.Context(), profile)
		if err != nil {
			databaseError(w, err)
			return
		}

		json(w, profiles, 200)
	}
}

type ProfilesProvider interface {
	ProfileWithStats(context context.Context, order string, orderDirection database.OrderByDirection, period database.Period, page int, limit int) ([]database.ProfileWithStats, error)
	ProfileCount(context context.Context) (int, error)
}

type ProfilesCache struct {
	Profiles  []database.ProfileWithStats
	Total     int
	CreatedAt time.Time
}

func Profiles(provider ProfilesProvider, cache *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orderBy := util.StringWithDefault(r.URL.Query().Get("orderBy"), "topics")
		orderDirStr := util.StringWithDefault(r.URL.Query().Get("orderDir"), "DESC")
		periodStr := util.StringWithDefault(r.URL.Query().Get("period"), "DESC")
		page := util.ToIntWithDefaultMin(r.URL.Query().Get("page"), 1)
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)

		orderDir := database.OrderByASC

		if strings.ToLower(orderDirStr) == "desc" {
			orderDir = database.OrderByDESC
		}

		period := database.PeriodFromString(periodStr)

		cacheKey := fmt.Sprintf("profiles_%s_%s_%d_%d_%s", orderBy, orderDir.Stringer(), page, limit, period.Stringer())
		profilesCache := cache.Get(cacheKey, func() interface{} {
			profiles, err := provider.ProfileWithStats(r.Context(), orderBy, orderDir, period, page, limit)
			if err != nil {
				return err
			}

			total, err := provider.ProfileCount(r.Context())
			if err != nil {
				return err
			}

			profilesCache := ProfilesCache{
				Profiles:  profiles,
				Total:     total,
				CreatedAt: time.Now(),
			}

			return profilesCache
		})

		data, castOK := profilesCache.(ProfilesCache)
		if !castOK {
			databaseError(w, profilesCache.(error))
			return
		}

		next := ""
		prev := ""

		if page != 1 {
			prev = fmt.Sprintf("%s/profiles?limit=%d&page=%d&orderBy=%s&orderDir=%s&period=%s", os.Getenv("APP_API_URL"), limit, page-1, orderBy, orderDir.Stringer(), period.URLString())
		}

		if page*limit < data.Total {
			next = fmt.Sprintf("%s/profiles?limit=%d&page=%d&orderBy=%s&orderDir=%s&period=%s", os.Getenv("APP_API_URL"), limit, page+1, orderBy, orderDir.Stringer(), period.URLString())
		}

		// NOTE(Pedro): To calculate the correct position for the records
		// when in asc order
		if orderDir == database.OrderByASC {
			for i := range data.Profiles {
				data.Profiles[i].Position = (data.Total + 1) - data.Profiles[i].Position
			}
		}

		pagination(w, data.Profiles, 200, PaginationMeta{
			Prev:     prev,
			Current:  fmt.Sprintf("%s/profiles?limit=%d&page=%d&orderBy=%s&orderDir=%s&period=%s", os.Getenv("APP_API_URL"), limit, page, orderBy, orderDir.Stringer(), period.URLString()),
			Next:     next,
			Total:    data.Total,
			CachedAt: &data.CreatedAt,
		})
	}
}

package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"

	"github.com/crossworth/painel-cartolafc/cache"
	"github.com/crossworth/painel-cartolafc/database"
	"github.com/crossworth/painel-cartolafc/httputil"
	"github.com/crossworth/painel-cartolafc/model"
	"github.com/crossworth/painel-cartolafc/util"
	"github.com/crossworth/painel-cartolafc/vk"
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
			httputil.SendJSON(w, httputil.NewError("link de perfil não informado"), 400)
			return
		}

		screenNameOrID, err := vk.ProfileScreenNameOrIDFromURL(link)
		if err != nil {
			httputil.SendErrorCode(w, err, 400)
			return
		}

		id, screenName, err := provider.ScreenNameToProfileID(r.Context(), screenNameOrID)
		if err != nil {

			// NOTE(Pedro): Maybe group?
			id, screenName, err = provider.GroupScreenNameToProfileID(r.Context(), screenNameOrID)
			if err != nil {
				httputil.SendErrorCode(w, err, 400)
				return
			}
		}

		canonicalProfileLink := "https://vk.com/id" + util.ToString(id)

		if id < 0 {
			canonicalProfileLink = "https://vk.com/club" + util.ToString(int(math.Abs(float64(id))))
		}

		httputil.SendJSON(w, ProfileLinkResponse{
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
			httputil.SendJSON(w, httputil.NewError("id de perfil inválido"), 400)
			return
		}

		profile, err := provider.ProfileByID(r.Context(), id)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		httputil.SendJSON(w, profile, 200)
	}
}

type ProfileHistoryProvider interface {
	ProfileHistoryByProfileID(context context.Context, id int) ([]model.ProfileNames, error)
}

func ProfileHistoryByID(provider ProfileHistoryProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))

		if id == 0 {
			httputil.SendJSON(w, httputil.NewError("id de perfil inválido"), 400)
			return
		}

		profileHistory, err := provider.ProfileHistoryByProfileID(r.Context(), id)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		httputil.SendJSON(w, profileHistory, 200)
	}
}

type ProfileTopicsProvider interface {
	TopicsByProfileID(context context.Context, id int, before int, limit int) ([]model.TopicWithLikes, error)
	TopicsCountByProfileID(context context.Context, id int) (int, error)
	TopicsPaginationTimestampByProfileID(context context.Context, id int, before int, limit int) (database.PaginationTimestamps, error)
}

func TopicsByProfileID(provider ProfileTopicsProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))
		before := util.ToIntWithDefault(r.URL.Query().Get("before"), int(time.Now().Unix()))
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)

		if id == 0 {
			httputil.SendJSON(w, httputil.NewError("id de perfil inválido"), 400)
			return
		}

		total, err := provider.TopicsCountByProfileID(r.Context(), id)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		topics, err := provider.TopicsByProfileID(r.Context(), id, before, limit)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		paginationTimestamps, err := provider.TopicsPaginationTimestampByProfileID(r.Context(), id, before, limit)
		if err != nil {
			httputil.SendDatabaseError(w, err)
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

		httputil.SendPagination(w, topics, 200, httputil.PaginationMeta{
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
			httputil.SendJSON(w, httputil.NewError("id de perfil inválido"), 400)
			return
		}

		total, err := provider.CommentsCountByProfileID(r.Context(), id)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		comments, err := provider.CommentsByProfileID(r.Context(), id, before, limit)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		paginationTimestamps, err := provider.CommentsPaginationTimestampByProfileID(r.Context(), id, before, limit)
		if err != nil {
			httputil.SendDatabaseError(w, err)
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

		httputil.SendPagination(w, comments, 200, httputil.PaginationMeta{
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
	TotalTopics             int `json:"total_topics"`
	TotalComments           int `json:"total_comments"`
	TotalLikes              int `json:"total_likes"`
	TotalTopicsPlusComments int `json:"total_topics_plus_comments"`
	TotalProfileChanges     int `json:"total_profile_changes"`
}

func ProfileStatsByID(provider ProfileStatsAndHistoryProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))

		if id == 0 {
			httputil.SendJSON(w, httputil.NewError("id de perfil inválido"), 400)
			return
		}

		stats, err := provider.ProfileStatsByProfileID(r.Context(), id)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		totalProfileChanges, err := provider.ProfileHistoryByProfileID(r.Context(), id)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		httputil.SendJSON(w, ProfileStatsResponse{
			TotalTopics:             stats.Topics,
			TotalComments:           stats.Comments,
			TotalLikes:              stats.Likes,
			TotalTopicsPlusComments: stats.TopicsPlusComments,
			TotalProfileChanges:     len(totalProfileChanges),
		}, 200)
	}
}

type MyProfileProvider interface {
	ProfileStatsAndHistoryProvider
	ProfileByID(context context.Context, id int) (model.Profile, error)
	TopicsWithMoreCommentsByID(context context.Context, id int, limit int) ([]model.TopicWithComments, error)
	TopicsWithMoreLikesByID(context context.Context, id int, limit int) ([]model.TopicWithLikes, error)
	CommentsWithMoreLikes(context context.Context, id int, limit int) ([]model.Comment, error)
}

type MyProfileStats struct {
	TotalTopics             int `json:"total_topics"`
	TotalComments           int `json:"total_comments"`
	TotalLikes              int `json:"total_likes"`
	TotalTopicsPlusComments int `json:"total_topics_plus_comments"`
}

type MyProfileCache struct {
	User                  model.Profile             `json:"user"`
	Stats                 MyProfileStats            `json:"stats"`
	TopicWithMoreLikes    []model.TopicWithLikes    `json:"topic_with_more_likes"`
	TopicWithMoreComments []model.TopicWithComments `json:"topic_with_more_comments"`
	CommentsWithMoreLikes []model.Comment           `json:"comments_with_more_likes"`
}

func MyProfile(provider MyProfileProvider, cache *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vkID, found := model.VKIDFromRequest(r)
		if !found || vkID == 0 {
			httputil.SendError(w, fmt.Errorf("usuário não logado!?"))
			return
		}

		cacheKey := fmt.Sprintf("profile_stats_%d", vkID)
		profilesCache := cache.GetShortCache(cacheKey, func() interface{} {
			stats, err := provider.ProfileStatsByProfileID(r.Context(), vkID)
			if err != nil {
				return err
			}

			user, err := provider.ProfileByID(r.Context(), vkID)
			if err != nil {
				return err
			}

			topicsWithMoreComments, err := provider.TopicsWithMoreCommentsByID(r.Context(), vkID, 5)
			if err != nil {
				return err
			}

			topicsWithMoreLikes, err := provider.TopicsWithMoreLikesByID(r.Context(), vkID, 5)
			if err != nil {
				return err
			}

			commentsWithMoreLikes, err := provider.CommentsWithMoreLikes(r.Context(), vkID, 5)
			if err != nil {
				return err
			}

			profilesCache := MyProfileCache{
				User: user,
				Stats: MyProfileStats{
					TotalTopics:             stats.Topics,
					TotalComments:           stats.Comments,
					TotalLikes:              stats.Likes,
					TotalTopicsPlusComments: stats.TopicsPlusComments,
				},
				TopicWithMoreLikes:    topicsWithMoreLikes,
				TopicWithMoreComments: topicsWithMoreComments,
				CommentsWithMoreLikes: commentsWithMoreLikes,
			}

			return profilesCache
		})

		data, castOK := profilesCache.(MyProfileCache)
		if !castOK {
			httputil.SendDatabaseError(w, profilesCache.(error))
			return
		}

		httputil.SendJSON(w, data, 200)
	}
}

type BotQuoteProvider interface {
	QuotesByBotByID(context context.Context, botID int, id int, page int, limit int) ([]database.QuotesByBot, error)
	QuotesByBotByIDCount(context context.Context, botID int, id int) (int, error)
}

type BotQuoteCache struct {
	Quotes    []database.QuotesByBot `json:"quotes"`
	Total     int                    `json:"total"`
	CreatedAt time.Time              `json:"created_at"`
}

func MyProfileBotQuotes(provider BotQuoteProvider, cache *cache.Cache, botQuoteID int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vkID, found := model.VKIDFromRequest(r)
		page := util.ToIntWithDefaultMin(r.URL.Query().Get("page"), 1)
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)
		if !found || vkID == 0 {
			httputil.SendError(w, fmt.Errorf("usuário não logado!?"))
			return
		}

		cacheKey := fmt.Sprintf("profile_bot_quotes_%d_%d_%d_%d", vkID, botQuoteID, page, limit)
		quotesCache := cache.GetShortCache(cacheKey, func() interface{} {
			quotes, err := provider.QuotesByBotByID(r.Context(), botQuoteID, vkID, page, limit)
			if err != nil {
				return err
			}

			total, err := provider.QuotesByBotByIDCount(r.Context(), botQuoteID, vkID)
			if err != nil {
				return err
			}

			return BotQuoteCache{
				Quotes:    quotes,
				Total:     total,
				CreatedAt: time.Now(),
			}
		})

		data, castOK := quotesCache.(BotQuoteCache)
		if !castOK {
			httputil.SendDatabaseError(w, quotesCache.(error))
			return
		}

		next := ""
		prev := ""

		if page != 1 {
			prev = fmt.Sprintf("%s/my-profile/bot-quotes?limit=%d&page=%d", os.Getenv("APP_API_URL"), limit, page-1)
		}

		if page*limit < data.Total {
			next = fmt.Sprintf("%s/my-profile/bot-quotes?limit=%d&page=%d", os.Getenv("APP_API_URL"), limit, page+1)
		}

		httputil.SendPagination(w, data.Quotes, 200, httputil.PaginationMeta{
			Prev:     prev,
			Current:  fmt.Sprintf("%s/my-profile/bot-quotes?limit=%d&page=%d", os.Getenv("APP_API_URL"), limit, page),
			Next:     next,
			Total:    data.Total,
			CachedAt: &data.CreatedAt,
		})
	}
}

type LastTopicsProvider interface {
	LastTopicsByID(context context.Context, id int, limit int) ([]model.Topic, error)
}

func LastTopics(provider LastTopicsProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vkID, found := model.VKIDFromRequest(r)
		if !found || vkID == 0 {
			httputil.SendError(w, fmt.Errorf("usuário não logado!?"))
			return
		}

		topics, err := provider.LastTopicsByID(r.Context(), vkID, 5)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		httputil.SendJSON(w, topics, 200)
	}
}

type GraphValueCache struct {
	Result    []database.GraphValue `json:"result"`
	CreatedAt time.Time             `json:"created_at"`
}

type TopicsGraphProvider interface {
	TopicsNumberByIDGraph(context context.Context, id int) ([]database.GraphValue, error)
}

func TopicsGraph(provider TopicsGraphProvider, cache *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vkID, found := model.VKIDFromRequest(r)
		if !found || vkID == 0 {
			httputil.SendError(w, fmt.Errorf("usuário não logado!?"))
			return
		}

		cacheKey := fmt.Sprintf("profile_topics_graph_%d", vkID)
		topicsGraphCache := cache.Get(cacheKey, func() interface{} {
			topicsGraph, err := provider.TopicsNumberByIDGraph(r.Context(), vkID)
			if err != nil {
				return err
			}

			return GraphValueCache{
				Result:    topicsGraph,
				CreatedAt: time.Now(),
			}
		})

		data, castOK := topicsGraphCache.(GraphValueCache)
		if !castOK {
			httputil.SendDatabaseError(w, topicsGraphCache.(error))
			return
		}

		httputil.SendJSON(w, data, 200)
	}
}

type CommentsGraphProvider interface {
	CommentsNumberByIDGraph(context context.Context, id int) ([]database.GraphValue, error)
}

func CommentsGraph(provider CommentsGraphProvider, cache *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vkID, found := model.VKIDFromRequest(r)
		if !found || vkID == 0 {
			httputil.SendError(w, fmt.Errorf("usuário não logado!?"))
			return
		}

		cacheKey := fmt.Sprintf("profile_comments_graph_%d", vkID)
		commentsGraphCache := cache.Get(cacheKey, func() interface{} {
			topicsGraph, err := provider.CommentsNumberByIDGraph(r.Context(), vkID)
			if err != nil {
				return err
			}

			return GraphValueCache{
				Result:    topicsGraph,
				CreatedAt: time.Now(),
			}
		})

		data, castOK := commentsGraphCache.(GraphValueCache)
		if !castOK {
			httputil.SendDatabaseError(w, commentsGraphCache.(error))
			return
		}

		httputil.SendJSON(w, data, 200)
	}
}

type ProfileNameProvider interface {
	SearchProfileName(context context.Context, text string) ([]model.Profile, error)
}

func AutoCompleteProfileName(provider ProfileNameProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profile := chi.URLParam(r, "profile")

		if profile == "" {
			httputil.SendJSON(w, []model.Profile{}, 200)
			return
		}

		profiles, err := provider.SearchProfileName(r.Context(), profile)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		httputil.SendJSON(w, profiles, 200)
	}
}

type AdministratorProfileProvider interface {
	GetAdministratorProfiles(context context.Context) ([]model.Profile, error)
	SetAdministratorProfiles(context context.Context, ids []int) error
}

func GetAdministratorProfiles(provider AdministratorProfileProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		profiles, err := provider.GetAdministratorProfiles(r.Context())
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		httputil.SendJSON(w, profiles, 200)
	}
}

func SetAdministratorProfiles(provider AdministratorProfileProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ids []int
		err := json.NewDecoder(r.Body).Decode(&ids)
		if err != nil {
			httputil.SendError(w, err)
			return
		}

		err = provider.SetAdministratorProfiles(r.Context(), ids)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		httputil.SendJSON(w, []int{}, 200)
	}
}

type SettingsProvider interface {
	SettingByName(context context.Context, name string) (string, error)
	UpdateSetting(context context.Context, name string, value string) error
}

type MembersRule struct {
	Value string `json:"value"`
}

func GetMembersRule(provider SettingsProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		value, err := provider.SettingByName(r.Context(), model.MembersRuleSettingName)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		httputil.SendJSON(w, MembersRule{
			Value: value,
		}, 200)
	}
}

func SetMembersRule(provider SettingsProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var result MembersRule
		err := json.NewDecoder(r.Body).Decode(&result)
		if err != nil {
			httputil.SendError(w, err)
			return
		}

		err = provider.UpdateSetting(r.Context(), model.MembersRuleSettingName, result.Value)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		httputil.SendJSON(w, []int{}, 200)
	}
}

type HomePage struct {
	Value string `json:"value"`
}

func GetHomePage(provider SettingsProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		value, err := provider.SettingByName(r.Context(), model.HomePageSettingName)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		httputil.SendJSON(w, MembersRule{
			Value: value,
		}, 200)
	}
}

func SetHomePage(provider SettingsProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var result MembersRule
		err := json.NewDecoder(r.Body).Decode(&result)
		if err != nil {
			httputil.SendError(w, err)
			return
		}

		err = provider.UpdateSetting(r.Context(), model.HomePageSettingName, result.Value)
		if err != nil {
			httputil.SendDatabaseError(w, err)
			return
		}

		httputil.SendJSON(w, []int{}, 200)
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
			httputil.SendDatabaseError(w, profilesCache.(error))
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

		httputil.SendPagination(w, data.Profiles, 200, httputil.PaginationMeta{
			Prev:     prev,
			Current:  fmt.Sprintf("%s/profiles?limit=%d&page=%d&orderBy=%s&orderDir=%s&period=%s", os.Getenv("APP_API_URL"), limit, page, orderBy, orderDir.Stringer(), period.URLString()),
			Next:     next,
			Total:    data.Total,
			CachedAt: &data.CreatedAt,
		})
	}
}

type PublicProfileStatCache struct {
	Profile   database.ProfileWithStats
	CreatedAt time.Time
}

func PublicProfileStatsByID(provider ProfileStatsAndHistoryProvider, cache *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))

		if id == 0 {
			httputil.SendJSON(w, httputil.NewError("id de perfil inválido"), 400)
			return
		}

		cacheKey := fmt.Sprintf("profile_stat_%d", id)
		profileStatCache := cache.GetShortCache(cacheKey, func() interface{} {
			stats, err := provider.ProfileStatsByProfileID(r.Context(), id)
			if err != nil {
				return err
			}

			return PublicProfileStatCache{
				Profile:   stats,
				CreatedAt: time.Now(),
			}
		})

		data, castOK := profileStatCache.(PublicProfileStatCache)
		if !castOK {
			httputil.SendDatabaseError(w, profileStatCache.(error))
			return
		}

		httputil.SendJSON(w, data, 200)
	}
}

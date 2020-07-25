package handle

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"

	"github.com/crossworth/cartola-web-admin/model"
	"github.com/crossworth/cartola-web-admin/util"
	"github.com/crossworth/cartola-web-admin/vk"
)

type ScreeNameProvider interface {
	ScreenNameToUserID(screenNameOrID string) (int, string, error)
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

		screenName, err := vk.ProfileScreenNameOrIDFromURL(link)
		if err != nil {
			errorCode(w, err, 400)
			return
		}

		id, screenName, err := provider.ScreenNameToUserID(screenName)
		if err != nil {
			errorCode(w, err, 400)
			return
		}

		json(w, ProfileLinkResponse{
			ID:                   id,
			ScreenName:           screenName,
			ProfileLink:          "https://vk.com/" + screenName,
			CanonicalProfileLink: "https://vk.com/id" + util.ToString(id),
		}, 200)
	}
}

type ProfileByIDProvider interface {
	ProfileByUserID(id int) (model.Profile, error)
	ProfileHistoryByUserID(id int) ([]model.ProfileNames, error)
}

func ProfileByID(provider ProfileByIDProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		profile, err := provider.ProfileByUserID(id)
		if err != nil {
			databaseError(w, err)
			return
		}

		json(w, profile, 200)
	}
}

func ProfileHistoryByID(provider ProfileByIDProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		profileHistory, err := provider.ProfileHistoryByUserID(id)
		if err != nil {
			databaseError(w, err)
			return
		}

		json(w, profileHistory, 200)
	}
}

type UserTopicProvider interface {
	TopicsByUserID(id int, before int, limit int) ([]model.Topic, error)
	TopicsCountByUserID(id int) (int, error)
	PrevTopicTimestampByUserID(id int, before int, limit int) (int, error)
}

func TopicsByID(provider UserTopicProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))
		before := util.ToIntWithDefault(r.URL.Query().Get("before"), int(time.Now().Unix()))
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		total, err := provider.TopicsCountByUserID(id)
		if err != nil {
			databaseError(w, err)
			return
		}

		topics, err := provider.TopicsByUserID(id, before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		prevTimestamp, err := provider.PrevTopicTimestampByUserID(id, before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		prev := ""
		next := ""

		lenTopics := len(topics)

		if lenTopics == limit && lenTopics < total {
			next = fmt.Sprintf("%s/topics/%d?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, topics[lenTopics-1].CreatedAt)
		}

		if prevTimestamp != 0 {
			prev = fmt.Sprintf("%s/topics/%d?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, prevTimestamp)
		}

		pagination(w, topics, 200, PaginationMeta{
			Prev:    prev,
			Current: fmt.Sprintf("%s/topics/%d?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, before),
			Next:    next,
			Total:   total,
		})
	}
}

type UserCommentProvider interface {
	CommentsByUserID(id int, before int, limit int) ([]model.Comment, error)
	CommentsCountByUserID(id int) (int, error)
	PrevCommentTimestampByUserID(id int, before int, limit int) (int, error)
}

func CommentsByID(provider UserCommentProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))
		before := util.ToIntWithDefault(r.URL.Query().Get("before"), int(time.Now().Unix()))
		limit := util.ToIntWithDefaultMin(r.URL.Query().Get("limit"), 10)

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		total, err := provider.CommentsCountByUserID(id)
		if err != nil {
			databaseError(w, err)
			return
		}

		comments, err := provider.CommentsByUserID(id, before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		prevTimestamp, err := provider.PrevCommentTimestampByUserID(id, before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		next := ""
		prev := ""

		lenComments := len(comments)
		if lenComments == limit && lenComments < total {
			next = fmt.Sprintf("%s/comments/%d?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, comments[lenComments-1].Date)
		}

		if prevTimestamp != 0 {
			prev = fmt.Sprintf("%s/comments/%d?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, prevTimestamp)
		}

		pagination(w, comments, 200, PaginationMeta{
			Prev:    prev,
			Current: fmt.Sprintf("%s/comments/%d?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, before),
			Next:    next,
			Total:   total,
		})
	}
}

type UserTopicCommentProvider interface {
	UserTopicProvider
	UserCommentProvider
}

type ProfileStatsResponse struct {
	TotalTopics   int `json:"total_topics"`
	TotalComments int `json:"total_comments"`
}

func ProfileStatsByID(provider UserTopicCommentProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		totalTopics, err := provider.TopicsCountByUserID(id)
		if err != nil {
			databaseError(w, err)
			return
		}

		totalComments, err := provider.CommentsCountByUserID(id)
		if err != nil {
			databaseError(w, err)
			return
		}

		json(w, ProfileStatsResponse{
			TotalTopics:   totalTopics,
			TotalComments: totalComments,
		}, 200)
	}
}

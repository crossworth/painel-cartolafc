package handle

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"

	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/model"
	"github.com/crossworth/cartola-web-admin/util"
	"github.com/crossworth/cartola-web-admin/vk"
)

type ScreeNameProvider interface {
	ScreenNameToUserID(context context.Context, screenNameOrID string) (int, string, error)
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

		id, screenName, err := provider.ScreenNameToUserID(r.Context(), screenName)
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
	ProfileByUserID(context context.Context, id int) (model.Profile, error)
	ProfileHistoryByUserID(context context.Context, id int) ([]model.ProfileNames, error)
}

func ProfileByID(provider ProfileByIDProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		profile, err := provider.ProfileByUserID(r.Context(), id)
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

		profileHistory, err := provider.ProfileHistoryByUserID(r.Context(), id)
		if err != nil {
			databaseError(w, err)
			return
		}

		json(w, profileHistory, 200)
	}
}

type UserTopicProvider interface {
	TopicsByUserID(context context.Context, id int, before int, limit int) ([]model.Topic, error)
	TopicsCountByUserID(context context.Context, id int) (int, error)
	PaginationTimestampTopicByUserID(context context.Context, id int, before int, limit int) (database.PaginationTimestamps, error)
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

		total, err := provider.TopicsCountByUserID(r.Context(), id)
		if err != nil {
			databaseError(w, err)
			return
		}

		topics, err := provider.TopicsByUserID(r.Context(), id, before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		paginationTimestamps, err := provider.PaginationTimestampTopicByUserID(r.Context(), id, before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		next := ""
		prev := ""

		if paginationTimestamps.Next != 0 {
			next = fmt.Sprintf("%s/topics/%d?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, paginationTimestamps.Next)
		}
		if paginationTimestamps.Prev != 0 {
			prev = fmt.Sprintf("%s/topics/%d?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, paginationTimestamps.Prev)
		}

		if len(topics) != 0 {
			before = topics[0].CreatedAt
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
	CommentsByUserID(context context.Context, id int, before int, limit int) ([]model.Comment, error)
	CommentsCountByUserID(context context.Context, id int) (int, error)
	PaginationTimestampCommentByUserID(context context.Context, id int, before int, limit int) (database.PaginationTimestamps, error)
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

		total, err := provider.CommentsCountByUserID(r.Context(), id)
		if err != nil {
			databaseError(w, err)
			return
		}

		comments, err := provider.CommentsByUserID(r.Context(), id, before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		paginationTimestamps, err := provider.PaginationTimestampCommentByUserID(r.Context(), id, before, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		next := ""
		prev := ""

		if paginationTimestamps.Next != 0 {
			next = fmt.Sprintf("%s/comments/%d?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, paginationTimestamps.Next)
		}
		if paginationTimestamps.Prev != 0 {
			prev = fmt.Sprintf("%s/comments/%d?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, paginationTimestamps.Prev)
		}

		if len(comments) != 0 {
			before = comments[0].Date
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
	ProfileHistoryByUserID(context context.Context, id int) ([]model.ProfileNames, error)
}

type ProfileStatsResponse struct {
	TotalTopics         int `json:"total_topics"`
	TotalComments       int `json:"total_comments"`
	TotalProfileChanges int `json:"total_profile_changes"`
}

func ProfileStatsByID(provider UserTopicCommentProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		totalTopics, err := provider.TopicsCountByUserID(r.Context(), id)
		if err != nil {
			databaseError(w, err)
			return
		}

		totalComments, err := provider.CommentsCountByUserID(r.Context(), id)
		if err != nil {
			databaseError(w, err)
			return
		}

		totalProfileChanges, err := provider.ProfileHistoryByUserID(r.Context(), id)
		if err != nil {
			databaseError(w, err)
			return
		}

		json(w, ProfileStatsResponse{
			TotalTopics:         totalTopics,
			TotalComments:       totalComments,
			TotalProfileChanges: len(totalProfileChanges),
		}, 200)
	}
}

type UserProfileNameProvider interface {
	SearchProfileName(context context.Context, text string) ([]model.Profile, error)
}

func AutoCompleteProfileName(provider UserProfileNameProvider) func(w http.ResponseWriter, r *http.Request) {
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

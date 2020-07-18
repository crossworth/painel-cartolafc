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
	FindProfileByID(id int) (model.Profile, error)
}

func ProfileByID(provider ProfileByIDProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		profile, err := provider.FindProfileByID(id)
		if err != nil {
			databaseError(w, err)
			return
		}

		json(w, profile, 200)
	}
}

type TopicsByIDProvider interface {
	FindTopicByUser(id int, before int64, after int64, limit int) ([]model.Topic, int64, error)
}

func TopicsByID(provider TopicsByIDProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := util.ToInt(chi.URLParam(r, "profile"))
		after := util.ToInt64(r.URL.Query().Get("after"))
		before := util.ToInt64(r.URL.Query().Get("before"))
		limit := util.ToIntWithDefault(r.URL.Query().Get("limit"), 10)

		if id == 0 {
			json(w, NewError("id de perfil inválido"), 400)
			return
		}

		// TODO(Pedro): This is not working!
		if before == 0 && after == 0 {
			before = time.Now().Unix()
		}

		topics, total, err := provider.FindTopicByUser(id, before, after, limit)
		if err != nil {
			databaseError(w, err)
			return
		}

		next := ""
		prev := ""

		// TODO(Pedro): Handle this in a better, reusable way
		if len(topics) > 0 {
			next = fmt.Sprintf("%s/topics/%d?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, topics[len(topics)-1].CreatedAt)
			prev = fmt.Sprintf("%s/topics/%d?limit=%d&after=%d", os.Getenv("APP_API_URL"), id, limit, topics[0].CreatedAt)
		}

		pagination(w, topics, 200, PaginationMeta{
			Current: fmt.Sprintf("%s/topics/%d?limit=%d&before=%d", os.Getenv("APP_API_URL"), id, limit, before),
			Next:    next,
			Prev:    prev,
			Total:   total,
		})
	}
}

package handle

import (
	"net/http"

	"github.com/crossworth/cartola-web-admin/util"
	"github.com/crossworth/cartola-web-admin/vk"
)

type ScreeNameProvider interface {
	ScreenNameToUserID(screenNameOrID string) (int, error)
}

type ProfileLinkResponse struct {
	ID            int    `json:"id"`
	IDProfileLink string `json:"profile_link"`
}

func ProfileLinkToID(provider ScreeNameProvider) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		link := r.URL.Query().Get("link")

		screenName, err := vk.ProfileScreenNameOrIDFromURL(link)
		if err != nil {
			json(w, err, 400)
			return
		}

		id, err := provider.ScreenNameToUserID(screenName)
		if err != nil {
			json(w, err, 400)
			return
		}

		json(w, ProfileLinkResponse{
			ID:            id,
			IDProfileLink: "https://vk.com/id" + util.ToString(id),
		}, 200)
	}
}

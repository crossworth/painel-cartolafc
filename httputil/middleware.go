package httputil

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi"

	"github.com/crossworth/painel-cartolafc/util"
)

// This is a copy of https://github.com/go-chi/chi/blob/master/middleware/strip.go#L13
func RemoveDoubleSlashes(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var path string
		rctx := chi.RouteContext(r.Context())
		if rctx.RoutePath != "" {
			path = rctx.RoutePath
		} else {
			path = r.URL.Path
		}
		if len(path) > 1 && strings.Contains(path, "//") {
			path = strings.ReplaceAll(path, "//", "/")
			http.Redirect(w, r, path, 301)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func OnlyAllowedHost(next http.Handler) http.Handler {
	host := util.GetStringFromEnvOrDefault("APP_BASE_URL", "")
	hostUrl, _ := url.Parse(host)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hostUrl.Host != "" && r.Host != hostUrl.Host {
			r.Close = true
			return
		}

		next.ServeHTTP(w, r)
	})
}

package httputil

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi"

	"github.com/crossworth/cartola-web-admin/util"
)

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

var cfConnectingIp = http.CanonicalHeaderKey("CF-Connecting-IP")
var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

func RealIP(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if rip := realIP(r); rip != "" {
			r.RemoteAddr = rip
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func realIP(r *http.Request) string {
	var ip string

	if cfIp := r.Header.Get(cfConnectingIp); cfIp != "" {
		ip = cfIp
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}

	return ip
}

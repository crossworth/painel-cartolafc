package auth

import (
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"

	"github.com/crossworth/painel-cartolafc/httputil"
	"github.com/crossworth/painel-cartolafc/logger"
	"github.com/crossworth/painel-cartolafc/model"
	"github.com/crossworth/painel-cartolafc/util"
)

var (
	alwaysAllowedRoutes = []string{
		"/favicon.ico",
		"/android-chrome-192x192.png",
		"/android-chrome-512x512.png",
		"/apple-touch-icon.png",
		"/browserconfig.xml",
		"/favicon-16x16.png",
		"/favicon-32x32.png",
		"/mstile-150x150.png",
		"/safari-pinned-tab.svg",
		"/site.webmanifest",
	}
)

func OnlyAuthenticatedUsers(sessionStorage sessions.Store, userHandler *UserHandler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if isRouteAllowed(request.URL.String()) {
				next.ServeHTTP(writer, request)
				return
			}

			session, err := sessionStorage.Get(request, model.UserSession)
			if err != nil {
				logger.Log.Error().Err(err).Msg("erro ao conseguir a session de usuário no middleware")
				httputil.Redirect(writer, request, "/fazer-login?motivo-redirect="+httputil.TextToQueryString("erro ao conseguir sessão do membro"), http.StatusTemporaryRedirect)
				return
			}

			debugUser := util.GetStringFromEnvOrDefault("DEBUG_AS_USER", "")
			if debugUser != "" {
				session.Values["user_id"] = debugUser
			}

			userID, ok := session.Values["user_id"].(string)
			if !ok {
				httputil.Redirect(writer, request, "/fazer-login?motivo-redirect="+httputil.TextToQueryString("membro não autenticado"), http.StatusTemporaryRedirect)
				return
			}

			userIDInt, err := strconv.Atoi(userID)
			if err == nil {
				request = model.SetVKIDOnRequest(request, userIDInt)
			}

			userType := userHandler.GetUserType(userIDInt)
			request = model.SetVKTypeOnRequest(request, userType)

			isUserAllowed := userHandler.IsUserAllowed(userIDInt)

			if !isUserAllowed && userType != "super_admin" {
				httputil.Redirect(writer, request, "/fazer-login?motivo-redirect="+httputil.TextToQueryString("membro não autorizado"), http.StatusTemporaryRedirect)
				return
			}

			if userID != "" {
				next.ServeHTTP(writer, request)
				return
			}

			httputil.Redirect(writer, request, "/fazer-login?motivo-redirect="+httputil.TextToQueryString("membro não autenticado"), http.StatusTemporaryRedirect)
		})
	}
}

func isRouteAllowed(routeCheck string) bool {
	for _, route := range alwaysAllowedRoutes {
		if route == routeCheck {
			return true
		}
	}

	return false
}

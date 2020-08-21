package auth

import (
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"

	"github.com/crossworth/cartola-web-admin/logger"
	"github.com/crossworth/cartola-web-admin/model"
	"github.com/crossworth/cartola-web-admin/util"
)

var (
	alwaysAllowedRoutes = []string{"/favicon.ico"}
)

func OnlyAuthenticatedUsers(sessionStorage sessions.Store, userTypeHandler *UserTypeHandler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if isRouteAllowed(request.URL.String()) {
				next.ServeHTTP(writer, request)
				return
			}

			session, err := sessionStorage.Get(request, model.UserSession)
			if err != nil {
				logger.Log.Error().Err(err).Msg("erro ao conseguir a session de usuário no middleware")
				http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
				return
			}

			debugUser := util.GetStringFromEnvOrDefault("DEBUG_AS_USER", "")
			if debugUser != "" {
				session.Values["user_id"] = debugUser
			}

			userID, ok := session.Values["user_id"].(string)
			if !ok {
				http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
				return
			}

			userIDInt, err := strconv.Atoi(userID)
			if err == nil {
				request = model.SetVKIDOnRequest(request, userIDInt)
			}

			userType := userTypeHandler.GetUserType(userIDInt)
			request = model.SetVKTypeOnRequest(request, userType)

			if userID != "" {
				next.ServeHTTP(writer, request)
				return
			}

			http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
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

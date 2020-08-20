package auth

import (
	"context"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"

	"github.com/crossworth/cartola-web-admin/logger"
	"github.com/crossworth/cartola-web-admin/model"
	"github.com/crossworth/cartola-web-admin/util"
	"github.com/crossworth/cartola-web-admin/vk"
)

func LoginPage(appName string) func(http.ResponseWriter, *http.Request) {
	loginPageTemplate, err := template.New("loginPage").Parse(loginPage)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("erro ao criar template de login")
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(context.WithValue(request.Context(), "provider", "vk"))
		_, err := gothic.CompleteUserAuth(writer, request)
		if err == nil {
			http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
			return
		}

		_ = loginPageTemplate.Execute(writer, struct {
			Title string
		}{
			Title: appName,
		})
	}
}

func Login() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(context.WithValue(request.Context(), "provider", "vk"))
		gothic.BeginAuthHandler(writer, request)
	}
}

func LoginCallback(vkAPI *vk.VKClient, sessionStorage sessions.Store) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(context.WithValue(request.Context(), "provider", "vk"))
		user, err := gothic.CompleteUserAuth(writer, request)
		if err != nil {
			logger.Log.Error().Err(err).Msg("erro ao fazer segunda parte do login")
			http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
			return
		}

		session, err := sessionStorage.Get(request, model.UserSession)
		if err != nil {
			logger.Log.Error().Err(err).Msg("erro ao conseguir a session de usuário")
			http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
			return
		}

		session.Values["user_id"] = user.UserID
		isOnCommunity, err := vkAPI.IsUserIDOnGroup(request.Context(), user.UserID,
			util.GetIntFromEnvOrFatalError("APP_VK_GROUP_ID"))
		if err != nil {
			logger.Log.Error().Err(err).Msg("erro ao salvar a session de usuário")
			http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
			return
		}

		if !isOnCommunity {
			_ = gothic.Logout(writer, request)
			logger.Log.Warn().Str("vk_id", user.UserID).Msg("usuário não faz parte da comunidade")
			http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
			return
		}

		err = session.Save(request, writer)
		if err != nil {
			logger.Log.Error().Err(err).Msg("erro ao salvar a session de usuário")
			http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
			return
		}

		logger.Log.Info().Str("vk_id", user.UserID).Msg("login membro")
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
	}
}

func Logout(sessionStorage sessions.Store) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		request = request.WithContext(context.WithValue(request.Context(), "provider", "vk"))

		session, _ := sessionStorage.Get(request, model.UserSession)
		userID, _ := session.Values["user_id"].(string)
		delete(session.Values, "user_id")
		_ = session.Save(request, writer)

		_ = gothic.Logout(writer, request)
		logger.Log.Info().Str("vk_id", userID).Msg("logout membro")
		http.Redirect(writer, request, "/fazer-login", http.StatusTemporaryRedirect)
	}
}

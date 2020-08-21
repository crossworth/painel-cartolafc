package api

import (
	"net/http"

	"github.com/crossworth/cartola-web-admin/httputil"
	"github.com/crossworth/cartola-web-admin/logger"
	"github.com/crossworth/cartola-web-admin/model"
)

func OnlySuperAdmin() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			vkType, _ := model.VKTypeFromRequest(request)

			if vkType == "super_admin" {
				next.ServeHTTP(writer, request)
				return
			}

			logger.LogFromRequest(request).Info().Msg("usuário não é administrador")
			httputil.Redirect(writer, request, "/?motivo-redirect="+httputil.TextToQueryString("você precisa ser um super administrador para ver isso"), http.StatusTemporaryRedirect)
		})
	}
}

func OnlyAdmin() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			vkType, _ := model.VKTypeFromRequest(request)

			if vkType == "admin" || vkType == "super_admin" {
				next.ServeHTTP(writer, request)
				return
			}

			logger.LogFromRequest(request).Info().Msg("usuário não é super administrador")
			httputil.Redirect(writer, request, "/?motivo-redirect="+httputil.TextToQueryString("você precisa ser um administrador para ver isso"), http.StatusTemporaryRedirect)
		})
	}
}

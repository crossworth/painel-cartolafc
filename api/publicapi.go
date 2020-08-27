package api

import (
	"compress/flate"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/crossworth/cartola-web-admin/api/handle"
	"github.com/crossworth/cartola-web-admin/cache"
	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/httputil"
	"github.com/crossworth/cartola-web-admin/logger"
)

type PublicAPI struct {
	router chi.Router
	db     *database.PostgreSQL
	cache  *cache.Cache
}

func NewPublicAPI(db *database.PostgreSQL, cache *cache.Cache) *PublicAPI {
	api := &PublicAPI{
		router: chi.NewRouter(),
		db:     db,
		cache:  cache,
	}

	logger.SetupLoggerOnRouter(api.router)

	api.router.Use(middleware.NoCache)
	api.router.Use(middleware.Compress(flate.DefaultCompression))
	api.router.Get("/profile-stat/{profile}", handle.PublicProfileStatsByID(api.db, api.cache))
	api.router.Get("/download-video/{video}", downloadVideo)
	return api
}

func (a *PublicAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

const extraKey = `hxoWHtZpj6ampd868QLLFA`
const cKey = `5c7c282e_4803f2`
const token = `UWkutTwTkIZnCdU9X_Mmm1-doiipV`
const credentials = `1xKkwe5Azv4rgNeJo-UlA8JDle6YnAtpP10sJHZdI6yhXrFnGS0rAhcIeq95U28JR5JKuo_TBXiIRqwIuCB4Utt-nWgnE__5DcoCY8dzf15WkvZ8XZu19_ZvA2L2WPxX-unIORODHh3nGi799E9Ah27DyW9`

type Item struct {
	Title string            `json:"title"`
	Files map[string]string `json:"files"`
}

type Response struct {
	Count int    `json:"count"`
	Items []Item `json:"items"`
}

type DaxabResponse struct {
	Response Response `json:"response"`
}

type VideoDownloadResponse struct {
	Found   bool   `json:"found"`
	Message string `json:"message"`
	Link    string `json:"link"`
}

func downloadVideo(writer http.ResponseWriter, request *http.Request) {
	video := chi.URLParam(request, "video")

	if video == "" {
		writer.Header().Add("Content-Type", "text/html")
		sendResponse(writer, request, "vídeo não informado")
		return
	}

	req, err := http.NewRequest("GET", "https://psv78-3.daxab.com/method/video.get", nil)
	if err != nil {
		sendResponse(writer, request, fmt.Sprintf("erro ao criar a request, %v", err))
		return
	}

	params := req.URL.Query()
	params.Set("credentials", credentials)
	params.Set("token", token)
	params.Set("extra_key", extraKey)
	params.Set("ckey", cKey)
	params.Set("videos", video)
	req.URL.RawQuery = params.Encode()

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36")

	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		sendResponse(writer, request, fmt.Sprintf("erro ao fazer a request, %v", err))
		return
	}
	defer resp.Body.Close()

	var response DaxabResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		sendResponse(writer, request, fmt.Sprintf("erro ao fazer o decode da request, %v", err))
		return
	}

	if response.Response.Count == 0 {
		sendResponse(writer, request, "nenhum vídeo encontrado")
		return
	}

	bestVideoFormat := getBestVideoSizeLink(response.Response.Items[0].Files)

	if bestVideoFormat == "" {
		sendResponse(writer, request, "nenhum formato mp4 encontrado")
		return
	}

	downloadUrl := fmt.Sprintf("https://psv78-3.daxab.com/%s&extra_key=%s&videos=%s&dl=1", strings.TrimPrefix(bestVideoFormat, "https://"), extraKey, video)

	downloadResp, err := http.Get(downloadUrl)
	if err != nil {
		sendResponse(writer, request, fmt.Sprintf("erro ao fazer o pipe do vídeo, %v", err))
		return
	}
	defer downloadResp.Body.Close()

	if strings.Contains(downloadResp.Header.Get("Content-Type"), "mp4") {
		if httputil.ExpectsJson(request) {
			writer.Header().Add("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(VideoDownloadResponse{
				Found:   true,
				Message: "",
				Link:    fmt.Sprintf("/public/api/video-download/%s", video),
			})
		} else {
			writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.mp4", url.QueryEscape(response.Response.Items[0].Title)))
			writer.Header().Set("Content-Type", downloadResp.Header.Get("Content-Type"))
			_, _ = io.Copy(writer, downloadResp.Body)
		}
	} else {
		sendResponse(writer, request, fmt.Sprintf("não foi possível baixar o vídeo, VK retornou, %v", downloadResp.Header.Get("Content-Type")))
		return
	}
}

func sendResponse(writer http.ResponseWriter, request *http.Request, msg string) {
	if httputil.ExpectsJson(request) {
		writer.Header().Add("Content-Type", "application/json")

		_ = json.NewEncoder(writer).Encode(VideoDownloadResponse{
			Found:   false,
			Message: msg,
			Link:    "",
		})
	} else {
		writer.Header().Add("Content-Type", "text/html")
		_, _ = writer.Write([]byte(fmt.Sprintf(basicHtml, msg)))
	}
}

func getBestVideoSizeLink(videos map[string]string) string {
	for format, link := range videos {
		if format == "mp4_1080" {
			return link
		}

		if format == "mp4_720" {
			return link
		}

		if format == "mp4_480" {
			return link
		}

		if format == "mp4_360" {
			return link
		}

		if format == "mp4_240" {
			return link
		}
	}

	return ""
}

const basicHtml = `
<!doctype html>
<html lang=pt_br>
<head>
<meta charset=utf-8>
<title>Download Vídeo</title>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
<p>%s</p>
</body>
</html>
`

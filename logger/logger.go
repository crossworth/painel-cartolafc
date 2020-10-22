package logger

import (
	stdLog "log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/crossworth/multiwriter"
	"github.com/go-chi/chi"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/crossworth/painel-cartolafc/model"
)

type LogLevel int

const (
	LogPanic = iota
	LogFatal
	LogError
	LogWarn
	LogDebug
	LogInfo
)

var Log = log.Logger

func Setup(logLevel LogLevel, fileName string) {
	switch logLevel {
	case LogPanic:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case LogFatal:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case LogError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case LogWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LogDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LogInfo:
		fallthrough
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	Log = log.Output(zerolog.ConsoleWriter{Out: colorable.NewColorableStdout()})

	if fileName != "" {
		fileHandle, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			stdLog.Printf("não foi possível abrir o arquivo de log, %v", err)
		}

		if fileHandle != nil {
			wr := diode.NewWriter(
				zerolog.ConsoleWriter{Out: colorable.NewColorableStdout()},
				1000,
				0,
				func(missed int) {})

			mw := multiwriter.MultiWriter{
				IO1: wr,
				IO2: fileHandle,
			}

			Log = log.Output(&mw)
		}
	}
}

func SetupLoggerOnRouter(router chi.Router) {
	router.Use(hlog.NewHandler(Log))

	router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		l := hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration)

		if vkID, ok := model.VKIDFromRequest(r); ok {
			l.Str("vk_id", strconv.Itoa(vkID))
		}

		if vkType, ok := model.VKTypeFromRequest(r); ok {
			l.Str("vk_type", vkType)
		}

		l.Msg("")
	}))

	router.Use(hlog.UserAgentHandler("user_agent"))
	router.Use(hlog.RefererHandler("referer"))
	router.Use(hlog.RequestIDHandler("req_id", "request-id"))
}

func LogFromRequest(r *http.Request) *zerolog.Logger {
	return hlog.FromRequest(r)
}

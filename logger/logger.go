package logger

import (
	stdLog "log"
	"os"

	"github.com/crossworth/multiwriter"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
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

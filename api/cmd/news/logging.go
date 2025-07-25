package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sunba23/news/config"
)

const timestampFormat = "2006-01-02T15:04:05.999Z07:00"

func SetUpLogging(cnf *config.Config) {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	zerolog.TimeFieldFormat = timestampFormat
	zerolog.TimestampFieldName = "timestamp"
	zerolog.MessageFieldName = "event"

	setGlobalLevel(cnf.LoggingLevel)

	logger := zerolog.New(os.Stdout).
		With().Caller().Timestamp().Logger()

	if cnf.LoggingPretty {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	log.Logger = logger
}

func setGlobalLevel(level string) {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		panic(fmt.Sprintf("invalid logging level: %s", level))
	}
}

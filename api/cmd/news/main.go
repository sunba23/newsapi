package main

import (
	"github.com/k0kubun/pp/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sunba23/news/config"
	"github.com/sunba23/news/internal/news"
)

func main() {
	conf, err := config.InitConfig()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	SetUpLogging(conf)

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		log.Debug().Msg("logging level is DEBUG. printing configuration")
		pp.Printf("%v\n", conf)
	}

	app, err := news.NewApplication(conf)

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	err = RunServer(app)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

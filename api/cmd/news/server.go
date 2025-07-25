package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sunba23/news/api"
	"github.com/sunba23/news/internal/news"
)

func RunServer(app news.App) error {
	conf := app.Config()

	handler := api.NewHttpHandler(app)
	server := &http.Server{
		Addr:        conf.ServerHost,
		ReadTimeout: time.Second * time.Duration(conf.ServerReadTimeoutSeconds),
		Handler:     handler,
	}

	go func() {
		log.Info().Msg(fmt.Sprintf("Starting server available at %v", conf.ServerHost))
		if err := server.ListenAndServe(); err != nil {
			log.Error().Err(err).Send()
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Info().Msg("Received interrupt. Cleaning up...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	os.Exit(0)
	return nil
}

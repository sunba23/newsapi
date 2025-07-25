package news

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/sunba23/news/config"
	"github.com/sunba23/news/internal/database"
)

type App interface {
	Config() *config.Config
	Repository() *database.Repository
}

type Application struct {
	config     *config.Config
	repository *database.Repository
}

func (app *Application) Config() *config.Config {
	return app.config
}

func (app *Application) Repository() *database.Repository {
	return app.repository
}

func NewApplication(conf *config.Config) (*Application, error) {
	db, err := sqlx.Connect("postgres", conf.PostgresConnStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	repo := database.NewSQLRepository(db)

	app := Application{
		config:     conf,
		repository: &repo,
	}
	return &app, nil
}

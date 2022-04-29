package main

import (
	"context"
	"flag"

	"main.go/internal/api/tbot"
	"main.go/internal/config"
	"main.go/internal/database/postgress"

	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		logger         *log.Logger
		cfg            *config.Config
		err            error
		configFilePath = flag.String("config", "./config.json", "path to configuration file")
	)

	ctx := context.Background()

	logger = log.New()
	logger.SetLevel(log.DebugLevel)

	cfg, err = config.LoadConfig(*configFilePath)
	if err != nil {
		logger.WithError(err).Fatal("reading config file error")
	}

	db, err := postgress.NewClient(*cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("constructing database error")
	}

	log.Info("successfully connected to DB!!")

	telegramServer, err := tbot.NewTelegramServer(cfg.API.BotToken, logger, db)
	if err != nil {
		logger.WithError(err).Fatal("constructing bot error")
	}
	telegramServer.RunBot(ctx)
}

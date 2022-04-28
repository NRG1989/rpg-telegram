package main

import (
	"context"
	"flag"

	"main.go/internal/api/telegram_bot"
	"main.go/internal/config"
	"main.go/internal/database/postgress"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		logger         *log.Logger
		cfg            *config.Config
		configFilePath = flag.String("config", "./config.json", "path to configuration file")
		envPrefix      = flag.String("env_prefix", "go-aut-registration-user", "environment prefix for override config variables")
	)

	ctx := context.Background()

	logger = log.New()
	logger.SetLevel(log.DebugLevel)

	var err error
	cfg, err = config.LoadConfig(*configFilePath, *envPrefix, logger)
	if err != nil {
		logger.WithError(err).Fatal("reading config file error")
	}

	db, err := postgress.NewClient(*cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("constructing database error")
	}

	log.Info("successfully connected to DB!!")

	bot, err := tgbotapi.NewBotAPI(cfg.API.BotToken)
	if err != nil {
		logger.WithError(err).Fatal("not connected to bot")
	}

	stageHandler := telegram_bot.NewStage(db, bot)

	telegram_bot.RunBot(ctx, logger, bot, stageHandler)

}

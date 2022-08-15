package main

import (
	"context"
	"flag"

	"go-aut-registration-user-telegram/internal/api/grpc"
	"go-aut-registration-user-telegram/internal/api/tbot"
	"go-aut-registration-user-telegram/internal/config"
	"go-aut-registration-user-telegram/internal/database/postgress"

	"github.com/sirupsen/logrus"
)

func main() {
	var (
		logger         *logrus.Logger
		cfg            *config.Config
		err            error
		configFilePath = flag.String("config", "./config.json", "path to configuration file")
	)

	ctx := context.Background()

	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	cfg, err = config.LoadConfig(*configFilePath)
	if err != nil {
		logger.WithError(err).Fatal("reading config file error")
	}

	db, err := postgress.NewClient(*cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("constructing database error")
	}

	storage := postgress.NewStorage(db)

	logger.Info("successfully connected to DB!!")

	telegramServer, err := tbot.NewTelegramServer(cfg.API.BotToken, logger, storage)
	if err != nil {
		logger.WithError(err).Fatal("constructing bot error")
	}

	go func() {
		logger.Info("starting 5012 gRPC server")
		if err := grpc.NewGRPCService(logger, cfg.GrpcCfg.AddressTelegramSrv, storage, telegramServer).Run(); err != nil {
			logger.WithError(err).Fatal("connection to grpc error")
		}
	}()

	telegramServer.RunBot(ctx)
}

package main

import (
	"context"
	"flag"

	"tgbotapi/internal/api/grpc"
	"tgbotapi/internal/api/tbot"
	"tgbotapi/internal/config"
	"tgbotapi/internal/database/postgress"

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
		logger.Info("starting gRPC server")
		if err := grpc.NewGRPCService(logger, cfg.GrpcCfg.AddressTelegramSrv, storage, telegramServer).Run(); err != nil {
			logger.WithError(err).Fatal("connection to grpc error")
		}
	}()

	go func() {
		logger.Info("starting gRPC server")
		if err := grpc.NewGRPCService(logger, cfg.GrpcCfg.AddressUserSrv, storage, telegramServer).Run(); err != nil {
			logger.WithError(err).Fatal("connection to grpc error")
		}
	}()

	telegramServer.RunBot(ctx)
}

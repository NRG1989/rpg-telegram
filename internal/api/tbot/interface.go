package tbot

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Stager interface {
	IsInProcess(chat int64) bool
	StageUnknownCommand(logger *logrus.Logger, chat int64)
	StageStart(ctx context.Context, logger *logrus.Logger, c chan string, chat int64)
	StageInfo(logger *logrus.Logger, chat int64)
	SendCode(ctx context.Context, logger *logrus.Logger, phone string) error
	StageError(logger *logrus.Logger, chat int64, err error)
}

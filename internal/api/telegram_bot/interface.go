package telegram_bot

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Stager interface {
	IsInProcess(chat int64) bool
	StageUnknownCommand(chat int64)
	StageStart(ctx context.Context, logger *logrus.Logger, chat int64)
	StageInfo(chat int64, message string)
}

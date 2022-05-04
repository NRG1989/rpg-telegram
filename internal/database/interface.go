package database

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Storage interface {
	AddUser(ctx context.Context, logger *logrus.Logger, chat int64, phone string) error
	IsExist(ctx context.Context, logger *logrus.Logger, phone string) (bool, error)
	FindUserChatId(ctx context.Context, logger *logrus.Logger, phone string) (int64, error)
}

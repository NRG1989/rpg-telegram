package database

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Storage interface {
	AddUser(ctx context.Context, logger *logrus.Logger, chat int64, phone string) error
	IsUserExist(ctx context.Context, logger *logrus.Logger, phone string) (bool, error)
	FindUserChatId(ctx context.Context, logger *logrus.Logger, phone string) (string, int64, error)
	IsPhoneExist(ctx context.Context, logger *logrus.Logger, phone string) (bool, error)
	AddPhone(ctx context.Context, logger *logrus.Logger, phone string, id string) error
}

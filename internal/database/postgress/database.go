package postgress

import (
	"context"

	"tgbotapi/internal/config"

	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Client struct {
	*sqlx.DB
	logger logrus.FieldLogger

	schemaName string
}

func NewClient(cfg config.Config, logger logrus.FieldLogger) (*Client, error) {
	db, err := sqlx.Open("pgx", cfg.DB.URL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	return &Client{
		db,
		logger,
		cfg.DB.SchemaName,
	}, nil
}

type loggerAdapter struct {
	logger logrus.FieldLogger
}

func (l *loggerAdapter) Log(_ context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	switch level {
	case pgx.LogLevelTrace:
		l.logger.WithFields(data).Debugf(msg)
	case pgx.LogLevelDebug:
		l.logger.WithFields(data).Debugf(msg)
	case pgx.LogLevelInfo:
		l.logger.WithFields(data).Infof(msg)
	case pgx.LogLevelWarn:
		l.logger.WithFields(data).Warnf(msg)
	case pgx.LogLevelError:
		l.logger.WithFields(data).Errorf(msg)
	case pgx.LogLevelNone:
		l.logger.WithFields(data).Errorf(msg)
	}
}

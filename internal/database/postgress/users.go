package postgress

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

type database struct {
	client *Client
}

func NewStorage(client *Client) *database {
	return &database{
		client: client,
	}
}

func (db database) AddUser(ctx context.Context, logger *logrus.Logger, chat int64, phone string) error {
	qb := sq.
		Update("userservice.client_telegram").
		Set("chat_id", chat).
		Where(sq.Eq{"phone_number": phone})

	query, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logger.Error(err)
		return err
	}

	if _, err = db.client.ExecContext(ctx, query, args...); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (db database) IsUserExist(ctx context.Context, logger *logrus.Logger, phone string) (bool, error) {
	qb := sq.
		Select(
			"chat_id",
		).
		From("userservice.client_telegram").
		Where(sq.Eq{"phone_number": phone})

	query, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logger.Error(err)
		return false, err
	}
	var chatId int64
	if err = db.client.QueryRowxContext(ctx, query, args...).Scan(&chatId); err != nil {
		logger.Error(err)
		return false, err
	}
	if chatId == 0 {
		return false, nil
	}
	return true, nil
}

func (db database) FindUserChatId(ctx context.Context, logger *logrus.Logger, phone string) (int64, string, error) {
	qb := sq.
		Select(
			"id",
			"chat_id",
		).
		From("userservice.client_telegram").
		Where(sq.Eq{"phone_number": phone})

	query, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logger.Error(err)
		return 0, "", err
	}

	var (
		chatId   int64
		clientID string
	)

	if err = db.client.QueryRowxContext(ctx, query, args...).Scan(&chatId, &clientID); err != nil {
		logger.Error(err)
		return 0, "", err
	}
	return chatId, clientID, nil
}

func (db database) IsPhoneExist(ctx context.Context, logger *logrus.Logger, phone string) (bool, error) {
	qb := sq.
		Select(
			"phone_number",
		).
		From("userservice.client_telegram").
		Where(sq.Eq{"phone_number": phone})

	query, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logger.Error(err)
		return false, err
	}
	var scanedPhone string
	if err = db.client.QueryRowxContext(ctx, query, args...).Scan(&scanedPhone); err != nil {
		logger.Error(err)
		return false, err
	}
	return true, nil
}

func (db database) AddPhone(ctx context.Context, logger *logrus.Logger, phone string, id string) error {

	query, args, err := sq.Insert("userservice.client_telegram").
		SetMap(map[string]interface{}{
			"id":           id,
			"phone_number": phone,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Error(err)
		return err
	}

	if _, err = db.client.ExecContext(ctx, query, args...); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

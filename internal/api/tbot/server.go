package tbot

import (
	"context"
	"strings"

	"tgbotapi/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	getCommandStart = "/start"
	getCommandInfo  = "/info"
)

type telegramServer struct {
	*tgbotapi.BotAPI
	logger       *logrus.Logger
	stageHandler Stager
}

func NewTelegramServer(botToken string, logger *logrus.Logger, db database.Storage) (*telegramServer, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.WithError(err).Fatal("bot starting error")
		return nil, err
	}
	stage := NewStage(db, bot)
	return &telegramServer{
		bot,
		logger,
		stage,
	}, nil
}

var AdditionalChat = make(map[int64]chan string)

func (s *telegramServer) RunBot(ctx context.Context) {

	s.logger.Info("telegram bot started!")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		s.logger.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if !s.stageHandler.IsInProcess(update.Message.Chat.ID) {
			if c, ok := AdditionalChat[update.Message.Chat.ID]; ok {
				delete(AdditionalChat, update.Message.Chat.ID)
				close(c)
			}
		}

		if c, ok := AdditionalChat[update.Message.Chat.ID]; ok {
			switch {
			case update.Message.Contact == nil:
				text := []string{"Пришлите свой контакт - нажав на кнопку, ",
					"либо выбрав свой номер из списка контактов"}
				s.stageHandler.StageError(s.logger, update.Message.Chat.ID, errors.New(strings.Join(text, "")))
			case update.Message.Contact.PhoneNumber != "" && update.Message.From.ID == update.Message.Contact.UserID:
				c <- update.Message.Contact.PhoneNumber
			default:
				text := []string{"Вы не можете исспользовать номер телефона, ",
					"на который не зарегестрирован ваш телеграм пользователь."}
				s.stageHandler.StageError(s.logger, update.Message.Chat.ID, errors.New(strings.Join(text, "")))
			}
			continue
		}

		switch update.Message.Text {
		case getCommandStart:
			if _, ok := AdditionalChat[update.Message.Chat.ID]; !ok {
				c := make(chan string)
				AdditionalChat[update.Message.Chat.ID] = c
				go s.stageHandler.StageStart(ctx, s.logger, c, update.Message.Chat.ID)
			}
		case getCommandInfo:
			s.stageHandler.StageInfo(s.logger, update.Message.Chat.ID)
		default:
			s.stageHandler.StageUnknownCommand(s.logger, update.Message.Chat.ID)
		}
	}
}

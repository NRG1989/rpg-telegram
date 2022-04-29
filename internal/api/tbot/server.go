package tbot

import (
	"context"
	"main.go/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

const (
	getCommandStart = "/start"
	getCommandInfo  = "/info"
)

type telegramServer struct {
	bot          *tgbotapi.BotAPI
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

	updates := s.bot.GetUpdatesChan(u)

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
			c <- update.Message.Text
			continue
		}

		switch update.Message.Text {
		case getCommandStart:
			s.stageHandler.StageStart(ctx, s.logger, update.Message.Chat.ID)
		case getCommandInfo:
			s.stageHandler.StageInfo(update.Message.Chat.ID, update.Message.Text)
		default:
			s.stageHandler.StageUnknownCommand(update.Message.Chat.ID)
		}
	}
}

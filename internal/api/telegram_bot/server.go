package telegram_bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

const (
	getCommandStart = "/start"
	getCommandInfo  = "/info"
)

var AdditionalChat = make(map[int64]chan string)

func RunBot(ctx context.Context, logger *logrus.Logger, bot *tgbotapi.BotAPI, stageHandler Stager) {

	logger.Info("telegram bot started!")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		logger.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if !stageHandler.IsInProcess(update.Message.Chat.ID) {
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
		default:
			stageHandler.StageUnknownCommand(update.Message.Chat.ID)
		case getCommandStart:
			stageHandler.StageStart(ctx, logger, update.Message.Chat.ID)
		case getCommandInfo:
			stageHandler.StageInfo(update.Message.Chat.ID, update.Message.Text)

		}
	}
}

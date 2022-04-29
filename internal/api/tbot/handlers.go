package tbot

import (
	"context"
	"log"

	"main.go/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

var myKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("отправить код повторно"),
	),
)

func NewStage(db database.Storage, bot *tgbotapi.BotAPI) *stage {
	return &stage{
		statusStageMap: make(map[int64]bool),
		db:             db,
		bot:            bot,
	}
}

type stage struct {
	statusStageMap map[int64]bool
	db             database.Storage
	bot            *tgbotapi.BotAPI
}

func (s *stage) IsInProcess(chat int64) bool {
	_, ok := s.statusStageMap[chat]
	return ok
}

func (s *stage) StageStart(ctx context.Context, logger *logrus.Logger, chat int64) {

	//Определяем существует ли пользователь, если нет, то добавляем его в таблицу   //TODO: future logic connecting with DB
	//flag, err := s.db.IsExist(ctx, logger, chat)
	//if err != nil {
	//	logger.Warning("problems with DB, when trying to check user: %s", err)
	//	return
	//}
	//if !flag {
	//	if err := s.db.AddUser(ctx, logger, chat); err != nil {
	//		logger.Warning("user was not added: %s", err)
	//		return
	//	}
	//}

	if _, err := s.bot.Send(tgbotapi.NewMessage(chat, "Добрый день, вам отправлен код активации для регистрации на сайте RPG-BANK")); err != nil {
		logger.Warning("message 'Добрый день....' was not send: %s", err)
	}
}

func (s *stage) StageUnknownCommand(chat int64) {
	if _, err := s.bot.Send(tgbotapi.NewMessage(chat, "Команда не распознана, выберите из списка /info")); err != nil {
		log.Printf("message 'Команда не распознана....' was not send: %s", err)
	}
}

func (s *stage) StageInfo(chat int64, message string) {
	msg := tgbotapi.NewMessage(chat, message)
	msg.ReplyMarkup = myKeyboard
	if _, err := s.bot.Send(msg); err != nil {
		log.Printf("some problems with main keyboard: %s", err)
	}
}

package tbot

import (
	"context"
	"strings"

	"tgbotapi/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

var (
	phoneNumberKeyboard = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{Text: "отправить ваш личный номер телефона", RequestContact: true},
		),
	)
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

func (s *stage) StageStart(ctx context.Context, logger *logrus.Logger, c chan string, chat int64) {
	s.statusStageMap[chat] = true
	defer delete(s.statusStageMap, chat)

	text := []string{"Добрый день, спасибо что присоединились к телеграм боту RPG-bank, для вашей ",
		"верификации нам необходимо подтвердить ваш номер телефона. Если вы согласны с обработкой ваших персональных данных - нажмите ",
		"на кнопку ниже, если нет, то можете ничего не делать."}

	msg := tgbotapi.NewMessage(chat, strings.Join(text, ""))
	msg.ReplyMarkup = phoneNumberKeyboard
	if _, err := s.bot.Send(msg); err != nil {
		logger.Printf("some problems with phone keyboard: %s", err)
	}

	phoneNumber := <-c
	if phoneNumber[0] != '+' {
		phoneNumber = "+" + phoneNumber
	}
	logger.Printf("we get phone number: " + phoneNumber)

	// Определяем существует ли номер чата пользователья если нет, то добавляем его в таблицу
	flag, err := s.db.IsUserExist(ctx, logger, phoneNumber)
	if err != nil {
		msg = tgbotapi.NewMessage(chat, "Этот номер телефона не принадлежит клиенту банка")
		if _, err := s.bot.Send(msg); err != nil {
			logger.Printf("some problems with sending message: %s", err)
		}
		logger.Warning("problems with DB, when trying to check user: %s", err)
		return
	}
	if flag {
		msg = tgbotapi.NewMessage(chat, "Пользователь с этим номером телефона уже зарегистрирован")
		if _, err := s.bot.Send(msg); err != nil {
			logger.Printf("some problems with sending message: %s", err)
		}
		return
	}
	if !flag {
		if err := s.db.AddUser(ctx, logger, chat, phoneNumber); err != nil {
			logger.Warning("user was not added: %s", err)
			return
		}
		msg = tgbotapi.NewMessage(chat, "Спасибо за оказанное доверие. Теперь при помощи Телеграм у вас будет возможность управлять вашими личными данными RPG-bank.")
		if _, err := s.bot.Send(msg); err != nil {
			logger.Printf("some problems with sending message: %s", err)
		}
	}
}

func (s *stage) StageUnknownCommand(logger *logrus.Logger, chat int64) {
	if _, err := s.bot.Send(tgbotapi.NewMessage(chat, "Данный бот не способен отвечать на подобные запросы")); err != nil {
		logger.Printf("message 'Данный бот не способен....' was not send: %s", err)
	}
}

func (s *stage) StageInfo(logger *logrus.Logger, chat int64) {
	if _, err := s.bot.Send(tgbotapi.NewMessage(chat, "В случае необходимости телеграмм бот отправит вам сообщение, касающееся RPG-BANK")); err != nil {
		logger.Printf("message 'В случае необходимости....' was not send: %s", err)
	}
}

func (s *stage) StageError(logger *logrus.Logger, chat int64, err error) {
	if _, err := s.bot.Send(tgbotapi.NewMessage(chat, err.Error())); err != nil {
		logger.Printf("message was not send: %s", err)
	}
}

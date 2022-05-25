package grpc

import (
	"context"
	"crypto/rand"
	"database/sql"

	pbBas "git.andersenlab.com/Andersen/rpg-new/go-aut-registration-user-grpc.git/protofiles/basic/.basic_server"
	pbTg "git.andersenlab.com/Andersen/rpg-new/go-aut-registration-user-grpc.git/protofiles/telegram/.telegram_server"
	"tgbotapi/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Sender interface {
	Send(message tgbotapi.Chattable) (tgbotapi.Message, error)
}

type handler struct {
	pbBas.UnimplementedGoAuthBasicServer
	pbTg.UnimplementedGoAuthRegistrationUserTelegramServer
	Logger *logrus.Logger
	DB     database.Storage
	Bot    Sender
}

func (h *handler) SendCode(ctx context.Context, request *pbTg.SendCodeRequest) (*pbTg.SendCodeResponse, error) {
	chat, clientID, err := h.DB.FindUserChatId(ctx, h.Logger, request.Phone)
	if err != nil {
		h.Logger.Printf("imposibble to send message to this user: %s", err)
		return nil, err
	}

	RandomCode, _ := rand.Prime(rand.Reader, 18)
	code := RandomCode.String()

	msg := tgbotapi.NewMessage(chat, code)
	if _, err := h.Bot.Send(msg); err != nil {
		h.Logger.Printf("some problems with sending message: %s", err)
		return nil, err
	}
	return &pbTg.SendCodeResponse{
		Id:   clientID,
		Code: code,
	}, nil
}

func (h *handler) SendPhoneNumber(ctx context.Context, request *pbBas.SendPhoneNumberRequest) (*pbBas.SendPhoneNumberResponse, error) {
	flag, err := h.DB.IsPhoneExist(ctx, h.Logger, request.Phone)

	if err == sql.ErrNoRows && !flag {
		if err := h.DB.AddPhone(ctx, h.Logger, request.Phone, request.Id); err != nil {
			h.Logger.Printf("phone was not add: %s", err)
			return &pbBas.SendPhoneNumberResponse{Result: false}, err
		}
		h.Logger.Info("phone added to DB")
		return &pbBas.SendPhoneNumberResponse{Result: true}, nil
	}
	h.Logger.Info("phone was already present at DB")
	return &pbBas.SendPhoneNumberResponse{Result: false}, err
}

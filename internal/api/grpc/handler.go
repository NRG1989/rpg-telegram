package grpc

import (
	"context"
	"crypto/rand"

	pb "tgbotapi/internal/api/grpc/proto"
	"tgbotapi/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Sender interface {
	Send(message tgbotapi.Chattable) (tgbotapi.Message, error)
}

type handler struct {
	pb.UnimplementedGoAutRegistrationUserTelegramServer
	Logger *logrus.Logger
	DB     database.Storage
	Bot    Sender
}

func (h *handler) SendCode(ctx context.Context, request *pb.SendCodeRequest) (*pb.SendCodeResponse, error) {
	chat, err := h.DB.FindUserChatId(ctx, h.Logger, request.Phone)
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
	return &pb.SendCodeResponse{
		Phone: request.Phone,
		Code:  code,
	}, nil
}

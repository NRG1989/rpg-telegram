package grpc

import (
	"net"

	"tgbotapi/internal/database"

	pbBas "git.andersenlab.com/Andersen/rpg-new/go-aut-registration-user-grpc.git/protofiles/basic/.basic_server"
	pbTg "git.andersenlab.com/Andersen/rpg-new/go-aut-registration-user-grpc.git/protofiles/telegram/.telegram_server"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type server struct {
	logger      *logrus.Logger
	GRPCAddress string
	db          database.Storage
	bot         Sender
}

func NewGRPCService(logger *logrus.Logger,
	GRPCAddress string,
	db database.Storage,
	bot Sender) *server {
	return &server{
		logger,
		GRPCAddress,
		db,
		bot,
	}
}

func (s server) Run() error {
	conn, err := net.Listen("tcp", s.GRPCAddress)
	if err != nil {
		s.logger.WithError(err).Errorln("no connection to grpc server")
		return err
	}
	grpcServer := grpc.NewServer()
	pbTg.RegisterGoAuthRegistrationUserTelegramServer(grpcServer, &handler{
		Logger: s.logger,
		Bot:    s.bot,
		DB:     s.db,
	})
	pbBas.RegisterGoAuthBasicServer(grpcServer, &handler{
		Logger: s.logger,
		Bot:    s.bot,
		DB:     s.db,
	})

	if err := grpcServer.Serve(conn); err != nil {
		s.logger.WithError(err).Errorln("connection to grpc error")
		return err
	}
	return nil
}

package grpc

import (
	"net"

	pbTg "go-aut-registration-user-telegram/internal/protofiles/telegram"

	"go-aut-registration-user-telegram/internal/database"

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
	grpcServer := grpc.NewServer()
	pbTg.RegisterGoAuthRegistrationUserTelegramServer(grpcServer, &handler{
		Logger: s.logger,
		Bot:    s.bot,
		DB:     s.db,
	})
	conn, err := net.Listen("tcp", s.GRPCAddress)
	if err != nil {
		s.logger.WithError(err).Errorln("no connection to grpc server")
		return err
	}

	if err := grpcServer.Serve(conn); err != nil {
		s.logger.WithError(err).Errorln("connection to grpc error: %#v", conn.Addr().String())
		return err
	}
	return nil
}

package grpc

import (
	"net"

	pb "tgbotapi/internal/api/grpc/proto"
	"tgbotapi/internal/database"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type server struct {
	logger      *logrus.Logger
	GRPCaddress string
	db          database.Storage
	bot         Sender
}

func NewGRPCService(logger *logrus.Logger,
	GRPCaddress string,
	db database.Storage,
	bot Sender) *server {
	return &server{
		logger,
		GRPCaddress,
		db,
		bot,
	}
}

func (s server) Run() error {
	conn, err := net.Listen("tcp", s.GRPCaddress)
	if err != nil {
		s.logger.WithError(err).Errorln("no connection to grpc server")
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGoAutRegistrationUserTelegramServer(grpcServer, &handler{
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

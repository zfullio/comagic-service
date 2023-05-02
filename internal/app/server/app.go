package server

import (
	"Comagic/internal/config"
	comagicGRPC "Comagic/internal/controllers/api/grpc/v1"
	"Comagic/pb"
	"context"
	"fmt"
	"github.com/nikoksr/notify"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	cfg           *config.ServerConfig
	grpcServer    *grpc.Server
	comagicServer pb.ComagicServiceServer
	logger        *zerolog.Logger
	Notify        notify.Notifier
}

func NewApp(cfg *config.ServerConfig, logger *zerolog.Logger, notify notify.Notifier) *App {
	grpcServer := comagicGRPC.NewServer(*cfg, logger, pb.UnimplementedComagicServiceServer{})

	return &App{
		cfg:           cfg,
		grpcServer:    nil,
		comagicServer: grpcServer,
		logger:        logger,
		Notify:        notify,
	}
}

func (a App) Run(ctx context.Context) (err error) {
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.StartGRPC(a.comagicServer)
	})

	return grp.Wait()
}

func (a App) StartGRPC(server pb.ComagicServiceServer) error {
	a.logger.Info().Msg(fmt.Sprintf("GRPC запущен на %s:%d", a.cfg.GRPC.IP, a.cfg.GRPC.Port))

	err := a.Notify.Send(context.Background(), "Comagic Service", fmt.Sprintf("gRPC запущен на %v:%v", a.cfg.GRPC.IP, a.cfg.GRPC.Port))
	if err != nil {
		a.logger.Fatal().Err(err).Msg("ошибка отправки уведомления")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.GRPC.IP, a.cfg.GRPC.Port))
	if err != nil {
		a.logger.Fatal().Err(err).Msg("failed to create listener")
	}

	a.grpcServer = grpc.NewServer()
	pb.RegisterComagicServiceServer(a.grpcServer, server)

	if err := a.grpcServer.Serve(lis); err != nil {
		a.logger.Fatal().Err(err).Msg("failed to serve")
	}

	return nil
}

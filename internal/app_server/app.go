package app_server

import (
	"Comagic/internal/config"
	comagicGRPC "Comagic/internal/controllers/api/grpc/v1"
	"Comagic/pb"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Entity int

type App struct {
	cfg           *config.ServerConfig
	grpcServer    *grpc.Server
	comagicServer pb.ComagicServiceServer
	logger        *zerolog.Logger
}

func NewApp(cfg *config.ServerConfig, logger *zerolog.Logger) *App {
	grpcServer := comagicGRPC.NewServer(*cfg, logger, pb.UnimplementedComagicServiceServer{})

	return &App{
		cfg:           cfg,
		grpcServer:    nil,
		comagicServer: grpcServer,
		logger:        logger,
	}
}

func (a App) Run(ctx context.Context) (err error) {
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.StartGRPC(a.comagicServer)
	})
	return grp.Wait()

}

func (a App) StartGRPC(server pb.ComagicServiceServer) (err error) {
	a.logger.Info().Msg(fmt.Sprintf(" GRPC запущен на %s:%d", a.cfg.GRPC.IP, a.cfg.GRPC.Port))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.GRPC.IP, a.cfg.GRPC.Port))
	if err != nil {
		log.Fatal("failed to create listener")
	}

	a.grpcServer = grpc.NewServer()
	pb.RegisterComagicServiceServer(a.grpcServer, server)

	if err := a.grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

	return nil
}

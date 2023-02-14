package v1

import (
	"Comagic/internal/config"
	"Comagic/pb"
	"github.com/rs/zerolog"
)

type Server struct {
	cfg    config.ServerConfig
	logger *zerolog.Logger
	pb.UnimplementedComagicServiceServer
}

func NewServer(cfg config.ServerConfig, logger *zerolog.Logger, srv pb.UnimplementedComagicServiceServer) *Server {
	apiLogger := logger.With().Str("api", "grpc").Logger()
	return &Server{
		cfg:                               cfg,
		logger:                            &apiLogger,
		UnimplementedComagicServiceServer: srv,
	}
}

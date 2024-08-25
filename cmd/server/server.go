package server

import (
	"github.com/PaBah/GophKeeper/internal/config"
	pb "github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"go.uber.org/zap"
)

type GrpcServer struct {
	pb.UnimplementedGophKeeperServiceServer
	logger *zap.Logger
	config config.ServerConfig
}

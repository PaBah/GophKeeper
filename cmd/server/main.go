package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"

	"github.com/PaBah/GophKeeper/internal/middlewares"
	"github.com/PaBah/GophKeeper/internal/tls"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	//"github.com/PaBah/GophKeeper/cmd/server"
	"github.com/PaBah/GophKeeper/internal/config"
	pb "github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"github.com/PaBah/GophKeeper/internal/logger"
	"github.com/PaBah/GophKeeper/internal/storage"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	serverConfig := &config.ServerConfig{}
	ParseFlags(serverConfig)

	if err := logger.Initialize(serverConfig.LogsLevel); err != nil {
		fmt.Printf("Logger can not be initialized %s", err)
		return
	}

	var store storage.Repository
	dbStore, _ := storage.NewDBStorage(context.Background(), serverConfig.DatabaseDSN)

	store = &dbStore
	defer dbStore.Close()

	newGRPCServer := NewGrpcServer(serverConfig, store)

	logger.Log().Info("Start gRPC server on", zap.String("address", serverConfig.GRPCAddress))
	interceptors := middlewares.NewGRPCServerMiddleware(serverConfig.Secret)
	authInterceptor := []grpc.UnaryServerInterceptor{interceptors.AuthInterceptor}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		listen, err := net.Listen("tcp", serverConfig.GRPCAddress)
		if err != nil {
			log.Fatal(err)
		}
		var s *grpc.Server

		const (
			certFilePath = "cert.pem" // certFilePath - path to TLS certificate
			keyFilePath  = "key.pem"  // keyFilePath - path to TLS key
		)
		_ = tls.CreateTLSCert(certFilePath, keyFilePath)
		creds, _ := credentials.NewServerTLSFromFile(certFilePath, keyFilePath)
		s = grpc.NewServer(grpc.Creds(creds),
			grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(authInterceptor...)),
			grpc.MaxConcurrentStreams(20),
			grpc.StreamInterceptor(interceptors.StreamAuthInterceptor),
		)

		pb.RegisterGophKeeperServiceServer(s, newGRPCServer)

		if err := s.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}

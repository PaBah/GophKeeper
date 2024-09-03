package middlewares

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/PaBah/GophKeeper/internal/auth"
	"github.com/PaBah/GophKeeper/internal/config"
	"github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"github.com/PaBah/GophKeeper/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCServerMiddleware struct {
	secret string
}

var (
	//nolint:gochecknoglobals
	once sync.Once
	//nolint:gochecknoglobals
	instance *GRPCServerMiddleware
)

// NewGRPCServerMiddleware initializes a MyMiddleware instance with provided logger and secret.
func NewGRPCServerMiddleware(secret string) *GRPCServerMiddleware {
	once.Do(func() {
		instance = &GRPCServerMiddleware{
			secret: secret,
		}
	})

	return instance
}

// AuthInterceptor provides a gRPC unary interceptor for authentication.
func (m GRPCServerMiddleware) AuthInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	switch info.FullMethod {
	case proto.GophKeeperService_SignUp_FullMethodName,
		proto.GophKeeperService_SignIn_FullMethodName:
		logger.Log().Debug("No protected method", zap.String("method", info.FullMethod))

		return handler(ctx, req)
	}

	logger.Log().Debug("Protected method")
	logger.Log().Debug(info.FullMethod)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Log().Debug("couldn't extract metadata from req")

		return nil, fmt.Errorf("%w", status.Error(codes.Internal, "couldn't extract metadata from req"))
	}

	authHeaders, ok := md[config.AUTHORIZATIONHEADER]
	if !ok || len(authHeaders) != 1 {
		logger.Log().Debug("authorization not exists")

		return nil, status.Errorf(codes.Unauthenticated, "authorization not exists")
	}

	token := strings.TrimPrefix(authHeaders[0], config.TOKENPREFIX)
	if token == "" {
		logger.Log().Debug("token empty or not valid")

		return nil, status.Errorf(codes.Unauthenticated, "token empty or not valid")
	}

	if isValid, err := auth.IsValidToken(token, m.secret); err != nil || !isValid {
		logger.Log().Debug("token is not valid")

		return nil, status.Errorf(codes.Unauthenticated, "token empty or not valid")
	}

	userID := auth.GetUserID(token, m.secret)
	if userID == "" {
		logger.Log().Debug("cannot get userID")

		return nil, status.Errorf(codes.Unauthenticated, "token empty or not valid")
	}

	//nolint:staticcheck
	newCtx := context.WithValue(ctx, config.USERIDCONTEXTKEY, userID)

	return handler(newCtx, req)
}

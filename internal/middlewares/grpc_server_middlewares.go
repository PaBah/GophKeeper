package middlewares

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/PaBah/GophKeeper/internal/auth"
	"github.com/PaBah/GophKeeper/internal/config"
	"github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCServerMiddleware struct {
	logger *zap.Logger
	secret string
}

var (
	//nolint:gochecknoglobals
	once sync.Once
	//nolint:gochecknoglobals
	instance *GRPCServerMiddleware
)

// NewGRPCServerMiddleware initializes a MyMiddleware instance with provided logger and secret.
func NewGRPCServerMiddleware(logger *zap.Logger, secret string) *GRPCServerMiddleware {
	once.Do(func() {
		instance = &GRPCServerMiddleware{
			logger: logger,
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
		m.logger.Debug("No protected method", zap.String("method", info.FullMethod))

		return handler(ctx, req)
	}

	m.logger.Debug("Protected method")
	m.logger.Debug(info.FullMethod)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		m.logger.Debug("couldn't extract metadata from req")

		return nil, fmt.Errorf("%w", status.Error(codes.Internal, "couldn't extract metadata from req"))
	}

	authHeaders, ok := md[config.AUTHORIZATIONHEADER]
	if !ok || len(authHeaders) != 1 {
		m.logger.Debug("authorization not exists")

		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "authorization not exists"))
	}

	token := strings.TrimPrefix(authHeaders[0], config.TOKENPREFIX)
	if token == "" {
		m.logger.Debug("token empty or not valid")

		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "token empty or not valid"))
	}

	if isValid, err := auth.IsValidToken(token, m.secret); err != nil || !isValid {
		m.logger.Debug("token is not valid")

		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "token empty or not valid"))
	}
	username, err := auth.GetUserEmailFromToken(token, m.secret)
	if err != nil {
		m.logger.Debug("cannot get username")

		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "token empty or not valid"))
	}

	userID := auth.GetUserID(token, m.secret)
	if userID == "" {
		m.logger.Debug("cannot get userID")

		return nil, fmt.Errorf("%w", status.Error(codes.Unauthenticated, "token empty or not valid"))
	}

	//nolint:staticcheck
	newCtx := context.WithValue(ctx, config.USEREMAILCONTEXTKEY, username)
	//nolint:staticcheck
	newCtx = context.WithValue(newCtx, config.USERIDCONTEXTKEY, userID)

	return handler(newCtx, req)
}

package middlewares

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/PaBah/GophKeeper/internal/auth"
	"github.com/PaBah/GophKeeper/internal/config"
	proto "github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"github.com/PaBah/GophKeeper/internal/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestNewGRPCServerMiddleware(t *testing.T) {

	tests := []struct {
		name      string
		secret    string
		wantToken bool
	}{
		{
			name:      "Valid Secret",
			secret:    "valid_secret",
			wantToken: true,
		},
		{
			name:      "Invalid Secret",
			secret:    "",
			wantToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = NewGRPCServerMiddleware(tt.secret)
		})
	}
}

func TestGRPCServerMiddleware_AuthInterceptor(t *testing.T) {
	validToken, _ := auth.BuildJWTString("1", "1", "valid_secret")
	tests := []struct {
		name                string
		fullMethod          string
		authorizationHeader string
		validToken          bool
		errExpected         bool
	}{
		{
			name:        "No Protected Method",
			fullMethod:  proto.GophKeeperService_SignUp_FullMethodName,
			errExpected: false,
		},
		{
			name:        "No Metadata",
			errExpected: true,
		},
		{
			name:        "No Auth Header",
			fullMethod:  proto.GophKeeperService_CreateCredentials_FullMethodName,
			errExpected: true,
		},
		{
			name:                "Auth Header Invalid",
			fullMethod:          proto.GophKeeperService_CreateCredentials_FullMethodName,
			authorizationHeader: validToken,
			errExpected:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			md := metadata.Pairs(string(config.AUTHORIZATIONHEADER), tt.authorizationHeader)
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			ctx = metadata.NewIncomingContext(ctx, md)

			m := GRPCServerMiddleware{secret: "valid_secret"}
			_, err := m.AuthInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: tt.fullMethod}, func(ctx context.Context, req any) (a any, e error) { return })

			if (err != nil) != tt.errExpected {
				t.Errorf("AuthInterceptor() error = %v, errExpected %v", err, tt.errExpected)
				return
			}

			if tt.validToken {
				token := strings.TrimPrefix(tt.authorizationHeader, string(config.TOKENPREFIX))
				if isValid, _ := auth.IsValidToken(token, m.secret); !isValid {
					t.Errorf("AuthInterceptor() token is not valid")
					return
				}
			}
		})
	}
}

func TestAuthInterceptor(t *testing.T) {
	var ErrContext = errors.New("context error")
	secret := "secret"

	middleware := NewGRPCServerMiddleware(secret)

	userID := "1"
	userName := "testuser"
	validToken, err := auth.BuildJWTString(userID, userName, secret)
	require.NoError(t, err)

	invalidToken, err := auth.BuildJWTString(userID, userName, secret+"fake")
	require.NoError(t, err)

	tests := []struct {
		name        string
		fullMethod  string
		token       string
		secret      string
		wantErr     bool
		isProtected bool
	}{
		{
			name:        "InvalidToken",
			fullMethod:  proto.GophKeeperService_CreateCard_FullMethodName,
			token:       invalidToken,
			secret:      secret,
			wantErr:     true,
			isProtected: true,
		},
		{
			name:        "NoMetadata",
			fullMethod:  proto.GophKeeperService_CreateCard_FullMethodName,
			token:       "",
			secret:      secret,
			wantErr:     true,
			isProtected: true,
		},
		{
			name:        "NoAuthHeader",
			fullMethod:  proto.GophKeeperService_CreateCard_FullMethodName,
			token:       "",
			secret:      secret,
			wantErr:     true,
			isProtected: true,
		},
		{
			name:        "EmptyToken",
			fullMethod:  proto.GophKeeperService_CreateCard_FullMethodName,
			token:       "Bearer ",
			secret:      secret,
			wantErr:     true,
			isProtected: true,
		},
		{
			name:        "UnprotectedMethod",
			fullMethod:  proto.GophKeeperService_SignIn_FullMethodName,
			token:       validToken,
			secret:      secret,
			wantErr:     false,
			isProtected: false,
		},
		{
			name:        "Empty token 2",
			fullMethod:  proto.GophKeeperService_CreateCard_FullMethodName,
			token:       "Bearer ",
			secret:      secret,
			wantErr:     true,
			isProtected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				if tt.isProtected {
					if ctx.Value(config.USERIDCONTEXTKEY) != userID {
						return "", ErrContext
					}
				}

				return "testResponse", nil
			}

			md := metadata.New(map[string]string{
				string(config.AUTHORIZATIONHEADER): tt.token,
			})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			switch tt.name {
			case "NoMetadata":
				ctx = context.Background()
			case "NoAuthHeader":
				md := metadata.New(map[string]string{})
				ctx = metadata.NewIncomingContext(context.Background(), md)
			case "EmptyToken":
				md := metadata.New(map[string]string{
					string(config.AUTHORIZATIONHEADER): "Bearer ",
				})
				ctx = metadata.NewIncomingContext(context.Background(), md)
			case "Empty token 2":
				md := metadata.New(map[string]string{
					string(config.AUTHORIZATIONHEADER): "",
				})
				ctx = metadata.NewIncomingContext(context.Background(), md)
			}

			info := &grpc.UnaryServerInfo{
				FullMethod: tt.fullMethod,
			}

			_, err := middleware.AuthInterceptor(ctx, nil, info, testHandler)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestStreamAuthInterceptor(t *testing.T) {
	secret := "valid_secret"
	middleware := NewGRPCServerMiddleware("valid_secret")
	userID := "1"
	userName := "testuser"
	validToken, _ := auth.BuildJWTString(userID, userName, secret)
	invalidToken, _ := auth.BuildJWTString(userID, userName, "fake")

	tests := []struct {
		name       string
		fullMethod string
		token      string
		wantErr    bool
		wantCode   codes.Code
		setup      func(mockStream *mock.MockMockServerStream)
	}{
		{
			name:       "ValidToken",
			fullMethod: "/package.Service/StreamMethod",
			token:      validToken,
			wantErr:    false,
			setup: func(mockStream *mock.MockMockServerStream) {
				md := metadata.New(map[string]string{
					string(config.AUTHORIZATIONHEADER): string(config.TOKENPREFIX) + validToken,
				})
				ctx := metadata.NewIncomingContext(context.Background(), md)
				mockStream.EXPECT().Context().Return(ctx).AnyTimes()
			},
		},
		{
			name:       "InvalidToken",
			fullMethod: "/package.Service/StreamMethod",
			token:      invalidToken,
			wantErr:    true,
			wantCode:   codes.Unauthenticated,
			setup: func(mockStream *mock.MockMockServerStream) {
				md := metadata.New(map[string]string{
					string(config.AUTHORIZATIONHEADER): invalidToken,
				})
				ctx := metadata.NewIncomingContext(context.Background(), md)
				mockStream.EXPECT().Context().Return(ctx).AnyTimes()
			},
		},
		{
			name:       "MetadataNotFound",
			fullMethod: "/package.Service/StreamMethod",
			token:      "",
			wantErr:    true,
			wantCode:   codes.Internal,
			setup: func(mockStream *mock.MockMockServerStream) {
				ctx := context.Background()
				mockStream.EXPECT().Context().Return(ctx).AnyTimes()
			},
		},
		{
			name:       "AuthorizationHeaderMissing",
			fullMethod: "/package.Service/StreamMethod",
			token:      "",
			wantErr:    true,
			setup: func(mockStream *mock.MockMockServerStream) {
				md := metadata.New(map[string]string{})
				ctx := metadata.NewIncomingContext(context.Background(), md)
				mockStream.EXPECT().Context().Return(ctx).AnyTimes()
			},
			wantCode: codes.Unauthenticated,
		},
		{
			name:       "Empty token",
			fullMethod: "/package.Service/StreamMethod",
			token:      "",
			wantErr:    true,
			wantCode:   codes.Unauthenticated,
			setup: func(mockStream *mock.MockMockServerStream) {
				md := metadata.New(map[string]string{
					string(config.AUTHORIZATIONHEADER): "",
				})
				ctx := metadata.NewIncomingContext(context.Background(), md)
				mockStream.EXPECT().Context().Return(ctx).AnyTimes()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStream := mock.NewMockMockServerStream(ctrl)
			tt.setup(mockStream)

			info := &grpc.StreamServerInfo{
				FullMethod: tt.fullMethod,
			}

			err := middleware.StreamAuthInterceptor(nil, mockStream,
				info, func(srv interface{}, stream grpc.ServerStream) error {
					return nil
				})

			if tt.wantErr {
				require.Error(t, err)
				if err != nil {
					require.Equal(t, tt.wantCode, status.Code(err))
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

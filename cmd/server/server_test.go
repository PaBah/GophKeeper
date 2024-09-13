package main

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/PaBah/GophKeeper/internal/config"
	pb "github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"github.com/PaBah/GophKeeper/internal/mocks"
	"github.com/PaBah/GophKeeper/internal/models"
	"github.com/PaBah/GophKeeper/internal/storage"
	"github.com/PaBah/GophKeeper/internal/utils"
	"go.uber.org/mock/gomock"
)

func TestNewGrpcServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	rm := mocks.NewMockRepository(ctrl)
	tests := []struct {
		name    string
		config  *config.ServerConfig
		storage storage.Repository
		wantErr bool
	}{
		{
			name:    "ValidConfiguration",
			config:  &config.ServerConfig{MinIOAddress: "localhost:9001", MinIOLogin: "admin", MinIOPassword: "password"},
			storage: rm,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grpc := NewGrpcServer(tt.config, tt.storage)
			if (grpc == nil) != tt.wantErr {
				t.Errorf("NewGrpcServer() not exists, wantErr %v", tt.wantErr)
			}
		})
	}
}

func TestSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	srv := &GrpcServer{
		storage:     repo,
		config:      &config.ServerConfig{Secret: "testing secret"},
		syncClients: make(map[string]map[string]pb.GophKeeperService_SubscribeToChangesServer),
		rwMutex:     &sync.RWMutex{},
	}

	tests := []struct {
		name    string
		request *pb.SignInRequest
		mock    func()
		wantErr bool
	}{
		{
			name:    "Successful",
			request: &pb.SignInRequest{Email: "email@example.com", Password: "password"},
			mock: func() {
				password := utils.PasswordHash("password")
				user := &models.User{ID: "user1", Email: "email@example.com", Password: password}
				repo.EXPECT().AuthorizeUser(gomock.Any(), "email@example.com").Return(*user, nil)
			},
			wantErr: false,
		},
		{
			name:    "UnknownEmail",
			request: &pb.SignInRequest{Email: "unknown@example.com", Password: "password"},
			mock: func() {
				repo.EXPECT().AuthorizeUser(gomock.Any(), "unknown@example.com").Return(models.User{}, errors.New("user not found"))
			},
			wantErr: true,
		},
		{
			name:    "IncorrectPassword",
			request: &pb.SignInRequest{Email: "email@example.com", Password: "incorrect"},
			mock: func() {
				password := utils.PasswordHash("password")
				user := &models.User{ID: "user1", Email: "email@example.com", Password: password}
				repo.EXPECT().AuthorizeUser(gomock.Any(), "email@example.com").Return(*user, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := srv.SignIn(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignIn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockRepository(ctrl)
	srv := &GrpcServer{
		storage:     repo,
		config:      &config.ServerConfig{Secret: "testing secret"},
		syncClients: make(map[string]map[string]pb.GophKeeperService_SubscribeToChangesServer),
		rwMutex:     &sync.RWMutex{},
	}

	tests := []struct {
		name        string
		request     *pb.CreateCredentialsRequest
		mock        func()
		wantErr     bool
		expectedRes *pb.CreateCredentialsResponse
	}{
		{
			name:    "SuccessfulCreation",
			request: &pb.CreateCredentialsRequest{ServiceName: "aws", Identity: "myIdentity", Password: "myPassword"},
			mock: func() {
				credentials := models.NewCredentials("aws", "myIdentity", "myPassword")
				createCredentials := models.Credentials{
					ID:          "1",
					ServiceName: credentials.ServiceName,
					Identity:    credentials.Identity,
					Password:    credentials.Password,
					UploadedAt:  time.Now(),
				}
				repo.EXPECT().CreateCredentials(gomock.Any(), gomock.Any()).Return(createCredentials, nil)
			},
			wantErr: false,
			expectedRes: &pb.CreateCredentialsResponse{
				Id:          "1",
				ServiceName: "aws",
				UploadedAt:  time.Now().Format(time.RFC3339),
			},
		},
		{
			name:    "DuplicateEntry",
			request: &pb.CreateCredentialsRequest{ServiceName: "aws", Identity: "myIdentity", Password: "myPassword"},
			mock: func() {
				_ = models.NewCredentials("aws", "myIdentity", "myPassword")
				repo.EXPECT().CreateCredentials(gomock.Any(), gomock.Any()).Return(models.Credentials{}, storage.ErrAlreadyExists)
			},
			wantErr:     true,
			expectedRes: &pb.CreateCredentialsResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			res, err := srv.CreateCredentials(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(res, tt.expectedRes) {
				t.Errorf("CreateCredentials() = %v, want %v", res, tt.expectedRes)
			}
		})
	}
}

func TestGetCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	srv := &GrpcServer{
		storage:     repo,
		config:      &config.ServerConfig{Secret: "testing secret"},
		syncClients: make(map[string]map[string]pb.GophKeeperService_SubscribeToChangesServer),
		rwMutex:     &sync.RWMutex{},
	}

	tests := []struct {
		name        string
		request     *pb.GetCredentialsRequest
		mock        func()
		wantErr     bool
		expectedRes *pb.GetCredentialsResponse
	}{
		{
			name:    "ValidGet",
			request: &pb.GetCredentialsRequest{},
			mock: func() {
				credentials := []models.Credentials{
					models.Credentials{
						ID:          "1",
						ServiceName: "aws",
						Identity:    "myIdentity",
						Password:    "myPassword",
						UploadedAt:  time.Now(),
					},
				}
				repo.EXPECT().GetCredentials(gomock.Any()).Return(credentials, nil)
			},
			wantErr: false,
			expectedRes: &pb.GetCredentialsResponse{
				Credentials: []*pb.GetCredentialsResponse_Credential{
					{
						Id:          "1",
						ServiceName: "aws",
						Identity:    "myIdentity",
						Password:    "myPassword",
						UploadedAt:  time.Now().Format(time.RFC3339),
					},
				},
			},
		},
		{
			name:    "StorageError",
			request: &pb.GetCredentialsRequest{},
			mock: func() {
				repo.EXPECT().GetCredentials(gomock.Any()).Return(nil, errors.New("storage error"))
			},
			wantErr:     true,
			expectedRes: &pb.GetCredentialsResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			res, err := srv.GetCredentials(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(res, tt.expectedRes) {
				t.Errorf("GetCredentials() = %v, want %v", res, tt.expectedRes)
			}
		})
	}
}

func TestUpdateCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	srv := &GrpcServer{
		storage:     repo,
		config:      &config.ServerConfig{Secret: "testing secret"},
		syncClients: make(map[string]map[string]pb.GophKeeperService_SubscribeToChangesServer),
		rwMutex:     &sync.RWMutex{},
	}

	tests := []struct {
		name        string
		request     *pb.UpdateCredentialsRequest
		mock        func()
		wantErr     bool
		expectedRes *pb.UpdateCredentialsResponse
	}{
		{
			name:    "SuccessfulUpdate",
			request: &pb.UpdateCredentialsRequest{Id: "1", ServiceName: "aws", Identity: "newIdentity", Password: "newPassword"},
			mock: func() {
				credentials := &models.Credentials{
					ID:          "1",
					ServiceName: "aws",
					Identity:    "newIdentity",
					Password:    "newPassword",
					UploadedAt:  time.Now(),
				}
				repo.EXPECT().UpdateCredentials(gomock.Any(), gomock.Any()).Return(*credentials, nil)
			},
			wantErr:     false,
			expectedRes: &pb.UpdateCredentialsResponse{Id: "1", ServiceName: "aws", UploadedAt: time.Now().Format(time.RFC3339)},
		},
		{
			name:    "UpdateFail",
			request: &pb.UpdateCredentialsRequest{Id: "notExists", ServiceName: "aws", Identity: "newIdentity", Password: "newPassword"},
			mock: func() {
				repo.EXPECT().UpdateCredentials(gomock.Any(), gomock.Any()).Return(models.Credentials{}, errors.New("update error"))
			},
			wantErr:     true,
			expectedRes: &pb.UpdateCredentialsResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			res, err := srv.UpdateCredentials(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(res, tt.expectedRes) {
				t.Errorf("UpdateCredentials() = %v, want %v", res, tt.expectedRes)
			}
		})
	}
}

func TestDeleteCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	srv := &GrpcServer{
		storage:     repo,
		config:      &config.ServerConfig{Secret: "testing secret"},
		syncClients: make(map[string]map[string]pb.GophKeeperService_SubscribeToChangesServer),
		rwMutex:     &sync.RWMutex{},
	}

	tests := []struct {
		name    string
		request *pb.DeleteCredentialsRequest
		mock    func()
		wantErr bool
	}{
		{
			name:    "SuccessfulDelete",
			request: &pb.DeleteCredentialsRequest{Id: "1"},
			mock: func() {
				repo.EXPECT().DeleteCredentials(gomock.Any(), "1").Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "NonExistentId",
			request: &pb.DeleteCredentialsRequest{Id: "2"},
			mock: func() {
				repo.EXPECT().DeleteCredentials(gomock.Any(), "2").Return(errors.New("credentials not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := srv.DeleteCredentials(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	srv := &GrpcServer{
		storage:     repo,
		config:      &config.ServerConfig{Secret: "testing secret"},
		syncClients: make(map[string]map[string]pb.GophKeeperService_SubscribeToChangesServer),
		rwMutex:     &sync.RWMutex{},
	}
	tests := []struct {
		name       string
		request    *pb.CreateCardRequest
		mock       func()
		wantErr    bool
		wantDigits string
	}{
		{
			name: "ValidCard",
			request: &pb.CreateCardRequest{
				Number:         "9426455762927963",
				ExpirationDate: "02/27",
				HolderName:     "John Doe",
				Cvv:            "123",
			},
			wantDigits: "7963",
			mock: func() {
				card := models.NewCard("9426455762927963", "02/27", "John Doe", "123")
				createdCard := models.Card{
					ID:             "1",
					Number:         "9426455762927963",
					ExpirationDate: "02/20",
					HolderName:     "John Doe",
					CVV:            "123",
					UploadedAt:     time.Now(),
				}
				repo.EXPECT().CreateCard(gomock.Any(), card).Return(createdCard, nil)
			},
			wantErr: false,
		},
		{
			name: "InValidCard",
			request: &pb.CreateCardRequest{
				Number:         "1234567890123456",
				ExpirationDate: "02/27",
				HolderName:     "John Doe",
				Cvv:            "123",
			},
			wantErr: true,
		},
		{
			name: "FailedToSave",
			request: &pb.CreateCardRequest{
				Number:         "9426455762927963",
				ExpirationDate: "02/27",
				HolderName:     "John Doe",
				Cvv:            "123",
			},
			mock: func() {
				_ = models.NewCard("9426455762927963", "02/27", "John Doe", "123")
				repo.EXPECT().CreateCard(gomock.Any(), gomock.Any()).Return(models.Card{}, errors.New("Unknown error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}
			resp, err := srv.CreateCard(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && resp.LastDigits != tt.wantDigits {
				t.Errorf("CreateCard() = %v, want %v", resp, tt.wantDigits)
			}
		})
	}
}

func TestGrpcServer_GetCards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	srv := &GrpcServer{
		storage:     repo,
		config:      &config.ServerConfig{Secret: "testing secret"},
		syncClients: make(map[string]map[string]pb.GophKeeperService_SubscribeToChangesServer),
		rwMutex:     &sync.RWMutex{},
	}

	tests := []struct {
		name        string
		request     *pb.GetCardsRequest
		mock        func()
		wantErr     bool
		expectedRes *pb.GetCardsResponse
	}{
		{
			name:    "ValidGet",
			request: &pb.GetCardsRequest{},
			mock: func() {
				cards := []models.Card{
					models.Card{
						ID:             "1",
						Number:         "aws",
						ExpirationDate: "myIdentity",
						CVV:            "myPassword",
						HolderName:     "myPassword",
						UploadedAt:     time.Now(),
					},
				}
				repo.EXPECT().GetCards(gomock.Any()).Return(cards, nil)
			},
			wantErr: false,
			expectedRes: &pb.GetCardsResponse{
				Cards: []*pb.GetCardsResponse_Card{
					{
						Id:             "1",
						Number:         "aws",
						ExpirationDate: "myIdentity",
						Cvv:            "myPassword",
						HolderName:     "myPassword",
						UploadedAt:     time.Now().Format(time.RFC3339),
					},
				},
			},
		},
		{
			name:    "StorageError",
			request: &pb.GetCardsRequest{},
			mock: func() {
				repo.EXPECT().GetCards(gomock.Any()).Return(nil, errors.New("storage error"))
			},
			wantErr:     true,
			expectedRes: &pb.GetCardsResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			res, err := srv.GetCards(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(res, tt.expectedRes) {
				t.Errorf("GetCredentials() = %v, want %v", res, tt.expectedRes)
			}
		})
	}
}

func TestUpdateCards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	srv := &GrpcServer{
		storage:     repo,
		config:      &config.ServerConfig{Secret: "testing secret"},
		syncClients: make(map[string]map[string]pb.GophKeeperService_SubscribeToChangesServer),
		rwMutex:     &sync.RWMutex{},
	}

	tests := []struct {
		name        string
		request     *pb.UpdateCardRequest
		mock        func()
		wantErr     bool
		expectedRes *pb.UpdateCardResponse
	}{
		{
			name: "SuccessfulUpdate",
			request: &pb.UpdateCardRequest{Id: "1",
				Number:         "9426455762927963",
				ExpirationDate: "myIdentity",
				Cvv:            "myPassword",
				HolderName:     "myPassword"},
			mock: func() {
				card := &models.Card{
					ID:             "1",
					Number:         "9426455762927963",
					ExpirationDate: "myIdentity",
					CVV:            "myPassword",
					HolderName:     "myPassword",
					UploadedAt:     time.Now(),
				}
				repo.EXPECT().UpdateCard(gomock.Any(), gomock.Any()).Return(*card, nil)
			},
			wantErr: false,
			expectedRes: &pb.UpdateCardResponse{LastDigits: "7963",
				ExpirationDate: "myIdentity",
				UploadedAt:     time.Now().Format(time.RFC3339)},
		},
		{
			name:    "UpdateFail",
			request: &pb.UpdateCardRequest{Id: "notExists"},
			mock: func() {
				repo.EXPECT().UpdateCard(gomock.Any(), gomock.Any()).Return(models.Card{}, errors.New("update error"))
			},
			wantErr:     true,
			expectedRes: &pb.UpdateCardResponse{},
		},
		{
			name:    "UpdateFail",
			request: &pb.UpdateCardRequest{Id: "1111111111111", Number: "1111111111111"},
			mock: func() {
			},
			wantErr:     true,
			expectedRes: &pb.UpdateCardResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			res, err := srv.UpdateCard(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(res, tt.expectedRes) {
				t.Errorf("UpdateCredentials() = %v, want %v", res, tt.expectedRes)
			}
		})
	}
}

func TestDeleteCards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	srv := &GrpcServer{
		storage:     repo,
		config:      &config.ServerConfig{Secret: "testing secret"},
		syncClients: make(map[string]map[string]pb.GophKeeperService_SubscribeToChangesServer),
		rwMutex:     &sync.RWMutex{},
	}

	tests := []struct {
		name    string
		request *pb.DeleteCardRequest
		mock    func()
		wantErr bool
	}{
		{
			name:    "SuccessfulDelete",
			request: &pb.DeleteCardRequest{Id: "1"},
			mock: func() {
				repo.EXPECT().DeleteCard(gomock.Any(), "1").Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "NonExistentId",
			request: &pb.DeleteCardRequest{Id: "2"},
			mock: func() {
				repo.EXPECT().DeleteCard(gomock.Any(), "2").Return(errors.New("credentials not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := srv.DeleteCard(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

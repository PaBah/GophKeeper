package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	pb "github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"github.com/PaBah/GophKeeper/internal/mocks"
	"github.com/PaBah/GophKeeper/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestClientService_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name             string
		email            string
		password         string
		mockResponse     pb.SignUpResponse
		mock             func()
		mockError        error
		expectedErrorMsg string
	}{
		{
			name:     "Valid Inputs",
			email:    "test@example.com",
			password: "password",
			mockResponse: pb.SignUpResponse{
				Token: "testToken",
			},
			mock: func() {
				client.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&pb.SignUpResponse{Token: "testToken"}, nil)
			},
			mockError:        nil,
			expectedErrorMsg: "",
		},
		{
			name:         "Empty Inputs",
			email:        "",
			password:     "",
			mockResponse: pb.SignUpResponse{},
			mock: func() {
				client.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&pb.SignUpResponse{}, errors.New("SignUp: invalid inputs"))
			},
			mockError:        fmt.Errorf("invalid inputs"),
			expectedErrorMsg: "SignUp: invalid inputs",
		},
		{
			name:         "Invalid Email",
			email:        "invalidEmail",
			password:     "password",
			mockResponse: pb.SignUpResponse{},
			mock: func() {
				client.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&pb.SignUpResponse{}, errors.New("SignUp: invalid email"))
			},
			mockError:        fmt.Errorf("invalid email"),
			expectedErrorMsg: "SignUp: invalid email",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			err := c.SignUp(tt.email, tt.password)

			if tt.expectedErrorMsg != "" {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tt.expectedErrorMsg)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestNewClientService(t *testing.T) {
	testTable := []struct {
		name          string
		serverAddress string
	}{
		{
			name:          "Valid Server Address",
			serverAddress: "localhost:8080",
		},
		{
			name:          "Empty Server Address",
			serverAddress: "",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			service := NewClientService(tt.serverAddress)
			require.NotNil(t, service)
			if tt.serverAddress == "" {
				require.Empty(t, service.serverAddress)
			} else {
				require.Equal(t, tt.serverAddress, service.serverAddress)
			}
		})
	}
}

func TestClientService_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name             string
		email            string
		password         string
		mockResponse     pb.SignInResponse
		mock             func()
		mockError        error
		expectedErrorMsg string
	}{
		{
			name:     "Correct Inputs",
			email:    "test@example.com",
			password: "password",
			mockResponse: pb.SignInResponse{
				Token: "testToken",
			},
			mock: func() {
				client.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(&pb.SignInResponse{Token: "testToken"}, nil)
			},
			mockError:        nil,
			expectedErrorMsg: "",
		},
		{
			name:         "Empty Inputs",
			email:        "",
			password:     "",
			mockResponse: pb.SignInResponse{},
			mock: func() {
				client.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(&pb.SignInResponse{}, errors.New("SignIn: invalid inputs"))
			},
			mockError:        fmt.Errorf("invalid inputs"),
			expectedErrorMsg: "SignIn: invalid inputs",
		},
		{
			name:         "Invalid Email",
			email:        "invalidEmail",
			password:     "password",
			mockResponse: pb.SignInResponse{},
			mock: func() {
				client.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(&pb.SignInResponse{}, errors.New("SignIn: invalid email"))
			},
			mockError:        fmt.Errorf("invalid email"),
			expectedErrorMsg: "SignIn: invalid email",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			err := c.SignIn(tt.email, tt.password)

			if tt.expectedErrorMsg != "" {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tt.expectedErrorMsg)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestClientService_CreateCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name             string
		serviceName      string
		identity         string
		password         string
		mock             func()
		mockError        error
		expectedErrorMsg string
	}{
		{
			name:        "Valid Inputs",
			serviceName: "testService",
			identity:    "identity",
			password:    "password",
			mock: func() {
				client.EXPECT().CreateCredentials(gomock.Any(), gomock.Any()).Return(&pb.CreateCredentialsResponse{}, nil)
			},
			mockError:        nil,
			expectedErrorMsg: "",
		},
		{
			name:        "Empty Inputs",
			serviceName: "",
			identity:    "",
			password:    "",
			mock: func() {
				client.EXPECT().CreateCredentials(gomock.Any(), gomock.Any()).Return(nil, errors.New("CreateCredentials: invalid inputs"))
			},
			mockError:        fmt.Errorf("invalid inputs"),
			expectedErrorMsg: "CreateCredentials: invalid inputs",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			err := c.CreateCredentials(context.Background(), tt.serviceName, tt.identity, tt.password)

			if tt.expectedErrorMsg != "" {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tt.expectedErrorMsg)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestClientService_GetCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name             string
		mockResp         *pb.GetCredentialsResponse
		mock             func()
		mockError        error
		expectedErrorMsg string
	}{
		{
			name: "Valid response",
			mockResp: &pb.GetCredentialsResponse{
				Credentials: []*pb.GetCredentialsResponse_Credential{
					{
						Id:          "id_1",
						ServiceName: "service_1",
						Identity:    "identity_1",
						Password:    "password_1",
						UploadedAt:  time.Now().Format(time.RFC3339),
					},
					{
						Id:          "id_2",
						ServiceName: "service_2",
						Identity:    "identity_2",
						Password:    "password_2",
						UploadedAt:  time.Now().Format(time.RFC3339),
					},
				},
			},
			mock: func() {
				client.EXPECT().GetCredentials(gomock.Any(), gomock.Any()).Return(&pb.GetCredentialsResponse{
					Credentials: []*pb.GetCredentialsResponse_Credential{
						{
							Id:          "id_1",
							ServiceName: "service_1",
							Identity:    "identity_1",
							Password:    "password_1",
							UploadedAt:  time.Now().Format(time.RFC3339),
						},
						{
							Id:          "id_2",
							ServiceName: "service_2",
							Identity:    "identity_2",
							Password:    "password_2",
							UploadedAt:  time.Now().Format(time.RFC3339),
						},
					},
				}, nil)
			},
			mockError:        nil,
			expectedErrorMsg: "",
		},
		{
			name:     "Empty response",
			mockResp: &pb.GetCredentialsResponse{},
			mock: func() {
				client.EXPECT().GetCredentials(gomock.Any(), gomock.Any()).Return(&pb.GetCredentialsResponse{}, nil)
			},
			mockError:        nil,
			expectedErrorMsg: "",
		},
		{
			name:      "Response Error",
			mockResp:  nil,
			mockError: fmt.Errorf("GetCredentials: test error"),
			mock: func() {
				client.EXPECT().GetCredentials(gomock.Any(), gomock.Any()).Return(nil, errors.New("GetCredentials: test error"))
			},
			expectedErrorMsg: "GetCredentials: test error",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			res, err := c.GetCredentials(context.Background())

			if tt.expectedErrorMsg != "" {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tt.expectedErrorMsg)
			} else {
				require.Nil(t, err)
				require.Equal(t, len(tt.mockResp.Credentials), len(res))
				for i, cred := range res {
					require.Equal(t, tt.mockResp.Credentials[i].Id, cred.ID)
					require.Equal(t, tt.mockResp.Credentials[i].ServiceName, cred.ServiceName)
					require.Equal(t, tt.mockResp.Credentials[i].Identity, cred.Identity)
					require.Equal(t, tt.mockResp.Credentials[i].Password, cred.Password)
					uploadedAt, _ := time.Parse(time.RFC3339, tt.mockResp.Credentials[i].UploadedAt)
					require.Equal(t, uploadedAt, cred.UploadedAt)
				}
			}
		})
	}
}

func TestClientService_UpdateCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name             string
		creds            models.Credentials
		mockResponse     *pb.UpdateCredentialsResponse
		mock             func()
		mockError        error
		expectedErrorMsg string
		expectedCreds    models.Credentials
	}{
		{
			name: "Valid Inputs",
			creds: models.Credentials{
				ID:          "id1",
				ServiceName: "serviceName1",
				Identity:    "identity1",
				Password:    "password1",
			},
			mockResponse: &pb.UpdateCredentialsResponse{
				ServiceName: "updatedServiceName",
				UploadedAt:  time.Now().Add(time.Hour).Format(time.RFC3339),
			},
			mock: func() {
				client.EXPECT().UpdateCredentials(gomock.Any(), gomock.Any()).
					Return(&pb.UpdateCredentialsResponse{
						ServiceName: "updatedServiceName",
						UploadedAt:  time.Now().Add(time.Hour).Format(time.RFC3339),
					}, nil)
			},
			mockError:        nil,
			expectedErrorMsg: "",
			expectedCreds: models.Credentials{
				ID:          "id1",
				ServiceName: "updatedServiceName",
				Identity:    "identity1",
				Password:    "password1",
				UploadedAt:  time.Now().Add(time.Hour),
			},
		},
		{
			name: "Error Updating",
			creds: models.Credentials{
				ID:          "id1",
				ServiceName: "serviceName1",
				Identity:    "identity1",
				Password:    "password1",
			},
			mockResponse: nil,
			mock: func() {
				client.EXPECT().UpdateCredentials(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("UpdateCredentials: test error"))
			},
			mockError:        fmt.Errorf("UpdateCredentials: test error"),
			expectedErrorMsg: "UpdateCredentials: test error",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			updatedCreds, err := c.UpdateCredentials(context.Background(), tt.creds)

			if tt.expectedErrorMsg != "" {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tt.expectedErrorMsg)
			} else {
				require.Nil(t, err)
				require.Equal(t, tt.expectedCreds.ID, updatedCreds.ID)
				require.Equal(t, tt.expectedCreds.ServiceName, updatedCreds.ServiceName)
				require.Equal(t, tt.expectedCreds.Identity, updatedCreds.Identity)
				require.Equal(t, tt.expectedCreds.Password, updatedCreds.Password)
				require.Equal(t, tt.expectedCreds.UploadedAt.Format(time.RFC3339), updatedCreds.UploadedAt.Format(time.RFC3339))
			}
		})
	}
}

func TestClientService_DeleteCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name                 string
		credentialsID        string
		mock                 func()
		mockError            error
		expectedErrorMessage string
	}{
		{
			name:          "Valid Credential Deletion",
			credentialsID: "id1",
			mock: func() {
				client.EXPECT().DeleteCredentials(gomock.Any(), gomock.Any()).
					Return(&pb.DeleteCredentialsResponse{}, nil)
			},
			mockError:            nil,
			expectedErrorMessage: "",
		},
		{
			name:          "Error Deleting Credential",
			credentialsID: "id1",
			mock: func() {
				client.EXPECT().DeleteCredentials(gomock.Any(), gomock.Any()).
					Return(&pb.DeleteCredentialsResponse{}, errors.New("DeleteCredentials: test error"))
			},
			mockError:            fmt.Errorf("DeleteCredentials: test error"),
			expectedErrorMessage: "DeleteCredentials: DeleteCredentials: test error",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			err := c.DeleteCredentials(context.Background(), tt.credentialsID)

			if tt.expectedErrorMessage != "" {
				require.Error(t, err)
				require.Equal(t, tt.expectedErrorMessage, err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestClientService_CreateCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name                 string
		number               string
		expirationDate       string
		holderName           string
		cvv                  string
		mock                 func()
		expectedErrorMessage string
	}{
		{
			name:           "Valid Card Information",
			number:         "1111222233334444",
			expirationDate: "11/24",
			holderName:     "John Doe",
			cvv:            "123",
			mock: func() {
				client.EXPECT().CreateCard(gomock.Any(), gomock.Any()).
					Return(&pb.CreateCardResponse{}, nil)
			},
			expectedErrorMessage: "",
		},
		{
			name:           "Error Creating Card",
			number:         "1111222233334444",
			expirationDate: "11/24",
			holderName:     "John Doe",
			cvv:            "123",
			mock: func() {
				client.EXPECT().CreateCard(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("CreateCard: test error"))
			},
			expectedErrorMessage: "CreateCard: CreateCard: test error",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			err := c.CreateCard(context.Background(), tt.number, tt.expirationDate, tt.holderName, tt.cvv)

			if tt.expectedErrorMessage != "" {
				require.Error(t, err)
				require.Equal(t, tt.expectedErrorMessage, err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestClientService_GetCards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name                 string
		mock                 func()
		expectedErrorMessage string
		expectedCards        []models.Card
	}{
		{
			name: "Valid Inputs",
			expectedCards: []models.Card{
				{
					ID:             "id1",
					Number:         "number1",
					ExpirationDate: "expirationDate1",
					HolderName:     "holderName1",
					CVV:            "cvv1",
					UploadedAt:     time.Now().Add(time.Hour),
				},
				{
					ID:             "id2",
					Number:         "number2",
					ExpirationDate: "expirationDate2",
					HolderName:     "holderName2",
					CVV:            "cvv2",
					UploadedAt:     time.Now().Add(time.Hour * 2),
				},
			},
			mock: func() {
				client.EXPECT().GetCards(gomock.Any(), gomock.Any()).
					Return(&pb.GetCardsResponse{
						Cards: []*pb.GetCardsResponse_Card{
							{
								Id:             "id1",
								Number:         "number1",
								ExpirationDate: "expirationDate1",
								HolderName:     "holderName1",
								Cvv:            "cvv1",
								UploadedAt:     time.Now().Add(time.Hour).Format(time.RFC3339),
							},
							{
								Id:             "id2",
								Number:         "number2",
								ExpirationDate: "expirationDate2",
								HolderName:     "holderName2",
								Cvv:            "cvv2",
								UploadedAt:     time.Now().Add(time.Hour * 2).Format(time.RFC3339),
							},
						},
					}, nil)
			},
			expectedErrorMessage: "",
		},
		{
			name: "Error Getting Cards",
			mock: func() {
				client.EXPECT().GetCards(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("GetCards: test error"))
			},
			expectedErrorMessage: "GetCards: GetCards: test error",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			cards, err := c.GetCards(context.Background())

			if tt.expectedErrorMessage != "" {
				require.Error(t, err)
				require.Equal(t, tt.expectedErrorMessage, err.Error())
			} else {
				require.NoError(t, err)
				require.Len(t, cards, len(tt.expectedCards))
				for i, card := range cards {
					require.Equal(t, tt.expectedCards[i].ID, card.ID)
					require.Equal(t, tt.expectedCards[i].Number, card.Number)
					require.Equal(t, tt.expectedCards[i].ExpirationDate, card.ExpirationDate)
					require.Equal(t, tt.expectedCards[i].HolderName, card.HolderName)
					require.Equal(t, tt.expectedCards[i].CVV, card.CVV)
					require.Equal(t, tt.expectedCards[i].UploadedAt.Format(time.RFC3339), card.UploadedAt.Format(time.RFC3339))
				}
			}
		})
	}
}

func TestClientService_GetFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name                 string
		mock                 func()
		expectedErrorMessage string
		expectedFiles        []models.File
	}{
		{
			name: "Valid Inputs",
			expectedFiles: []models.File{
				{
					Name:       "file1",
					Size:       "1234",
					UploadedAt: time.Now().Add(time.Hour),
				},
				{
					Name:       "file2",
					Size:       "5678",
					UploadedAt: time.Now().Add(time.Hour * 2),
				},
			},
			mock: func() {
				client.EXPECT().GetFiles(gomock.Any(), gomock.Any()).
					Return(&pb.GetFilesResponse{
						Files: []*pb.GetFilesResponse_File{
							{
								Name:       "file1",
								Size:       "1234",
								UploadedAt: time.Now().Add(time.Hour).Format(time.RFC3339),
							},
							{
								Name:       "file2",
								Size:       "5678",
								UploadedAt: time.Now().Add(time.Hour * 2).Format(time.RFC3339),
							},
						},
					}, nil)
			},
			expectedErrorMessage: "",
		},
		{
			name: "Error Getting Files",
			mock: func() {
				client.EXPECT().GetFiles(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("GetFiles: test error"))
			},
			expectedErrorMessage: "GetFiles: GetFiles: test error",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			files, err := c.GetFiles(context.Background())

			if tt.expectedErrorMessage != "" {
				require.Error(t, err)
				require.Equal(t, tt.expectedErrorMessage, err.Error())
			} else {
				require.NoError(t, err)
				require.Len(t, files, len(tt.expectedFiles))
				for i, file := range files {
					require.Equal(t, tt.expectedFiles[i].Name, file.Name)
					require.Equal(t, tt.expectedFiles[i].Size, file.Size)
					require.Equal(t, tt.expectedFiles[i].UploadedAt.Format(time.RFC3339), file.UploadedAt.Format(time.RFC3339))
				}
			}
		})
	}
}

func TestClientService_DeleteFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name                 string
		mock                 func()
		expectedErrorMessage string
	}{
		{
			name: "File Deleted Successfully",
			mock: func() {
				client.EXPECT().DeleteFile(gomock.Any(), gomock.Any()).
					Return(&pb.DeleteFileResponse{}, nil)
			},
			expectedErrorMessage: "",
		},
		{
			name: "Error Deleting File",
			mock: func() {
				client.EXPECT().DeleteFile(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("DeleteFile: test error"))
			},
			expectedErrorMessage: "DeleteFile: DeleteFile: test error",
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			err := c.DeleteFile(context.Background(), "testFile.txt")

			if tt.expectedErrorMessage != "" {
				require.Error(t, err)
				require.Equal(t, tt.expectedErrorMessage, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestClientService_UpdateCards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name                 string
		mock                 func()
		inputCard            models.Card
		expectedErrorMessage string
		expectedCard         models.Card
	}{
		{
			name: "Valid update",
			inputCard: models.Card{
				ID:             "inputCardID",
				Number:         "inputCardNumber",
				ExpirationDate: "inputCardExpirationDate",
				HolderName:     "inputCardHolderName",
				CVV:            "inputCardCVV",
				UploadedAt:     time.Now(),
			},
			expectedCard: models.Card{
				ID:             "inputCardID",
				Number:         "inputCardNumber",
				ExpirationDate: "updatedExpirationDate",
				HolderName:     "inputCardHolderName",
				CVV:            "inputCardCVV",
				UploadedAt:     time.Now(),
			},
			mock: func() {
				client.EXPECT().UpdateCard(gomock.Any(), &pb.UpdateCardRequest{
					Id:             "inputCardID",
					Number:         "inputCardNumber",
					ExpirationDate: "inputCardExpirationDate",
					HolderName:     "inputCardHolderName",
					Cvv:            "inputCardCVV",
				}).Return(&pb.UpdateCardResponse{
					ExpirationDate: "updatedExpirationDate",
					UploadedAt:     time.Now().Format(time.RFC3339),
				}, nil)
			},
		},
		{
			name: "Error updating card",
			inputCard: models.Card{
				ID:             "errorCardID",
				Number:         "errorCardNumber",
				ExpirationDate: "errorCardExpirationDate",
				HolderName:     "errorCardHolderName",
				CVV:            "errorCardCVV",
			},
			mock: func() {
				client.EXPECT().UpdateCard(gomock.Any(), &pb.UpdateCardRequest{
					Id:             "errorCardID",
					Number:         "errorCardNumber",
					ExpirationDate: "errorCardExpirationDate",
					HolderName:     "errorCardHolderName",
					Cvv:            "errorCardCVV",
				}).Return(nil, errors.New("UpdateCards test error"))
			},
			expectedErrorMessage: "UpdateCards: UpdateCards test error",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			card, err := c.UpdateCards(context.Background(), tt.inputCard)

			if tt.expectedErrorMessage != "" {
				require.Error(t, err)
				require.Equal(t, tt.expectedErrorMessage, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedCard.ID, card.ID)
				require.Equal(t, tt.expectedCard.Number, card.Number)
				require.Equal(t, tt.expectedCard.ExpirationDate, card.ExpirationDate)
				require.Equal(t, tt.expectedCard.HolderName, card.HolderName)
				require.Equal(t, tt.expectedCard.CVV, card.CVV)
				require.Equal(t, tt.expectedCard.UploadedAt.Format(time.RFC3339), card.UploadedAt.Format(time.RFC3339))
			}
		})
	}
}
func TestClientService_DeleteCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name                 string
		mock                 func()
		inputCardID          string
		expectedErrorMessage string
	}{
		{
			name:        "Valid delete",
			inputCardID: "validID",
			mock: func() {
				client.EXPECT().DeleteCard(gomock.Any(), &pb.DeleteCardRequest{
					Id: "validID",
				}).Return(&pb.DeleteCardResponse{}, nil)
			},
		},
		{
			name:        "Error deleting card",
			inputCardID: "errorID",
			mock: func() {
				client.EXPECT().DeleteCard(gomock.Any(), &pb.DeleteCardRequest{
					Id: "errorID",
				}).Return(nil, errors.New("DeleteCard test error"))
			},
			expectedErrorMessage: "DeleteCard: DeleteCard test error",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			err := c.DeleteCard(context.Background(), tt.inputCardID)

			if tt.expectedErrorMessage != "" {
				require.Error(t, err)
				require.Equal(t, tt.expectedErrorMessage, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestClientService_UploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tempFile, _ := os.CreateTemp("", "tempfile_")
	defer tempFile.Close()
	_, _ = tempFile.WriteString("test")

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name  string
		mock  func()
		input string
	}{
		{
			name:  "Valid file upload",
			input: tempFile.Name(),
			mock: func() {
				stream := mocks.NewMockGophKeeperService_UploadFileClient(ctrl)
				client.EXPECT().UploadFile(gomock.Any()).Return(stream, nil).Times(1)
				stream.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
				stream.EXPECT().CloseSend().Return(nil).Times(1)
				stream.EXPECT().Recv().Return(&pb.UploadFileResponse{Message: "File uploaded successful."}, nil).Times(1)
			},
		},
		{
			name:  "Invalid file path",
			input: "tempFile.Name()",
			mock: func() {
				stream := mocks.NewMockGophKeeperService_UploadFileClient(ctrl)
				stream.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
			},
		},
		{
			name:  "Can not upload file",
			input: tempFile.Name(),
			mock: func() {
				stream := mocks.NewMockGophKeeperService_UploadFileClient(ctrl)
				client.EXPECT().UploadFile(gomock.Any()).Return(stream, errors.New("can not upload")).Times(1)
				stream.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
			},
		},
		{
			name:  "Can not Send",
			input: tempFile.Name(),
			mock: func() {
				stream := mocks.NewMockGophKeeperService_UploadFileClient(ctrl)
				client.EXPECT().UploadFile(gomock.Any()).Return(stream, nil).Times(1)
				stream.EXPECT().Send(gomock.Any()).Return(errors.New("test err")).AnyTimes()
			},
		},
		{
			name:  "Can not CloseSend",
			input: tempFile.Name(),
			mock: func() {
				stream := mocks.NewMockGophKeeperService_UploadFileClient(ctrl)
				client.EXPECT().UploadFile(gomock.Any()).Return(stream, nil).Times(1)
				stream.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
				stream.EXPECT().CloseSend().Return(errors.New("test")).Times(1)
			},
		},
		{
			name:  "Can not Recv",
			input: tempFile.Name(),
			mock: func() {
				stream := mocks.NewMockGophKeeperService_UploadFileClient(ctrl)
				client.EXPECT().UploadFile(gomock.Any()).Return(stream, nil).Times(1)
				stream.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
				stream.EXPECT().CloseSend().Return(nil).Times(1)
				stream.EXPECT().Recv().Return(&pb.UploadFileResponse{}, errors.New("test")).Times(1)
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			c.UploadFile(context.Background(), tt.input)
		})
	}
}

func TestClientService_DownloadsFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)

	testTable := []struct {
		name                 string
		mock                 func()
		input                string
		expectedErrorMessage string
	}{
		{
			name:  "Valid download",
			input: "validFile",
			mock: func() {
				stream := mocks.NewMockGophKeeperService_DownloadFileClient(ctrl)
				client.EXPECT().DownloadFile(gomock.Any(), &pb.DownloadFileRequest{Name: "validFile"}).Return(stream, nil)
				stream.EXPECT().Recv().Return(&pb.DownloadFileResponse{Data: []byte("test data")}, nil).Times(1)
				stream.EXPECT().Recv().Return(nil, io.EOF)
			},
		},
		{
			name:  "Error when calling DownloadFile",
			input: "nonExistentFile",
			mock: func() {
				client.EXPECT().DownloadFile(gomock.Any(), &pb.DownloadFileRequest{Name: "nonExistentFile"}).Return(nil, errors.New("error in DownloadFile"))
			},
		},
		{
			name:  "Error when receiving chunks",
			input: "validFile",
			mock: func() {
				stream := mocks.NewMockGophKeeperService_DownloadFileClient(ctrl)
				client.EXPECT().DownloadFile(gomock.Any(), &pb.DownloadFileRequest{Name: "validFile"}).Return(stream, nil)
				stream.EXPECT().Recv().Return(nil, errors.New("error in Recv"))
			},
		},
		{
			name:  "Can not Create",
			input: "",
			mock: func() {
				stream := mocks.NewMockGophKeeperService_DownloadFileClient(ctrl)
				client.EXPECT().DownloadFile(gomock.Any(), &pb.DownloadFileRequest{Name: ""}).Return(stream, nil)
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			c := ClientService{client: client}
			err := testChdir(t, func() error {
				c.DownloadsFile(context.Background(), tt.input)
				_, err := os.Stat(tt.input)
				if os.IsNotExist(err) {
					return errors.New("file was not created")
				}
				return nil
			})
			if tt.expectedErrorMessage != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrorMessage)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func testChdir(t *testing.T, fn func() error) (err error) {
	tDir, err := os.MkdirTemp("", "")
	if err != nil {
		return
	}
	defer func() {
		err = os.RemoveAll(tDir)
	}()
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	if err = os.Chdir(tDir); err != nil {
		return
	}
	defer func() {
		_ = os.Chdir(pwd)
	}()
	return fn()
}

func TestClientService_SubscribeToChanges_Valid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGophKeeperServiceClient(ctrl)
	stream := mocks.NewMockGophKeeperService_SubscribeToChangesClient(ctrl)

	client.EXPECT().SubscribeToChanges(gomock.Any(), &pb.SubscribeToChangesRequest{}).Return(stream, nil)

	c := ClientService{client: client}
	_, err := c.SubscribeToChanges(context.Background())
	require.NoError(t, err)
}

func TestClientService_TryToConnect_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockGophKeeperServiceClient(ctrl)
	c := ClientService{client: client, serverAddress: ""}
	isConnected := c.TryToConnect()
	require.True(t, isConnected)
}

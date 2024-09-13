package main

import (
	"testing"

	"github.com/PaBah/GophKeeper/internal/client"
	"github.com/PaBah/GophKeeper/internal/mock"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/mock/gomock"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Valid Email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Email without domain",
			email:   "test@",
			wantErr: true,
		},
		{
			name:    "Email without user",
			email:   "@example.com",
			wantErr: true,
		},
		{
			name:    "Empty Email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "Email with multiple '@'",
			email:   "test@example@.com",
			wantErr: true,
		},
		{
			name:    "Email with invalid characters",
			email:   "test@@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleKeyMsg(t *testing.T) {
	form := NewAuthForm("title", nil)
	model := new(Model)

	tests := []struct {
		name string
		key  tea.KeyType
	}{
		{name: "Enter key", key: tea.KeyEnter},
		{name: "Up key", key: tea.KeyUp},
		{name: "Down key", key: tea.KeyDown},
		{name: "Random key", key: tea.KeyBackspace},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = form.handleKeyMsg(model, tea.KeyMsg{Type: tt.key})
		})
	}
}

func TestHandleEnterKey(t *testing.T) {
	authForm := NewAuthForm("title", nil)
	model := NewModel(Initial)
	var clientMock client.GRPCClientProvider
	ctrl := gomock.NewController(t)
	gm := mock.NewMockGRPCClientProvider(ctrl)
	gm.EXPECT().
		SignIn(gomock.Eq("test@example.com"), gomock.Any()).
		Return(nil).
		AnyTimes()
	gm.EXPECT().
		SignUp(gomock.Eq("test@example.com"), gomock.Any()).
		Return(nil).
		AnyTimes()

	clientMock = gm
	model.clientService = clientMock
	model.signInScreen = authForm
	tests := []struct {
		name   string
		state  State
		email  string
		errNil bool
	}{
		{
			name:   "SignIn with invalid email",
			state:  SignIn,
			email:  "invalid email",
			errNil: false,
		},
		{
			name:   "SignUp with invalid email",
			state:  SignUp,
			email:  "invalid email",
			errNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model.state = tt.state
			authForm.emailInput.SetValue(tt.email)
			authForm.focusIndex = 2
			_, _ = authForm.handleEnterKey(&model)
			if (model.err == nil) != tt.errNil {
				t.Errorf("handleEnterKey() error = %v, wantNil %v", model.err, tt.errNil)
			}
		})
	}
}

func TestUpdateDashboardOnChange(t *testing.T) {
	tests := []struct {
		name        string
		source      menuItem
		expectedStr string
	}{
		{
			name:        "credentials source",
			source:      credentials,
			expectedStr: "GophKeeper: credentials changed, shift → to refresh",
		},
		{
			name:        "cards source",
			source:      cards,
			expectedStr: "GophKeeper: cards changed, shift → to refresh",
		},
		{
			name:        "files source",
			source:      files,
			expectedStr: "GophKeeper: files changed, shift → to refresh",
		},
		{
			name:        "unknown source",
			source:      5, // Assuming 5 is not defined in menuItem
			expectedStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{dashboardScreen: NewDashboardScreen()}
			authForm := NewAuthForm("title", nil)
			m.dashboardScreen.cursor = tt.source
			authForm.updateDashboardOnChange(m, tt.source)
			if m.dashboardScreen.updateMsg != tt.expectedStr {
				t.Errorf("updateDashboardOnChange() got = %v, want = %v", m.dashboardScreen.updateMsg, tt.expectedStr)
			}
		})
	}
}

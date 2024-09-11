package main

import (
	"errors"
	"reflect"
	"testing"

	"github.com/PaBah/GophKeeper/internal/client"
	"github.com/PaBah/GophKeeper/internal/mocks"
	"github.com/PaBah/GophKeeper/internal/models"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/mock/gomock"
)

func TestDashboardScreen_HandleKeyMsg(t *testing.T) {
	tests := []struct {
		name     string
		keyInput tea.KeyType
		wantErr  bool
	}{
		{
			name:     "handle up key msg",
			keyInput: tea.KeyUp,
			wantErr:  false,
		},
		{
			name:     "handle down key msg",
			keyInput: tea.KeyDown,
			wantErr:  false,
		},
		{
			name:     "handle shift+right key msg",
			keyInput: tea.KeyShiftRight,
			wantErr:  false,
		},
		{
			name:     "handle left key",
			keyInput: tea.KeyLeft,
			wantErr:  false,
		},
		{
			name:     "handle unknown key",
			keyInput: tea.KeyShiftLeft,
			wantErr:  false,
		},
		{
			name:     "f1",
			keyInput: tea.KeyF1,
			wantErr:  false,
		},
		{
			name:     "enter",
			keyInput: tea.KeyEnter,
			wantErr:  false,
		},
	}
	var clientMock client.GRPCClientProvider
	ctrl := gomock.NewController(t)
	gm := mocks.NewMockGRPCClientProvider(ctrl)
	gm.EXPECT().
		GetCredentials(gomock.Any()).
		Return([]models.Credentials{models.Credentials{}}, nil).
		AnyTimes()
	gm.EXPECT().
		GetCards(gomock.Any()).
		Return([]models.Card{models.Card{}}, nil).
		AnyTimes()
	gm.EXPECT().
		GetFiles(gomock.Any()).
		Return([]models.File{models.File{}}, nil).
		AnyTimes()

	clientMock = gm
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(Dashboard)
			m.dashboardScreen.loadActual(&m)
			m.clientService = clientMock
			_, err := m.dashboardScreen.handleKeyMsg(&m, tea.KeyMsg{Type: tt.keyInput})

			if (err != nil) != tt.wantErr {
				t.Errorf("DashboardScreen.handleKeyMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestDashboardScreen_LoadActual(t *testing.T) {
	ctrl := gomock.NewController(t)
	gm := mocks.NewMockGRPCClientProvider(ctrl)
	gm.EXPECT().GetCredentials(gomock.Any()).Return([]models.Credentials{models.Credentials{}}, nil).AnyTimes()
	gm.EXPECT().GetCards(gomock.Any()).Return([]models.Card{models.Card{}}, nil).AnyTimes()
	gm.EXPECT().GetFiles(gomock.Any()).Return([]models.File{models.File{}}, nil).AnyTimes()

	tests := []struct {
		name       string
		mockClient client.GRPCClientProvider
		cursor     menuItem
	}{
		{name: "credentials", mockClient: gm, cursor: credentials},
		{name: "cards", mockClient: gm, cursor: cards},
		{name: "files", mockClient: gm, cursor: files},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DashboardScreen{
				cursor: tt.cursor,
			}
			m := NewModel(Dashboard)
			m.clientService = tt.mockClient
			ds.loadActual(&m)
		})
	}
}

func TestDashboardScreen_handleEnterKey(t *testing.T) {
	tests := []struct {
		name            string
		cursor          menuItem
		tableNavigation bool
		tableCursor     int
		wantQuit        bool
	}{
		{
			name:            "switch to table navigation",
			cursor:          credentials,
			tableNavigation: false,
			tableCursor:     0,
			wantQuit:        false,
		},
		{
			name:            "remain in table navigation",
			cursor:          credentials,
			tableNavigation: true,
			tableCursor:     1,
			wantQuit:        false,
		},
		{
			name:            "exit dashboard",
			cursor:          exit,
			tableNavigation: false,
			tableCursor:     0,
			wantQuit:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DashboardScreen{
				cursor:          tt.cursor,
				tableNavigation: tt.tableNavigation,
				tableCursor:     tt.tableCursor,
			}
			m := NewModel(Dashboard)
			_, cmd := ds.handleEnterKey(&m)

			if quitOk := reflect.DeepEqual(cmd, tea.Quit); quitOk != tt.wantQuit {
				t.Errorf("DashboardScreen.handleEnterKey() quit = %v, wantQuit %v", quitOk, tt.wantQuit)
			}
		})
	}
}

func TestDashboardScreen_handleShiftRightKey(t *testing.T) {
	tests := []struct {
		name    string
		errVal  error
		wantErr bool
	}{
		{
			name:    "no error during ShiftRightKey",
			errVal:  nil,
			wantErr: false,
		},
		{
			name:    "error during ShiftRightKey",
			errVal:  errors.New("error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DashboardScreen{}
			m := NewModel(Dashboard)
			m.err = tt.errVal
			ds.handleShiftRightKey(&m)

			if (m.err != nil) != tt.wantErr {
				t.Errorf("DashboardScreen.handleShiftRightKey() error = %v, wantErr %v", m.err, tt.wantErr)
			}
		})
	}
}

func TestDashboardScreen_handleF7Key(t *testing.T) {
	tests := []struct {
		name   string
		cursor menuItem
	}{
		{
			name:   "cursor is on cards",
			cursor: cards,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DashboardScreen{
				cursor: tt.cursor,
				cardsState: []models.Card{
					{
						CVV: "123",
					},
				},
			}
			m := NewModel(Dashboard)
			_, _ = ds.handleF7Key(&m)
			if tt.cursor == cards {
				copied, err := clipboard.ReadAll()
				if err != nil || copied != "123" {
					t.Errorf("Expected %s in clipboard but got %s", "", copied)
				}
			} else {
				copied, err := clipboard.ReadAll()
				if err != nil || copied == "123" {
					t.Errorf("Expected anything except %s in clipboard but got %s", "", copied)
				}
			}
		})
	}
}
func TestDashboardScreen_handleF6Key(t *testing.T) {
	tests := []struct {
		name     string
		cursor   menuItem
		selected int
		want     string
	}{
		{
			name:     "cursor on cards",
			cursor:   cards,
			selected: 0,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DashboardScreen{
				cursor:      tt.cursor,
				tableCursor: tt.selected,
				cardsState: []models.Card{
					{
						HolderName: "",
					},
				},
			}
			m := NewModel(Dashboard)
			_, _ = ds.handleF6Key(&m)
			pasted, _ := clipboard.ReadAll()
			if tt.cursor == cards && pasted != tt.want {
				t.Errorf("Expected %s in clipboard but got %s", tt.want, pasted)
			} else if tt.cursor != cards && pasted == tt.want {
				t.Errorf("Expected anything except %s in clipboard but got %s", tt.want, pasted)
			}
		})
	}
}

func TestDashboardScreen_handleF5Key(t *testing.T) {
	tests := []struct {
		name        string
		cursor      menuItem
		tableCursor int
		want        string
	}{
		{
			name:        "cursor on credentials",
			cursor:      credentials,
			tableCursor: 0,
			want:        "",
		},
		{
			name:        "cursor on cards",
			cursor:      cards,
			tableCursor: 1,
			want:        "",
		},
		{name: "cursor on other", want: "test", cursor: files},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DashboardScreen{
				cursor:      tt.cursor,
				tableCursor: tt.tableCursor,
				credentialsState: []models.Credentials{
					{Password: ""},
					{Password: "password2"},
				},
				cardsState: []models.Card{
					{ExpirationDate: "04/24"},
					{ExpirationDate: ""},
				},
			}
			m := NewModel(Dashboard)
			_, _ = ds.handleF5Key(&m)
			if tt.cursor == credentials || tt.cursor == cards {
				pasted, _ := clipboard.ReadAll()
				if pasted != tt.want {
					t.Errorf("Expected %s in clipboard but got %s", tt.want, pasted)
				}
			} else {
				pasted, _ := clipboard.ReadAll()
				if pasted == tt.want {
					t.Errorf("Expected anything except %s in clipboard but got %s", tt.want, pasted)
				}
			}
		})
	}
}

func TestDashboardScreen_handleF4Key(t *testing.T) {
	tests := []struct {
		name        string
		cursor      menuItem
		tableCursor int
		want        string
	}{
		{
			name:        "cursor on credentials",
			cursor:      credentials,
			tableCursor: 0,
			want:        "",
		},
		{
			name:        "cursor on cards",
			cursor:      cards,
			tableCursor: 1,
			want:        "",
		},
		{name: "cursor on other", want: "test", cursor: files},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DashboardScreen{
				cursor:      tt.cursor,
				tableCursor: tt.tableCursor,
				credentialsState: []models.Credentials{
					{Identity: ""},
					{Identity: "admin"},
				},
				cardsState: []models.Card{
					{Number: "5424003791772490"},
					{Number: ""},
				},
			}
			m := NewModel(Dashboard)
			_, _ = ds.handleF4Key(&m)
			if tt.cursor == credentials || tt.cursor == cards {
				pasted, _ := clipboard.ReadAll()
				if pasted != tt.want {
					t.Errorf("Expected %s in clipboard but got %s", tt.want, pasted)
				}
			} else {
				pasted, _ := clipboard.ReadAll()
				if pasted == tt.want {
					t.Errorf("Expected anything except %s in clipboard but got %s", tt.want, pasted)
				}
			}
		})
	}
}

func TestDashboardScreen_handleF3Key(t *testing.T) {
	tests := []struct {
		name        string
		cursor      menuItem
		tableCursor int
	}{
		{
			name:        "cursor on credentials",
			cursor:      credentials,
			tableCursor: 0,
		},
		{
			name:        "cursor on cards",
			cursor:      cards,
			tableCursor: 1,
		},
		{
			name:        "cursor on files",
			cursor:      files,
			tableCursor: 0,
		},
		{
			name:        "cursor on default",
			cursor:      exit,
			tableCursor: 0,
		},
	}
	var clientMock client.GRPCClientProvider
	ctrl := gomock.NewController(t)
	gm := mocks.NewMockGRPCClientProvider(ctrl)
	gm.EXPECT().GetCredentials(gomock.Any()).Return([]models.Credentials{models.Credentials{}}, nil).AnyTimes()
	gm.EXPECT().GetCards(gomock.Any()).Return([]models.Card{models.Card{Number: "5424003791772490"}}, nil).AnyTimes()
	gm.EXPECT().GetFiles(gomock.Any()).Return([]models.File{models.File{}}, nil).AnyTimes()
	gm.EXPECT().DeleteCredentials(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	gm.EXPECT().DeleteCard(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	gm.EXPECT().DeleteFile(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	clientMock = gm

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DashboardScreen{
				cursor:      tt.cursor,
				tableCursor: tt.tableCursor,
				credentialsState: []models.Credentials{
					{ID: "id1"},
					{ID: "id2"},
				},
				cardsState: []models.Card{
					{ID: "id1"},
					{ID: "id2"},
				},
				filesState: []models.File{
					{Name: "file1"},
					{Name: "file2"},
				},
			}
			m := NewModel(Dashboard)
			m.clientService = clientMock
			ds.handleF3Key(&m)
		})
	}
}

func TestDashboardScreen_handleF2Key(t *testing.T) {
	tests := []struct {
		name     string
		cursor   menuItem
		wantMode State
	}{
		{
			name:     "cursor on credentials",
			cursor:   credentials,
			wantMode: CredentialsForm,
		},
		{
			name:     "cursor on cards",
			cursor:   cards,
			wantMode: CardForm,
		},
		{
			name:     "cursor on files",
			cursor:   files,
			wantMode: Dashboard,
		},
		{
			name:     "cursor on default",
			cursor:   exit,
			wantMode: Dashboard,
		},
	}
	var clientMock client.GRPCClientProvider
	ctrl := gomock.NewController(t)
	gm := mocks.NewMockGRPCClientProvider(ctrl)
	gm.EXPECT().DownloadsFile(gomock.Any(), gomock.Any()).Return().AnyTimes()
	clientMock = gm
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DashboardScreen{
				cursor: tt.cursor,
			}
			m := NewModel(Dashboard)
			m.clientService = clientMock
			m.credentialsScreen.inputs = []textinput.Model{textinput.New(), textinput.New(), textinput.New()}
			m.cardsScreen.inputs = []textinput.Model{textinput.New(), textinput.New(), textinput.New(), textinput.New()}
			ds.credentialsState = []models.Credentials{models.Credentials{ServiceName: "test", Identity: "test", Password: "<PASSWORD>"}}
			ds.cardsState = []models.Card{models.Card{Number: "1111111111111111", HolderName: "test", ExpirationDate: "12/34", CVV: "123"}}
			ds.filesState = []models.File{models.File{Name: "test.test", Size: "12 Kb"}}
			_, _ = ds.handleF2Key(&m)
			if m.state != tt.wantMode {
				t.Errorf("Expected mode %v but got %v", tt.wantMode, m.state)
			}
		})
	}
}

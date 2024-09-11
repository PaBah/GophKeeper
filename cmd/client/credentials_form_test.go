package main

import (
	"context"
	"testing"

	"github.com/PaBah/GophKeeper/internal/models"
	"github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/textinput"
)

type mockClientService struct{}

func (m *mockClientService) CreateCredentials(ctx context.Context, service, identity, password string) error {
	return nil
}

func (m *mockClientService) UpdateCredentials(ctx context.Context, credential models.Credentials) (bool, error) {
	return true, nil
}

func TestCredentialsFormUpdate(t *testing.T) {
	tables := []struct {
		name       string
		inputs     []textinput.Model
		updateID   string
		createMode bool
		focused    credentialsFormInput
		title      string
	}{
		{name: "Case 1", inputs: []textinput.Model{textinput.NewModel()}, updateID: "1234", createMode: false, focused: 0, title: "Test Case 1"},
		{name: "Case 2", inputs: []textinput.Model{textinput.NewModel(), textinput.NewModel()}, updateID: "5678", createMode: true, focused: 1, title: "Test Case 2"},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			m := &Model{
				service:         &mockClientService{},
				dashboardScreen: &DashboardScreen{},
				err:             nil,
				state:           Dashboard,
			}
			form := &CredentialsScreen{
				inputs:     table.inputs,
				updateID:   table.updateID,
				createMode: table.createMode,
				focused:    table.focused,
				title:      table.title,
			}

			msg := bubbletea.KeyMsg{
				Type: bubbletea.KeyTypeRune,
				Key:  "enter",
			}

			_, _ = form.Update(m, msg)
		})
	}
}

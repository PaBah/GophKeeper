package main

import (
	"testing"

	"github.com/PaBah/GophKeeper/internal/client"
	"github.com/PaBah/GophKeeper/internal/mocks"
	"github.com/PaBah/GophKeeper/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/mock/gomock"
)

func TestUpdate(t *testing.T) {
	var clientMock client.GRPCClientProvider
	ctrl := gomock.NewController(t)
	gm := mocks.NewMockGRPCClientProvider(ctrl)
	gm.EXPECT().
		SignIn(gomock.Eq("test@example.com"), gomock.Any()).
		Return(nil).
		AnyTimes()
	gm.EXPECT().
		SignUp(gomock.Eq("test@example.com"), gomock.Any()).
		Return(nil).
		AnyTimes()

	clientMock = gm

	tests := []struct {
		name       string
		inputs     string
		updateID   string
		createMode bool
		msg        tea.Msg
	}{
		{"EnterKeyComplete", "validInput", "123", false, tea.KeyMsg{Type: tea.KeyEnter}},
		{"HandleDownKey", "validInput", "123", true, tea.KeyMsg{Type: tea.KeyDown}},
		{"HandleUpKey", "validInput", "123", true, tea.KeyMsg{Type: tea.KeyUp}},
		{"UpdateInputs", "validInput", "456", true, tea.KeyMsg{Type: tea.KeyCtrlShiftHome}},
		{"CreateMode", "validInput", "123", true, tea.KeyMsg{Type: tea.KeyEnter}},
		{"UpdateMode", "validInput", "456", false, tea.KeyMsg{Type: tea.KeyEnter}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{
				clientService: clientMock,
			}
			cs := NewCardScreen()
			cs.inputs[cardNumber].SetValue(tt.inputs)
			cs.inputs[expiryDate].SetValue(tt.inputs)
			cs.inputs[cardHolder].SetValue(tt.inputs)
			cs.inputs[cvv].SetValue(tt.inputs)
			_, _ = cs.Update(m, tt.msg)
			if m.err != nil {
				t.Errorf("Update returned error: %v", m.err)
			}
		})
	}
}

func TestCardHandleEnterKey(t *testing.T) {
	var clientMock client.GRPCClientProvider
	ctrl := gomock.NewController(t)
	gm := mocks.NewMockGRPCClientProvider(ctrl)
	gm.EXPECT().
		CreateCard(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	gm.EXPECT().
		UpdateCards(gomock.Any(), gomock.Any()).
		Return(models.Card{}, nil).
		AnyTimes()
	gm.EXPECT().
		GetCredentials(gomock.Any()).
		Return([]models.Credentials{models.Credentials{}}, nil).
		AnyTimes()
	clientMock = gm

	tests := []struct {
		name        string
		inputs      []string
		updateID    string
		focused     cardFormInput
		createMode  bool
		errExpected bool
	}{
		{"ValidCardDetailsCreateMode", []string{"4242424242424242", "01/23", "John Doe", "737"}, "", 4, true, false},
		{"ValidCardDetailsUpdateMode", []string{"4242424242424242", "01/23", "John Doe", "737"}, "123", 4, false, false},
		{"InvalidCardNumber", []string{"1234567890123456", "01/23", "John Doe", "737"}, "", 4, true, true},
		{"InvalidExpiryDate", []string{"4242424242424242", "13/23", "John Doe", "737"}, "", 4, true, true},
		{"InvalidHolderName", []string{"4242424242424242", "01/23", "1234567890", "737"}, "", 4, true, true},
		{"InvalidCVV", []string{"4242424242424242", "01/23", "John Doe", "12345"}, "", 4, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(Initial)
			m.clientService = clientMock
			cs := NewCardScreen()
			for i, inp := range tt.inputs {
				cs.inputs[i].SetValue(inp)
			}
			cs.focused = tt.focused
			cs.createMode = tt.createMode
			_, _ = cs.handleEnterKey(&m)
			if (m.err != nil) != tt.errExpected {
				t.Errorf("handleEnterKey returned error: %v, expected error: %v", m.err, tt.errExpected)
			}
		})
	}
}

func TestHandleCardUpKey(t *testing.T) {
	var tests = []struct {
		name       string
		focused    cardFormInput
		updateID   string
		createMode bool
	}{
		{"HandleAboveFirstElement", 0, "", false},
		{"HandleAtSomeElement", 2, "", true},
		{"HandleAtLastElement", totalCardFields, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(Initial)
			cs := NewCardScreen()
			cs.focused = tt.focused
			cs.createMode = tt.createMode

			// Test case when Update is called with `up` KeyMsg
			_, _ = cs.Update(&m, tea.KeyMsg{Type: tea.KeyUp})

			// Check that focused input shifted down
			if cs.focused != tt.focused-1 && tt.focused > 0 {
				t.Errorf("handleUpKey did not decrease focused input. Expected: %d, Actual: %d", tt.focused-1, cs.focused)
			}
		})
	}
}

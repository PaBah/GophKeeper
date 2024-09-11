package main

import (
	"testing"

	"github.com/PaBah/GophKeeper/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
)

func TestUpdate(t *testing.T) {
	clientService := models.CreateTestClientService() // Create a test ClientService for the model
	tests := []struct {
		name       string
		inputs     []textinput.Model
		updateID   string
		createMode bool
		msg        tea.Msg
	}{
		{"EnterKeyComplete", "validInput", "123", false, tea.KeyMsg{"enter"}},
		{"HandleDownKey", "validInput", "123", true, tea.KeyMsg{"down"}},
		{"HandleUpKey", "validInput", "123", true, tea.KeyMsg{"up"}},
		{"UpdateInputs", "validInput", "456", true, tea.KeyMsg{"chars"}},
		{"CreateMode", "validInput", "123", true, tea.KeyMsg{"enter"}},
		{"UpdateMode", "validInput", "456", false, tea.KeyMsg{"enter"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{
				clientService: clientService,
			}
			cs := &CardScreen{
				inputs:     tt.inputs,
				updateID:   tt.updateID,
				createMode: tt.createMode,
				focused:    totalCardFields,
				title:      "Test Screen",
			}
			cs.inputs[cardNumber] = makeCardNumberInput()
			cs.inputs[expiryDate] = makeExpiryDateInput()
			cs.inputs[cardHolder] = makeCardHolderInput()
			cs.inputs[cvv] = makeCVVInput()
			m, _ = cs.Update(m, tt.msg)
			if m.err != nil {
				t.Errorf("Update returned error: %v", m.err)
			}
		})
	}
}

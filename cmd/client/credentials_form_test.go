package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func TestCredentialsFormUpdate(t *testing.T) {
	tables := []struct {
		name       string
		inputs     []textinput.Model
		updateID   string
		createMode bool
		focused    credentialsFormInput
		title      string
	}{
		{name: "Case 1", inputs: []textinput.Model{textinput.New()}, updateID: "1234", createMode: false, focused: 0, title: "Test Case 1"},
		{name: "Case 2", inputs: []textinput.Model{textinput.New(), textinput.New()}, updateID: "5678", createMode: true, focused: 1, title: "Test Case 2"},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			m := NewModel(Dashboard)
			form := &CredentialsScreen{
				inputs:     table.inputs,
				updateID:   table.updateID,
				createMode: table.createMode,
				focused:    table.focused,
				title:      table.title,
			}

			msg := tea.KeyMsg{
				Type: tea.KeyEnter,
			}

			_, _ = form.Update(&m, msg)
		})
	}
}

// Additional tests are to be placed here:
func TestCredentialsScreenHandleKeyMsg(t *testing.T) {
	tables := []struct {
		name        string
		inputs      []textinput.Model
		updateID    string
		createMode  bool
		focused     credentialsFormInput
		title       string
		keyMsg      tea.KeyMsg
		expectedErr string
	}{
		{
			name:        "Test Case: Enter key focused on non-final field",
			inputs:      []textinput.Model{textinput.New(), textinput.New()},
			updateID:    "1234",
			createMode:  false,
			focused:     0,
			title:       "Test Case 1",
			keyMsg:      tea.KeyMsg{Type: tea.KeyEnter},
			expectedErr: "no value: Service name is required",
		},
		{
			name:        "Test Case: Enter key focused on final field, but missing values",
			inputs:      []textinput.Model{textinput.New(), textinput.New(), textinput.New()},
			updateID:    "1234",
			createMode:  true,
			focused:     totalCredentialsFields,
			title:       "Test Case 2",
			keyMsg:      tea.KeyMsg{Type: tea.KeyEnter},
			expectedErr: "no value: Service name is required",
		},
		{
			name:       "Test Case: Down key",
			inputs:     []textinput.Model{textinput.New(), textinput.New()},
			updateID:   "",
			createMode: false,
			focused:    0,
			title:      "Test Case 3",
			keyMsg:     tea.KeyMsg{Type: tea.KeyDown},
		},
		{
			name:       "Test Case: Up key",
			inputs:     []textinput.Model{textinput.New(), textinput.New()},
			updateID:   "",
			createMode: false,
			focused:    1,
			title:      "Test Case 4",
			keyMsg:     tea.KeyMsg{Type: tea.KeyUp},
		},
		{
			name:       "Test Case: CtrlC key",
			inputs:     []textinput.Model{textinput.New(), textinput.New()},
			updateID:   "",
			createMode: false,
			focused:    1,
			title:      "Test Case 4",
			keyMsg:     tea.KeyMsg{Type: tea.KeyCtrlC},
		},
	}
	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			m := NewModel(Initial)
			form := &CredentialsScreen{
				inputs:     table.inputs,
				updateID:   table.updateID,
				createMode: table.createMode,
				focused:    table.focused,
				title:      table.title,
			}
			_, _ = form.handleKeyMsg(&m, table.keyMsg)
			if m.err != nil && m.err.Error() != table.expectedErr {
				t.Errorf("expected error: %q, got: %q", table.expectedErr, m.err.Error())
			}
		})
	}
}

func TestCredentialsScreenValidateAndSubmit(t *testing.T) {
	tables := []struct {
		name       string
		inputs     []string
		updateID   string
		createMode bool
		focused    credentialsFormInput
		title      string
		errMessage string
	}{
		{
			name:       "Empty Service Name",
			inputs:     []string{"", "test", "test"},
			updateID:   "",
			createMode: true,
			focused:    0,
			title:      "Test Case 1",
			errMessage: "no value: Service name is required",
		},
		{
			name:       "Empty Identity",
			inputs:     []string{"test", "", "test"},
			updateID:   "",
			createMode: true,
			focused:    1,
			title:      "Test Case 2",
			errMessage: "no value: Identity is required",
		},
		{
			name:       "Empty Password",
			inputs:     []string{"test", "test", ""},
			updateID:   "",
			createMode: true,
			focused:    2,
			title:      "Test Case 3",
			errMessage: "no value: Password is required",
		},
		{
			name:       "Create Mode Test",
			inputs:     []string{"test", "test", "test"},
			updateID:   "",
			createMode: true,
			focused:    3,
			title:      "Test Case 4",
			errMessage: "",
		},
		{
			name:       "Update Mode Test",
			inputs:     []string{"test", "test", "test"},
			updateID:   "1234",
			createMode: false,
			focused:    3,
			title:      "Test Case 5",
			errMessage: "",
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			m := NewModel(Dashboard)
			form := NewCredentialsScreen()
			form.updateID = table.updateID
			form.createMode = table.createMode
			form.focused = table.focused
			form.title = table.title
			form.inputs[serviceName].SetValue(table.inputs[serviceName])
			form.inputs[identity].SetValue(table.inputs[identity])
			form.inputs[password].SetValue(table.inputs[password])

			_ = form.validateAndSubmit(&m)
			if m.err != nil && m.err.Error() != table.errMessage || m.err == nil && table.errMessage != "" {
				t.Errorf("expected error: %q, got: %q", table.errMessage, m.err.Error())
			}
		})
	}
}

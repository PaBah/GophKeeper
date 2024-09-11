package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInitialForm_UpdateUpdate(t *testing.T) {
	var initialMod = &Model{state: Initial}

	tables := []struct {
		InitForm InitialForm
		Key      string
		Expected Model
	}{
		{InitialForm{AuthThroughSignIn: false, SelectedOption: 0}, "down", Model{state: Initial}},
		{InitialForm{AuthThroughSignIn: false, SelectedOption: 0}, "up", Model{state: Initial}},
		{InitialForm{AuthThroughSignIn: false, SelectedOption: 0}, "enter", Model{state: SignIn}},
		{InitialForm{AuthThroughSignIn: false, SelectedOption: 1}, "enter", Model{state: SignUp}},
	}

	for _, table := range tables {
		t.Run("asd", func(t *testing.T) {
			result, _ := table.InitForm.Update(initialMod, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(table.Key)})

			if result.state != table.Expected.state {
				t.Errorf("Invalid Model state after Update, got: %d, want: %d", int(result.state), int(table.Expected.state))
			}
		})
	}
}

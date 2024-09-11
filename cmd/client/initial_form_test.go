package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdate(t *testing.T) {
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
		{InitialForm{AuthThroughSignIn: false, SelectedOption: 0}, "other", Model{state: Initial}},
	}

	for _, table := range tables {
		result, _ := table.InitForm.Update(initialMod, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(table.Key)})

		if result.state != table.Expected.state {
			t.Errorf("Invalid Model state after Update, got: %s, want: %s", result.state, table.Expected.state)
		}
	}
}

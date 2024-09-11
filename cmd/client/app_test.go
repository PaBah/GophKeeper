package main

import (
	"errors"
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestModel_Init(t *testing.T) {
	tests := []struct {
		name string
		m    Model
		want tea.Cmd
	}{
		{
			name: "InitTest",
			m:    NewModel(Initial),
			want: tea.Batch(textinput.Blink, spinner.New().Tick, nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Init(); got == nil {
				t.Errorf("Model.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_Update(t *testing.T) {
	tests := []struct {
		name string
		m    Model
		msg  tea.Msg
		want tea.Cmd
	}{
		{
			name: "Shift Tab State Update - Initial State",
			m:    NewModel(Initial),
			msg:  tea.KeyMsg{Type: tea.KeyShiftTab},
			want: nil,
		},

		{
			name: "CtrlC Key Message - Initial State",
			m:    NewModel(Initial),
			msg:  tea.KeyMsg{Type: tea.KeyCtrlC},
			want: tea.Quit,
		},
		{
			name: "Esc Key Message - Initial State",
			m:    NewModel(Initial),
			msg:  tea.KeyMsg{Type: tea.KeyEsc},
			want: tea.Quit,
		},
		{
			name: "WindowSizeMsg - Initial State",
			m:    NewModel(Initial),
			msg:  tea.WindowSizeMsg{Width: 50, Height: 20},
			want: nil,
		},
		{
			name: "Shift Tab State Update - SignIn State",
			m:    NewModel(SignIn),
			msg:  tea.KeyMsg{Type: tea.KeyShiftTab},
			want: nil,
		},
		{
			name: "Shift Tab State Update - SignUp State",
			m:    NewModel(SignUp),
			msg:  tea.KeyMsg{Type: tea.KeyShiftTab},
			want: nil,
		},
		{
			name: "Shift Tab State Update - CredentialsForm State",
			m:    NewModel(CredentialsForm),
			msg:  tea.KeyMsg{Type: tea.KeyShiftTab},
			want: nil,
		},
		{
			name: "Shift Tab State Update - CardForm State",
			m:    NewModel(CardForm),
			msg:  tea.KeyMsg{Type: tea.KeyShiftTab},
			want: nil,
		},
		{
			name: "Shift Tab State Update - FileLoad State",
			m:    NewModel(FileLoad),
			msg:  tea.KeyMsg{Type: tea.KeyShiftTab},
			want: nil,
		},
		{
			name: "Shift Tab State Update - Dashboard State",
			m:    NewModel(Dashboard),
			msg:  tea.KeyMsg{Type: tea.KeyShiftTab},
			want: nil,
		},
		{
			name: "WindowSizeMsg - SignIn State",
			m:    NewModel(SignIn),
			msg:  tea.WindowSizeMsg{Width: 50, Height: 20},
			want: nil,
		},
		{
			name: "WindowSizeMsg - CredentialsForm State",
			m:    NewModel(CredentialsForm),
			msg:  tea.WindowSizeMsg{Width: 50, Height: 20},
			want: nil,
		},
		{
			name: "WindowSizeMsg - CardForm State",
			m:    NewModel(CardForm),
			msg:  tea.WindowSizeMsg{Width: 50, Height: 20},
			want: nil,
		},
		{
			name: "WindowSizeMsg - FileLoad State",
			m:    NewModel(FileLoad),
			msg:  tea.WindowSizeMsg{Width: 50, Height: 20},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := tt.m.Update(tt.msg)
			assert.Equal(t, reflect.TypeOf(got), reflect.TypeOf(tt.want), "Model.Update() = %v, want %v", got, tt.want)
		})
	}
}

func TestModel_View(t *testing.T) {
	errored := NewModel(12)
	errored.err = errors.New("test")
	tests := []struct {
		name string
		m    Model
		want string
	}{
		{
			name: "InitialState",
			m:    NewModel(Initial),
			want: NewModel(Initial).View(),
		},
		{
			name: "SignInState",
			m:    NewModel(SignIn),
			want: NewModel(SignIn).View(),
		},
		{
			name: "SignUpState",
			m:    NewModel(SignUp),
			want: NewModel(SignUp).View(),
		},
		{
			name: "DashboardState",
			m:    NewModel(Dashboard),
			want: NewModel(Dashboard).View(),
		},
		{
			name: "CredentialsFormState",
			m:    NewModel(CredentialsForm),
			want: NewModel(CredentialsForm).View(),
		},
		{
			name: "CardFormState",
			m:    NewModel(CardForm),
			want: NewModel(CardForm).View(),
		},
		{
			name: "FileLoadState",
			m:    NewModel(FileLoad),
			want: NewModel(FileLoad).View(),
		},
		{
			name: "ErrorState",
			m:    NewModel(12),
			want: NewModel(12).View(),
		},
		{
			name: "ErrorHeaderState",
			m:    errored,
			want: errored.View(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.View()
			assert.Equal(t, reflect.TypeOf(got), reflect.TypeOf(tt.want), "Model.View() = %v, want %v", got, tt.want)
		})
	}
}

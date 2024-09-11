package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

func TestModel_Init(t *testing.T) {
	tests := []struct {
		name string
		m    Model
		want tea.Cmd
	}{
		{
			name: "InitTest",
			m:    Model{},
			want: tea.Batch(textinput.Blink, spinner.Tick, nil),
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

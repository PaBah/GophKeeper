package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type FilePicker struct {
	filepicker filepicker.Model
	quitting   bool
}

func (fp *FilePicker) Update(m *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	fp.filepicker, cmd = fp.filepicker.Update(msg)

	if didSelect, path := fp.filepicker.DidSelectFile(msg); didSelect {
		log.Println("Selected file:", path)
		m.clientService.UploadFile(context.Background(), path)
		m.dashboardScreen.loadActual(m)
		m.dashboardScreen.content = m.dashboardScreen.drawContent(m)
		m.state = Dashboard
	}

	return m, cmd
}

func (fp *FilePicker) View(m *Model) string {
	if fp.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	s.WriteString(fp.filepicker.View() + "\n")
	return s.String()
}

func NewFilePicker() *FilePicker {
	fp := filepicker.New()
	fp.FileAllowed = true
	fp.Height = 20
	fp.CurrentDirectory, _ = os.UserHomeDir()

	return &FilePicker{
		filepicker: fp,
	}
}

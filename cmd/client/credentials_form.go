package main

import (
	"context"
	"errors"

	"github.com/PaBah/GophKeeper/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CredentialsScreen struct {
	inputs     []textinput.Model
	updateID   string
	createMode bool
	focused    int
	title      string
}

const (
	serviceName            = 0
	identity               = 1
	password               = 2
	totalCredentialsFields = 3
)

func NewCredentialsScreen() *CredentialsScreen {

	serviceName := textinput.New()
	serviceName.Placeholder = "Service name"
	serviceName.Width = 200
	serviceName.Focus()

	identity := textinput.New()
	identity.Placeholder = "Identity (email, login, phone, etc.)"
	identity.Width = 200

	password := textinput.New()
	password.Placeholder = "Password"
	password.Width = 200

	return &CredentialsScreen{
		inputs:  []textinput.Model{serviceName, identity, password},
		focused: 0,
		title:   "Please, enter credentials of service I should keep",
	}
}

// Update обрабатывает действия пользователя
func (form *CredentialsScreen) Update(m *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return form.handleKeyMsg(m, msg)
	default:
		return m, nil
	}
}

func (form *CredentialsScreen) handleKeyMsg(m *Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		return form.handleEnterKey(m)
	case "down":
		return form.handleDownKey(m)
	case "up":
		return form.handleUpKey(m)
	default:
		return form.updateInputs(m, msg)
	}
}

func (form *CredentialsScreen) handleEnterKey(m *Model) (tea.Model, tea.Cmd) {
	if form.focused == totalCredentialsFields {
		m.err = form.validateAndSubmit(m)
		if m.err == nil {
			m.dashboardScreen.loadActual(m)
			m.dashboardScreen.content = m.dashboardScreen.drawContent(m)
			m.state = Dashboard
		}
	} else {
		form.moveFocusForward()
	}
	return form.updateInputs(m, tea.KeyMsg{})
}

// Валидация полей и отправка данных на сервер
func (form *CredentialsScreen) validateAndSubmit(m *Model) error {
	m.err = form.validateFields()
	if m.err != nil {
		return m.err
	}

	if form.createMode {
		return m.clientService.CreateCredentials(
			context.Background(),
			form.inputs[serviceName].Value(),
			form.inputs[identity].Value(),
			form.inputs[password].Value(),
		)
	}

	_, err := m.clientService.UpdateCredentials(
		context.Background(),
		models.Credentials{
			ID:          form.updateID,
			ServiceName: form.inputs[serviceName].Value(),
			Identity:    form.inputs[identity].Value(),
			Password:    form.inputs[password].Value(),
		},
	)
	return err
}

func (form *CredentialsScreen) moveFocusForward() {
	if form.focused < totalCredentialsFields {
		form.inputs[form.focused].Blur()
		form.focused++
		if len(form.inputs)-1 >= form.focused {
			form.inputs[form.focused].Focus()
		}
	}
}

func (form *CredentialsScreen) handleDownKey(m *Model) (tea.Model, tea.Cmd) {
	form.moveFocusForward()
	return form.updateInputs(m, tea.KeyMsg{})
}

func (form *CredentialsScreen) handleUpKey(m *Model) (tea.Model, tea.Cmd) {
	if form.focused > 0 {
		if len(form.inputs)-1 >= form.focused {
			form.inputs[form.focused].Blur()
		}
		form.focused--
		form.inputs[form.focused].Focus()
	}
	return form.updateInputs(m, tea.KeyMsg{})
}

func (form *CredentialsScreen) updateInputs(m *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	for i := range form.inputs {
		var cmd tea.Cmd
		form.inputs[i], cmd = form.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (form *CredentialsScreen) View(m *Model) string {
	var submitButton string

	if form.focused == totalCredentialsFields {
		submitButton = buttonStyle.Render("Submit")
	} else {
		submitButton = buttonBlurredStyle.Render("Submit")
	}
	fields := []string{titleStyle.Render(form.title)}
	for _, field := range form.inputs {
		fields = append(fields, field.View())
	}
	fields = append(fields, submitButton)
	ui := lipgloss.JoinVertical(lipgloss.Left,
		fields...,
	)
	return lipgloss.NewStyle().Align(lipgloss.Center).Padding(1, 2).Render(ui)
}

func (form *CredentialsScreen) validateFields() error {
	if form.inputs[serviceName].Value() == "" {
		return errors.New("no value: Service name is required")
	}

	if form.inputs[identity].Value() == "" {
		return errors.New("no value: Identity is required")
	}

	if form.inputs[password].Value() == "" {
		return errors.New("no value: Password is required")
	}

	return nil
}

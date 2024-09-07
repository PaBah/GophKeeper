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
		switch msg.String() {
		case "enter":
			if msg.String() == "enter" && form.focused == totalCredentialsFields {
				m.err = form.validateFields()
				if m.err == nil {
					if form.createMode {
						m.err = m.clientService.CreateCredentials(
							context.Background(),
							form.inputs[serviceName].Value(),
							form.inputs[identity].Value(),
							form.inputs[password].Value(),
						)
					} else {
						_, m.err = m.clientService.UpdateCredentials(
							context.Background(),
							models.Credentials{
								ID:          form.updateID,
								ServiceName: form.inputs[serviceName].Value(),
								Identity:    form.inputs[identity].Value(),
								Password:    form.inputs[password].Value(),
							},
						)
					}
				}
				if m.err == nil {
					m.dashboardScreen.content = m.dashboardScreen.drawContent(m)
					m.state = Dashboard
				}
			}
			if form.focused < totalCredentialsFields {
				form.inputs[form.focused].Blur()
				form.focused++
				if len(form.inputs)-1 >= form.focused {
					form.inputs[form.focused].Focus()
				}
			}
		case "down":
			if form.focused < totalCredentialsFields {
				form.inputs[form.focused].Blur()
				form.focused++
				if len(form.inputs)-1 >= form.focused {
					form.inputs[form.focused].Focus()
				}
			}
		case "up":
			if form.focused > 0 {
				if len(form.inputs)-1 >= form.focused {
					form.inputs[form.focused].Blur()
				}
				form.focused--
				form.inputs[form.focused].Focus()
			}
		}

		var cmds []tea.Cmd
		for i := range form.inputs {
			var cmd tea.Cmd
			form.inputs[i], cmd = form.inputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

// View отвечает за рендеринг интерфейса
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

// validateFields проверяет все поля на наличие ошибок
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

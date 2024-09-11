package main

import (
	"context"
	"io"
	"log"
	"net/mail"

	pb "github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProcessCallback func(email, password string) error

type AuthForm struct {
	emailInput      textinput.Model
	passwordInput   textinput.Model
	processCallback ProcessCallback
	focusIndex      int
	title           string
}

func NewAuthForm(title string, processCallback ProcessCallback) *AuthForm {
	email := textinput.New()
	email.Placeholder = "Email"
	email.Focus()
	email.CharLimit = 64
	email.Width = 30

	password := textinput.New()
	password.Placeholder = "Password"
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'
	password.Width = 30

	return &AuthForm{
		processCallback: processCallback,
		title:           title,
		emailInput:      email,
		passwordInput:   password,
		focusIndex:      0,
	}
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// Update processes incoming messages and updates the state of the AuthForm.
// It handles key messages and dispatches them to the appropriate handlers based on the type of key.
func (form *AuthForm) Update(m *Model, msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return form.handleKeyMsg(m, msg)
	default:
		return m, nil
	}
}

func (form *AuthForm) handleKeyMsg(m *Model, msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		return form.handleEnterKey(m)
	case tea.KeyUp:
		return form.handleUpKey(m)
	case tea.KeyDown:
		return form.handleDownKey(m)
	default:
		return form.updateInputs(m, msg)
	}
}

func (form *AuthForm) handleEnterKey(m *Model) (*Model, tea.Cmd) {
	if form.focusIndex == 2 {
		m.err = form.validateInputs(m)
		if m.err == nil {
			form.subscribeToChanges(m)
			m.state = Dashboard
			m.dashboardScreen.tableNavigation = false
		}
	} else {
		form.moveFocusForward()
	}
	return form.updateInputs(m, tea.KeyMsg{})
}

func (form *AuthForm) validateInputs(m *Model) error {
	m.err = validateEmail(form.emailInput.Value())
	if m.err != nil {
		return m.err
	}

	if m.state == SignIn {
		return m.clientService.SignIn(form.emailInput.Value(), form.passwordInput.Value())
	}
	return m.clientService.SignUp(form.emailInput.Value(), form.passwordInput.Value())
}

func (form *AuthForm) subscribeToChanges(m *Model) {
	stream, err := m.clientService.SubscribeToChanges(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	go form.handleSubscription(stream, m)
}

func (form *AuthForm) handleSubscription(stream pb.GophKeeperService_SubscribeToChangesClient, m *Model) {
	for {
		var resp pb.SubscribeToChangesResponse
		err := stream.RecvMsg(&resp)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		form.updateDashboardOnChange(m, menuItem(resp.Source))
	}
}

func (form *AuthForm) updateDashboardOnChange(m *Model, source menuItem) {
	switch source {
	case credentials:
		if m.dashboardScreen.cursor == credentials {
			m.dashboardScreen.updateMsg = "GophKeeper: credentials changed, shift → to refresh"
		}
	case cards:
		if m.dashboardScreen.cursor == cards {
			m.dashboardScreen.updateMsg = "GophKeeper: cards changed, shift → to refresh"
		}
	case files:
		if m.dashboardScreen.cursor == files {
			m.dashboardScreen.updateMsg = "GophKeeper: files changed, shift → to refresh"
		}
	default:
		m.dashboardScreen.updateMsg = ""
	}
}

func (form *AuthForm) moveFocusForward() {
	if form.focusIndex < 2 {
		form.focusIndex++
	}
}

func (form *AuthForm) handleUpKey(m *Model) (*Model, tea.Cmd) {
	if form.focusIndex > 0 {
		form.focusIndex--
	}
	return form.updateInputs(m, tea.KeyMsg{})
}

func (form *AuthForm) handleDownKey(m *Model) (*Model, tea.Cmd) {
	if form.focusIndex < 2 {
		form.focusIndex++
	}
	return form.updateInputs(m, tea.KeyMsg{})
}

func (form *AuthForm) updateInputs(m *Model, msg tea.Msg) (*Model, tea.Cmd) {
	form.updateFocus()
	return form.updateFocusedInput(msg)
}

func (form *AuthForm) updateFocus() {
	if form.focusIndex == 0 {
		form.emailInput.Focus()
		form.passwordInput.Blur()
	} else if form.focusIndex == 1 {
		form.passwordInput.Focus()
		form.emailInput.Blur()
	} else {
		form.emailInput.Blur()
		form.passwordInput.Blur()
	}
}

func (form *AuthForm) updateFocusedInput(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	if form.focusIndex == 0 {
		form.emailInput, cmd = form.emailInput.Update(msg)
	} else if form.focusIndex == 1 {
		form.passwordInput, cmd = form.passwordInput.Update(msg)
	}
	return nil, cmd
}

// View renders the AuthForm by composing the title, email input, password input, and submit button.
func (form *AuthForm) View(m *Model) string {
	var submitButton string

	if form.focusIndex == 2 {
		submitButton = buttonStyle.Render("Submit")
	} else {
		submitButton = buttonBlurredStyle.Render("Submit")
	}

	ui := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(form.title),
		form.emailInput.View(),
		form.passwordInput.View(),
		submitButton,
	)

	return lipgloss.NewStyle().Align(lipgloss.Center).Padding(1, 2).Render(ui)
}

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

func (form *AuthForm) Update(m *Model, msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if form.focusIndex == 2 {
				m.err = validateEmail(form.emailInput.Value())
				if m.err == nil {
					if m.state == SignIn {
						m.err = m.clientService.SignIn(form.emailInput.Value(), form.passwordInput.Value())
					} else {
						m.err = m.clientService.SignUp(form.emailInput.Value(), form.passwordInput.Value())
					}
				}
				if m.err == nil {
					stream, err := m.clientService.SubscribeToChanges(context.Background())
					if err != nil {
						log.Fatal(err)
					}
					go func() {
						for {
							var resp pb.SubscribeToChangesResponse
							err := stream.RecvMsg(&resp)
							if err == io.EOF {
								break
							}
							if err != nil {
								log.Fatal(err)
							}
							switch resp.Source {
							case credentials:
								if m.dashboardScreen.cursor == credentials {
									m.dashboardScreen.updateMsg = "GophKeeper: credentials changed, shift → to refresh"
								}
							case cards:
								if m.dashboardScreen.cursor == cards {
									m.dashboardScreen.updateMsg = "GophKeeper: cards changed, shift → to refresh"
								}

							}
						}
					}()
					m.state = Dashboard
					m.dashboardScreen.tableNavigation = false
				}
			} else {
				if form.focusIndex < 2 {
					form.focusIndex++
				}
			}

		case tea.KeyUp:
			if form.focusIndex > 0 {
				form.focusIndex--
			}
		case tea.KeyDown:
			if form.focusIndex < 2 {
				form.focusIndex++
			}
		}
	}

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

	var cmd tea.Cmd
	if form.focusIndex == 0 {
		form.emailInput, cmd = form.emailInput.Update(msg)
	} else if form.focusIndex == 1 {
		form.passwordInput, cmd = form.passwordInput.Update(msg)
	}

	return m, cmd
}

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

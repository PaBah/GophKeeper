package main

import (
	"context"
	"errors"
	"strings"
	"unicode"

	"github.com/PaBah/GophKeeper/internal/models"
	"github.com/PaBah/GophKeeper/internal/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CardScreen struct {
	inputs     []textinput.Model
	updateID   string
	createMode bool
	focused    int
	title      string
}

const (
	cardNumber      = 0
	expiryDate      = 1
	cardHolder      = 2
	cvv             = 3
	totalCardFields = 4
)

func NewCardScreen() *CardScreen {

	cardNumber := textinput.New()
	cardNumber.Placeholder = "1234 5678 9012 345"
	cardNumber.Width = 200
	cardNumber.Focus()

	expiryDate := textinput.New()
	expiryDate.Placeholder = "MM/YY"
	expiryDate.Width = 200

	cardHolder := textinput.New()
	cardHolder.Placeholder = "John Doe"
	cardHolder.Width = 200

	cvv := textinput.New()
	cvv.Placeholder = "123"
	cvv.Width = 200

	return &CardScreen{
		inputs:  []textinput.Model{cardNumber, expiryDate, cardHolder, cvv},
		focused: 0,
		title:   "Please, enter your payment card's details",
	}
}

func (form *CardScreen) Update(m *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if msg.String() == "enter" && form.focused == totalCardFields {
				m.err = form.validateInputs()
				if m.err == nil {
					if form.createMode {
						m.err = m.clientService.CreateCard(
							context.Background(),
							strings.ReplaceAll(form.inputs[cardNumber].Value(), " ", ""),
							form.inputs[expiryDate].Value(),
							form.inputs[cardHolder].Value(),
							form.inputs[cvv].Value(),
						)
					} else {
						_, m.err = m.clientService.UpdateCards(
							context.Background(),
							models.Card{
								ID:             form.updateID,
								Number:         strings.ReplaceAll(form.inputs[cardNumber].Value(), " ", ""),
								ExpirationDate: form.inputs[expiryDate].Value(),
								HolderName:     form.inputs[cardHolder].Value(),
								CVV:            form.inputs[cvv].Value(),
							},
						)
					}
				}
				if m.err == nil {
					m.dashboardScreen.content = m.dashboardScreen.drawContent(m)
					m.state = Dashboard
				}
			}
			if form.focused < totalCardFields {
				form.inputs[form.focused].Blur()
				form.focused++
				if len(form.inputs)-1 >= form.focused {
					form.inputs[form.focused].Focus()
				}
			}
		case "down":
			if form.focused < totalCardFields {
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
func (form *CardScreen) View(m *Model) string {
	var submitButton string

	if form.focused == totalCardFields {
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

func (form *CardScreen) validateInputs() error {
	cardNum := strings.ReplaceAll(form.inputs[cardNumber].Value(), " ", "")
	if len(cardNum) != 16 || !isNumeric(cardNum) || utils.ValidateLuhn(cardNum) != nil {
		return errors.New("incorrect value: Invalid card number")
	}

	if len(form.inputs[expiryDate].Value()) != 5 || !isValidExpiryDate(form.inputs[expiryDate].Value()) {
		return errors.New("incorrect value: Invalid Expiration date")
	}

	if !isAlphabetic(form.inputs[cardHolder].Value()) {
		return errors.New("incorrect value: Invalid Holder name")
	}

	if len(form.inputs[cvv].Value()) != 3 || !isNumeric(form.inputs[cvv].Value()) {
		return errors.New("incorrect value: Invalid CVV")
	}
	return nil
}

func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isAlphabetic(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func isValidExpiryDate(s string) bool {
	if len(s) != 5 || s[2] != '/' {
		return false
	}

	month := s[:2]
	year := s[3:]

	return isNumeric(month) && isNumeric(year) && month >= "01" && month <= "12"
}

package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InitialForm struct {
	AuthThroughSignIn bool
	SelectedOption    int
}

func (form *InitialForm) Update(model *Model, message tea.Msg) (*Model, tea.Cmd) {
	switch message.(type) {
	case tea.KeyMsg:
		switch message.(tea.KeyMsg).String() {
		case "down":
			form.SelectedOption = (form.SelectedOption + 1) % 2
		case "up":
			form.SelectedOption = (form.SelectedOption + 1) % 2
		case "enter":
			if form.SelectedOption == 0 {
				form.AuthThroughSignIn = true
				model.state = SignIn
			} else {
				form.AuthThroughSignIn = false
				model.state = SignUp
			}
		}
	default:
		return model, nil
	}
	return model, nil
}

func (form *InitialForm) View(model Model) string {
	var options [2]string
	options[0] = "SignIn"
	options[1] = "SignUp"

	var view string
	for i, option := range options {
		if i == form.SelectedOption {
			view += buttonStyle.Render(option) + "\n\n"
		} else {
			view += buttonBlurredStyle.Render(option) + "\n\n"
		}
	}

	return model.styles.Base.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			titleStyle.Render("Hi! Welcome to GophKeeper!"),
			titleStyle.Render("Please, choose authorisation method before continue."),
		),
		view,
	)
}

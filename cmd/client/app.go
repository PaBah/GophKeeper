package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	selectedOption int
	err            error
	lg             *lipgloss.Renderer
	state          State
	styles         *Styles
	spinner        spinner.Model
}

const (
	Option1 = iota
	Option2
)

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type State int

const (
	Initial State = iota
	SignIn
	SignUp
	CredentialsForm
	CardForm
	FileLoad
	Dashboard
)

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().Foreground(indigo).Bold(true).Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(indigo).PaddingLeft(1).MarginTop(1)
	s.StatusHeader = lg.NewStyle().Foreground(green).Bold(true)
	s.Highlight = lg.NewStyle().Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.Copy().Foreground(red)
	s.Help = lg.NewStyle().Foreground(lipgloss.Color("240"))
	return &s
}

func NewModel(state State) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := Model{state: state, spinner: s}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	return m
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message.(type) {
	case tea.KeyMsg:
		switch message.(tea.KeyMsg).String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "down":
			m.selectedOption = (m.selectedOption + 1) % 2
		case "up":
			m.selectedOption = (m.selectedOption + 1) % 2
		default:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(message)
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}
	selectedOptionStyle := lipgloss.NewStyle().Background(lipgloss.Color("229"))
	switch m.state {
	case Initial:
		var options [2]string
		options[Option1] = "SignIn"
		options[Option2] = "SignUp"

		var view string
		for i, option := range options {
			if i == m.selectedOption {
				view += selectedOptionStyle.Render(option) + "\n\n"
			} else {
				view += m.styles.Base.Render(option) + "\n\n"
			}
		}

		return m.styles.Base.Render(
			"Hi! Welcome to GophKeeper!\n\n",
			"Please, choose authorisation method before continue.\n\n",
			view,
		)
	default:
		return m.styles.Base.Render(m.spinner.View(), "Oh-oh, something crashed... press esc to quit")
	}
}

func main() {
	m := NewModel(SignUp)
	m.selectedOption = Option1
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

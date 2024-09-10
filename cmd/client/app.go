package main

import (
	"log"
	"os"
	"strings"

	"github.com/PaBah/GophKeeper/internal/client"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	err               error
	width             int
	lg                *lipgloss.Renderer
	state             State
	styles            *Styles
	spinner           spinner.Model
	initialScreen     *InitialForm
	signInScreen      *AuthForm
	signUpScreen      *AuthForm
	dashboardScreen   *DashboardScreen
	clientService     *client.ClientService
	credentialsScreen *CredentialsScreen
	cardsScreen       *CardScreen
	filesScreen       *FilePicker
}

var (
	normalFg      = lipgloss.AdaptiveColor{Light: "235", Dark: "252"}
	indigo        = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	cream         = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}
	fuchsia       = lipgloss.Color("#F780E2")
	help          = lipgloss.Color("240")
	green         = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
	red           = lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}
	blurBg        = lipgloss.Color("235")
	blurText      = lipgloss.Color("240")
	helpSeparator = " · "

	titleStyle = lipgloss.NewStyle().
			Foreground(normalFg).
			Bold(true).
			MarginBottom(1)

	buttonStyle = lipgloss.NewStyle().
			Foreground(cream).
			Background(fuchsia).
			Padding(0, 3).
			MarginTop(1)

	buttonBlurredStyle = lipgloss.NewStyle().
				Foreground(blurText).
				Background(blurBg).
				Padding(0, 3).
				MarginTop(1)
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
	return tea.Batch(m.filesScreen.filepicker.Init(), m.spinner.Tick, textinput.Blink)
}

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().Foreground(indigo).Bold(true).Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(indigo).PaddingLeft(1).MarginTop(1)
	s.StatusHeader = lg.NewStyle().Foreground(green).Bold(true)
	s.ErrorHeaderText = lg.NewStyle().Foreground(red).Bold(true).Padding(0, 1, 0, 2)
	s.Help = lg.NewStyle().Foreground(help)
	return &s
}

func NewModel(state State) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(normalFg)

	m := Model{state: state, spinner: s}
	m.clientService = client.NewClientService(":3200")
	m.clientService.TryToConnect()
	m.initialScreen = &InitialForm{SelectedOption: 0}
	m.signInScreen = NewAuthForm("Please, enter your credentials to SignIn:", func(email, password string) error {
		return m.clientService.SignIn(email, password)
	})
	m.signUpScreen = NewAuthForm("Please, enter your email and create a password to SignUp:", func(email, password string) error {
		return m.clientService.SignUp(email, password)
	})
	m.dashboardScreen = NewDashboardScreen()
	m.credentialsScreen = NewCredentialsScreen()
	m.cardsScreen = NewCardScreen()
	m.filesScreen = NewFilePicker()
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	return m
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message.(type) {
	case tea.WindowSizeMsg:
		m.width = message.(tea.WindowSizeMsg).Width
	case tea.KeyMsg:
		switch message.(tea.KeyMsg).Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyShiftTab:
			switch m.state {
			case SignIn, SignUp:
				m.state = Initial
			case CredentialsForm:
				m.state = Dashboard
			case CardForm:
				m.state = Dashboard
			case FileLoad:
				m.state = Dashboard
			case Dashboard:
				if m.initialScreen.AuthThroughSignIn {
					m.state = SignIn
				} else {
					m.state = SignUp
				}
			}
		}
	}
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.spinner, cmd = m.spinner.Update(message)
	cmds = append(cmds, cmd)
	//if m.state == FileLoad {
	//	_, cmd = m.filesScreen.Update(&m, message)
	//	cmds = append(cmds, cmd)
	//} else {
	//	cmds = append(cmds, m.filesScreen.filepicker.Init())
	//}
	switch m.state {
	case Initial:
		var cmd tea.Cmd
		_, cmd = m.initialScreen.Update(&m, message)
		cmds = append(cmds, cmd)
	case SignIn:
		var cmd tea.Cmd
		_, cmd = m.signInScreen.Update(&m, message)
		cmds = append(cmds, cmd)
	case SignUp:
		var cmd tea.Cmd
		_, cmd = m.signUpScreen.Update(&m, message)
		cmds = append(cmds, cmd)
	case Dashboard:
		var cmd tea.Cmd
		_, cmd = m.dashboardScreen.Update(&m, message)
		cmds = append(cmds, cmd)
	case CredentialsForm:
		var cmd tea.Cmd
		_, cmd = m.credentialsScreen.Update(&m, message)
		cmds = append(cmds, cmd)
	case CardForm:
		var cmd tea.Cmd
		_, cmd = m.cardsScreen.Update(&m, message)
		cmds = append(cmds, cmd)
	case FileLoad:
		var cmd tea.Cmd
		_, cmd = m.filesScreen.Update(&m, message)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := m.appBoundaryView("GophKeeper")
	if m.err != nil {
		header = m.appErrorBoundaryView(m.err.Error())
	}
	body := ""
	footer := ""
	switch m.state {
	case Initial:
		body = m.initialScreen.View(m)
	case SignIn:
		body = m.signInScreen.View(&m)
		footer = strings.Join([]string{"shft+tab back"}, helpSeparator) + " "
	case SignUp:
		body = m.signUpScreen.View(&m)
		footer = strings.Join([]string{"shft+tab back"}, helpSeparator) + " "
	case Dashboard:
		body = m.dashboardScreen.View(&m)
		lines := []string{}
		switch m.dashboardScreen.cursor {
		case credentials:
			lines = []string{
				"shft+tab back", "← menu", "F1 new", "F2 update", "F3 delete", "F4 copy identity", "F5 copy password",
			}
		case cards:
			lines = []string{
				"shft+tab back", "← menu", "F1 new", "F2 update", "F3 delete", "F4 copy number", "F5 copy expiration", "F6 copy holder", "F7 copy CVV",
			}
		case files:
			lines = []string{
				"shft+tab back", "← menu", "F1 upload", "F2 download", "F3 delete",
			}
		}
		if !m.dashboardScreen.tableNavigation {
			lines = []string{"shft+tab back"}
		}
		footer = strings.Join(lines, helpSeparator) + " "
	case CredentialsForm:
		body = m.credentialsScreen.View(&m)
		lines := []string{"shft+tab back"}
		footer = strings.Join(lines, helpSeparator) + " "
	case CardForm:
		body = m.cardsScreen.View(&m)
		lines := []string{"shft+tab back"}
		footer = strings.Join(lines, helpSeparator) + " "
	case FileLoad:
		body = m.filesScreen.View(&m)
		lines := []string{"shft+tab back"}
		footer = strings.Join(lines, helpSeparator) + " "
	default:
		return m.styles.Base.Render(m.spinner.View(), "Oh-oh, something crashed... press ctrl+c to quit")
	}
	footer = m.appFooterView(footer)
	return m.styles.Base.Render(header + "\n" + body + "\n\n" + footer)
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m Model) appFooterView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.Help.Render(text),
		lipgloss.WithWhitespaceChars("\\"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func main() {
	logFile, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()

	// Настроить логгер для записи в файл
	log.SetOutput(logFile)

	m := NewModel(Initial)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
